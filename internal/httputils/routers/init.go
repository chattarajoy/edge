package routers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ErrorHandler interface {
	NotFound(http.Handler)
}

type Router interface {
	ErrorHandler
	http.Handler
	Handle(string, string, http.Handler)
	Name() string
}


func CreateRouter(routerName string) Router {
	if routerName == "httprouter" {
		return NewHTTPRouter(httprouter.New())
	}
	return nil
}
