package users

import (
	"io"
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
		[]string{"users.create"},
		"POST",
		uc.create,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/list",
		[]string{"users.information"},
		"GET",
		uc.list,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/{email}/remove",
		[]string{"users.remove"},
		"DELETE",
		uc.remove,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/add-teams",
		[]string{"users.update"},
		"PUT",
		uc.addTeams,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/remove-teams",
		[]string{"users.update"},
		"PUT",
		uc.removeTeams,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/add-permissions",
		[]string{"users.update"},
		"PUT",
		uc.addPermissions,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/remove-permissions",
		[]string{"users.update"},
		"PUT",
		uc.removePermissions,
	)
}

func (uc *UserController) create(w http.ResponseWriter, r *http.Request) {
	order := &CreateUserOrder{}
	if data, err := io.ReadAll(r.Body); err != nil {
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

func (uc *UserController) addTeams(w http.ResponseWriter, r *http.Request) {
	order := &UpdateTeams{}
	if data, err := io.ReadAll(r.Body); err != nil {
		uc.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		uc.BadRequest(w, err)
	} else if err = uc.service.AddTeams(order); err != nil {
		uc.ServiceFailure(w, err)
	} else {
		uc.OKNoResponse(w)
	}
}

func (uc *UserController) removeTeams(w http.ResponseWriter, r *http.Request) {
	order := &UpdateTeams{}
	if data, err := io.ReadAll(r.Body); err != nil {
		uc.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		uc.BadRequest(w, err)
	} else if err = uc.service.RemoveTeams(order); err != nil {
		uc.ServiceFailure(w, err)
	} else {
		uc.OKNoResponse(w)
	}
}

func (uc *UserController) addPermissions(w http.ResponseWriter, r *http.Request) {
	order := &UpdatePermissions{}
	if data, err := io.ReadAll(r.Body); err != nil {
		uc.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		uc.BadRequest(w, err)
	} else if err = uc.service.AddPermissions(order); err != nil {
		uc.ServiceFailure(w, err)
	} else {
		uc.OKNoResponse(w)
	}
}

func (uc *UserController) removePermissions(w http.ResponseWriter, r *http.Request) {
	order := &UpdatePermissions{}
	if data, err := io.ReadAll(r.Body); err != nil {
		uc.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		uc.BadRequest(w, err)
	} else if err = uc.service.RemovePermissions(order); err != nil {
		uc.ServiceFailure(w, err)
	} else {
		uc.OKNoResponse(w)
	}
}
