// Package netpoll contains server and client implementation for netpoll.
package netpoll

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/cloudwego/kitex/internal/utils"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/netpoll"
)

// NewTransServerFactory creates a default netpoll transport server factory.
func NewTransServerFactory() remote.TransServerFactory {
	return &netpollTransServerFactory{}
}

type netpollTransServerFactory struct {
}

// NewTransServer implements the remote.TransServerFactory interface.
func (f *netpollTransServerFactory) NewTransServer(opt *remote.ServerOption, transHdlr remote.ServerTransHandler) remote.TransServer {
	return &transServer{opt: opt, transHdlr: transHdlr}
}

type transServer struct {
	opt       *remote.ServerOption
	transHdlr remote.ServerTransHandler

	evl       netpoll.EventLoop
	ln        netpoll.Listener
	connCount utils.AtomicInt
	sync.Mutex
}

var _ remote.TransServer = &transServer{}

// CreateListener implements the remote.TransServer interface.
func (ts *transServer) CreateListener(addr net.Addr) (net.Listener, error) {
	ln, err := netpoll.CreateListener(addr.Network(), addr.String())
	ts.ln = ln
	return ts.ln, err
}

// BootstrapServer implements the remote.TransServer interface.
func (ts *transServer) BootstrapServer() (err error) {
	if ts.ln == nil {
		return errors.New("listener is nil in netpoll transport server")
	}
	opts := []netpoll.Option{
		netpoll.WithIdleTimeout(ts.opt.MaxConnectionIdleTime),
		netpoll.WithReadTimeout(ts.opt.ReadWriteTimeout),
	}

	ts.Lock()
	opts = append(opts, netpoll.WithOnPrepare(ts.onConnActive))
	ts.evl, err = netpoll.NewEventLoop(ts.onConnRead, opts...)
	ts.Unlock()

	if err != nil {
		return err
	}
	return ts.evl.Serve(ts.ln)
}

// Shutdown implements the remote.TransServer interface.
func (ts *transServer) Shutdown() (err error) {
	ts.Lock()
	defer ts.Unlock()
	if ts.evl != nil {
		ctx, cancel := context.WithTimeout(context.Background(), ts.opt.ExitWaitTime)
		defer cancel()
		err = ts.evl.Shutdown(ctx)
	}
	if ts.ln != nil {
		err = ts.ln.Close()
	}
	return
}

// ConnCount implements the remote.TransServer interface.
func (ts *transServer) ConnCount() utils.AtomicInt {
	return ts.connCount
}

// TODO: Is it necessary to init RPCInfo when connection is first created?
// Now the strategy is the same as the original one.
// 1. There's no request when the connection is first created.
// 2. Doesn't need to init RPCInfo if it's not RPC request, such as heartbeat.
func (ts *transServer) onConnActive(conn netpoll.Connection) context.Context {
	ctx := context.Background()
	conn.AddCloseCallback(func(connection netpoll.Connection) error {
		ts.onConnInactive(ctx, conn)
		return nil
	})
	ts.connCount.Inc()
	ctx, err := ts.transHdlr.OnActive(ctx, conn)
	if err != nil {
		ts.onError(ctx, err, conn)
		conn.Close()
		return ctx
	}
	return ctx
}

func (ts *transServer) onConnRead(ctx context.Context, conn netpoll.Connection) error {
	err := ts.transHdlr.OnRead(ctx, conn)
	if err != nil {
		ts.onError(ctx, err, conn)
		conn.Close()
	}
	return nil
}

func (ts *transServer) onConnInactive(ctx context.Context, conn netpoll.Connection) {
	ts.connCount.Dec()
	ts.transHdlr.OnInactive(ctx, conn)
}

func (ts *transServer) onError(ctx context.Context, err error, conn netpoll.Connection) {
	ts.transHdlr.OnError(ctx, err, conn)
}
