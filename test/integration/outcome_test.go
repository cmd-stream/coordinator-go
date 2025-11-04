//go:build integration
// +build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	coord "github.com/cmd-stream/coordinator-go"
	coordman "github.com/cmd-stream/coordinator-go/manager"
	"github.com/cmd-stream/coordinator-go/support"
	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/testkit-go/fixtures/coordinator/cmds"
	"github.com/cmd-stream/testkit-go/fixtures/coordinator/codecs"
	rcvr "github.com/cmd-stream/testkit-go/fixtures/coordinator/receiver"
	"github.com/cmd-stream/testkit-go/fixtures/coordinator/results"
	"github.com/cmd-stream/testkit-go/helpers"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestOutcome(t *testing.T) {
	const (
		addr  = "127.0.0.1:9002"
		delta = 100 * time.Millisecond
	)

	// Start coordinator.
	var (
		storage  = support.NewMemoryStorage(support.NewBatcher(coord.CheckpointSize))
		receiver = &rcvr.Receiver{}
	)
	err := StartServer(addr, receiver, storage)
	assertfatal.EqualError(err, nil, t)
	time.Sleep(100 * time.Millisecond)

	client, err := makeClient(addr)
	assertfatal.EqualError(err, nil, t)

	cmd := cmds.Cmd{OutcomesCount: 3}
	var (
		cmdSeq   core.Seq = 1
		wantSend          = helpers.WantSend{
			Seq: cmdSeq,
			N:   codecs.CmdSize(cmdSeq, cmd),
		}
		wantReceive = helpers.WantReceive{
			AsyncResult: codecs.AsyncResult(cmdSeq, results.Result(3)),
		}
	)
	err = helpers.Exchange(cmd, NewSendFn(client), ReceiveFn, wantSend, wantReceive)
	assertfatal.EqualError(err, nil, t)

	// TODO Check storage.
	storage.Load(func(partID coordman.PartID, partSeq coordman.PartSeq,
		bs []byte,
	) error {
		fmt.Println(partSeq)
		fmt.Println(bs)
		return nil
	})
	fmt.Println(storage.Checkpoint())
}
