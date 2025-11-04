package coordinator

import (
	srv "github.com/cmd-stream/cmd-stream-go/server"
	man "github.com/cmd-stream/coordinator-go/manager"
	"github.com/cmd-stream/coordinator-go/support"
	csrv "github.com/cmd-stream/core-go/server"
	dsrv "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
)

const (
	MaxCmds              = 100
	CheckpointSize int64 = 1
)

func Make[T any](invoker handler.Invoker[T], codec srv.Codec[T],
	outcomeSer man.OutcomeSerializer,
	ops ...SetOption[T],
) (server *csrv.Server, err error) {
	o := Options[T]{
		Info:                 srv.ServerInfo,
		MaxCmds:              MaxCmds, // TODO
		ImmediateResultProxy: ProxyFactoryFn[T](support.SendBackProxyFactory[T]),
		DelayedCmdProxy:      support.DummyProxy{},
	}
	Apply(ops, &o)
	// TODO Remove srv.Options{}
	service, err := man.NewManager(srv.AdaptCodec(codec, srv.Options{}),
		outcomeSer, o.MaxCmds, o.Manager...)
	if err != nil {
		return
	}
	var (
		inv = NewInvoker(service, invoker)
		// TODO Remove srv.Options{}
		f = NewTransportFactory(srv.AdaptCodec(codec, srv.Options{}),
			o.Transport...)
		h  = NewHandler(service, inv, o.Handler...)
		d  = dsrv.New(o.Info, f, h, o.Delegate...)
		wf = NewWorkersFactory(service, inv)
	)
	server = csrv.NewWithWorkers(d, wf, o.Core...)
	return
}
