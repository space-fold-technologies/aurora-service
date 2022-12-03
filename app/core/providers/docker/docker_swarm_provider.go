package docker

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/space-fold-technologies/aurora-service/app/core/plugins"
	"github.com/space-fold-technologies/aurora-service/app/core/providers"
)

type DockerSwarmProvider struct {
	pluginRegistry plugins.PluginRegistry
	operator       *SwarmOperator
}

func NewProvider(pluginRegistry plugins.PluginRegistry, agent providers.AgentClient, configuration DockerServiceConfigurations) providers.Provider {
	instance := new(DockerSwarmProvider)
	instance.pluginRegistry = pluginRegistry
	instance.operator = NewSwarmOperator(agent, configuration)
	return instance
}

func (dsp *DockerSwarmProvider) Initialize(listenAddr, avertiseAddr string) (string, error) {
	ctx := context.Background()
	defer ctx.Done()
	return dsp.operator.InitializeManager(ctx, listenAddr, avertiseAddr)
}

func (dsp *DockerSwarmProvider) Details() (*providers.ManagerDetails, error) {
	ctx := context.Background()
	defer ctx.Done()
	return dsp.operator.ParentNode(ctx)
}

func (dsp *DockerSwarmProvider) DeployService(ws *websocket.Conn, order *providers.DeploymentOrder, reporter *providers.Reporter) error {
	ctx := context.Background()
	defer ctx.Done()
	if dsp.operator.HasWorkers(ctx) {
		return dsp.operator.DeployToWorker(ctx, dsp.pluginRegistry, order, reporter)
	}
	return dsp.operator.DeployToManager(ctx, dsp.pluginRegistry, order, reporter)
}
func (dsp *DockerSwarmProvider) Stop(serviceId string) error {
	ctx := context.Background()
	defer ctx.Done()
	return dsp.operator.StopService(ctx, serviceId)
}

func (dsp *DockerSwarmProvider) Log(ws *websocket.Conn, properties *providers.LogProperties) error {
	ctx := context.Background()
	defer ctx.Done()
	return dsp.operator.ForwardServiceLogs(ctx, ws, properties)
}
func (dsp *DockerSwarmProvider) Shell(ws *websocket.Conn, properties *providers.ShellProperties) error {
	ctx := context.Background()
	defer ctx.Done()
	if dsp.operator.IsLocalNode(properties.Host) {
		return dsp.operator.ForwardLocalShell(ctx, ws, properties.ContainerID, properties.Terminal)
	}
	return dsp.operator.ForwardRemoteShell(ctx, ws, properties.ContainerID, properties.Terminal, properties.Host)
}

func (dsp *DockerSwarmProvider) Join(order *providers.JoinOrder) (*providers.NodeDetails, error) {
	ctx := context.Background()
	defer ctx.Done()
	return dsp.operator.JoinNodeToCluster(ctx, order)
}
func (dsp *DockerSwarmProvider) Leave(order *providers.LeaveOrder) error {
	ctx := context.Background()
	defer ctx.Done()
	return dsp.operator.RemoveNodeFromCluster(ctx, order)
}
func (dsp *DockerSwarmProvider) DeployDependency(order *providers.DependencyOrder) (string, error) {
	ctx := context.Background()
	defer ctx.Done()
	return dsp.operator.DeployInternalDependency(ctx, order)
}

func (dsp *DockerSwarmProvider) FetchContainers(identifiers []string, status providers.ContainersCallback) error {
	ctx := context.Background()
	defer ctx.Done()
	return dsp.operator.FetchContainers(ctx, identifiers, status)
}

// internals
