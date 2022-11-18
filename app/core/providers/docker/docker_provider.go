package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/plugins"
	"github.com/space-fold-technologies/aurora-service/app/core/providers"
	terminal "golang.org/x/term"
)

var (
	pongWait      = 60 * time.Second
	pingInterval  = 20 * time.Second
	defaultAppDir = "/home/application/current"
	MAX_RETRIES   = 10
)

type DockerServiceConfigurations struct {
	NetworkName   string
	NetworkPrefix string
}

type DockerProvider struct {
	dkr            *client.Client
	pluginRegistry plugins.PluginRegistry
	agent          providers.AgentClient
	configuration  DockerServiceConfigurations
}

func NewProvider(pluginRegistry plugins.PluginRegistry, agent providers.AgentClient, configuration DockerServiceConfigurations) providers.Provider {
	instance := new(DockerProvider)
	instance.pluginRegistry = pluginRegistry
	instance.agent = agent
	instance.configuration = configuration
	if err := instance.initialize(); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	}
	return instance
}

func (dp *DockerProvider) initialize() error {
	if cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()); err != nil {
		return err
	} else {
		dp.dkr = cli
	}
	return nil
}

func (dp *DockerProvider) Deploy(ws *websocket.Conn, properties *providers.TerminalProperties, order *providers.Order, callback providers.DeploymentCallback) error {
	ctx := context.Background()
	quit := make(chan struct{})
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	defer close(quit)
	defer ctx.Done()
	defer ws.WriteControl(websocket.CloseMessage, nil, time.Now().Add(time.Second))
	go func() {
		for {
			select {
			case <-quit:
				return
			case <-time.After(pingInterval):
			}
			ws.WriteControl(websocket.PingMessage, nil, time.Now().Add(2*time.Second))
		}
	}()
	retries := 0
	if digest, err := dp.pullImage(ctx, ws, order.URI, order.Username, order.Password, &retries); err != nil {
		return err
	} else if serviceId, err := dp.createService(ctx, ws, order, digest); err != nil {
		return err
	} else if err := dp.success(ctx, serviceId, callback); err != nil {
		logging.GetInstance().Error(err)
		return err
	}
	return nil
}

func (dp *DockerProvider) Stop(container string) error {
	ctx := context.Background()
	defer ctx.Done()
	return dp.dkr.ServiceRemove(ctx, container)
}

func (dp *DockerProvider) Nuke(serviceId string) error {
	ctx := context.Background()
	defer ctx.Done()
	return dp.dkr.ServiceRemove(ctx, serviceId)
}

func (dp *DockerProvider) Fetch(name string, callback providers.StatusCallback) error {
	ctx := context.Background()
	defer ctx.Done()
	retries := 0
	containers := make(map[string]*providers.Instance)
	if err := dp.containers(ctx, name, containers, &retries); err != nil {
		return err
	} else if len(containers) == 0 {
		return fmt.Errorf("no running containers found matching %s service", name)
	}

	return callback(ctx, &providers.CurrentState{
		Status:    "RUNNING",
		Message:   "currently running containers",
		Instances: containers,
	})
}

func (dp *DockerProvider) LogContainer(ws *websocket.Conn, properties *providers.TerminalProperties, container string) error {
	ctx := context.Background()
	quit := make(chan struct{})
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	defer close(quit)
	defer ctx.Done()
	go func() {
		for {
			select {
			case <-quit:
				return
			case <-time.After(pingInterval):
			}
			ws.WriteControl(websocket.PingMessage, nil, time.Now().Add(2*time.Second))
		}
	}()
	return dp.logContainer(ctx, ws, properties, container)
}

func (dp *DockerProvider) LogService(ws *websocket.Conn, properties *providers.TerminalProperties, serviceId string) error {
	ctx := context.Background()
	quit := make(chan struct{})
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	defer close(quit)
	defer ctx.Done()
	go func() {
		for {
			select {
			case <-quit:
				return
			case <-time.After(pingInterval):
			}
			ws.WriteControl(websocket.PingMessage, nil, time.Now().Add(2*time.Second))
		}
	}()
	return dp.logService(ctx, ws, properties, serviceId)
}

func (dp *DockerProvider) Shell(ws *websocket.Conn, properties *providers.TerminalProperties, container string) error {
	ctx := context.Background()
	quit := make(chan struct{})
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	defer close(quit)
	defer ctx.Done()
	go func() {
		for {
			select {
			case <-quit:
				return
			case <-time.After(pingInterval):
			}
			ws.WriteControl(websocket.PingMessage, nil, time.Now().Add(2*time.Second))
		}
	}()
	return dp.attach(ctx, ws, properties, container)
}

func (dp *DockerProvider) Initialize(ListenAddr, AvertiseAddr string) (string, error) {
	ctx := context.Background()
	defer ctx.Done()
	return dp.dkr.SwarmInit(ctx, swarm.InitRequest{
		ListenAddr:      ListenAddr,
		AdvertiseAddr:   AvertiseAddr,
		DataPathAddr:    AvertiseAddr,
		DataPathPort:    0,
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

func (dp *DockerProvider) Details() (*providers.ManagerDetails, error) {
	ctx := context.Background()
	defer ctx.Done()

	if nodes, err := dp.dkr.NodeList(ctx, types.NodeListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "role", Value: "manager",
		}),
	}); err != nil {
		return nil, err
	} else if len(nodes) == 0 {
		return nil, fmt.Errorf("no manager nodes found")
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

func (dp *DockerProvider) Join(order *providers.JoinOrder) (*providers.NodeDetails, error) {
	logging.GetInstance().Infof("ADDING NODE WITH IP: %s TO CLUSTER WITH ADDR: %s", order.WorkerAddress, order.CaptainAddress)
	ctx := context.Background()
	retries := 0
	if swarm, err := dp.dkr.SwarmInspect(ctx); err != nil {
		return nil, err
	} else if err := dp.join(ctx, order, swarm.JoinTokens.Worker); err != nil {
		logging.GetInstance().Errorf("Failed to add: %s to cluster : %s", order.WorkerAddress, order.CaptainAddress)
		return nil, err
	} else if node, err := dp.worker(ctx, order.WorkerAddress, &retries); err != nil {
		return nil, err
	} else if err := dp.label(ctx, node.ID, order.Name, swarm.Version); err != nil {
		return nil, err
	} else {
		logging.GetInstance().Infof("ADDED NODE WITH ID: %s AND ADDR: %s", node.ID, node.Status.Addr)
		return &providers.NodeDetails{ID: node.ID, IP: node.Status.Addr}, nil
	}
}

func (dp *DockerProvider) Leave(order *providers.LeaveOrder) error {
	ctx := context.Background()
	defer ctx.Done()
	return dp.agent.Leave(ctx, &providers.RemoveAgent{Id: order.NodeID}, order.Address, order.Token)
}

func (dp *DockerProvider) CreateApplication(order *providers.ApplicationOrder) (string, error) {
	ctx := context.Background()
	defer ctx.Done()
	labels := make(map[string]string)
	labels["aurora.plugin.service.name"] = order.Name
	labels["aurora.plugin.service.id"] = order.ID
	mounts := make([]mount.Mount, 0)
	for source, target := range order.Volumes {
		mounts = append(mounts, mount.Mount{Source: source, Target: target, Type: mount.TypeBind})
	}
	networks := make([]swarm.NetworkAttachmentConfig, 0)
	if ok, err := dp.hasNetwork(ctx); err != nil {
		logging.GetInstance().Error(err)
		return "", err
	} else if !ok {
		if networkId, err := dp.createNetwork(ctx); err != nil {
			logging.GetInstance().Error(err)
			return "", err
		} else {
			logging.GetInstance().Infof("created default network with id :%s", networkId)
		}
	}
	networks = append(networks, swarm.NetworkAttachmentConfig{
		Target: fmt.Sprintf("%s-%s", dp.configuration.NetworkPrefix, dp.configuration.NetworkName),
	})

	ports := make([]swarm.PortConfig, 0)
	for port := range order.Ports {
		ports = append(ports, swarm.PortConfig{TargetPort: uint32(port), PublishedPort: uint32(port)})
	}
	service := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name:   order.Name,
			Labels: labels,
		},
		TaskTemplate: swarm.TaskSpec{
			Placement: &swarm.Placement{
				Constraints: []string{"node.role == manager"},
				MaxReplicas: 1,
			},
			ContainerSpec: &swarm.ContainerSpec{
				Hostname: order.Name,
				Image:    order.Image(order.Digest),
				Env:      []string{},
				Command:  order.Command,
				Mounts:   mounts,
				TTY:      true,
				Labels:   map[string]string{},
			},
			RestartPolicy: &swarm.RestartPolicy{Condition: swarm.RestartPolicyConditionAny},
			Networks:      networks,
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: ports,
		},
	}
	options := types.ServiceCreateOptions{}
	if resp, err := dp.dkr.ServiceCreate(ctx, service, options); err != nil {
		logging.GetInstance().Error(err)
		return "", err
	} else {
		return resp.ID, nil
	}
}

func (dp *DockerProvider) fetchImage(ctx context.Context, uri string, retries *int) error {
	logging.GetInstance().Infof("Pulling docker image : %s", uri)
	if stream, err := dp.dkr.ImagePull(ctx, uri, types.ImagePullOptions{}); err != nil {
		logging.GetInstance().Error(err)
		if strings.Contains(strings.ToLower(err.Error()), "timeout") && *retries <= MAX_RETRIES {
			time.Sleep(5 * time.Second)
			*retries++
			return dp.fetchImage(ctx, uri, retries)
		}
		return err
	} else {
		io.Copy(os.Stdout, stream)
	}
	return nil
}

func (dp *DockerProvider) join(ctx context.Context, order *providers.JoinOrder, workerToken string) error {
	logging.GetInstance().Infof("WORKER TOKEN: %s WORKER IP: %s", workerToken, order.WorkerAddress)
	return dp.agent.Join(ctx, &providers.RegisterAgent{
		Token:   workerToken,
		Name:    order.Name,
		Address: order.CaptainAddress,
	}, order.WorkerAddress, order.Token)
}

func (dp *DockerProvider) worker(ctx context.Context, address string, retries *int) (swarm.Node, error) {
	if nodes, err := dp.dkr.NodeList(ctx, types.NodeListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "role", Value: "worker",
		}),
	}); err != nil {
		return swarm.Node{}, err
	} else if len(nodes) == 0 && *retries < MAX_RETRIES {
		*retries++
		return dp.worker(ctx, address, retries)
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
			return dp.worker(ctx, address, retries)
		}
		return swarm.Node{}, fmt.Errorf("could not find worker node on time")
	}
}

func (dp *DockerProvider) label(ctx context.Context, identifier, name string, version swarm.Version) error {
	return dp.dkr.NodeUpdate(ctx, identifier, version, swarm.NodeSpec{
		Annotations: swarm.Annotations{
			Name: name,
		},
		Role:         swarm.NodeRoleWorker,
		Availability: swarm.NodeAvailabilityActive,
	})
}

func (dp *DockerProvider) pullImage(ctx context.Context, ws *websocket.Conn, image, username, password string, retries *int) (string, error) {
	logging.GetInstance().Infof("Pulling docker image : %s", image)
	if stream, err := dp.dkr.ImagePull(ctx, image, types.ImagePullOptions{
		RegistryAuth: dp.encodeCredentials(username, password),
	}); err != nil {
		logging.GetInstance().Error(err)
		if strings.Contains(strings.ToLower(err.Error()), "timeout") && *retries <= MAX_RETRIES {
			time.Sleep(5 * time.Second)
			*retries++
			return dp.pullImage(ctx, ws, image, username, password, retries)
		}
		return "", err
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
				} else if err = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s\n\r", string(line)))); err != nil {
					logging.GetInstance().Error("write:", err)
				}
			}

		}
		if digest == "" {
			return "", errors.New("failed to report docker image digest")
		}
		return digest, nil
	}
}

func (dp *DockerProvider) fetchImageDetailsByDigest(ctx context.Context, url, repoId string, retries *int) (*types.ImageInspect, error) {
	filterArguments := filters.NewArgs()
	filterArguments.Add("reference", fmt.Sprintf("%s*", url))
	if summary, err := dp.dkr.ImageList(ctx, types.ImageListOptions{
		All:     true,
		Filters: filterArguments,
	}); err != nil {
		logging.GetInstance().Error(err)
		return nil, err
	} else {
		if len(summary) == 0 && *retries < MAX_RETRIES {
			time.Sleep(5 * time.Second)
			*retries++
			return dp.fetchImageDetailsByDigest(ctx, url, repoId, retries)
		}
		for _, image := range summary {

			if (image.RepoDigests != nil && image.RepoDigests[0] == repoId) || (image.RepoTags != nil && image.RepoTags[0] == url) {
				if inspect, _, err := dp.dkr.ImageInspectWithRaw(ctx, image.ID); err != nil && *retries < MAX_RETRIES {
					logging.GetInstance().Error(err)
					*retries++
					return dp.fetchImageDetailsByDigest(ctx, url, repoId, retries)
				} else {
					return &inspect, err
				}
			}
		}
	}
	return nil, errors.New("no images found")
}

func (dp *DockerProvider) createService(ctx context.Context, ws *websocket.Conn, order *providers.Order, digest string) (string, error) {
	labels := make(map[string]string)
	labels["aurora.service.name"] = order.Name
	labels["aurora.service.id"] = order.Identifier

	mounts := make([]mount.Mount, 0)
	for _, volume := range order.Volumes {
		mounts = append(mounts, mount.Mount{
			Source: volume.Source,
			Target: volume.Target,
			Type:   mount.TypeBind,
		})
	}
	ws.WriteJSON(map[string]string{"status": "service-setup", "step": "Mounting volumes"})
	networks := make([]swarm.NetworkAttachmentConfig, 0)
	if ok, err := dp.hasNetwork(ctx); err != nil {
		logging.GetInstance().Error(err)
		return "", err
	} else if !ok {
		if networkId, err := dp.createNetwork(ctx); err != nil {
			logging.GetInstance().Error(err)
			return "", err
		} else {
			logging.GetInstance().Infof("created default network with id :%s", networkId)
		}
	}
	ws.WriteJSON(map[string]string{"status": "service-setup", "step": "setting up network"})
	networks = append(networks, swarm.NetworkAttachmentConfig{
		Target: fmt.Sprintf("%s-%s", dp.configuration.NetworkPrefix, dp.configuration.NetworkName),
	})

	ports := make([]swarm.PortConfig, 0)
	portList := make([]uint, 0)
	commands := make([]string, 0)
	retries := 0
	if image, err := dp.fetchImageDetailsByDigest(ctx, order.URI, order.Image(digest), &retries); err != nil {
		logging.GetInstance().Error(err)
		return "", err
	} else {
		for port := range image.Config.ExposedPorts {
			ports = append(ports, swarm.PortConfig{TargetPort: uint32(port.Int()), PublishedPort: uint32(port.Int()), Protocol: swarm.PortConfigProtocol(port.Proto())})
			portList = append(portList, uint(port.Int()))
		}
		for _, command := range image.Config.Cmd {
			commands = append(commands, command)
		}
	}

	if err := dp.pluginRegistry.Invoke(plugins.REVERSE_PROXY, func(p plugins.Plugin) error {
		return p.Call(
			plugins.REVERSE_PROXY_REGISTRATION,
			&plugins.ProxyRequest{Hostname: order.Hostname(), Port: portList[0]},
			&plugins.ProxyResponse{Labels: labels})
	}); err != nil {
		logging.GetInstance().Errorf("plugin registration for [%s] failed", plugins.REVERSE_PROXY_REGISTRATION)
		logging.GetInstance().Error(err)
	} else {
		logging.GetInstance().Infof("loaded : %s with reverse proxy plugin", order.Hostname())
		ws.WriteJSON(map[string]string{"status": "service-setup", "step": "registered service route with reverse proxy plugin"})
	}

	logging.GetInstance().Infof("LOCAL IMAGE+DIGEST: %s", order.Image(digest))
	ws.WriteJSON(map[string]string{"status": "service-setup", "step": "injecting traefik labels"})
	service := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name:   order.Name,
			Labels: labels,
		},
		Mode: swarm.ServiceMode{Replicated: &swarm.ReplicatedService{Replicas: order.Replicas()}},
		TaskTemplate: swarm.TaskSpec{
			Placement: &swarm.Placement{
				Constraints: []string{},
				MaxReplicas: uint64(order.Scale),
			},
			ContainerSpec: &swarm.ContainerSpec{
				Hostname: order.Hostname(),
				Image:    order.Image(digest),
				Env:      order.Env(),
				Command:  commands,
				Mounts:   mounts,
				TTY:      true,
				Labels:   map[string]string{"aurora.container.service.name": order.Name, "aurora.container.service.id": order.Identifier},
				Hosts:    []string{"host.docker.internal"}, //<< Some nonsese like this
			},
			RestartPolicy: &swarm.RestartPolicy{Condition: swarm.RestartPolicyConditionAny},
			Networks:      networks,
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: ports,
		},
	}
	options := types.ServiceCreateOptions{}
	ws.WriteJSON(map[string]string{"status": "service-setup", "step": "creating service"})
	if resp, err := dp.dkr.ServiceCreate(ctx, service, options); err != nil {
		logging.GetInstance().Error(err)
		return "", err
	} else {
		ws.WriteJSON(map[string]string{"status": "service-setup", "step": fmt.Sprintf("created service with id %s", resp.ID)})
		return resp.ID, nil
	}
}

func (dp *DockerProvider) createNetwork(ctx context.Context) (string, error) {
	if response, err := dp.dkr.NetworkCreate(ctx, fmt.Sprintf("%s-%s", dp.configuration.NetworkPrefix, dp.configuration.NetworkName), types.NetworkCreate{
		Driver: "overlay",
		Scope:  "swarm",
		Labels: map[string]string{
			dp.configuration.NetworkPrefix: "true",
		},
	}); err != nil {
		return "", err
	} else {
		return response.ID, nil
	}
}

func (dp *DockerProvider) hasNetwork(ctx context.Context) (bool, error) {
	if networks, err := dp.dkr.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "name", Value: fmt.Sprintf("%s-%s", dp.configuration.NetworkPrefix, dp.configuration.NetworkName),
		}),
	}); err != nil {
		return false, err
	} else if len(networks) > 0 {
		if networks[0].Driver == "overlay" && networks[0].Scope == "swarm" && networks[0].Labels[dp.configuration.NetworkPrefix] == "true" {
			logging.GetInstance().Info("network already exists and is configured correctly")
			return true, nil
		}
	}
	return false, nil
}

func (dp *DockerProvider) success(ctx context.Context, serviceId string, callback providers.DeploymentCallback) error {
	logger := logging.GetInstance()
	retries := 0
	containers := make(map[string]*providers.Instance)
	if info, _, err := dp.dkr.ServiceInspectWithRaw(ctx, serviceId, types.ServiceInspectOptions{InsertDefaults: true}); err != nil {
		return err
	} else if err := dp.containers(ctx, info.Spec.Name, containers, &retries); err != nil {
		return err
	} else {
		logger.Infof("DEPLOYED SERVICE NAME: %s", info.Spec.Name)
		logger.Infof("REPLICATON FACTOR: %d", *info.Spec.Mode.Replicated.Replicas)
		return callback(ctx, &providers.Report{
			Status:    "DEPLOYED",
			Message:   "successfully deployed service to cluster",
			ServiceID: info.ID,
			Instances: containers,
		})
	}
}

func (dp *DockerProvider) encodeCredentials(username, password string) string {
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

func (dp *DockerProvider) containers(ctx context.Context, name string, pack map[string]*providers.Instance, retries *int) error {
	logging.GetInstance().Infof("FETCHING CONTAINER DETAILS FROM SERVICE :%s ATTEMPT :%d", name, *retries)
	filterArguments := filters.NewArgs()
	filterArguments.Add("label", fmt.Sprintf("com.docker.swarm.service.name=%s", strings.Trim(name, " ")))
	if containers, err := dp.dkr.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filterArguments,
	}); err != nil {
		logging.GetInstance().Error(err)
		return err
	} else {
		if len(containers) == 0 && *retries <= MAX_RETRIES {
			time.Sleep(5 * time.Second)
			*retries++
			return dp.containers(ctx, name, pack, retries)
		}

		logging.GetInstance().Infof("CONTAINERS FOUND: %d", len(containers))
		for _, container := range containers {
			pack[container.ID] = &providers.Instance{
				ID:        container.ID,
				IP:        container.NetworkSettings.Networks["aurora-default"].IPAddress,
				Family:    4,
				Node:      container.Labels["com.docker.swarm.node.id"],
				ServiceID: container.Labels["com.docker.swarm.service.id"],
				TaskID:    container.Labels["com.docker.swarm.task.id"],
			}
		}
		return nil
	}
}

func (dp *DockerProvider) attach(ctx context.Context, ws *websocket.Conn, properties *providers.TerminalProperties, container string) error {
	commands := dp.commandsForExec("[ $(command -v bash) ] && exec bash -l || exec sh -l")
	if properties.ClientTerminal != "" {
		commands = append([]string{"/usr/bin/env", "TERM=" + properties.ClientTerminal}, commands...)
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
	if resp, err := dp.dkr.ContainerExecCreate(ctx, container, options); err != nil {
		return err
	} else if connection, err := dp.dkr.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{Detach: false, Tty: true}); err != nil {
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

func (dp *DockerProvider) logContainer(ctx context.Context, ws *websocket.Conn, properties *providers.TerminalProperties, container string) error {
	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Details: true, Timestamps: true}
	if out, err := dp.dkr.ContainerLogs(ctx, container, options); err != nil {
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

func (dp *DockerProvider) logService(ctx context.Context, ws *websocket.Conn, properties *providers.TerminalProperties, serviceId string) error {
	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Details: true, Timestamps: true}
	if out, err := dp.dkr.ServiceLogs(ctx, serviceId, options); err != nil {
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

func (dp *DockerProvider) commandsForExec(cmd string) []string {
	source := "[ -f /home/application/apprc ] && source /home/application/apprc"
	cd := fmt.Sprintf("[ -d %s ] && cd %s", defaultAppDir, defaultAppDir)
	return []string{"/bin/sh", "-c", fmt.Sprintf("%s; %s; %s", source, cd, cmd)}
}
