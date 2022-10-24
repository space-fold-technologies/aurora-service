package controllers

import "github.com/space-fold-technologies/aurora-service/app/core/server/registry"

type HTTPController interface {
	Name() string
	Initialize(RouteRegistry registry.RouterRegistry)
}
