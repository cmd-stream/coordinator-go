package coordinator

import (
	"net"

	man "github.com/cmd-stream/coordinator-go/manager"
	csrv "github.com/cmd-stream/core-go/server"
	"github.com/ymz-ncnk/jointwork-go"
)

func NewWorkersFactory[T any](manager *man.Manager[T],
	invoker Invoker[T],
) WorkersFactory[T] {
	return WorkersFactory[T]{manager, invoker}
}

type WorkersFactory[T any] struct {
	manager *man.Manager[T]
	invoker Invoker[T]
}

func (f WorkersFactory[T]) New(count int, conns <-chan net.Conn,
	delegate csrv.Delegate,
	callback csrv.LostConnCallback,
) (workers []jointwork.Task) {
	workers = make([]jointwork.Task, count+1)
	for i := range count {
		workers[i] = csrv.NewWorker(conns, delegate, callback)
	}
	workers[count] = NewRetryWorker(f.manager, f.invoker)
	return
}
