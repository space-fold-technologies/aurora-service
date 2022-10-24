package nodes

import (
	"io/ioutil"
	"net/http"

	"github.com/space-fold-technologies/aurora-service/app/core/server/http/controllers"
	"github.com/space-fold-technologies/aurora-service/app/core/server/registry"
	"google.golang.org/protobuf/proto"
)

const BASE_PATH = "/api/v1/aurora-service/nodes"

type NodeController struct {
	*controllers.ControllerBase
	service *NodeService
}

func NewController(service *NodeService) controllers.HTTPController {
	return &NodeController{service: service}
}

func (nc *NodeController) Name() string {
	return "node-controller"
}

func (nc *NodeController) Initialize(RouteRegistry registry.RouterRegistry) {
	RouteRegistry.AddRestricted(
		BASE_PATH+"/create",
		[]string{"node.create"},
		"POST",
		nc.create,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/update",
		[]string{"node.update"},
		"PUT",
		nc.update,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/{cluster-name}/list",
		[]string{"node.update"},
		"GET",
		nc.list,
	)

	RouteRegistry.AddRestricted(
		BASE_PATH+"/{node-name}",
		[]string{"node.remove"},
		"DELETE",
		nc.remove,
	)
}

func (nc *NodeController) create(w http.ResponseWriter, r *http.Request) {
	order := &CreateNodeOrder{}
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		nc.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		nc.BadRequest(w, err)
	} else if err := nc.service.Create(order); err != nil {
		nc.ServiceFailure(w, err)
	} else {
		nc.OKNoResponse(w)
	}
}

func (nc *NodeController) update(w http.ResponseWriter, r *http.Request) {
	order := &UpdateNodeOrder{}
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		nc.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, order); err != nil {
		nc.BadRequest(w, err)
	} else if err := nc.service.Update(order); err != nil {
		nc.ServiceFailure(w, err)
	} else {
		nc.OKNoResponse(w)
	}
}

func (nc *NodeController) list(w http.ResponseWriter, r *http.Request) {
	cluster := nc.GetVar("cluster-name", r)
	if results, err := nc.service.List(cluster); err != nil {
		nc.ServiceFailure(w, err)
	} else {
		nc.OK(w, results)
	}
}

func (nc *NodeController) remove(w http.ResponseWriter, r *http.Request) {
	if err := nc.service.Remove(nc.GetVar("node-name", r)); err != nil {
		nc.ServiceFailure(w, err)
	} else {
		nc.OKNoResponse(w)
	}
}
