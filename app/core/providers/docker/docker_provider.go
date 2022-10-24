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
	} else if serviceId, err := dp.createService(ctx, order, digest); err != nil {
		return err
	} else if err := dp.success(ctx, serviceId, callback); err != nil {
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

func (dp *DockerProvider) Initialize() (string, error) {
	ctx := context.Background()
	defer ctx.Done()
	return dp.dkr.SwarmInit(ctx, swarm.InitRequest{
		ListenAddr:       "",
		AdvertiseAddr:    "",
		DataPathAddr:     "",
		DataPathPort:     0,
		ForceNewCluster:  true,
		Spec:             swarm.Spec{},
		AutoLockManagers: true,
		Availability:     swarm.NodeAvailabilityActive,
		DefaultAddrPool:  []string{},
		SubnetSize:       0,
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
			ListenAddr:    "",
			AdvertiseAddr: "",
			DataPathAddr:  "",
			RemoteAddrs:   []string{},
			JoinToken:     order.Token,
			Availability:  swarm.NodeAvailabilityActive,
		}); err != nil {
			return nil, err
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
	logging.GetInstance().Info("Pulling docker image : %s", image)
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
				} else if err = ws.WriteMessage(websocket.BinaryMessage, []byte(chunk.Status)); err != nil {
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

func (dp *DockerProvider) createService(ctx context.Context, order *providers.Order, digest string) (string, error) {
	labels := make(map[string]string)
	host := fmt.Sprintf("%s.%s", order.Hostname(), dp.domain)
	if err := dp.injector.Labels(labels, network, host, order.Hostname(), order.Ports()); err != nil {
		return "", err
	}
	mounts := make([]mount.Mount, 0)
	for _, volume := range order.Volumes {
		mounts = append(mounts, mount.Mount{
			Source: volume.Source,
			Target: volume.Target,
			Type:   mount.TypeBind,
		})
	}
	networks := make([]swarm.NetworkAttachmentConfig, 0)
	if ok, err := dp.hasNetwork(ctx, network); err != nil {
		return "", err
	} else if !ok {
		if networkId, err := dp.createNetwork(ctx, network); err != nil {
			return "", err
		} else {
			logging.GetInstance().Infof("created default network with id :%s", networkId)
		}
	}
	networks = append(networks, swarm.NetworkAttachmentConfig{
		Target: fmt.Sprintf("%s-%s", networkPrefix, network),
	})

	ports := make([]swarm.PortConfig, 0)
	for _, port := range order.Ports() {
		ports = append(ports, swarm.PortConfig{TargetPort: uint32(port), PublishedPort: uint32(port)})
	}

	service := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name:   order.Name,
			Labels: labels,
		},
		Mode: swarm.ServiceMode{Replicated: &swarm.ReplicatedService{Replicas: order.Replicas()}},
		TaskTemplate: swarm.TaskSpec{
			Placement: &swarm.Placement{
				Constraints: []string{},
			},
			ContainerSpec: &swarm.ContainerSpec{
				Hostname: order.Hostname(),
				Image:    order.ImageName(),
				Env:      order.Env(),
				Command:  []string{},
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
	if resp, err := dp.dkr.ServiceCreate(ctx, service, options); err != nil {
		return "", err
	} else {
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

func (dp *DockerProvider) success(ctx context.Context, serviceId string, callback providers.DeploymentCallback) error {
	logger := logging.GetInstance()
	containers := make(map[string]*providers.Instance)
	if info, _, err := dp.dkr.ServiceInspectWithRaw(ctx, serviceId, types.ServiceInspectOptions{InsertDefaults: true}); err != nil {
		return err
	} else if err := dp.containers(ctx, info.Spec.Name, containers); err != nil {
		return err
	} else if err := dp.nodes(ctx, info.Spec.Name, containers); err != nil {
		return err
	} else {
		logger.Infof("fetched : %d containers for service: %s", len(containers), info.Spec.Name)
		logger.Infof("DEPLOYED SERVICE NAME: %s", info.Spec.Name)
		logger.Infof("REPLICATON FACTOR: %d", *info.Spec.Mode.Replicated.Replicas)
		return callback(ctx, &providers.Report{
			Status:    "DEPLOYED",
			Message:   "successfully deployed service to cluster",
			ServiceID: serviceId,
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

func (dp *DockerProvider) containers(ctx context.Context, name string, pack map[string]*providers.Instance) error {
	// the patter of the containers is `service-name`.{index}.`container-identifier`
	// filters.Add("label", fmt.Sprintf("com.docker.swarm.node.id=%s", nodeID))
	filterArguments := filters.NewArgs()
	filterArguments.Add("status", "running")
	filterArguments.Add("status", "paused")
	filterArguments.Add("status", "restarting")
	filterArguments.Add("label", fmt.Sprintf("com.docker.swarm.service.name=%s", name))
	if containers, err := dp.dkr.ContainerList(ctx, types.ContainerListOptions{
		Filters: filterArguments,
	}); err != nil {
		logging.GetInstance().Error(err)
		return err
	} else {

		for _, container := range containers {
			pack[container.ID] = &providers.Instance{
				ID:     container.ID,
				IP:     container.NetworkSettings.Networks[""].IPAddress,
				Family: 4,
				Node:   "",
			}
		}
		return nil
	}
}

func (dp *DockerProvider) nodes(ctx context.Context, name string, pack map[string]*providers.Instance) error {
	filterArguments := filters.NewArgs()
	filterArguments.Add("status", "running")
	filterArguments.Add("label", fmt.Sprintf("com.docker.swarm.service.name=%s", name))
	if tasks, err := dp.dkr.TaskList(ctx, types.TaskListOptions{}); err != nil {
		return err
	} else {
		for _, entry := range tasks {
			instance := pack[entry.Status.ContainerStatus.ContainerID]
			instance.Node = entry.NodeID
			instance.ServiceID = entry.ServiceID
			instance.TaskID = entry.ID

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
