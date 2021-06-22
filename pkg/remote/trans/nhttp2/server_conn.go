package nhttp2

import (
	"io"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/remote/trans/nhttp2/codes"
	"github.com/cloudwego/kitex/pkg/remote/trans/nhttp2/grpc"
	"github.com/cloudwego/kitex/pkg/remote/trans/nhttp2/status"
)

type serverConn struct {
	tr grpc.ServerTransport
	s  *grpc.Stream
}

var _ net.Conn = (*serverConn)(nil)

func newServerConn(tr grpc.ServerTransport, s *grpc.Stream) *serverConn {
	return &serverConn{
		tr: tr,
		s:  s,
	}
}

// impl net.Conn
func (c *serverConn) Read(b []byte) (n int, err error) {
	n, err = c.s.Read(b)
	return n, convertErrorFromGrpcToKitex(err)
}

func (c *serverConn) Write(b []byte) (n int, err error) {
	if len(b) < 5 {
		return 0, io.ErrShortWrite
	}
	err = c.tr.Write(c.s, b[:5], b[5:], nil)
	return len(b), convertErrorFromGrpcToKitex(err)
}

func (c *serverConn) LocalAddr() net.Addr                { return c.tr.LocalAddr() }
func (c *serverConn) RemoteAddr() net.Addr               { return c.tr.RemoteAddr() }
func (c *serverConn) SetDeadline(t time.Time) error      { return nil }
func (c *serverConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *serverConn) SetWriteDeadline(t time.Time) error { return nil }
func (c *serverConn) Close() error {
	return c.tr.WriteStatus(c.s, status.New(codes.OK, ""))
}
