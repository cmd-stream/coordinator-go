package integration

import (
	"errors"
	"time"

	srv "github.com/cmd-stream/cmd-stream-go/server"
	coord "github.com/cmd-stream/coordinator-go"
	coordman "github.com/cmd-stream/coordinator-go/manager"
	"github.com/cmd-stream/coordinator-go/support"
	"github.com/cmd-stream/core-go"
	ccln "github.com/cmd-stream/core-go/client"
	"github.com/cmd-stream/testkit-go/fixtures/coordinator/codecs"
	rcvr "github.com/cmd-stream/testkit-go/fixtures/coordinator/receiver"
	"github.com/cmd-stream/testkit-go/helpers"
)

func StartServer(addr string, receiver *rcvr.Receiver,
	storage *support.MemoryStorage,
) (err error) {
	var (
		invoker    = srv.NewInvoker(receiver)
		outcomeSer = support.OutcomeSerializerJSON{}
	)
	coordinator, err := coord.Make(invoker, codecs.ServerCodec{}, outcomeSer,
		coord.WithManager[*rcvr.Receiver](
			coordman.WithStorage[*rcvr.Receiver](storage),
		),
	)
	if err != nil {
		return
	}
	go func() {
		coordinator.ListenAndServe(addr)
	}()
	return
}

func NewSendFn(client *ccln.Client[*rcvr.Receiver]) helpers.SendFn[*rcvr.Receiver] {
	return func(cmd core.Cmd[*rcvr.Receiver], cmdResults chan<- core.AsyncResult) (
		seq core.Seq, n int, err error,
	) {
		return client.Send(cmd, cmdResults)
	}
}

var (
	ErrTimeout = errors.New("timeout")

	ReceiveFn = func(results <-chan core.AsyncResult) (asyncResult core.AsyncResult,
		err error,
	) {
		select {
		case <-time.NewTimer(coordman.SlowRetryInterval).C:
			err = ErrTimeout
		case asyncResult = <-results:
		}
		return
	}
)
