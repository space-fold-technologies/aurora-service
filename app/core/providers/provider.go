package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
)

type Order struct {
	Identifier string
	Name       string
	URI        string
	Username   string
	Password   string
	Temporary  bool
	Scale      uint
	Variables  []*Variable
	Volumes    []Mount
}

type Instance struct {
	ID        string
	IP        string
	Family    uint
	Node      string
	ServiceID string
	TaskID    string
}

func (o *Order) ImageName() string {
	parts := strings.Split(o.URI, "/")
	if len(parts) > 0 {
		return strings.Trim(parts[len(parts)-1], " ")
	}
	return ""
}

func (o *Order) Env() []string {
	envs := make([]string, 0)
	for _, entry := range o.Variables {
		logging.GetInstance().Infof("ENV-VARS : KEY %s VAL %s", entry.Key, entry.Value)
		envs = append(envs, fmt.Sprintf("%s=%s", entry.Key, entry.Value))
	}
	return envs
}

func (o *Order) Hostname() string {
	return strings.ReplaceAll(strings.ToLower(o.Name), " ", "-")
}

func (o *Order) Image(digest string) string {
	var noTag string
	if i := strings.LastIndex(o.URI, ":"); i >= 0 {
		noTag = o.URI[0:i]
	} else {
		noTag = o.URI
	}
	return noTag + "@" + digest
}

func (o *Order) Ports() []uint {
	return []uint{}
}

func (o *Order) Replicas() *uint64 {
	value := uint64(o.Scale)
	if value == 0 {
		value = 1
	}
	return &value
}

type Report struct {
	Status      string
	Message     string
	ServiceID   string
	ImageDigest string
	Instances   map[string]*Instance
}

type CurrentState struct {
	Status    string
	Message   string
	Instances map[string]*Instance
}

type DeploymentCallback func(ctx context.Context, report *Report) error
type StatusCallback func(ctx context.Context, state *CurrentState) error
type JoinOrder struct {
	Name           string
	WorkerAddress  string
	CaptainAddress string
	Token          string
}

type LeaveOrder struct {
	NodeID  string
	Address string
	Token   string
}

type NodeDetails struct {
	ID string
	IP string
}

type ApplicationOrder struct {
	ID      string
	Name    string
	URI     string
	Digest  string
	Ports   []int
	Volumes map[string]string
	Command []string
}

func (o *ApplicationOrder) Image(digest string) string {
	var noTag string
	if i := strings.LastIndex(o.URI, ":"); i >= 0 {
		noTag = o.URI[0:i]
	} else {
		noTag = o.URI
	}
	return noTag + "@" + digest
}

type Provider interface {
	Deploy(ws *websocket.Conn, properties *TerminalProperties, order *Order, callback DeploymentCallback) error
	Stop(serviceId string) error
	Nuke(serviceId string) error
	Fetch(name string, callback StatusCallback) error
	LogContainer(ws *websocket.Conn, properties *TerminalProperties, container string) error
	LogService(ws *websocket.Conn, properties *TerminalProperties, service string) error
	Shell(ws *websocket.Conn, properties *TerminalProperties, container string) error
	Initialize(ListenAddr, AvertiseAddr string) (string, error)
	Details() (*ManagerDetails, error)
	Join(order *JoinOrder) (*NodeDetails, error)
	Leave(order *LeaveOrder) error
	CreateApplication(order *ApplicationOrder) (string, error)
}
