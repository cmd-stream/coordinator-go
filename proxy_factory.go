package coordinator

import (
	"github.com/cmd-stream/core-go"
	dsrv "github.com/cmd-stream/delegate-go/server"
)

type ProxyFactory[T any] interface {
	New(transport dsrv.Transport[T]) core.Proxy
}

type ProxyFactoryFn[T any] func(transport dsrv.Transport[T]) core.Proxy

func (f ProxyFactoryFn[T]) New(transport dsrv.Transport[T]) core.Proxy {
	return f(transport)
}
