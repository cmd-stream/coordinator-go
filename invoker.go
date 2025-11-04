package coordinator

import (
	"context"
	"fmt"
	"time"

	man "github.com/cmd-stream/coordinator-go/manager"
	"github.com/cmd-stream/core-go"

	"github.com/cmd-stream/handler-go"
)

func NewInvoker[T any](manager *man.Manager[T],
	invoker handler.Invoker[T],
) Invoker[T] {
	return Invoker[T]{manager, invoker}
}

type Invoker[T any] struct {
	manager *man.Manager[T]
	invoker handler.Invoker[T]
}

func (i Invoker[T]) Invoke(ctx context.Context, seq core.Seq, at time.Time,
	bytesRead int,
	cmd core.Cmd[T],
	proxy core.Proxy,
) (err error) {
	err = i.invoker.Invoke(ctx, seq, at, bytesRead, cmd, proxy)
	switch err {
	case ErrCmdDelayed:
		i.manager.DelayCmd(man.CmdRecord[T]{Seq: seq, At: at, Cmd: cmd})
		err = nil
	case ErrCmdBlocked:
		i.manager.Suspend()
		i.manager.DelayCmd(man.CmdRecord[T]{Seq: seq, At: at, Cmd: cmd})
		err = nil
	default:
		compCmd, ok := cmd.(man.Cmd[T])
		if ok {
			svcErr := i.manager.CompleteCmd(compCmd.Workflow().ID())
			if svcErr != nil {
				// TODO
				err = fmt.Errorf("orchestrator: failed to complete cmd, cause: %w",
					svcErr)
			}
		}
	}
	return
}
