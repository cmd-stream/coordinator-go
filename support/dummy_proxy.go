package support

import (
	"net"
	"time"

	"github.com/cmd-stream/core-go"
)

type DummyProxy struct{}

func (p DummyProxy) LocalAddr() (addr net.Addr) {
	return
}

func (p DummyProxy) RemoteAddr() (addr net.Addr) {
	return
}

func (p DummyProxy) Send(seq core.Seq, result core.Result) (n int, err error) {
	return
}

func (p DummyProxy) SendWithDeadline(seq core.Seq, result core.Result,
	deadline time.Time,
) (n int, err error) {
	return
}
