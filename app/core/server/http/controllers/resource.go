package controllers

import (
	"encoding/json"
	"log"

	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/security"
	"google.golang.org/protobuf/proto"
)

type ErrorMessage struct {
	Error      string `json:"error"`
	Detail     string `json:"details"`
	Code       int    `json:"code"`
	OccurredAt string `json:"occurred_at"`
}

type ControllerBase struct {
}

type DebugTransport struct{}

func (DebugTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	b, err := httputil.DumpRequestOut(r, false)
	if err != nil {
		return nil, err
	}
	log.Println(string(b))
	return http.DefaultTransport.RoundTrip(r)
}

func (a *ControllerBase) GetAccessToken(r *http.Request) string {
	return strings.Split(r.Header.Get("Authorization"), " ")[1]
}

func (a *ControllerBase) GetPrincipals(r *http.Request) *security.Claims {
	claimsMap := context.Get(r, "principals").(map[string]interface{})
	claims := &security.Claims{}
	claims.FromMap(claimsMap)
	return claims
}

func (a *ControllerBase) GetHeader(key string, r *http.Request) string {
	return r.Header.Get(key)
}

func (a *ControllerBase) ProxyRequest(w http.ResponseWriter, r *http.Request, Base, Path string) {
	url, _ := url.Parse(Base)
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host
	r.URL.Path = Path
	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(w, r)
}

func (a *ControllerBase) Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func (a *ControllerBase) Respond(w http.ResponseWriter, data interface{}, statusCode int) int {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
	return 0
}

func (a *ControllerBase) BadRequest(w http.ResponseWriter, Err error) int {
	payload := &ErrorMessage{Error: "Bad Request", Detail: Err.Error(), Code: http.StatusBadRequest, OccurredAt: time.Now().Format(time.RFC3339)}
	return a.Respond(w, payload, http.StatusBadRequest)
}

func (a *ControllerBase) ServiceFailure(w http.ResponseWriter, Err error) int {
	payload := &ErrorMessage{Error: "Service Expectation Failure", Detail: Err.Error(), Code: http.StatusExpectationFailed, OccurredAt: time.Now().Format(time.RFC3339)}
	return a.Respond(w, payload, http.StatusExpectationFailed)
}

func (a *ControllerBase) WebsocketFailure(ws *websocket.Conn, err error) error {
	payload := &ErrorMessage{Error: "Service Expectation Failure", Detail: err.Error(), Code: http.StatusExpectationFailed, OccurredAt: time.Now().Format(time.RFC3339)}
	if data, err := json.Marshal(payload); err != nil {
		return err
	} else {
		ws.WriteMessage(websocket.TextMessage, data)
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		return ws.Close()
	}
}

func (a *ControllerBase) WebSocketClose(ws *websocket.Conn) error {
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Closing off"))
	return ws.Close()
}

func (a *ControllerBase) Created(w http.ResponseWriter, Msg string) int {
	payload := map[string]string{}
	payload["occurred_at"] = time.Now().Format(time.RFC3339)
	payload["message"] = Msg
	return a.Respond(w, payload, http.StatusCreated)
}

func (a *ControllerBase) OK(w http.ResponseWriter, data proto.Message) {
	w.Header().Add("Content-Type", "application/x-protobuf")
	w.WriteHeader(200)
	if content, err := proto.Marshal(data); err != nil {
		logging.GetInstance().Error(err)
		a.ServiceFailure(w, err)
	} else {
		w.Write(content)
	}
}

func (a *ControllerBase) OKNoResponse(w http.ResponseWriter) int {
	w.WriteHeader(http.StatusNoContent)
	return 0
}

func (a *ControllerBase) GetVar(Key string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[Key]
}

func (a *ControllerBase) GetQueryInt(name string, r *http.Request) int64 {
	keys, ok := r.URL.Query()[name]
	if !ok || len(keys) < 1 {
		log.Println("Url Param 'key' is missing")
		return 0
	}
	i, err := strconv.Atoi(keys[0])
	if err != nil {
		return 0
	}
	return int64(i)
}

func (a *ControllerBase) GetQueryString(name string, r *http.Request) string {
	keys, ok := r.URL.Query()[name]
	if !ok || len(keys) < 1 {
		log.Println("Url Param 'key' is missing")
		return ""
	}

	return string(keys[0])
}

func (a *ControllerBase) GetBool(name string, r *http.Request) bool {
	keys, ok := r.URL.Query()[name]
	if !ok || len(keys) < 1 {
		log.Println("Url Param 'key' is missing")
		return false
	}
	b, err := strconv.ParseBool(keys[0])
	if err != nil {
		return false
	}
	return b
}
