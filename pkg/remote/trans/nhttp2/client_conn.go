package nhttp2

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/remote/trans/nhttp2/grpc"
	"github.com/cloudwego/kitex/pkg/remote/trans/nhttp2/metadata"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

type clientConn struct {
	tr grpc.ClientTransport
	s  *grpc.Stream
}

var _ net.Conn = (*clientConn)(nil)

func newClientConn(ctx context.Context, tr grpc.ClientTransport, addr string) (*clientConn, error) {
	ri := rpcinfo.GetRPCInfo(ctx)
	s, err := tr.NewStream(ctx, &grpc.CallHdr{
		Host: addr,
		// grpc method format /package.Service/Method
		Method: fmt.Sprintf("/%s.%s/%s", ri.Invocation().PackageName(), ri.Invocation().ServiceName(), ri.Invocation().MethodName()),
	})
	if err != nil {
		return nil, err
	}
	return &clientConn{
		tr: tr,
		s:  s,
	}, nil
}

// impl net.Conn
func (c *clientConn) Read(b []byte) (n int, err error) {
	n, err = c.s.Read(b)
	if err == io.EOF {
		if statusErr := c.s.Status().Err(); statusErr != nil {
			err = statusErr
		}
	}
	return n, convertErrorFromGrpcToKitex(err)
}

func (c *clientConn) Write(b []byte) (n int, err error) {
	if len(b) < 5 {
		return 0, io.ErrShortWrite
	}
	err = c.tr.Write(c.s, b[:5], b[5:], &grpc.Options{})
	return len(b), convertErrorFromGrpcToKitex(err)
}

func (c *clientConn) LocalAddr() net.Addr                { return c.tr.LocalAddr() }
func (c *clientConn) RemoteAddr() net.Addr               { return c.tr.RemoteAddr() }
func (c *clientConn) SetDeadline(t time.Time) error      { return nil }
func (c *clientConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *clientConn) SetWriteDeadline(t time.Time) error { return nil }
func (c *clientConn) Close() error {
	c.tr.Write(c.s, nil, nil, &grpc.Options{Last: true})
	// Always return nil; io.EOF is the only error that might make sense
	// instead, but there is no need to signal the client to call Read
	// as the only use left for the stream after Close is to call
	// Read. This also matches historical behavior.
	return nil
}

func (c *clientConn) Header() (metadata.MD, error) { return c.s.Header() }
func (c *clientConn) Trailer() metadata.MD         { return c.s.Trailer() }
