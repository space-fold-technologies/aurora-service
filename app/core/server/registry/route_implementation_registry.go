package registry

import (
	ctx "context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/security"
)

type RouteEntry struct {
	Path         string
	IsRestricted bool
	Permissions  []string
	Method       string
	HttpHandler  func(http.ResponseWriter, *http.Request)
}

type RegistryRoute struct {
	Path              string
	ExpressionValue   string
	Method            string
	RegularExpression *regexp.Regexp
	IsRestricted      bool
	Permissions       []string
}

func NewRegistryRoute(Path, ExpressionValue, Method string, IsRestricted bool, Permissions []string) *RegistryRoute {
	registryRoute := &RegistryRoute{
		Path:            Path,
		ExpressionValue: ExpressionValue,
		Method:          Method,
		IsRestricted:    IsRestricted,
		Permissions:     Permissions,
	}
	registryRoute.init()
	return registryRoute
}

func (aRx *RegistryRoute) init() *RegistryRoute {
	if len(strings.TrimSpace(aRx.ExpressionValue)) > 0 {
		aRx.RegularExpression = regexp.MustCompile(aRx.ExpressionValue)
	} else {
		aRx.RegularExpression = regexp.MustCompile(aRx.Path)
	}
	return aRx
}

func (aRx *RegistryRoute) IsMatch(Path, Method string) bool {
	submatches := aRx.RegularExpression.FindStringSubmatch(Path)
	if len(submatches) == 0 {
		return false
	}
	return Method == aRx.Method
}

type AuthorizedRouterRegistry struct {
	Router         *mux.Router
	routeEntries   []*RouteEntry
	registryRoutes []*RegistryRoute
	TokenHandler   security.TokenHandler
}

func NewRouteRegistry(Router *mux.Router, TokenHandler security.TokenHandler) RouterRegistry {
	authorizedRouterRegistry := &AuthorizedRouterRegistry{
		Router:       Router,
		TokenHandler: TokenHandler,
	}
	authorizedRouterRegistry.setUp()
	return authorizedRouterRegistry
}

func (ar *AuthorizedRouterRegistry) setUp() *AuthorizedRouterRegistry {
	ar.routeEntries = make([]*RouteEntry, 0)
	ar.registryRoutes = make([]*RegistryRoute, 0)
	return ar
}

func (ar *AuthorizedRouterRegistry) Add(
	Path string,
	IsRestricted bool,
	Method string,
	HttpHandler func(http.ResponseWriter, *http.Request)) {
	ar.routeEntries = append(ar.routeEntries, &RouteEntry{
		Path:         Path,
		IsRestricted: IsRestricted,
		Permissions:  []string{},
		Method:       Method,
		HttpHandler:  HttpHandler,
	})
}

func (ar *AuthorizedRouterRegistry) AddRestricted(
	Path string,
	Permissions []string,
	Method string,
	HttpHandler func(http.ResponseWriter, *http.Request)) {
	ar.routeEntries = append(ar.routeEntries, &RouteEntry{
		Path:         Path,
		IsRestricted: true,
		Permissions:  Permissions,
		Method:       Method,
		HttpHandler:  HttpHandler,
	})
}

func (ar *AuthorizedRouterRegistry) Initialize() {
	logging.GetInstance().Info("Setting Up controller routes")
	for _, entry := range ar.routeEntries {
		r := ar.Router.HandleFunc(entry.Path, entry.HttpHandler).Methods(entry.Method)
		expr, _ := r.GetPathRegexp()
		ar.registryRoutes = append(ar.registryRoutes, NewRegistryRoute(entry.Path, expr, entry.Method, entry.IsRestricted, entry.Permissions))
		logging.GetInstance().Info("ROUTE ->->", entry.Path)
	}
	ar.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ar.verifyCall(next, w, r)
		})
	})
}

func (ar *AuthorizedRouterRegistry) fetchAuthorizationHeader(r *http.Request) string {
	tokenHeader := r.Header.Get("Authorization")
	if len(tokenHeader) == 0 {
		tokenHeader = r.Header.Values("Authorization")[0]
	}
	return tokenHeader
}

func (ar *AuthorizedRouterRegistry) verifyCall(next http.Handler, w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path
	for _, value := range ar.registryRoutes {
		if value.IsMatch(requestPath, r.Method) && !value.IsRestricted {
			r.WithContext(ctx.Background())
			next.ServeHTTP(w, r)
			return
		}
	}

	tokenHeader := ar.fetchAuthorizationHeader(r) //Grab the token from the header
	splitTokenHeader := strings.Split(tokenHeader, " ")
	if len(splitTokenHeader) < 2 { //Token is missing, returns with error code 401 Unauthorized
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":       "AUTHORIZATION ERROR",
			"details":     "This request is not valid",
			"occurred_at": time.Now().Format(time.RFC3339),
			"code":        401,
		})
		return
	}
	claims, err := ar.TokenHandler.VerifyToken(splitTokenHeader[1])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":       "AUTHORIZATION ERROR",
			"details":     "Could not validate this request",
			"occurred_at": time.Now().Format(time.RFC3339),
			"code":        401,
		})
		return
	}
	permissions := ar.matchRoute(requestPath, r.Method) // Then check if the role in the claim is in the list of roles

	logging.GetInstance().Info("MATCHED PATH : " + requestPath)
	if ar.isClaimValid(permissions, claims) {
		context.Set(r, "principals", claims.ToMap())
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":       "AUTHORIZATION ERROR",
			"details":     "Could not validate this request",
			"occurred_at": time.Now().Format(time.RFC3339),
			"code":        401,
		})
	}
}

func (ar *AuthorizedRouterRegistry) matchRoute(Path, Method string) (Permissions []string) {
	for _, value := range ar.registryRoutes {
		if value.IsMatch(Path, Method) {
			return value.Permissions
		}
	}
	return []string{}
}

func (ar *AuthorizedRouterRegistry) isClaimValid(Permissions []string, Claim *security.Claims) bool {
	if len(Permissions) == 0 {
		return true
	}

	for _, role := range Claim.Permissions {
		for _, permission := range Permissions {
			if strings.TrimSpace(role) == strings.TrimSpace(permission) {
				return true
			}
		}
	}
	return false
}
