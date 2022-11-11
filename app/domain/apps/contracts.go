package apps

import (
	"net/http"
	"strconv"
	"time"

	"github.com/space-fold-technologies/aurora-service/app/core/providers"
)

func Parse(r *http.Request) *providers.TerminalProperties {
	return &providers.TerminalProperties{
		Identifier:     r.Header.Values("identifier")[0],
		Name:           r.Header.Values("name")[0],
		Width:          ToInt(r.Header.Values("width")[0]),
		Heigth:         ToInt(r.Header.Values("height")[0]),
		ClientTerminal: SafeFetch("term", r),
		Token:          SafeFetch("token", r),
		Isolated:       ToBool(r.Header.Values("isolated")[0]),
	}
}

func ToInt(val string) int {
	if value, err := strconv.Atoi(val); err != nil {
		panic(err)
	} else {
		return value
	}
}

func ToBool(val string) bool {
	if ok, err := strconv.ParseBool(val); err != nil {
		panic(err)
	} else {
		return ok
	}
}

func SafeFetch(key string, r *http.Request) string {
	values := r.Header.Values(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
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
	ImageURI   string
	UpdatedAt  *time.Time
}
