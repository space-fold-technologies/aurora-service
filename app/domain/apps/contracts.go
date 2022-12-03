package apps

import (
	"net/http"
	"strconv"
	"time"
)

type LogProperties struct {
	Name           string
	Width          int
	Height         int
	ShowDetails    bool
	ShowTimeStamps bool
	Token          string
}

func ParseLogProperties(r *http.Request) *LogProperties {
	return &LogProperties{
		Name:           r.Header.Values("name")[0],
		Width:          ToInt(r.Header.Values("width")[0]),
		Height:         ToInt(r.Header.Values("height")[0]),
		ShowDetails:    ToBool(r.Header.Values("show-details")[0]),
		ShowTimeStamps: ToBool(r.Header.Values("show-timestamps")[0]),
		Token:          SafeFetch("token", r),
	}
}

type ShellProperties struct {
	Name       string
	Identifier string
	Width      int
	Height     int
	Term       string
	Token      string
}

func ParseShellProperties(r *http.Request) *ShellProperties {
	return &ShellProperties{
		Identifier: r.Header.Values("identifier")[0],
		Name:       r.Header.Values("name")[0],
		Width:      ToInt(r.Header.Values("width")[0]),
		Height:     ToInt(r.Header.Values("height")[0]),
		Term:       SafeFetch("term", r),
		Token:      SafeFetch("token", r),
	}
}

type DeploymentProperties struct {
	Name       string
	Identifier string
	Token      string
}

func ParseDeploymentProperties(r *http.Request) *DeploymentProperties {
	return &DeploymentProperties{
		Identifier: r.Header.Values("identifier")[0],
		Name:       r.Header.Values("name")[0],
		Token:      SafeFetch("token", r),
	}
}

type RollbackProperties struct {
	Name       string
	Identifier string
}

func ParseRollbackProperties(r *http.Request) *RollbackProperties {
	return &RollbackProperties{Name: r.Header.Values("name")[0], Identifier: r.Header.Values("identifier")[0]}
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
