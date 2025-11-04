package support

import "net"

func NewRecConn(conn net.Conn) *RecConn {
	return &RecConn{conn, []byte{}}
}

type RecConn struct {
	net.Conn
	bs []byte
}

func (c *RecConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	c.bs = append(c.bs, b...)
	return
}

func (c *RecConn) Bytes() (bs []byte) {
	return c.bs
}

func (c *RecConn) ClearRecordedBytes() {
	c.bs = c.bs[:0]
}
