//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	coord "github.com/cmd-stream/coordinator-go"
	coordman "github.com/cmd-stream/coordinator-go/manager"
	"github.com/cmd-stream/coordinator-go/support"
	"github.com/cmd-stream/testkit-go/fixtures/coordinator/cmds"
	"github.com/cmd-stream/testkit-go/fixtures/coordinator/codecs"
	rcvr "github.com/cmd-stream/testkit-go/fixtures/coordinator/receiver"

	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestRetry(t *testing.T) {
	const (
		addr  = "127.0.0.1:9003"
		delta = 100 * time.Millisecond
	)

	storage := support.NewMemoryStorage(support.NewBatcher(coord.CheckpointSize))
	err := fillStorage(storage)
	assertfatal.EqualError(err, nil, t)

	var (
		receiver = &rcvr.Receiver{}
		start    = time.Now()
	)
	err = StartServer(addr, receiver, storage)
	assertfatal.EqualError(err, nil, t)
	time.Sleep(time.Second)

	// Before each Command retry coordinator waits for a while.
	want := start.Add(coordman.SlowRetryInterval)
	assertfatal.SameTime(receiver.DoneAt(), want, delta, t)
}

func fillStorage(storage coordman.Storage) (err error) {
	cmd := cmds.Cmd{OutcomesCount: 3}

	w := NewWriter()
	w.WriteByte(1) // seq
	codecs.ClientCodec{}.Encode(cmd, w)
	buf := w.Buffer

	cmdEnvelope := coordman.CmdEnvelope{At: time.Now(), Bs: buf.Bytes()}
	bs := make([]byte, coordman.CmdEnvelopeDTS.Size(cmdEnvelope))
	coordman.CmdEnvelopeDTS.Marshal(cmdEnvelope, bs)

	partID, partSeq, err := storage.Save(bs)
	if err != nil {
		return
	}

	outcome := cmds.Outcome(0)
	bs, err = json.Marshal(outcome)
	if err != nil {
		return
	}

	outcomeEnvelope := coordman.OutcomeEnvelope{
		WorkflowID: coordman.WorkflowID{PartID: partID, PartSeq: partSeq},
		Bs:         bs,
	}
	bs = make([]byte, coordman.OutcomeEnvelopeDTS.Size(outcomeEnvelope))
	coordman.OutcomeEnvelopeDTS.Marshal(outcomeEnvelope, bs)

	partID, partSeq, err = storage.Save(bs)
	if err != nil {
		return
	}
	storage.SetCompleted(partID, partSeq)
	return
}

func NewWriter() Writer {
	return Writer{Buffer: &bytes.Buffer{}}
}

type Writer struct {
	*bytes.Buffer
}

func (w Writer) Flush() error {
	return nil
}
