package remotesvr

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/cloudwego/kitex/internal/mocks"
	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/internal/utils"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/remote"
)

func TestServerStart(t *testing.T) {
	isCreateListener := false
	isBootstrapped := false
	transSvr := &mocks.MockTransServer{
		CreateListenerFunc: func(addr net.Addr) (listener net.Listener, err error) {
			isCreateListener = true
			return nil, nil
		},
		BootstrapServerFunc: func() (err error) {
			isBootstrapped = true
			return nil
		},
	}
	var opt = &remote.ServerOption{
		Address:            utils.NewNetAddr("tcp", "test"),
		TransServerFactory: mocks.NewMockTransServerFactory(transSvr),
		Logger:             klog.DefaultLogger(),
	}
	inkHdlrFunc := func(ctx context.Context, req, resp interface{}) (err error) {
		return nil
	}
	transHdrl := &mocks.MockSvrTransHandler{}
	svr, err := NewServer(opt, inkHdlrFunc, transHdrl)
	test.Assert(t, err == nil, err)

	err = <-svr.Start()
	test.Assert(t, err == nil, err)
	test.Assert(t, isBootstrapped)
	test.Assert(t, isCreateListener)

	err = svr.Stop()
	test.Assert(t, err == nil, err)
}

func TestServerStartListenErr(t *testing.T) {
	mockErr := errors.New("test")
	transSvr := &mocks.MockTransServer{
		CreateListenerFunc: func(addr net.Addr) (listener net.Listener, err error) {
			return nil, mockErr
		},
	}
	var opt = &remote.ServerOption{
		Address:            utils.NewNetAddr("tcp", "test"),
		TransServerFactory: mocks.NewMockTransServerFactory(transSvr),
		Logger:             klog.DefaultLogger(),
	}
	inkHdlrFunc := func(ctx context.Context, req, resp interface{}) (err error) {
		return nil
	}
	transHdrl := &mocks.MockSvrTransHandler{}
	svr, err := NewServer(opt, inkHdlrFunc, transHdrl)
	test.Assert(t, err == nil, err)

	err = <-svr.Start()
	test.Assert(t, err != nil, err)
	test.Assert(t, err == mockErr)

	err = svr.Stop()
	test.Assert(t, err == nil, err)
}
