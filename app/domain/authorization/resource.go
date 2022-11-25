package authorization

import (
	"io"
	"net/http"

	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/server/http/controllers"
	"github.com/space-fold-technologies/aurora-service/app/core/server/registry"
	"google.golang.org/protobuf/proto"
)

const BASE_PATH = "/api/v1/aurora-service/authorization"

type AuthorizationController struct {
	*controllers.ControllerBase
	service *AuthorizationService
}

func NewController(service *AuthorizationService) controllers.HTTPController {
	return &AuthorizationController{service: service}
}

func (ac *AuthorizationController) Name() string {
	return "authorization-controller"
}

func (ac *AuthorizationController) Initialize(RouteRegistry registry.RouterRegistry) {
	RouteRegistry.Add(
		BASE_PATH+"/session",
		false,
		"POST",
		ac.session,
	)
}

func (ac *AuthorizationController) session(w http.ResponseWriter, r *http.Request) {
	credentials := &Credentials{}
	if data, err := io.ReadAll(r.Body); err != nil {
		ac.BadRequest(w, err)
	} else if err = proto.Unmarshal(data, credentials); err != nil {
		ac.BadRequest(w, err)
	} else if session, err := ac.service.RequestSession(credentials); err != nil {
		logging.GetInstance().Error(err)
		ac.ServiceFailure(w, err)
	} else {
		ac.OK(w, session)
	}
}
