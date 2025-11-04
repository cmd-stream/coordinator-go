package coordinator

import (
	"context"
	"sync"
	"time"

	man "github.com/cmd-stream/coordinator-go/manager"
	"github.com/cmd-stream/core-go"
	csrv "github.com/cmd-stream/core-go/server"
	dsrv "github.com/cmd-stream/delegate-go/server"
	"github.com/cmd-stream/handler-go"
)

// NewHandler creates a new Handler.
func NewHandler[T any](manager *man.Manager[T], invoker Invoker[T],
	ops ...handler.SetOption,
) *Handler[T] {
	h := Handler[T]{manager: manager, invoker: invoker}
	handler.Apply(ops, &h.options)
	return &h
}

// Handler implements the delegate.ServerTransportHandler interface.
//
// It receives Commands sequentially and executes each in a separate goroutine
// using the Invoker.
//
// If an error occurs, the Handler closes the transport connection.
type Handler[T any] struct {
	manager *man.Manager[T]
	invoker Invoker[T]
	options handler.Options
}

func (h *Handler[T]) Handle(ctx context.Context, transport dsrv.Transport[T]) (
	err error,
) {
	var (
		wg             = &sync.WaitGroup{}
		ownCtx, cancel = context.WithCancel(ctx)
		errs           = make(chan error, 1)
	)
	wg.Add(1)
	go receiveCmdAndInvoke(ownCtx, transport.(Transport[T]), h.manager,
		h.invoker, errs, wg, h.options)

	select {
	case <-ownCtx.Done():
		err = context.Canceled
	case err = <-errs:
		if err == man.ErrAnotherOrchestratorDetected {
			err = csrv.ErrClosed
		}
	}
	cancel()
	if err := transport.Close(); err != nil {
		panic(err)
	}
	wg.Wait()
	return
}

func receiveCmdAndInvoke[T any](ctx context.Context,
	transport Transport[T],
	manager *man.Manager[T],
	invoker Invoker[T],
	errs chan<- error,
	wg *sync.WaitGroup,
	options handler.Options,
) {
	defer wg.Done()
	var (
		seq   core.Seq
		cmd   core.Cmd[T]
		bs    []byte
		err   error
		at    time.Time
		proxy = handler.NewProxy(transport)
	)
	for {
		if options.CmdReceiveDuration != 0 {
			deadline := time.Now().Add(options.CmdReceiveDuration)
			if err = transport.SetReceiveDeadline(deadline); err != nil {
				queueErr(err, errs)
				return
			}
		}
		seq, cmd, bs, err = transport.ReceiveWithBytes()
		if err != nil {
			queueErr(err, errs)
			return
		}
		if options.At {
			at = time.Now()
		}
		cmd, err = prepareCmd(manager, at, cmd, bs)
		if err != nil {
			if err == man.ErrServiceSuspended {
				continue
			}
			queueErr(err, errs)
			return
		}
		wg.Add(1)
		go invokeCmd(ctx, seq, at, len(bs), cmd, invoker, proxy, errs, wg)
	}
}

func prepareCmd[T any](manager *man.Manager[T], at time.Time, cmd core.Cmd[T],
	bs []byte,
) (core.Cmd[T], error) {
	ocmd, ok := cmd.(man.Cmd[T])
	if ok {
		id, err := manager.SaveCmd(at, bs)
		if err != nil {
			return nil, err
		}
		cmd = ocmd.SetWorkflow(man.NewWorkflow(id, manager))
	}
	return cmd, nil
}

func invokeCmd[T any](ctx context.Context, seq core.Seq, at time.Time,
	bytesRead int,
	cmd core.Cmd[T],
	invoker Invoker[T],
	proxy core.Proxy,
	errs chan<- error,
	wg *sync.WaitGroup,
) {
	err := invoker.Invoke(ctx, seq, at, bytesRead, cmd, proxy)
	if err != nil {
		queueErr(err, errs)
	}
	wg.Done()
}

func queueErr(err error, errs chan<- error) {
	select {
	case errs <- err:
	default:
	}
}
