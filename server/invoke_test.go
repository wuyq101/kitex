package server

import (
	"testing"

	"github.com/apache/thrift/lib/go/thrift"

	"github.com/cloudwego/kitex/internal/mocks"
	internal_server "github.com/cloudwego/kitex/internal/server"
	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/pkg/remote/trans/invoke"
	"github.com/cloudwego/kitex/pkg/utils"
)

func TestInvokerCall(t *testing.T) {
	var opts []Option
	opts = append(opts, internal_server.WithMetaHandler(noopMetahandler{}))
	invoker := NewInvoker(opts...)

	err := invoker.RegisterService(mocks.ServiceInfo(), mocks.MyServiceHandler())
	if err != nil {
		t.Fatal(err)
	}
	err = invoker.Init()
	if err != nil {
		t.Fatal(err)
	}
	args := mocks.NewMockArgs()
	codec := utils.NewThriftMessageCodec()
	b, _ := codec.Encode("mock", thrift.CALL, 0, args.(thrift.TStruct))

	msg := invoke.NewMessage(nil, nil)
	msg.SetRequestBytes(b)

	err = invoker.Call(msg)
	if err != nil {
		t.Fatal(err)
	}
	b, err = msg.GetResponseBytes()
	if err != nil {
		t.Fatal(err)
	}
	test.Assert(t, len(b) > 0)
}
