package coordinator

import (
	"net"

	"github.com/cmd-stream/coordinator-go/support"
	"github.com/cmd-stream/core-go"
	trn "github.com/cmd-stream/transport-go"
	strn "github.com/cmd-stream/transport-go/server"
)

func NewTransport[T any](conn net.Conn,
	codec trn.Codec[core.Result, core.Cmd[T]],
	ops ...trn.SetOption,
) Transport[T] {
	return Transport[T]{support.NewRecConn(conn), strn.New(conn,
		codec, ops...), codec}
}

type Transport[T any] struct {
	conn *support.RecConn
	*strn.Transport[T]
	codec trn.Codec[core.Result, core.Cmd[T]]
}

func (t Transport[T]) Codec() trn.Codec[core.Result, core.Cmd[T]] {
	return t.codec
}

func (t Transport[T]) ReceiveWithBytes() (seq core.Seq, cmd core.Cmd[T],
	bs []byte, err error,
) {
	t.conn.ClearRecordedBytes()
	seq, cmd, _, err = t.Receive()
	bs = t.conn.Bytes()
	return
}
