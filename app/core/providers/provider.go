package providers

import (
	"context"

	"github.com/gorilla/websocket"
)

type DoneCallback func(report Report) error
type ProgressCallback func(line []byte)
type Reporter struct {
	Done     DoneCallback
	Progress ProgressCallback
}
type StatusCallback func(ctx context.Context, state *CurrentState) error

type Provider interface {
	Initialize(ListenAddr, AvertiseAddr string) (string, error)
	Details() (*ManagerDetails, error)
	DeployService(ws *websocket.Conn, order *DeploymentOrder, reporter *Reporter) error
	Stop(serviceId string) error
	Log(ws *websocket.Conn, properties *LogProperties) error
	Shell(ws *websocket.Conn, properties *ShellProperties) error
	Join(order *JoinOrder) (*NodeDetails, error)
	Leave(order *LeaveOrder) error
	DeployDependency(order *DependencyOrder) (string, error)
}
