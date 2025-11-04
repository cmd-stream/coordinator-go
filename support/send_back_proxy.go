package support

import (
	"github.com/cmd-stream/core-go"
	dsrv "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
)

func SendBackProxyFactory[T any](transport dsrv.Transport[T]) core.Proxy {
	return handler.NewProxy(transport)
}
