package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	h "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/security"
	"github.com/space-fold-technologies/aurora-service/app/core/server/http/controllers"
	"github.com/space-fold-technologies/aurora-service/app/core/server/registry"
)

type ServerCore struct {
	tokenVerifier                  security.TokenVerifier
	routerRegistry                 registry.RouterRegistry
	router                         *mux.Router
	server                         *http.Server
	controllerRegistry             *controllers.HTTPControllerRegistry
	host                           string
	port                           int
	startupCallback                func() bool
	shutdownCallback               func() bool
	healthCheckCallback            func() (interface{}, error)
	middlewareRegistrationCallback func(router *mux.Router)
	details                        Details
}

func New(details Details, host string, port int, tokenVerifier security.TokenVerifier) *ServerCore {
	instance := &ServerCore{
		details:       details,
		host:          host,
		port:          port,
		tokenVerifier: tokenVerifier,
	}
	instance.initialize()
	return instance
}

func (sc *ServerCore) initialize() {
	sc.router = mux.NewRouter()

	sc.routerRegistry = registry.NewRouteRegistry(sc.router, sc.tokenVerifier)
	sc.controllerRegistry = controllers.New(sc.routerRegistry)
}

func (sc *ServerCore) OnMiddlewareRegistration(middlewareRegistrationCallback func(router *mux.Router)) {
	sc.middlewareRegistrationCallback = middlewareRegistrationCallback
}

// GetWebControllerRegistry : Get an instance of the web controller registry
func (sc *ServerCore) GetRegistry() *controllers.HTTPControllerRegistry {
	return sc.controllerRegistry
}

// OnStartUp : Triggered on a start up event
func (sc *ServerCore) OnStartUp(startupCallback func() bool) {
	sc.startupCallback = startupCallback
}

// OnShutDown : Triggered on a shutdown event
func (sc *ServerCore) OnShutDown(shutdownCallback func() bool) {
	sc.shutdownCallback = shutdownCallback
}

// OnHealthCheck : Triggered on a call to the HC endpoint
func (sc *ServerCore) OnHealthCheck(healthCheckCallback func() (interface{}, error)) {
	sc.healthCheckCallback = healthCheckCallback
}

// Start Server
func (sc *ServerCore) Start() {
	if sc.startupCallback() {
		sc.middlewareRegistrationCallback(sc.router)
		sc.startUpServer()
	}
}

func (sc *ServerCore) startUpServer() {
	sc.server = &http.Server{
		Addr: fmt.Sprintf("%s:%d", sc.host, sc.port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: h.CORS(h.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			h.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			h.AllowedOrigins([]string{"*"}))(sc.router),
	}

	sc.routerRegistry.Initialize()
	go func() {
		sc.logoPrint()
		LOG := logging.GetInstance()
		LOG.Infof("Starting    %s ", sc.details.Name)
		LOG.Infof("Version     %s ", sc.details.Version)
		LOG.Infof("Environment %s ", sc.details.Environment)
		LOG.Infof("Repository  %s ", sc.details.Repository)
		LOG.Infof("Commit Hash %s ", sc.details.Hash)
		LOG.Infof("Build Date  %s ", sc.details.BuildDate)
		LOG.Infof("Build Epoch %d ", sc.details.BuildEpoch)

		if err := sc.server.ListenAndServe(); err != nil {
			LOG.Error(err)
		}
	}()
}

func (sc *ServerCore) Stop() {
	if sc.shutdownCallback() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		sc.server.Shutdown(ctx)
		// Optionally, you could run srv.Shutdown in a goroutine and block on
		// <-ctx.Done() if your application should wait for other services
		// to finalize based on context cancellation.
		fmt.Println("shutting down")
		os.Exit(0)
	}
}

func (sc *ServerCore) logoPrint() {
	logoData, err := Asset("resources/boot.txt")
	if err != nil {
		logging.GetInstance().Error(err)
	} else {
		fmt.Print(string(logoData))
		fmt.Println()
	}
}
