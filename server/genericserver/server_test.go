package genericserver

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/server"
)

func TestNewServer(t *testing.T) {
	addr, _ := net.ResolveTCPAddr("tcp", ":9999")
	g := generic.BinaryThriftGeneric()
	svr := NewServer(new(mockImpl), g, server.WithServiceAddr(addr))
	time.AfterFunc(time.Second, func() {
		err := svr.Stop()
		test.Assert(t, err == nil, err)
	})
	err := svr.Run()
	test.Assert(t, err == nil, err)
}

func TestNewServerWithServiceInfo(t *testing.T) {
	g := generic.BinaryThriftGeneric()
	svr := NewServerWithServiceInfo(new(mockImpl), g, nil)
	err := svr.Run()
	test.Assert(t, strings.Contains(err.Error(), "no service"))

	svr = NewServerWithServiceInfo(nil, g, generic.ServiceInfo(g.PayloadCodecType()))
	err = svr.Run()
	test.Assert(t, strings.Contains(err.Error(), "handler is nil"))

	test.PanicAt(t, func() {
		NewServerWithServiceInfo(nil, nil, nil)
	}, func(err interface{}) bool {
		if errMsg, ok := err.(string); ok {
			return strings.Contains(errMsg, "invalid Generic")
		}
		return false
	})
}

// GenericServiceImpl ...
type mockImpl struct {
}

// GenericCall ...
func (g *mockImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	return nil, nil
}
