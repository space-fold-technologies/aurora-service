package apps

import (
	"net/http"
	"time"

	"github.com/space-fold-technologies/aurora-service/app/core/providers"
)

func Parse(r *http.Request) *providers.TerminalProperties {
	return nil
}

type ApplicationEntry struct {
	Identifier  string
	Name        string
	Description string
	Team        string
	Cluster     string
	Scale       uint32
}

type ContainerOrder struct {
	Identifier string
	IP         string
	Family     uint
	ServiceID  string
	Node       string
}
type DeploymentOrder struct {
	Identifier      string
	ApplicationName string
	Status          string
	ImageURI        string
	CreatedAt       *time.Time
}
type DeploymentUpdate struct {
	Identifier string
	Status     string
	Report     string
	ServiceID  string
	UpdatedAt  *time.Time
}
