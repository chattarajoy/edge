package routers

import (
	"bitbucket.org/qubole/edge/pkg/helpers"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type HTTPRouter struct {
	mux   http.Handler
	muxFn func() *httprouter.Router
}

func (router *HTTPRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *HTTPRouter) NotFound(handler http.Handler) {
	router.muxFn().NotFound = handler
}

func (router *HTTPRouter) Handle(method, path string, handler http.Handler) {
	allowedMethods := []string{"head", "get", "post", "put", "patch", "delete", "options"}
	curMethod := strings.ToLower(method)
	if helpers.StringInSlice(curMethod, allowedMethods) {
		router.muxFn().Handler(method, path, handler)
	}
	return
}

func (router *HTTPRouter) Name() string {
	return "HTTP Router"
}

func NewHTTPRouter(mux http.Handler) *HTTPRouter {
	g := &HTTPRouter{mux: mux}
	g.muxFn = func() *httprouter.Router {
		return g.mux.(*httprouter.Router)
	}
	return g
}
