package coordinator

import (
	"context"
	"sync"

	man "github.com/cmd-stream/coordinator-go/manager"
	"github.com/cmd-stream/coordinator-go/support"
)

func NewRetryWorker[T any](manager *man.Manager[T],
	invoker Invoker[T],
) RetryWorker[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return RetryWorker[T]{ctx, cancel, manager, invoker}
}

type RetryWorker[T any] struct {
	ctx     context.Context
	cancel  context.CancelFunc
	manager *man.Manager[T]
	invoker Invoker[T]
}

func (w RetryWorker[T]) Run() (err error) {
	var (
		errs = make(chan error, 1)
		wg   = &sync.WaitGroup{}
	)
	wg.Add(1)
	go retryCmds(w.ctx, w.manager, w.invoker, errs, wg)

	select {
	case <-w.ctx.Done():
		err = context.Canceled
	case err = <-errs:
	}
	w.cancel()
	wg.Wait()
	return
}

func (w RetryWorker[T]) Stop() (err error) {
	w.cancel()
	return
}

func retryCmds[T any](ctx context.Context, manager *man.Manager[T],
	invoker Invoker[T],
	errs chan<- error,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	var (
		record man.CmdRecord[T]
		proxy  = support.DummyProxy{}
		err    error
	)
	for {
		record, err = manager.NextDelayedCmd(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			queueErr(err, errs)
			return
		}
		wg.Add(1)
		go invokeCmd(ctx, record.Seq, record.At, record.N, record.Cmd, invoker,
			proxy, errs, wg)
	}
}
