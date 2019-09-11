package httputils

import "github.com/go-kit/kit/transport"

type Router interface {
	transport.ErrorHandler
}