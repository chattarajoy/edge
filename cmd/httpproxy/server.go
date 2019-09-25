package httpproxy

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/heptio/workgroup"

	"github.com/chattarajoy/edge/internal/httputils/routers"
)

type serverInput struct {
	Port            int
	Logger          log.Logger
	Router          routers.Router
	ServerDrainTime int
	NotFoundHandler http.Handler
	Db              *sql.DB
}

type httpServer struct {
	*serverInput
}

func (server *httpServer) routes() {
	server.handle("GET", "/", http.HandlerFunc(server.homeHandler))
}

func (server *httpServer) handle(method, path string, handler http.Handler) {
	server.Router.Handle(method, path, server.wrappedHandlers(handler))
	server.Router.NotFound(server.wrappedHandlers(server.NotFoundHandler))
}

func (server *httpServer) homeHandler(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `
	{
		"name": "httpproxy",
		"status": "OK"
	}`
	_, _ = fmt.Fprintf(w, htmlIndex)
}

func createAndRunServer(group *workgroup.Group, inp *serverInput) error {
	hServer := &httpServer{serverInput: inp}
	hServer.routes()

	server := &http.Server{Addr: fmt.Sprintf(":%d", inp.Port), Handler: hServer.Router}
	_ = inp.Logger.Log("Listening on port", inp.Port)
	serveErr := server.ListenAndServe()
	_ = inp.Db.Close()
	return serveErr
}
