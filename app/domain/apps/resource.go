package apps

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/server/http/controllers"
	"github.com/space-fold-technologies/aurora-service/app/core/server/registry"
	"google.golang.org/protobuf/proto"
)

const BASE_PATH = "/api/v1/aurora-service/apps"

type AppController struct {
	*controllers.ControllerBase
	service  *AppService
	upgrader websocket.Upgrader
}

func NewController(service *AppService) controllers.HTTPController {
	return &AppController{service: service, upgrader: websocket.Upgrader{
		HandshakeTimeout: 0,
		ReadBufferSize:   0,
		WriteBufferSize:  0,
		WriteBufferPool:  nil,
		Subprotocols:     []string{},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}}
}

func (ac *AppController) Name() string {
	return "apps-controller"
}

func (ac *AppController) Initialize(RouteRegistry registry.RouterRegistry) {
	RouteRegistry.AddRestricted(
		BASE_PATH+"/create",
		[]string{"apps.create"},
		"POST",
		ac.create,
	)
	RouteRegistry.AddRestricted(
		BASE_PATH+"/setup-deployment",
		[]string{"apps.deploy"},
		"POST",
		ac.setupDeployment,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/deploy",
		[]string{"apps.deploy"},
		"GET",
		ac.deploy,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/{application-name}/information",
		[]string{"apps.information"},
		"GET",
		ac.information,
	)
	RouteRegistry.AddRestricted(
		BASE_PATH+"/{cluster}/list",
		[]string{"apps.information"},
		"GET",
		ac.list,
	)
	RouteRegistry.AddRestricted(
		BASE_PATH+"/update",
		[]string{"apps.update"},
		"PUT",
		ac.update,
	)
	RouteRegistry.AddRestricted(
		BASE_PATH+"/{application-name}/logs",
		[]string{"apps.information"},
		"GET",
		ac.logs,
	)
	RouteRegistry.AddRestricted(
		BASE_PATH+"/{application-name}/shell",
		[]string{"apps.shell"},
		"GET",
		ac.shell,
	)
}

func (ac *AppController) create(w http.ResponseWriter, r *http.Request) {
	order := &CreateAppOrder{}
	if data, err := io.ReadAll(r.Body); err != nil {
		ac.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		ac.BadRequest(w, err)
	} else if err := ac.service.Create(order); err != nil {
		ac.ServiceFailure(w, err)
	} else {
		ac.OKNoResponse(w)
	}
}

func (ac *AppController) deploy(w http.ResponseWriter, r *http.Request) {
	if ws, err := ac.upgrader.Upgrade(w, r, nil); err != nil {
		ac.ServiceFailure(w, err)
	} else if err := ac.service.Deploy(ws, Parse(r)); err != nil {
		ac.ServiceFailure(w, err)
	}
}

func (ac *AppController) setupDeployment(w http.ResponseWriter, r *http.Request) {
	order := &DeployAppOrder{}
	if data, err := io.ReadAll(r.Body); err != nil {
		ac.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		ac.BadRequest(w, err)
	} else if pass, err := ac.service.SetupDeployment(order); err != nil {
		ac.ServiceFailure(w, err)
	} else {
		ac.OK(w, pass)
	}
}

func (ac *AppController) information(w http.ResponseWriter, r *http.Request) {
	name := ac.GetVar("application-name", r)
	if result, err := ac.service.Information(name); err != nil {
		ac.ServiceFailure(w, err)
	} else {
		ac.OK(w, result)
	}
}

func (ac *AppController) list(w http.ResponseWriter, r *http.Request) {
	//Returns Application names against instances and nodes with IP
	cluster := ac.GetVar("cluster", r)
	if result, err := ac.service.List(cluster); err != nil {
		logging.GetInstance().Error(err)
		ac.ServiceFailure(w, err)
	} else {
		ac.OK(w, result)
	}
}

func (ac *AppController) update(w http.ResponseWriter, r *http.Request) {
	order := &UpdateAppOrder{}
	if data, err := io.ReadAll(r.Body); err != nil {
		ac.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		ac.BadRequest(w, err)
	} else if err := ac.service.Update(order); err != nil {
		ac.ServiceFailure(w, err)
	} else {
		ac.OKNoResponse(w)
	}
}

func (ac *AppController) logs(w http.ResponseWriter, r *http.Request) {
	//Requires a web socket connetion
	name := ac.GetVar("application-name", r)

	if ws, err := ac.upgrader.Upgrade(w, r, nil); err != nil {
		ac.ServiceFailure(w, err)
	} else if err := ac.service.Log(ws, name, Parse(r)); err != nil {
		ac.ServiceFailure(w, err)
	}
}

func (ac *AppController) shell(w http.ResponseWriter, r *http.Request) {
	//Requires a web socket connetion
	name := ac.GetVar("application-name", r)
	if ws, err := ac.upgrader.Upgrade(w, r, nil); err != nil {
		ac.ServiceFailure(w, err)
	} else if err := ac.service.Shell(ws, name, Parse(r)); err != nil {
		ac.ServiceFailure(w, err)
	}
}
