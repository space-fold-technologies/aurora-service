package users

import (
	"io/ioutil"
	"net/http"

	"github.com/space-fold-technologies/aurora-service/app/core/server/http/controllers"
	"github.com/space-fold-technologies/aurora-service/app/core/server/registry"
	"google.golang.org/protobuf/proto"
)

const BASE_PATH = "/api/v1/aurora-service/users"

type UserController struct {
	*controllers.ControllerBase
	service *UserService
}

func NewController(service *UserService) controllers.HTTPController {
	return &UserController{service: service}
}

func (uc *UserController) Name() string {
	return "user-controller"
}

func (uc *UserController) Initialize(RouteRegistry registry.RouterRegistry) {
	RouteRegistry.AddRestricted(
		BASE_PATH+"/create",
		[]string{"user.create"},
		"POST",
		uc.create,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/list",
		[]string{"user.read"},
		"POST",
		uc.list,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/{email}/remove",
		[]string{"user.remove"},
		"DELETE",
		uc.remove,
	)
}

func (uc *UserController) create(w http.ResponseWriter, r *http.Request) {
	order := &CreateUserOrder{}
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		uc.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		uc.BadRequest(w, err)
	} else if err := uc.service.Create(order); err != nil {
		uc.ServiceFailure(w, err)
	} else {
		uc.OKNoResponse(w)
	}
}

func (uc *UserController) list(w http.ResponseWriter, r *http.Request) {
	principals := uc.GetPrincipals(r)
	if results, err := uc.service.List(principals); err != nil {
		uc.ServiceFailure(w, err)
	} else {
		uc.OK(w, results)
	}
}

func (uc *UserController) remove(w http.ResponseWriter, r *http.Request) {
	email := uc.GetVar("email", r)
	if err := uc.service.Remove(email); err != nil {
		uc.ServiceFailure(w, err)
	} else {
		uc.OKNoResponse(w)
	}
}
