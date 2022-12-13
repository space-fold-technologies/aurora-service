package providers

import (
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

type Variable struct {
	Key   string
	Value string
}

type Mount struct {
	Source   string
	Target   string
	TypeBind string
}

type ShellProperties struct {
	ContainerID string
	Host        string
	Name        string
	Width       int
	Heigth      int
	Terminal    string
}

type LogProperties struct {
	ServiceID      string
	Details        bool
	ShowTimestamps bool
}

type ManagerDetails struct {
	ID      string
	Address string
}

type DeploymentOrder struct {
	ID                   string
	Name                 string
	ImageURI             string
	Username             string
	Password             string
	Temporary            bool
	Scale                uint
	Constraint           []string
	EnvironmentVariables map[string]string
	Volumes              []Mount
	Token                string
}

type Instance struct {
	ID        string
	IP        string
	Family    uint
	Node      string
	ServiceID string
	TaskID    string
}

func (o *DeploymentOrder) ImageName() string {
	parts := strings.Split(o.ImageURI, "/")
	if len(parts) > 0 {
		return strings.Trim(parts[len(parts)-1], " ")
	}
	return ""
}

func (o *DeploymentOrder) Env() []string {
	unique := mapset.NewSet[string]()
	for key, value := range o.EnvironmentVariables {
		unique.Add(fmt.Sprintf("%s=%s", key, value))
	}
	return unique.ToSlice()
}

func (o *DeploymentOrder) Hostname() string {
	return strings.ReplaceAll(strings.ToLower(o.Name), " ", "-")
}

func (o *DeploymentOrder) Image(digest string) string {
	var noTag string
	if i := strings.LastIndex(o.ImageURI, ":"); i >= 0 {
		noTag = o.ImageURI[0:i]
	} else {
		noTag = o.ImageURI
	}
	return noTag + "@" + digest
}

func (o *DeploymentOrder) Ports() []uint {
	return []uint{}
}

func (o *DeploymentOrder) Replicas() *uint64 {
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

type DependencyOrder struct {
	ID                   string
	Name                 string
	URI                  string
	Digest               string
	Ports                []int
	Volumes              map[string]string
	Command              []string
	EnvironmentVariables map[string]string
}

func (o *DependencyOrder) Image(digest string) string {
	var noTag string
	if i := strings.LastIndex(o.URI, ":"); i >= 0 {
		noTag = o.URI[0:i]
	} else {
		noTag = o.URI
	}
	return noTag + "@" + digest
}

func (o *DependencyOrder) Env() []string {
	unique := mapset.NewSet[string]()
	for key, value := range o.EnvironmentVariables {
		unique.Add(fmt.Sprintf("%s=%s", key, value))
	}
	return unique.ToSlice()
}
