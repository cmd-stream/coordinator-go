//go:build integration
// +build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	coord "github.com/cmd-stream/coordinator-go"
	"github.com/cmd-stream/core-go"

	"github.com/cmd-stream/testkit-go/fixtures/coordinator/cmds"
	"github.com/cmd-stream/testkit-go/fixtures/coordinator/codecs"
	"github.com/cmd-stream/testkit-go/helpers"

	rcvr "github.com/cmd-stream/testkit-go/fixtures/coordinator/receiver"

	coordman "github.com/cmd-stream/coordinator-go/manager"
	"github.com/cmd-stream/coordinator-go/support"

	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestDelay(t *testing.T) {
	const (
		addr  = "127.0.0.1:9001"
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

	// Make client.
	client, err := makeClient(addr)
	assertfatal.EqualError(err, nil, t)

	// Send Command.
	cmd := cmds.DelayCmd{DelaysCount: 2}
	var (
		cmdSeq   core.Seq = 1
		wantSend          = helpers.WantSend{
			Seq: cmdSeq,
			N:   codecs.DelayCmdSize(cmdSeq, cmd),
		}
		wantReceive = helpers.WantReceive{
			Err: ErrTimeout,
		}

		start    = time.Now()
		execTime = coordman.SlowRetryInterval * time.Duration(cmd.DelaysCount)
	)
	err = helpers.Exchange(cmd, NewSendFn(client), ReceiveFn, wantSend, wantReceive)
	assertfatal.EqualError(err, nil, t)

	// Wait for Command execution.
	time.Sleep(execTime + 50*time.Millisecond)

	// Check Command execution time.
	var (
		end  = receiver.DoneAt()
		want = start.Add(execTime)
	)
	assertfatal.SameTime(end, want, delta, t)

	// TODO Check storage.
	storage.Load(func(partID coordman.PartID, partSeq coordman.PartSeq, bs []byte) error {
		fmt.Println(partSeq)
		fmt.Println(bs)
		return nil
	})
	fmt.Println(storage.Checkpoint())
}
