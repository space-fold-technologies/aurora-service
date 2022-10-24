package controllers

import (
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/server/registry"
)

type HTTPControllerRegistry struct {
	routeRegistry registry.RouterRegistry
	controllers   map[string]HTTPController
}

func New(routeRegistry registry.RouterRegistry) *HTTPControllerRegistry {
	return &HTTPControllerRegistry{
		routeRegistry: routeRegistry,
		controllers:   make(map[string]HTTPController),
	}
}

func (cr *HTTPControllerRegistry) AddController(controller HTTPController) {
	cr.controllers[controller.Name()] = controller
}

func (cr *HTTPControllerRegistry) InitializeControllers() {
	LOG := logging.GetInstance()
	for name, controller := range cr.controllers {
		LOG.Infof("CONTROLLER -> %s", name)
		controller.Initialize(cr.routeRegistry)
	}
}
