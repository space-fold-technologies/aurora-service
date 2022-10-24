package clusters

import (
	"io/ioutil"
	"net/http"

	"github.com/space-fold-technologies/aurora-service/app/core/server/http/controllers"
	"github.com/space-fold-technologies/aurora-service/app/core/server/registry"
	"google.golang.org/protobuf/proto"
)

const BASE_PATH = "/api/v1/aurora-service/clusters"

type ClusterController struct {
	*controllers.ControllerBase
	service *ClusterService
}

func NewController(service *ClusterService) controllers.HTTPController {
	return &ClusterController{service: service}
}

func (cc *ClusterController) Name() string {
	return "cluster-controller"
}

func (cc *ClusterController) Initialize(RouteRegistry registry.RouterRegistry) {
	RouteRegistry.AddRestricted(
		BASE_PATH+"/create",
		[]string{"cluster.create"},
		"POST",
		cc.create,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/update",
		[]string{"cluster.update"},
		"PUT",
		cc.update,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/list",
		[]string{"cluster.read"},
		"GET",
		cc.list,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/{name}/remove",
		[]string{"cluster.remove"},
		"DELETE",
		cc.remove,
	)

}

func (cc *ClusterController) create(w http.ResponseWriter, r *http.Request) {
	order := &CreateClusterOrder{}
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		cc.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		cc.BadRequest(w, err)
	} else if err := cc.service.Create(order); err != nil {
		cc.ServiceFailure(w, err)
	} else {
		cc.OKNoResponse(w)
	}
}

func (cc *ClusterController) update(w http.ResponseWriter, r *http.Request) {
	order := &UpdateClusterOrder{}
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		cc.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		cc.BadRequest(w, err)
	} else if err := cc.service.Update(order); err != nil {
		cc.ServiceFailure(w, err)
	} else {
		cc.OKNoResponse(w)
	}
}

func (cc *ClusterController) list(w http.ResponseWriter, r *http.Request) {
	principals := cc.GetPrincipals(r)
	if results, err := cc.service.List(principals); err != nil {
		cc.ServiceFailure(w, err)
	} else {
		cc.OK(w, results)
	}
}

func (cc *ClusterController) remove(w http.ResponseWriter, r *http.Request) {
	name := cc.GetVar("name", r)
	if err := cc.service.Remove(name); err != nil {
		cc.ServiceFailure(w, err)
	} else {
		cc.OKNoResponse(w)
	}
}
