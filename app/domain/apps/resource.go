package apps

import (
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/server/http/controllers"
	"github.com/space-fold-technologies/aurora-service/app/core/server/registry"
	"google.golang.org/protobuf/proto"
)

const BASE_PATH = "/api/v1/aurora-service/apps"

var (
	pongWait     = 60 * time.Second
	pingInterval = 20 * time.Second
)

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
		BASE_PATH+"/log",
		[]string{"apps.information"},
		"GET",
		ac.log,
	)
	RouteRegistry.AddRestricted(
		BASE_PATH+"/shell",
		[]string{"apps.shell"},
		"GET",
		ac.shell,
	)
	RouteRegistry.AddRestricted(
		BASE_PATH+"/{application-name}/deployments",
		[]string{"apps.information"},
		"GET",
		ac.deployments,
	)
	RouteRegistry.AddRestricted(
		BASE_PATH+"/rollback",
		[]string{"apps.information"},
		"GET",
		ac.rollback,
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
	} else if err := ac.service.Deploy(ws, ParseDeploymentProperties(r)); err != nil {
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
		logging.GetInstance().Error(err)
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

func (ac *AppController) log(w http.ResponseWriter, r *http.Request) {
	//Requires a web socket connetion
	if ws, err := ac.upgrader.Upgrade(w, r, nil); err != nil {
		ac.ServiceFailure(w, err)
	} else {
		quit := make(chan struct{})
		defer ac.WebSocketClose(ws)
		defer close(quit)
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
		if err := ac.service.Log(ac.setup(ws), ParseLogProperties(r)); err != nil {
			ac.WebsocketFailure(ws, err)
		}
	}
}

func (ac *AppController) shell(w http.ResponseWriter, r *http.Request) {
	//Requires a web socket connetion
	if ws, err := ac.upgrader.Upgrade(w, r, nil); err != nil {
		ac.ServiceFailure(w, err)
	} else {

		if err := ac.service.Shell(ac.setup(ws), ParseShellProperties(r)); err != nil {
			ac.WebsocketFailure(ws, err)
		}
	}
}

func (ac *AppController) deployments(w http.ResponseWriter, r *http.Request) {
	name := ac.GetVar("application-name", r)
	if list, err := ac.service.Deployments(name); err != nil {
		ac.ServiceFailure(w, err)
	} else {
		ac.OK(w, list)
	}
}

func (ac *AppController) rollback(w http.ResponseWriter, r *http.Request) {
	//Requires a web socket connetion
	if ws, err := ac.upgrader.Upgrade(w, r, nil); err != nil {
		ac.ServiceFailure(w, err)
	} else {
		quit := make(chan struct{})
		defer ac.WebSocketClose(ws)
		defer close(quit)
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
		if err := ac.service.Rollback(ac.setup(ws), ParseRollbackProperties(r)); err != nil {
			ac.WebsocketFailure(ws, err)
		}
	}
}

func (ac *AppController) setup(ws *websocket.Conn) *websocket.Conn {
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(appData string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	ws.SetCloseHandler(func(code int, text string) error {
		return ws.Close()
	})
	return ws
}
