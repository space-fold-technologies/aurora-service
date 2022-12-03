package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gorilla/websocket"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/plugins"
	"github.com/space-fold-technologies/aurora-service/app/core/providers"
	terminal "golang.org/x/term"
)

var (
	MAX_RETRIES   = 10
	defaultAppDir = "/home/application/current"
)

type DockerServiceConfigurations struct {
	NetworkName        string
	NetworkPrefix      string
	ListenerAddress    string
	AdvertisingAddress string
}

type ImageFetchOrder struct {
	uri      string
	username string
	password string
}

type InternalSpecOrder struct {
	Name     string
	Image    string
	Command  []string
	Mounts   []mount.Mount
	Networks []swarm.NetworkAttachmentConfig
	Ports    []swarm.PortConfig
}

type ServiceSpecOrder struct {
	Name        string
	HostName    string
	Image       string
	Command     []string
	Mounts      []mount.Mount
	Networks    []swarm.NetworkAttachmentConfig
	Ports       []swarm.PortConfig
	Contraints  []string
	Labels      map[string]string
	Envs        []string
	Replicas    *uint64
	MaxReplicas uint64
}

type ImageProgressCallback func(line []byte)

type SwarmOperator struct {
	dkr           *client.Client
	agent         providers.AgentClient
	configuration DockerServiceConfigurations
}

func NewSwarmOperator(agent providers.AgentClient, configuration DockerServiceConfigurations) *SwarmOperator {
	instance := new(SwarmOperator)
	instance.agent = agent
	instance.configuration = configuration
	return instance.setup()
}

func (so *SwarmOperator) setup() *SwarmOperator {
	logger := logging.GetInstance()
	if cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()); err != nil {
		logger.Errorf("failed to open connection to docker daemon: %s", err.Error())
	} else {
		so.dkr = cli
	}
	return so
}

func (so *SwarmOperator) HasWorkers(ctx context.Context) bool {
	if nodes, err := so.dkr.NodeList(ctx, types.NodeListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "role", Value: "worker",
		}),
	}); err != nil {
		return false
	} else {
		return len(nodes) > 0
	}
}

func (so *SwarmOperator) IsLocalNode(host string) bool {
	return host == so.configuration.AdvertisingAddress
}

func (so *SwarmOperator) LocalContainers(ctx context.Context, serviceID string, progress func(lines []byte)) ([]*ContainerDetails, error) {
	details := make([]*ContainerDetails, 0)
	filter := filters.NewArgs()
	filter.Add("label", fmt.Sprintf("com.docker.swarm.service.id=%s", serviceID))
	if containers, err := so.dkr.ContainerList(ctx, types.ContainerListOptions{Filters: filter}); err != nil {
		return nil, err
	} else {
		for _, container := range containers {
			details = append(details, &ContainerDetails{
				ID:            container.ID,
				NodeID:        container.Labels["com.docker.swarm.node.id"],
				TaskID:        container.Labels["com.docker.swarm.task.id"],
				ServiceID:     serviceID,
				IPAddress:     container.NetworkSettings.Networks["aurora-default"].IPAddress,
				AddressFamily: 4,
			})
		}
		return details, nil
	}
}

func (so *SwarmOperator) RemoteContainers(ctx context.Context, serviceID, token string, progress func(lines []byte)) ([]*ContainerDetails, error) {
	retries := 0
	details := make([]*ContainerDetails, 0)
	if tasks, err := so.serviceTasks(ctx, serviceID, &retries); err != nil {
		for _, task := range tasks {
			if node, _, err := so.dkr.NodeInspectWithRaw(ctx, task.NodeID); err != nil {
				return nil, err
			} else if report, err := so.agent.Containers(ctx, serviceID, node.Status.Addr, token); err != nil {
				return nil, err
			} else {
				for _, container := range report.GetContainers() {
					details = append(details, &ContainerDetails{
						ID:            container.GetIdentifier(),
						NodeID:        container.GetNodeIdentifier(),
						TaskID:        container.GetTaskIdentifier(),
						ServiceID:     serviceID,
						IPAddress:     container.GetIpAddress(),
						AddressFamily: uint(container.GetAddressFamily()),
					})
				}
				return details, nil
			}
		}
	}
	return nil, fmt.Errorf("no service tasks found matching service id: %s, retries: %d", serviceID, retries)
}

func (so *SwarmOperator) DeployToManager(ctx context.Context, registry plugins.PluginRegistry, order *providers.DeploymentOrder, reporter *providers.Reporter) error {
	retries := 0
	spec := swarm.ServiceSpec{}
	if digest, err := so.fetchImage(ctx, &ImageFetchOrder{uri: order.ImageURI},
		func(line []byte) { reporter.Progress([]byte(fmt.Sprintf("%s\n\r", string(line)))) },
		&retries); err != nil {
		return err
	} else if err := so.prepare(ctx, registry, order, digest, []string{"node.role==manager"}, &spec); err != nil {
		return err
	} else if resp, err := so.dkr.ServiceCreate(ctx, spec, types.ServiceCreateOptions{}); err != nil {
		return err
	} else if containers, err := so.LocalContainers(ctx, resp.ID, func(lines []byte) { reporter.Progress(lines) }); err != nil {
		return err
	} else {

		report := providers.Report{
			Status:      "DEPLOYED",
			Message:     "successfully deployed service to cluster",
			ServiceID:   resp.ID,
			ImageDigest: digest,
			Instances:   make(map[string]*providers.Instance),
		}
		for _, container := range containers {
			report.Instances[container.ID] = &providers.Instance{
				ID:        container.ID,
				Node:      container.NodeID,
				TaskID:    container.TaskID,
				IP:        container.IPAddress,
				Family:    container.AddressFamily,
				ServiceID: resp.ID,
			}
		}
		reporter.Done(report)
		return nil
	}
}

func (so *SwarmOperator) DeployToWorker(ctx context.Context, registry plugins.PluginRegistry, order *providers.DeploymentOrder, reporter *providers.Reporter) error {
	retries := 0
	spec := swarm.ServiceSpec{}
	if digest, err := so.fetchImage(ctx, &ImageFetchOrder{uri: order.ImageURI},
		func(line []byte) { reporter.Progress([]byte(fmt.Sprintf("%s\n\r", string(line)))) },
		&retries); err != nil {
		return err
	} else if err := so.prepare(ctx, registry, order, digest, []string{"node.role==worker"}, &spec); err != nil {
		return err
	} else if resp, err := so.dkr.ServiceCreate(ctx, spec, types.ServiceCreateOptions{}); err != nil {
		return err
	} else if containers, err := so.RemoteContainers(ctx, resp.ID, order.Token, func(lines []byte) { reporter.Progress(lines) }); err != nil {
		return err
	} else {
		report := providers.Report{
			Status:      "DEPLOYED",
			Message:     "successfully deployed service to cluster",
			ServiceID:   resp.ID,
			ImageDigest: digest,
			Instances:   make(map[string]*providers.Instance),
		}
		for _, container := range containers {
			report.Instances[container.ID] = &providers.Instance{
				ID:        container.ID,
				Node:      container.NodeID,
				TaskID:    container.TaskID,
				IP:        container.IPAddress,
				Family:    container.AddressFamily,
				ServiceID: resp.ID,
			}
		}
		reporter.Done(report)
		return nil
	}
}

func (so *SwarmOperator) ForwardLocalShell(ctx context.Context, ws *websocket.Conn, containerId, shell string) error {
	//nothing special to do here
	return so.localAttach(ctx, ws, shell, containerId)
}

func (so *SwarmOperator) ForwardRemoteShell(ctx context.Context, ws *websocket.Conn, containerId, shell, host string) error {
	if rc, err := so.remoteClient(host); err != nil {
		return err
	} else {
		return so.remoteAttach(ctx, rc, ws, shell, containerId)
	}
}

func (so *SwarmOperator) StopService(ctx context.Context, serviceId string) error {
	return so.dkr.ServiceRemove(ctx, serviceId)
}

func (so *SwarmOperator) ForwardServiceLogs(ctx context.Context, ws *websocket.Conn, properties *providers.LogProperties) error {
	// last battle and well, the only part i really want
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Details:    properties.Details,
		Timestamps: properties.ShowTimestamps,
	}
	if out, err := so.dkr.ServiceLogs(ctx, properties.ServiceID, options); err != nil {
		logging.GetInstance().Error(err)
		return err
	} else {
		defer out.Close()
		defer ws.Close()
		errs := make(chan error, 2)
		quit := make(chan bool)
		handler := &providers.WebSocketWriter{Conn: ws}
		go func() {
			defer close(quit)
			_, err := io.Copy(handler, out)
			if err != nil && err != io.EOF {
				errs <- err
			}
		}()
		<-quit
		close(errs)
		return <-errs
	}
}

func (so *SwarmOperator) InitializeManager(ctx context.Context, listenAddr, advertiseAddr string) (string, error) {
	cmd := fmt.Sprintf("docker swarm init --advertise-addr=%s --listen-addr=%s", advertiseAddr, listenAddr)
	if _, err := exec.Command(cmd).Output(); err != nil {
		return "", err
	} else if details, err := so.dkr.SwarmInspect(ctx); err != nil {
		return "", err
	} else {
		return details.ID, nil
	}

	// TODO: This generates a faulty swarm, the reason why is still not know but, will be uncovered
	// Looking into programmatically creating swarms and joining nodes will be done over the CLI as exec commands as
	//
	//return so.initialize(ctx, advertiseAddr, listenAddr)
}

func (so *SwarmOperator) ParentNode(ctx context.Context) (*providers.ManagerDetails, error) {
	if nodes, err := so.dkr.NodeList(
		ctx,
		types.NodeListOptions{Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "role",
			Value: "manager",
		})}); err != nil {
		return nil, err
	} else {
		for _, node := range nodes {
			if node.ManagerStatus != nil && node.ManagerStatus.Leader {
				return &providers.ManagerDetails{ID: node.ID, Address: node.Status.Addr}, nil
			}
			continue
		}
	}
	return nil, fmt.Errorf("no manager nodes found")
}

func (so *SwarmOperator) JoinNodeToCluster(ctx context.Context, order *providers.JoinOrder) (*providers.NodeDetails, error) {
	retries := 0
	if swarm, err := so.dkr.SwarmInspect(ctx); err != nil {
		return nil, err
	} else if err := so.agent.Join(ctx, &providers.RegisterAgent{
		Token:   swarm.JoinTokens.Worker,
		Name:    order.Name,
		Address: order.CaptainAddress,
	}, order.WorkerAddress, order.Token); err != nil {
		logging.GetInstance().Errorf("Failed to add: %s to cluster : %s", order.WorkerAddress, order.CaptainAddress)
		return nil, err
	} else if node, err := so.worker(ctx, order.WorkerAddress, &retries); err != nil {
		return nil, err
	} else {
		logging.GetInstance().Infof("ADDED NODE WITH ID: %s AND ADDR: %s", node.ID, node.Status.Addr)
		return &providers.NodeDetails{ID: node.ID, IP: node.Status.Addr}, nil
	}
}

func (so *SwarmOperator) RemoveNodeFromCluster(ctx context.Context, order *providers.LeaveOrder) error {
	return so.agent.Leave(ctx, &providers.RemoveAgent{Id: order.NodeID}, order.Address, order.Token)
}

func (so *SwarmOperator) DeployInternalDependency(ctx context.Context, order *providers.DependencyOrder) (string, error) {
	mounts := make([]mount.Mount, 0)
	for source, target := range order.Volumes {
		mounts = append(mounts, mount.Mount{Source: source, Target: target, Type: mount.TypeBind})
	}
	ports := make([]swarm.PortConfig, 0)
	for _, port := range order.Ports {
		ports = append(ports, swarm.PortConfig{TargetPort: uint32(port), PublishedPort: uint32(port)})
	}
	retries := 0
	if _, err := so.fetchImage(ctx, &ImageFetchOrder{uri: order.URI}, func(line []byte) { fmt.Println(string(line)) }, &retries); err != nil {
		return "", err
	}
	specification := so.internalDependencySpec(&InternalSpecOrder{
		Name:     order.Name,
		Image:    order.Image(order.Digest),
		Command:  order.Command,
		Mounts:   mounts,
		Networks: so.networks(ctx),
		Ports:    ports,
	})
	if resp, err := so.dkr.ServiceCreate(ctx, specification, types.ServiceCreateOptions{}); err != nil {
		logging.GetInstance().Error(err)
		return "", err
	} else {
		return resp.ID, nil
	}
}

func (so *SwarmOperator) FetchContainers(ctx context.Context, identifiers []string, status providers.ContainersCallback) error {
	//TODO: Implement call to update container state
	logger := logging.GetInstance()
	state := make(map[string][]*providers.Instance)
	for _, identifier := range identifiers {
		if service, _, err := so.dkr.ServiceInspectWithRaw(ctx, identifier, types.ServiceInspectOptions{}); err != nil {
			return err
		} else if service.Spec.TaskTemplate.Placement.Constraints[0] == "node.role==worker" {
			if containers, err := so.RemoteContainers(ctx, service.ID, "", func(lines []byte) { logger.Infof("%s", string(lines)) }); err != nil {
				return err
			} else {
				so.packInstanceDetails(service.ID, state, containers)
			}
		} else {
			if containers, err := so.LocalContainers(ctx, service.ID, func(lines []byte) { logger.Infof("%s", string(lines)) }); err != nil {
				return err
			} else {
				so.packInstanceDetails(service.ID, state, containers)
			}
		}
	}
	return nil
}

func (so *SwarmOperator) packInstanceDetails(serviceID string, state map[string][]*providers.Instance, containers []*ContainerDetails) {
	state[serviceID] = make([]*providers.Instance, 0)
	for _, container := range containers {
		state[serviceID] = append(state[serviceID], &providers.Instance{
			ID:        container.ID,
			Node:      container.NodeID,
			TaskID:    container.TaskID,
			IP:        container.IPAddress,
			Family:    container.AddressFamily,
			ServiceID: serviceID,
		})
	}
}

/*
func (so *SwarmOperator) initialize(ctx context.Context, advertiseAddr, listenAddr string) (string, error) {
	return so.dkr.SwarmInit(ctx, swarm.InitRequest{
		ListenAddr:      listenAddr,
		AdvertiseAddr:   advertiseAddr,
		ForceNewCluster: false,
		Spec: swarm.Spec{Orchestration: swarm.OrchestrationConfig{TaskHistoryRetentionLimit: func(n int64) *int64 { return &n }(5)},
			Raft: swarm.RaftConfig{
				SnapshotInterval: 10000,
				KeepOldSnapshots: func(n uint64) *uint64 { return &n }(0),
			},
			Dispatcher: swarm.DispatcherConfig{
				HeartbeatPeriod: 5 * time.Second,
			},
			EncryptionConfig: swarm.EncryptionConfig{
				AutoLockManagers: false,
			},
			CAConfig: swarm.CAConfig{NodeCertExpiry: 90 * 24 * time.Hour}},
		AutoLockManagers: false,
		Availability:     swarm.NodeAvailabilityActive,
		DefaultAddrPool:  []string{"10.0.0.0/8"},
		SubnetSize:       24,
	})
}
*/

// utilities
func (so *SwarmOperator) networks(ctx context.Context) []swarm.NetworkAttachmentConfig {
	networks := make([]swarm.NetworkAttachmentConfig, 0)
	if ok, err := so.hasNetwork(ctx); err != nil {
		logging.GetInstance().Error(err)
		return []swarm.NetworkAttachmentConfig{}
	} else if !ok {
		if networkId, err := so.createNetwork(ctx); err != nil {
			logging.GetInstance().Error(err)
			return []swarm.NetworkAttachmentConfig{}
		} else {
			logging.GetInstance().Infof("created default network with id :%s", networkId)
		}
	}
	networks = append(networks, swarm.NetworkAttachmentConfig{
		Target: fmt.Sprintf("%s-%s", so.configuration.NetworkPrefix, so.configuration.NetworkName),
	})
	return networks
}

func (so *SwarmOperator) portMapper(ctx context.Context, portSet nat.PortSet) ([]swarm.PortConfig, error) {
	logger := logging.GetInstance()
	if services, err := so.dkr.ServiceList(ctx, types.ServiceListOptions{Status: true}); err != nil {
		logger.Errorf("failed to get a list of services: %s", err.Error())
		return []swarm.PortConfig{}, err
	} else {
		publishPorts := mapset.NewSet[uint32]()
		for _, service := range services {
			for _, port := range service.Endpoint.Ports {
				publishPorts.Add(port.PublishedPort)
			}
		}
		existing := publishPorts.ToSlice()
		for j := 1; j < len(existing); j++ {
			if existing[0] < existing[j] {
				existing[0] = existing[j]
			}
		}
		ports := make([]swarm.PortConfig, 0)
		for port := range portSet {
			ports = append(ports, swarm.PortConfig{TargetPort: uint32(port.Int()), Protocol: swarm.PortConfigProtocol(port.Proto())})
		}
		return ports, nil
	}
}

func (so *SwarmOperator) createNetwork(ctx context.Context) (string, error) {
	if response, err := so.dkr.NetworkCreate(
		ctx,
		fmt.Sprintf(
			"%s-%s",
			so.configuration.NetworkPrefix,
			so.configuration.NetworkName), types.NetworkCreate{
			Driver: "overlay",
			Scope:  "swarm",
			Labels: map[string]string{
				so.configuration.NetworkPrefix: "true",
			},
		}); err != nil {
		return "", err
	} else {
		return response.ID, nil
	}
}

func (so *SwarmOperator) hasNetwork(ctx context.Context) (bool, error) {
	if networks, err := so.dkr.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "name", Value: fmt.Sprintf(
				"%s-%s",
				so.configuration.NetworkPrefix, so.configuration.NetworkName),
		}),
	}); err != nil {
		return false, err
	} else if len(networks) > 0 {
		if networks[0].Driver == "overlay" &&
			networks[0].Scope == "swarm" &&
			networks[0].Labels[so.configuration.NetworkPrefix] == "true" {
			logging.GetInstance().Info("network already exists and is configured correctly")
			return true, nil
		}
	}
	return false, nil
}

func (so *SwarmOperator) internalDependencySpec(order *InternalSpecOrder) swarm.ServiceSpec {
	return swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name:   order.Name,
			Labels: map[string]string{},
		},
		TaskTemplate: swarm.TaskSpec{
			Placement: &swarm.Placement{
				Constraints: []string{"node.role==manager"},
				MaxReplicas: 1,
			},
			ContainerSpec: &swarm.ContainerSpec{
				Hostname: order.Name,
				Image:    order.Image,
				Env:      []string{},
				Command:  order.Command,
				Mounts:   order.Mounts,
				TTY:      true,
				Labels:   map[string]string{},
				Hosts:    []string{}, //<< Some nonsese like this //"host.docker.internal"
			},
			RestartPolicy: &swarm.RestartPolicy{Condition: swarm.RestartPolicyConditionAny},
			Networks:      order.Networks,
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: order.Ports,
		},
	}
}

// internals
func (so *SwarmOperator) fetchImage(ctx context.Context, order *ImageFetchOrder, callback ImageProgressCallback, retries *int) (string, error) {
	if stream, err := so.dkr.ImagePull(ctx, order.uri, types.ImagePullOptions{
		RegistryAuth: so.encodeCredentials(order.username, order.password),
	}); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "timeout") && *retries < MAX_RETRIES {
			time.Sleep(5 * time.Second)
			*retries++
			return so.fetchImage(ctx, order, callback, retries)
		}
	} else {
		defer stream.Close()
		var digest string
		reader := bufio.NewReader(stream)
		for {
			if line, _, err := reader.ReadLine(); err != nil {
				if err == io.EOF {
					break
				} else {
					logging.GetInstance().Errorf("Failed to read line: %v", err)
					return "", err
				}
			} else {
				var chunk ProgressChunk
				if err = json.Unmarshal(line, &chunk); err != nil {
					return "", err
				} else if strings.HasPrefix(chunk.Status, "Digest:") {
					digest = strings.TrimSpace(strings.ReplaceAll(chunk.Status, "Digest: ", ""))
				}
				// intercept to get image digest
				callback(line)
			}
		}
		return digest, nil
	}
	return "", nil
}

func (so *SwarmOperator) imageSpecification(ctx context.Context, url, digest string, retries *int) (types.ImageInspect, error) {
	filter := filters.NewArgs()
	filter.Add("reference", url)
	if images, err := so.dkr.ImageList(ctx, types.ImageListOptions{All: true, Filters: filter}); err != nil {
		return types.ImageInspect{}, err
	} else if len(images) == 0 && *retries < MAX_RETRIES {
		time.Sleep(5 * time.Second)
		*retries++
		return so.imageSpecification(ctx, url, digest, retries)
	} else {
		for _, image := range images {
			if (image.RepoDigests != nil && image.RepoDigests[0] == digest) || (image.RepoTags != nil && image.RepoTags[0] == url) {
				if inspect, _, err := so.dkr.ImageInspectWithRaw(ctx, image.ID); err != nil {
					*retries++
					return so.imageSpecification(ctx, url, digest, retries)
				} else {
					return inspect, nil
				}

			}
		}
	}
	return types.ImageInspect{}, fmt.Errorf("no matching image found")
}

func (so *SwarmOperator) encodeCredentials(username, password string) string {
	authConfig := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		logging.GetInstance().Error(err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(encodedJSON)
}

func (so *SwarmOperator) worker(ctx context.Context, address string, retries *int) (swarm.Node, error) {
	if nodes, err := so.dkr.NodeList(ctx, types.NodeListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "role", Value: "worker",
		}),
	}); err != nil {
		return swarm.Node{}, err
	} else if len(nodes) == 0 && *retries < MAX_RETRIES {
		*retries++
		return so.worker(ctx, address, retries)
	} else {
		for _, node := range nodes {
			if node.Status.Addr != address {
				continue
			} else {
				return node, nil
			}
		}
		if *retries < MAX_RETRIES {
			*retries++
			return so.worker(ctx, address, retries)
		}
		return swarm.Node{}, fmt.Errorf("could not find worker node on time")
	}
}

func (so *SwarmOperator) prepare(ctx context.Context, registry plugins.PluginRegistry, order *providers.DeploymentOrder, digest string, constraints []string, spec *swarm.ServiceSpec) error {
	labels := make(map[string]string)
	commands := make([]string, 0)
	ports := make([]swarm.PortConfig, 0)
	mounts := make([]mount.Mount, 0)
	labels["com.docker.stack.namespace"] = so.configuration.NetworkPrefix
	for _, volume := range order.Volumes {
		mounts = append(mounts, mount.Mount{Source: volume.Source, Target: volume.Target, Type: mount.TypeBind})
	}
	retries := 0
	if image, err := so.imageSpecification(ctx, order.ImageURI, order.Image(digest), &retries); err != nil {
		logging.GetInstance().Error(err)
		return err
	} else {
		for _, command := range image.Config.Cmd {
			commands = append(commands, command)
		}
		if mappedPorts, err := so.portMapper(ctx, image.Config.ExposedPorts); err != nil {
			return err
		} else {
			ports = append(ports, mappedPorts...)
		}
		if err := registry.Invoke(plugins.REVERSE_PROXY, func(p plugins.Plugin) error {
			return p.Call(
				plugins.REVERSE_PROXY_REGISTRATION,
				&plugins.ProxyRequest{Hostname: order.Hostname(), Port: uint(ports[0].TargetPort)},
				&plugins.ProxyResponse{Labels: labels})
		}); err != nil {
			return fmt.Errorf("plugin registration for [%s] failed: %s", plugins.REVERSE_PROXY_REGISTRATION, err.Error())
		}
		spec.Annotations = swarm.Annotations{
			Labels: labels,
			Name:   order.Name,
		}
		spec.Mode = swarm.ServiceMode{Replicated: &swarm.ReplicatedService{Replicas: func(v uint64) *uint64 { return &v }(uint64(order.Scale))}}
		//spec.EndpointSpec = &swarm.EndpointSpec{Ports: ports} //
		spec.EndpointSpec = &swarm.EndpointSpec{Mode: swarm.ResolutionModeVIP}
		spec.TaskTemplate = swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Hostname: order.Hostname(),
				Image:    order.Image(digest),
				Env:      order.Env(),
				Command:  commands,
				Mounts:   mounts,
				TTY:      true,
				Init:     func(b bool) *bool { return &b }(false),
				Hosts:    []string{"host.docker.internal:host-gateway"},
			},
			Placement: &swarm.Placement{MaxReplicas: *order.Replicas(), Platforms: []swarm.Platform{{Architecture: "amd64", OS: "linux"}}, Constraints: constraints},
			RestartPolicy: &swarm.RestartPolicy{
				Condition: swarm.RestartPolicyConditionAny,
				Delay:     func(t time.Duration) *time.Duration { return &t }(2 * time.Second),
			},
		}
		spec.Networks = so.networks(ctx)
	}
	return nil
}

func (so *SwarmOperator) serviceTasks(ctx context.Context, serviceId string, retries *int) ([]swarm.Task, error) {
	filter := filters.NewArgs()
	filter.Add("label", fmt.Sprintf("com.docker.stack.namespace=%s", so.configuration.NetworkPrefix))
	results := make([]swarm.Task, 0)
	if tasks, err := so.dkr.TaskList(ctx, types.TaskListOptions{Filters: filter}); err != nil {
		return []swarm.Task{}, err
	} else if len(tasks) == 0 && *retries < MAX_RETRIES {
		time.Sleep(5 * time.Second)
		*retries++
		return so.serviceTasks(ctx, serviceId, retries)
	} else if len(tasks) == 0 && *retries >= MAX_RETRIES {
		return []swarm.Task{}, fmt.Errorf("no tasks found matching service id: %s", serviceId)
	} else {
		for _, task := range tasks {
			if task.ServiceID == serviceId {
				if target, _, err := so.dkr.TaskInspectWithRaw(ctx, task.ID); err != nil {
					return []swarm.Task{}, err
				} else if len(task.NodeID) != 0 && task.Status.ContainerStatus != nil {
					results = append(results, target)
				} else {
					time.Sleep(5 * time.Second)
					*retries++
					return so.serviceTasks(ctx, serviceId, retries)
				}
			}
		}
		logging.GetInstance().Infof("tasks found :%d", len(results))
		return results, nil
	}
}

func (so *SwarmOperator) localAttach(ctx context.Context, ws *websocket.Conn, shell, container string) error {
	commands := so.commandsForExec("[ $(command -v bash) ] && exec bash -l || exec sh -l")
	if shell != "" {
		commands = append([]string{"/usr/bin/env", "TERM=" + shell}, commands...)
	}

	var term *terminal.Terminal
	buffer := providers.OptionalWriter{}
	defer func() {
		for term != nil {
			buffer.Disable()
			if _, err := term.ReadLine(); err != nil {
				break
			} else {
				//fmt.Fprintf("", "> %s\n", line)
			}
		}
	}()
	term = terminal.NewTerminal(&buffer, "")
	options := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Detach:       false,
		Tty:          true,
		Cmd:          commands,
	}
	commander := providers.ShellLogger{Base: &providers.WebSocketWriter{Conn: ws}, Term: term}
	//commander := providers.WebSocketWriter{Conn: ws}
	if resp, err := so.dkr.ContainerExecCreate(ctx, container, options); err != nil {
		return err
	} else if connection, err := so.dkr.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{Detach: false, Tty: true}); err != nil {
		return err
	} else {
		defer connection.Close()
		defer commander.Close()
		errs := make(chan error, 2)
		quit := make(chan bool)
		go func() {
			if _, err := io.Copy(connection.Conn, &commander); err != nil {
				errs <- err
			}
		}()
		go func() {
			defer close(quit)
			_, err := io.Copy(&commander, connection.Conn)
			if err != nil && err != io.EOF {
				errs <- err
			}
		}()
		<-quit
		close(errs)
		return <-errs
	}
}

func (so *SwarmOperator) remoteAttach(ctx context.Context, rc *client.Client, ws *websocket.Conn, shell, container string) error {
	commands := so.commandsForExec("[ $(command -v bash) ] && exec bash -l || exec sh -l")
	if shell != "" {
		commands = append([]string{"/usr/bin/env", "TERM=" + shell}, commands...)
	}

	var term *terminal.Terminal
	buffer := providers.OptionalWriter{}
	defer func() {
		for term != nil {
			buffer.Disable()
			if _, err := term.ReadLine(); err != nil {
				break
			} else {
				//fmt.Fprintf("", "> %s\n", line)
			}
		}
	}()
	term = terminal.NewTerminal(&buffer, "")
	options := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Detach:       false,
		Tty:          true,
		Cmd:          commands,
	}
	commander := providers.ShellLogger{Base: &providers.WebSocketWriter{Conn: ws}, Term: term}
	//commander := providers.WebSocketWriter{Conn: ws}
	if resp, err := so.dkr.ContainerExecCreate(ctx, container, options); err != nil {
		return err
	} else if connection, err := so.dkr.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{Detach: false, Tty: true}); err != nil {
		return err
	} else {
		defer connection.Close()
		defer commander.Close()
		errs := make(chan error, 2)
		quit := make(chan bool)
		go func() {
			if _, err := io.Copy(connection.Conn, &commander); err != nil {
				errs <- err
			}
		}()
		go func() {
			defer close(quit)
			_, err := io.Copy(&commander, connection.Conn)
			if err != nil && err != io.EOF {
				errs <- err
			}
		}()
		<-quit
		close(errs)
		return <-errs
	}
}

func (so *SwarmOperator) commandsForExec(cmd string) []string {
	source := "[ -f /home/application/apprc ] && source /home/application/apprc"
	cd := fmt.Sprintf("[ -d %s ] && cd %s", defaultAppDir, defaultAppDir)
	return []string{"/bin/sh", "-c", fmt.Sprintf("%s; %s; %s", source, cd, cmd)}
}

func (so *SwarmOperator) remoteClient(workerIP string) (*client.Client, error) {
	options := make([]client.Opt, 0)
	options = append(options, client.WithTimeout(5*time.Second))
	options = append(options, client.WithHost(fmt.Sprintf("tcp://%s:2375", workerIP)))
	options = append(options, client.WithAPIVersionNegotiation())
	return client.NewClientWithOpts(options...)
}
