package teams

import (
	"io/ioutil"
	"net/http"

	"github.com/space-fold-technologies/aurora-service/app/core/server/http/controllers"
	"github.com/space-fold-technologies/aurora-service/app/core/server/registry"
	"google.golang.org/protobuf/proto"
)

const BASE_PATH = "/api/v1/aurora-service/teams"

type TeamController struct {
	*controllers.ControllerBase
	service *TeamService
}

func NewController(service *TeamService) controllers.HTTPController {
	return &TeamController{service: service}
}

func (tc *TeamController) Name() string {
	return "team-controller"
}

func (tc *TeamController) Initialize(RouteRegistry registry.RouterRegistry) {
	RouteRegistry.AddRestricted(
		BASE_PATH+"/create",
		[]string{"team.create"},
		"POST",
		tc.create,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/list",
		[]string{"team.read"},
		"GET",
		tc.list,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/{name}/remove",
		[]string{"team.remove"},
		"DELETE",
		tc.remove,
	)
}

func (tc *TeamController) create(w http.ResponseWriter, r *http.Request) {
	order := &CreateTeamOrder{}
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		tc.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		tc.BadRequest(w, err)
	} else if err := tc.service.Create(order); err != nil {
		tc.ServiceFailure(w, err)
	} else {
		tc.OKNoResponse(w)
	}
}

func (tc *TeamController) list(w http.ResponseWriter, r *http.Request) {
	principals := tc.GetPrincipals(r)
	if results, err := tc.service.List(principals); err != nil {
		tc.ServiceFailure(w, err)
	} else {
		tc.OK(w, results)
	}
}

func (tc *TeamController) remove(w http.ResponseWriter, r *http.Request) {
	name := tc.GetVar("name", r)
	if err := tc.service.Remove(name); err != nil {
		tc.ServiceFailure(w, err)
	} else {
		tc.OKNoResponse(w)
	}
}
