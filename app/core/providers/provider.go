package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
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
	}
	return noTag + "@" + digest
}

func (o *Order) Ports() []uint {
	return []uint{}
}

func (o *Order) Replicas() *uint64 {
	value := uint64(o.Scale)
	return &value
}

type Report struct {
	Status      string
	Message     string
	ServiceID   string
	ImageDigest string
	Instances   map[string]*Instance
}

type DeploymentCallback func(ctx context.Context, report *Report) error
type PluginParameterInjector interface {
	Labels(target map[string]string, network, host, hostname string, ports []uint) error
}

type JoinOrder struct {
	ListenAddress  string
	ClusterAddress string
	Token          string
}

type NodeDetails struct {
	ID string
	IP string
}

type Provider interface {
	Deploy(ws *websocket.Conn, properties *TerminalProperties, order *Order, callback DeploymentCallback) error
	Stop(container string) error
	Log(ws *websocket.Conn, properties *TerminalProperties, container string) error
	Shell(ws *websocket.Conn, properties *TerminalProperties, container string) error
	Initialize(ListenAddr, AvertiseAddr string) (string, error)
	Join(order *JoinOrder) (*NodeDetails, error)
	Leave(nodeId string) error
}
