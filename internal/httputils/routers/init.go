package routers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
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
