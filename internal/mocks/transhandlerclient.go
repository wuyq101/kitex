package mocks

import (
	"context"
	"net"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/remote"
)

type mockCliTransHandlerFactory struct {
	hdlr *MockCliTransHandler
}

func NewMockCliTransHandlerFactory(hdrl *MockCliTransHandler) remote.ClientTransHandlerFactory {
	return &mockCliTransHandlerFactory{hdrl}
}

func (f *mockCliTransHandlerFactory) NewTransHandler(opt *remote.ClientOption) (remote.ClientTransHandler, error) {
	f.hdlr.opt = opt
	return f.hdlr, nil
}

type MockCliTransHandler struct {
	opt       *remote.ClientOption
	transPipe *remote.TransPipeline

	WriteFunc func(ctx context.Context, conn net.Conn, send remote.Message) error

	// 调用decode
	ReadFunc func(ctx context.Context, conn net.Conn, msg remote.Message) error
}

func (t *MockCliTransHandler) Write(ctx context.Context, conn net.Conn, send remote.Message) (err error) {
	if t.WriteFunc != nil {
		return t.WriteFunc(ctx, conn, send)
	}
	return
}

// 阻塞等待
func (t *MockCliTransHandler) Read(ctx context.Context, conn net.Conn, msg remote.Message) (err error) {
	if t.ReadFunc != nil {
		return t.ReadFunc(ctx, conn, msg)
	}
	return
}

func (t *MockCliTransHandler) OnMessage(ctx context.Context, args, result remote.Message) error {
	// do nothing
	return nil
}

// 新连接建立时触发，主要用于服务端，对用netpoll onPrepare
func (t *MockCliTransHandler) OnActive(ctx context.Context, conn net.Conn) (context.Context, error) {
	// ineffective now and do nothing
	return ctx, nil
}

// 连接关闭时回调
func (t *MockCliTransHandler) OnInactive(ctx context.Context, conn net.Conn) {
	// ineffective now and do nothing
}

// 传输层扩展中panic 回调
func (t *MockCliTransHandler) OnError(ctx context.Context, err error, conn net.Conn) {
	if pe, ok := err.(*kerrors.DetailedError); ok {
		t.opt.Logger.Errorf("KITEX: send request error, remote=%s, err=%s\n%s", conn.RemoteAddr(), err.Error(), pe.Stack())
	} else {
		t.opt.Logger.Errorf("KITEX: send request error, remote=%s, err=%s", conn.RemoteAddr(), err.Error())
	}
}

func (t *MockCliTransHandler) SetPipeline(p *remote.TransPipeline) {
	t.transPipe = p
}
