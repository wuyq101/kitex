package trans

import (
	"context"
	"net"
	"testing"

	"github.com/cloudwego/kitex/internal/mocks"
	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

func TestDefaultSvrTransHandler(t *testing.T) {
	buf := remote.NewReaderWriterBuffer(1024)
	ext := &MockExtension{
		NewWriteByteBufferFunc: func(ctx context.Context, conn net.Conn, msg remote.Message) remote.ByteBuffer {
			return buf
		},
		NewReadByteBufferFunc: func(ctx context.Context, conn net.Conn, msg remote.Message) remote.ByteBuffer {
			return buf
		},
	}

	tagEncode, tagDecode := 0, 0
	opt := &remote.ServerOption{
		Codec: &MockCodec{
			EncodeFunc: func(ctx context.Context, msg remote.Message, out remote.ByteBuffer) error {
				tagEncode++
				test.Assert(t, out == buf)
				return nil
			},
			DecodeFunc: func(ctx context.Context, msg remote.Message, in remote.ByteBuffer) error {
				tagDecode++
				test.Assert(t, in == buf)
				return nil
			},
		},
		SvcInfo: mocks.ServiceInfo(),
	}

	handler, err := NewDefaultSvrTransHandler(opt, ext)
	test.Assert(t, err == nil)

	ctx := context.Background()
	conn := &mocks.Conn{}
	msg := &MockMessage{
		RPCInfoFunc: func() rpcinfo.RPCInfo {
			return newMockRPCInfo()
		},
	}
	err = handler.Write(ctx, conn, msg)
	test.Assert(t, err == nil, err)
	test.Assert(t, tagEncode == 1, tagEncode)
	test.Assert(t, tagDecode == 0, tagDecode)

	err = handler.Read(ctx, conn, msg)
	test.Assert(t, err == nil, err)
	test.Assert(t, tagEncode == 1, tagEncode)
	test.Assert(t, tagDecode == 1, tagDecode)
}
