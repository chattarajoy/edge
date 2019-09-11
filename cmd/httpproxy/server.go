package httpproxy

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/heptio/workgroup"
)

type serverInput struct {
	Port            int
	Logger          log.Logger
	ServerDrainTime int
}

type httpServer struct {
	*serverInput
}

func createServer(group *workgroup.Group, inp *serverInput) error {
	hServer := &httpServer{serverInput: inp}
	server := &http.Server{Addr: fmt.Sprintf(":%d", inp.Port), Hander: in.Router}
}


//func (server *httpServer) setupRoutes() {
//	server.handle("GET", )
//}