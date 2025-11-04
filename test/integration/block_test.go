//go:build integration
// +build integration

package integration

import (
	"fmt"
	"net"
	"testing"
	"time"

	cmdstream "github.com/cmd-stream/cmd-stream-go"
	coord "github.com/cmd-stream/coordinator-go"
	coordman "github.com/cmd-stream/coordinator-go/manager"
	"github.com/cmd-stream/coordinator-go/support"
	"github.com/cmd-stream/core-go"
	ccln "github.com/cmd-stream/core-go/client"
	"github.com/cmd-stream/testkit-go/fixtures/coordinator/cmds"
	"github.com/cmd-stream/testkit-go/fixtures/coordinator/codecs"
	rcvr "github.com/cmd-stream/testkit-go/fixtures/coordinator/receiver"
	"github.com/cmd-stream/testkit-go/fixtures/coordinator/results"
	helpers "github.com/cmd-stream/testkit-go/helpers"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestBlock(t *testing.T) {
	const (
		addr  = "127.0.0.1:9000"
		delta = 100 * time.Millisecond
	)

	var (
		storage  = support.NewMemoryStorage(support.NewBatcher(coord.CheckpointSize))
		receiver = &rcvr.Receiver{}
	)
	err := StartServer(addr, receiver, storage)
	assertfatal.EqualError(err, nil, t)
	time.Sleep(100 * time.Millisecond)

	client, err := makeClient(addr)
	assertfatal.EqualError(err, nil, t)
	sendFn := NewSendFn(client)

	// Send first cmd.
	cmd1 := cmds.DelayCmd{DelaysCount: 5, BlockOn: 2}
	var (
		cmdSeq1   core.Seq = 1
		wantSend1          = helpers.WantSend{
			Seq: cmdSeq1,
			N:   codecs.DelayCmdSize(cmdSeq1, cmd1),
		}
		wantReceive1 = helpers.WantReceive{Err: ErrTimeout}

		start1    = time.Now()
		execTime1 = coordman.SlowRetryInterval * time.Duration(cmd1.DelaysCount)
	)
	err = helpers.Exchange(cmd1, sendFn, ReceiveFn, wantSend1, wantReceive1)
	assertfatal.EqualError(err, nil, t)

	// Wait for the coordinator to block.
	time.Sleep(coordman.SlowRetryInterval * time.Duration(cmd1.BlockOn))

	// Try to send a Command while the coordinator is blocked.
	cmd2 := cmds.Cmd{OutcomesCount: 1}
	var (
		cmdSeq2   core.Seq = 2
		wantSend2          = helpers.WantSend{
			Seq: cmdSeq2,
			N:   codecs.CmdSize(cmdSeq2, cmd2),
		}
		wantReceive2 = helpers.WantReceive{
			Err: ErrTimeout,
		}
	)
	err = helpers.Exchange(cmd2, sendFn, ReceiveFn, wantSend2, wantReceive2)
	assertfatal.EqualError(err, nil, t)

	// No one Command was executed.
	time.Sleep(50 * time.Millisecond)
	assertfatal.SameTime(receiver.DoneAt(), time.Time{}, delta, t)

	// Wait for the coordinator to unblock - at least one already received
	// Command should succeed.
	time.Sleep(execTime1 + 50*time.Millisecond)

	// Check first Command execution time.
	var (
		end1  = receiver.DoneAt()
		want1 = start1.Add(execTime1)
	)
	assertfatal.SameTime(end1, want1, delta, t)

	// Send another Command. Now the execution is ok.
	cmd3 := cmds.Cmd{OutcomesCount: 1}
	var (
		cmdSeq3   core.Seq = 3
		wantSend3          = helpers.WantSend{
			Seq: cmdSeq3,
			N:   codecs.CmdSize(cmdSeq3, cmd3),
		}
		wantReceive3 = helpers.WantReceive{
			AsyncResult: codecs.AsyncResult(cmdSeq3, results.Result(cmd3.OutcomesCount)),
		}
		start3    = time.Now()
		execTime3 = time.Duration(0)
	)
	err = helpers.Exchange(cmd3, sendFn, ReceiveFn, wantSend3, wantReceive3)
	assertfatal.EqualError(err, nil, t)

	time.Sleep(execTime3 + 50*time.Millisecond)

	// Check third Command execution time.
	var (
		end3  = receiver.DoneAt()
		want3 = start3.Add(execTime3)
	)
	assertfatal.SameTime(end3, want3, delta, t)

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

func makeClient(addr string) (client *ccln.Client[*rcvr.Receiver], err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	return cmdstream.MakeClient(codecs.ClientCodec{}, conn)
}
