package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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
	"github.com/space-fold-technologies/aurora-service/app/core/providers"
)

var (
	pongWait      = 60 * time.Second
	pingInterval  = 20 * time.Second
	defaultAppDir = "/home/application/current"
	networkPrefix = "aurora"
	network       = "default"
	MAX_RETRIES   = 10
)

type DockerProvider struct {
	dkr      *client.Client
	injector providers.PluginParameterInjector
	domain   string
}

func NewProvider(injector providers.PluginParameterInjector) providers.Provider {
	instance := new(DockerProvider)
	instance.injector = injector
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
	if digest, err := dp.pullImage(ctx, ws, order.URI, order.Username, order.Password); err != nil {
		return err
	} else if serviceId, err := dp.createService(ctx, ws, order, digest); err != nil {
		return err
	} else if err := dp.success(ctx, serviceId, int(order.Scale), callback); err != nil {
		logging.GetInstance().Errorf("no peace, no deploy")
		logging.GetInstance().Error(err)
	}
	return nil
}

func (dp *DockerProvider) Stop(container string) error {
	ctx := context.Background()
	defer ctx.Done()
	return dp.dkr.ServiceRemove(ctx, container)
}

func (dp *DockerProvider) Log(ws *websocket.Conn, properties *providers.TerminalProperties, container string) error {
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
	return dp.log(ctx, ws, properties, container)
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
		ForceNewCluster: true,
		Spec: swarm.Spec{Orchestration: swarm.OrchestrationConfig{TaskHistoryRetentionLimit: func(n int64) *int64 { return &n }(5)},
			Raft: swarm.RaftConfig{
				SnapshotInterval: 10000,
				KeepOldSnapshots: func(n uint64) *uint64 { return &n }(0),
			},
			Dispatcher: swarm.DispatcherConfig{
				HeartbeatPeriod: 5 * time.Second,
			},
			EncryptionConfig: swarm.EncryptionConfig{
				AutoLockManagers: true,
			},
			CAConfig: swarm.CAConfig{NodeCertExpiry: 90 * 24 * time.Hour}},
		AutoLockManagers: true,
		Availability:     swarm.NodeAvailabilityActive,
		DefaultAddrPool:  []string{"10.0.0.0/8"},
		SubnetSize:       24,
	})
}

func (dp *DockerProvider) Join(order *providers.JoinOrder) (*providers.NodeDetails, error) {
	ctx := context.Background()
	defer ctx.Done()
	existingNodes := map[string]bool{}
	if nodes, err := dp.dkr.NodeList(ctx, types.NodeListOptions{}); err != nil {
		return nil, err
	} else {
		for _, node := range nodes {
			existingNodes[node.ID] = true
		}
		if err := dp.dkr.SwarmJoin(ctx, swarm.JoinRequest{
			ListenAddr:    order.ListenAddress,
			AdvertiseAddr: order.ClusterAddress,
			DataPathAddr:  "",
			RemoteAddrs:   []string{},
			JoinToken:     order.Token,
			Availability:  swarm.NodeAvailabilityActive,
		}); err != nil && !strings.Contains(err.Error(), "already part of") {
			return nil, err
		} else if strings.Contains(err.Error(), "already part of") {
			for k := range existingNodes {
				delete(existingNodes, k)
			}
		}
		if nodes, err := dp.dkr.NodeList(ctx, types.NodeListOptions{}); err != nil {
			return nil, err
		} else {
			for _, node := range nodes {
				if _, ok := existingNodes[node.ID]; ok {
					continue
				} else {
					return &providers.NodeDetails{ID: node.ID, IP: node.Status.Addr}, nil
				}
			}
			return nil, errors.New("no node created for some reason")
		}
	}
}

func (dp *DockerProvider) Leave(nodeId string) error {
	ctx := context.Background()
	defer ctx.Done()
	return dp.dkr.NodeRemove(ctx, nodeId, types.NodeRemoveOptions{Force: true})
}

func (dp *DockerProvider) pullImage(ctx context.Context, ws *websocket.Conn, image, username, password string) (string, error) {
	logging.GetInstance().Infof("Pulling docker image : %s", image)
	if stream, err := dp.dkr.ImagePull(ctx, image, types.ImagePullOptions{
		RegistryAuth: dp.encodeCredentials(username, password),
	}); err != nil {
		logging.GetInstance().Error(err)
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

func (dp *DockerProvider) fetchImageDetailsByDigest(ctx context.Context, url, digest string) (types.ImageInspect, error) {
	filterArguments := filters.NewArgs()
	filterArguments.Add("reference", fmt.Sprintf("%s*", url))
	if summary, err := dp.dkr.ImageList(ctx, types.ImageListOptions{
		Filters: filterArguments,
	}); err != nil {
		logging.GetInstance().Error(err)
		return types.ImageInspect{}, err
	} else {
		for _, image := range summary {
			target := fmt.Sprintf("%s@%s", url, digest)
			if image.RepoDigests[0] == target {
				if inspect, _, err := dp.dkr.ImageInspectWithRaw(ctx, image.ID); err != nil {
					logging.GetInstance().Error(err)
					return types.ImageInspect{}, err
				} else {
					return inspect, err
				}
			}
		}
	}
	return types.ImageInspect{}, nil
}

func (dp *DockerProvider) createService(ctx context.Context, ws *websocket.Conn, order *providers.Order, digest string) (string, error) {
	labels := make(map[string]string)
	host := fmt.Sprintf("%s.%s", order.Hostname(), dp.domain)
	ws.WriteJSON(map[string]string{"status": "service-setup", "step": fmt.Sprintf("creating service with host %s", host)})
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
	if ok, err := dp.hasNetwork(ctx, network); err != nil {
		logging.GetInstance().Error(err)
		return "", err
	} else if !ok {
		if networkId, err := dp.createNetwork(ctx, network); err != nil {
			logging.GetInstance().Error(err)
			return "", err
		} else {
			logging.GetInstance().Infof("created default network with id :%s", networkId)
		}
	}
	ws.WriteJSON(map[string]string{"status": "service-setup", "step": "setting up network"})
	networks = append(networks, swarm.NetworkAttachmentConfig{
		Target: fmt.Sprintf("%s-%s", networkPrefix, network),
	})

	ports := make([]swarm.PortConfig, 0)
	portList := make([]uint, 0)
	commands := make([]string, 0)
	if image, err := dp.fetchImageDetailsByDigest(ctx, order.URI, digest); err != nil {
		ws.WriteJSON(map[string]string{"status": "service-setup", "step": "fetching image details"})
	} else {
		for port, _ := range image.Config.ExposedPorts {
			ports = append(ports, swarm.PortConfig{TargetPort: uint32(port.Int()), PublishedPort: uint32(port.Int()), Protocol: swarm.PortConfigProtocol(port.Proto())})
			portList = append(portList, uint(port.Int()))
		}
		for _, command := range image.Config.Cmd {
			commands = append(commands, command)
		}
	}

	if err := dp.injector.Labels(labels, network, host, order.Hostname(), portList); err != nil {
		logging.GetInstance().Error(err)
		return "", err
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

func (dp *DockerProvider) createNetwork(ctx context.Context, name string) (string, error) {
	if response, err := dp.dkr.NetworkCreate(ctx, fmt.Sprintf("%s-%s", networkPrefix, name), types.NetworkCreate{
		Driver: "overlay",
		Scope:  "swarm",
		Labels: map[string]string{
			networkPrefix: "true",
		},
	}); err != nil {
		return "", err
	} else {
		return response.ID, nil
	}
}

func (dp *DockerProvider) hasNetwork(ctx context.Context, name string) (bool, error) {
	if networks, err := dp.dkr.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "name", Value: fmt.Sprintf("%s-%s", networkPrefix, name),
		}),
	}); err != nil {
		return false, err
	} else if len(networks) > 0 {
		if networks[0].Driver == "overlay" && networks[0].Scope == "swarm" && networks[0].Labels[networkPrefix] == "true" {
			logging.GetInstance().Info("network already exists and is configured correctly")
			return true, nil
		}
	}
	return false, nil
}

func (dp *DockerProvider) success(ctx context.Context, serviceId string, limit int, callback providers.DeploymentCallback) error {
	logger := logging.GetInstance()
	retries := 0
	containers := make(map[string]*providers.Instance)
	if info, _, err := dp.dkr.ServiceInspectWithRaw(ctx, serviceId, types.ServiceInspectOptions{InsertDefaults: true}); err != nil {
		return err
	} else if err := dp.containers(ctx, info.Spec.Name, limit, containers, &retries); err != nil {
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

func (dp *DockerProvider) containers(ctx context.Context, name string, limit int, pack map[string]*providers.Instance, retries *int) error {
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
			return dp.containers(ctx, name, limit, pack, retries)
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
	var inout chan []byte
	var output chan []byte
	commands := dp.commandsForExec("[ $(command -v bash) ] && exec bash -l || exec sh -l")
	if properties.ClientTerminal != "" {
		commands = append([]string{"/usr/bin/env", "TERM=" + properties.ClientTerminal}, commands...)
	}
	options := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Detach:       false,
		Tty:          true,
		Cmd:          commands,
	}

	if resp, err := dp.dkr.ContainerExecCreate(ctx, container, options); err != nil {
		return err
	} else if connection, err := dp.dkr.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{Detach: false, Tty: true}); err != nil {
		return err
	} else {
		defer connection.Close()
		go func(w io.WriteCloser) {
			for {
				data, ok := <-inout
				if !ok {
					fmt.Println("!ok")
					w.Close()
					return
				}

				fmt.Println(string(data))
				data = append(data, '\n')
				w.Write(append(data, '\n'))
			}
		}(connection.Conn)
		go func() {
			buffer := make([]byte, 4096)
			for {
				c, err := connection.Reader.Read(buffer)
				if c > 0 {

					output <- buffer[:c]
				}
				if err != nil {
					break
				}
			}
		}()
		for {

			mt, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			log.Printf("recv: %s", message)
			inout <- message
			data := <-output
			err = ws.WriteMessage(mt, data)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
	return nil
}

func (dp *DockerProvider) log(ctx context.Context, ws *websocket.Conn, properties *providers.TerminalProperties, container string) error {
	var output chan []byte
	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true}
	if out, err := dp.dkr.ServiceLogs(ctx, container, options); err != nil {
		return err
	} else {
		defer out.Close()
		go func() {
			buffer := make([]byte, 4096)
			for {
				c, err := out.Read(buffer)
				if c > 0 {

					output <- buffer[:c]
				}
				if err != nil {
					break
				}
			}
		}()
		for {
			data := <-output
			if err = ws.WriteMessage(websocket.BinaryMessage, data); err != nil {
				logging.GetInstance().Error("write:", err)
				break
			}
		}
	}

	return nil
}

func (dp *DockerProvider) commandsForExec(cmd string) []string {
	source := "[ -f /home/application/apprc ] && source /home/application/apprc"
	cd := fmt.Sprintf("[ -d %s ] && cd %s", defaultAppDir, defaultAppDir)
	return []string{"/bin/sh", "-c", fmt.Sprintf("%s; %s; %s", source, cd, cmd)}
}
