package connpool

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/cloudwego/kitex/internal/mocks"
	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/pkg/remote"
)

func TestShortPool(t *testing.T) {
	d := &remote.SynthesizedDialer{}
	p := NewShortPool("test")

	d.DialFunc = func(network, address string, timeout time.Duration) (net.Conn, error) {
		cost := time.Millisecond * 10
		if timeout < cost {
			return nil, errors.New("connection timeout")
		}
		return nil, nil
	}

	_, err := p.Get(nil, "tcp", "127.0.0.1:8080", &remote.ConnOption{Dialer: d, ConnectTimeout: time.Second})
	test.Assert(t, err == nil)

	_, err = p.Get(nil, "tcp", "127.0.0.1:8888", &remote.ConnOption{Dialer: d, ConnectTimeout: time.Millisecond})
	test.Assert(t, err != nil)
}

func TestShortPoolPut(t *testing.T) {
	d := &remote.SynthesizedDialer{}
	p := NewShortPool("test")

	c := &mocks.Conn{}
	d.DialFunc = func(network, address string, timeout time.Duration) (net.Conn, error) {
		return c, nil
	}

	x, err := p.Get(nil, "tcp", "127.0.0.1:8080", &remote.ConnOption{Dialer: d, ConnectTimeout: time.Second})
	test.Assert(t, err == nil)
	y, ok := x.(*shortConn)
	test.Assert(t, ok)
	test.Assert(t, y.Conn == c)
	test.Assert(t, !y.closed)

	err = p.Put(x)
	test.Assert(t, err == nil)
	test.Assert(t, y.closed)
}

func TestShortPoolPutClosed(t *testing.T) {
	d := &remote.SynthesizedDialer{}
	p := NewShortPool("test")

	var closed bool
	var cmto = errors.New("closed more than once")
	c := &mocks.Conn{
		CloseFunc: func() (err error) {
			if closed {
				err = cmto
			}
			closed = true
			return
		},
	}
	d.DialFunc = func(network, address string, timeout time.Duration) (net.Conn, error) {
		return c, nil
	}

	x, err := p.Get(nil, "tcp", "127.0.0.1:8080", &remote.ConnOption{Dialer: d, ConnectTimeout: time.Second})
	test.Assert(t, err == nil)
	y, ok := x.(*shortConn)
	test.Assert(t, ok)
	test.Assert(t, y.Conn == c)
	test.Assert(t, !y.closed)
	test.Assert(t, !closed)

	err = y.Close()
	test.Assert(t, err == nil)
	test.Assert(t, y.closed)
	test.Assert(t, closed)

	err = y.Close()
	test.Assert(t, err == nil)

	err = p.Put(x)
	test.Assert(t, err == nil)
}
