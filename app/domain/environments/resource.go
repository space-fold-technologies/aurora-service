package environments

import (
	"io"
	"net/http"

	"github.com/space-fold-technologies/aurora-service/app/core/server/http/controllers"
	"github.com/space-fold-technologies/aurora-service/app/core/server/registry"
	"google.golang.org/protobuf/proto"
)

const BASE_PATH = "/api/v1/aurora-service/environments"

type EnvironmentController struct {
	*controllers.ControllerBase
	service *EnvironmentService
}

func NewController(service *EnvironmentService) controllers.HTTPController {
	return &EnvironmentController{service: service}
}

func (ec *EnvironmentController) Name() string {
	return "environment-controller"
}

func (ec *EnvironmentController) Initialize(RouteRegistry registry.RouterRegistry) {
	RouteRegistry.AddRestricted(
		BASE_PATH+"/create",
		[]string{"environments.create"},
		"PUT",
		ec.create,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/{scope}/{target}/list",
		[]string{"environments.information"},
		"GET",
		ec.list,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/remove",
		[]string{"environments.remove"},
		"PUT",
		ec.remove,
	)
}

func (ec *EnvironmentController) create(w http.ResponseWriter, r *http.Request) {
	order := &CreateEnvEntryOrder{}
	if data, err := io.ReadAll(r.Body); err != nil {
		ec.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		ec.BadRequest(w, err)
	} else if err := ec.service.Create(order); err != nil {
		ec.ServiceFailure(w, err)
	} else {
		ec.OKNoResponse(w)
	}
}

func (ec *EnvironmentController) list(w http.ResponseWriter, r *http.Request) {
	principals := ec.GetPrincipals(r)
	scope := ec.GetVar("scope", r)
	target := ec.GetVar("target", r)
	if results, err := ec.service.List(principals, scope, target); err != nil {
		ec.ServiceFailure(w, err)
	} else {
		ec.OK(w, results)
	}
}

func (ec *EnvironmentController) remove(w http.ResponseWriter, r *http.Request) {
	order := &RemoveEnvEntryOrder{}
	if data, err := io.ReadAll(r.Body); err != nil {
		ec.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		ec.BadRequest(w, err)
	} else if err := ec.service.Remove(order); err != nil {
		ec.ServiceFailure(w, err)
	} else {
		ec.OKNoResponse(w)
	}
}
