package manager

import (
	"bytes"
	"time"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/transport-go"
)

type CmdRecord[T any] struct {
	Seq core.Seq
	At  time.Time
	N   int
	Cmd core.Cmd[T]
}

func unmarshalCmdRecord[T any](codec transport.Codec[core.Result, core.Cmd[T]],
	bs []byte) (cmdRecord CmdRecord[T], err error,
) {
	envelope, _, err := CmdEnvelopeDTS.Unmarshal(bs)
	if err != nil {
		return
	}
	seq, cmd, n, err := codec.Decode(bytes.NewBuffer(envelope.Bs))
	if err != nil {
		return
	}
	cmdRecord = CmdRecord[T]{
		Seq: seq,
		At:  envelope.At,
		N:   n,
		Cmd: cmd,
	}
	return
}
