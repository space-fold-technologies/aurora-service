package registry

import (
	"net/http"
)

type RouterRegistry interface {
	Add(Path string, IsRestricted bool, Method string, HttpHandler func(http.ResponseWriter, *http.Request))
	AddRestricted(Path string, Roles []string, Method string, HttpHandler func(http.ResponseWriter, *http.Request))
	Initialize()
}
