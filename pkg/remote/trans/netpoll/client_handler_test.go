package netpoll

import (
	"context"
	"testing"
	"time"

	"github.com/cloudwego/kitex/internal/mocks"
	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/trans"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/netpoll"
)

var cliTransHdlr remote.ClientTransHandler

func init() {
	var opt = &remote.ClientOption{
		Codec: &MockCodec{
			EncodeFunc: nil,
			DecodeFunc: nil,
		},
	}
	cliTransHdlr, _ = newCliTransHandler(opt)
}

func TestWrite(t *testing.T) {
	var isWriteBufFlushed bool
	conn := &MockNetpollConn{
		WriterFunc: func() (r netpoll.Writer) {
			writer := &MockNetpollWriter{
				FlushFunc: func() (err error) {
					isWriteBufFlushed = true
					return nil
				},
			}
			return writer
		},
	}
	ri := rpcinfo.NewRPCInfo(nil, nil, rpcinfo.NewInvocation("", ""), nil, rpcinfo.NewRPCStats())
	ctx := context.Background()
	var req interface{}
	var msg = remote.NewMessage(req, mocks.ServiceInfo(), ri, remote.Call, remote.Client)
	err := cliTransHdlr.Write(ctx, conn, msg)
	test.Assert(t, err == nil, err)
	test.Assert(t, isWriteBufFlushed)
}

func TestRead(t *testing.T) {
	rwTimeout := time.Second

	var readTimeout time.Duration
	var isReaderBufReleased bool
	conn := &MockNetpollConn{
		SetReadTimeoutFunc: func(timeout time.Duration) (e error) {
			readTimeout = timeout
			return nil
		},
		ReaderFunc: func() (r netpoll.Reader) {
			reader := &MockNetpollReader{
				ReleaseFunc: func() (err error) {
					isReaderBufReleased = true
					return nil
				},
			}
			return reader
		},
	}
	cfg := rpcinfo.NewRPCConfig()
	rpcinfo.AsMutableRPCConfig(cfg).SetReadWriteTimeout(rwTimeout)
	ri := rpcinfo.NewRPCInfo(nil, nil, nil, cfg, rpcinfo.NewRPCStats())
	ctx := context.Background()
	var req interface{}
	var msg = remote.NewMessage(req, mocks.ServiceInfo(), ri, remote.Reply, remote.Client)
	err := cliTransHdlr.Read(ctx, conn, msg)
	test.Assert(t, err == nil, err)
	test.Assert(t, readTimeout == trans.GetReadTimeout(ri.Config()))
	test.Assert(t, isReaderBufReleased)
}
