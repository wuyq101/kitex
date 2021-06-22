// Package detection protocol detection
package detection

import (
	"bytes"
	"context"
	"net"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/codec"
	transNetpoll "github.com/cloudwego/kitex/pkg/remote/trans/netpoll"
	"github.com/cloudwego/kitex/pkg/remote/trans/nhttp2"
	"github.com/cloudwego/kitex/pkg/remote/trans/nhttp2/grpc"
	"github.com/cloudwego/netpoll"
)

type detectionKey struct{}

type detectionHandler struct {
	handler remote.ServerTransHandler
}

// NewSvrTransHandlerFactory detection factory construction
func NewSvrTransHandlerFactory() remote.ServerTransHandlerFactory {
	return &svrTransHandlerFactory{
		nhttp2:  nhttp2.NewSvrTransHandlerFactory(),
		netpoll: transNetpoll.NewSvrTransHandlerFactory(),
	}
}

type svrTransHandlerFactory struct {
	nhttp2  remote.ServerTransHandlerFactory
	netpoll remote.ServerTransHandlerFactory
}

func (f *svrTransHandlerFactory) NewTransHandler(opt *remote.ServerOption) (remote.ServerTransHandler, error) {
	t := &svrTransHandler{}
	var err error
	if t.nhttp2, err = f.nhttp2.NewTransHandler(opt); err != nil {
		return nil, err
	}
	if t.netpoll, err = f.netpoll.NewTransHandler(opt); err != nil {
		return nil, err
	}
	return t, nil
}

type svrTransHandler struct {
	nhttp2  remote.ServerTransHandler
	netpoll remote.ServerTransHandler
}

func (t *svrTransHandler) Write(ctx context.Context, conn net.Conn, send remote.Message) error {
	return t.which(ctx).Write(ctx, conn, send)
}

func (t *svrTransHandler) Read(ctx context.Context, conn net.Conn, msg remote.Message) error {
	return t.which(ctx).Read(ctx, conn, msg)
}

var prefaceReadAtMost = func() int {
	// min(len(ClientPreface), len(flagBuf))
	// len(flagBuf) = 2 * codec.Size32
	if 2*codec.Size32 < grpc.ClientPrefaceLen {
		return 2 * codec.Size32
	}
	return grpc.ClientPrefaceLen
}()

func (t *svrTransHandler) OnRead(ctx context.Context, conn net.Conn) error {
	// only need detect once when connection is reused
	r := ctx.Value(detectionKey{}).(*detectionHandler)
	if r.handler != nil {
		return r.handler.OnRead(ctx, conn)
	}
	// Check the validity of client preface.
	zr := conn.(netpoll.Connection).Reader()
	// read at most avoid block
	preface, err := zr.Peek(prefaceReadAtMost)
	if err != nil {
		return err
	}
	// compare preface one by one
	which := t.netpoll
	if bytes.Equal(preface[:prefaceReadAtMost], grpc.ClientPreface[:prefaceReadAtMost]) {
		which = t.nhttp2
		ctx, err = which.OnActive(ctx, conn)
		if err != nil {
			return err
		}
	}
	r.handler = which
	return which.OnRead(ctx, conn)
}

func (t *svrTransHandler) OnInactive(ctx context.Context, conn net.Conn) {
	t.which(ctx).OnInactive(ctx, conn)
}

func (t *svrTransHandler) which(ctx context.Context) remote.ServerTransHandler {
	if r := ctx.Value(detectionKey{}).(*detectionHandler); r.handler != nil {
		return r.handler
	}
	// use noop transHandler
	return noopSvrTransHandler{}
}

func (t *svrTransHandler) OnError(ctx context.Context, err error, conn net.Conn) {
	t.which(ctx).OnError(ctx, err, conn)
}

func (t *svrTransHandler) OnMessage(ctx context.Context, args, result remote.Message) error {
	return t.which(ctx).OnMessage(ctx, args, result)
}

//
func (t *svrTransHandler) SetPipeline(pipeline *remote.TransPipeline) {
	t.nhttp2.SetPipeline(pipeline)
	t.netpoll.SetPipeline(pipeline)
}

func (t *svrTransHandler) SetInvokeHandleFunc(inkHdlFunc endpoint.Endpoint) {
	if t, ok := t.nhttp2.(remote.InvokeHandleFuncSetter); ok {
		t.SetInvokeHandleFunc(inkHdlFunc)
	}
	if t, ok := t.netpoll.(remote.InvokeHandleFuncSetter); ok {
		t.SetInvokeHandleFunc(inkHdlFunc)
	}
}

func (t *svrTransHandler) OnActive(ctx context.Context, conn net.Conn) (context.Context, error) {
	ctx, err := t.netpoll.OnActive(ctx, conn)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, detectionKey{}, &detectionHandler{}), nil
}
