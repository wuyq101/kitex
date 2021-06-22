package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/cloudwego/kitex/internal/mocks"
	internal_server "github.com/cloudwego/kitex/internal/server"
	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/internal/utils"
	"github.com/cloudwego/kitex/pkg/diagnosis"
	"github.com/cloudwego/kitex/pkg/endpoint"
)

func TestOptionDebugInfo(t *testing.T) {
	var opts []Option
	md := newMockDiagnosis()
	opts = append(opts, internal_server.WithDiagnosisService(md))
	buildMw := 0
	opts = append(opts, WithMiddleware(func(next endpoint.Endpoint) endpoint.Endpoint {
		buildMw++
		return func(ctx context.Context, req, resp interface{}) (err error) {
			return next(ctx, req, resp)
		}
	}))
	opts = append(opts, WithMiddlewareBuilder(func(ctx context.Context) endpoint.Middleware {
		return func(next endpoint.Endpoint) endpoint.Endpoint {
			buildMw++
			return next
		}
	}))

	svr := NewServer(opts...)
	// check middleware build
	test.Assert(t, buildMw == 2)

	// check probe result
	pp := md.ProbePairs()
	test.Assert(t, pp != nil)

	optFunc := pp[diagnosis.OptionsKey]
	test.Assert(t, optFunc != nil)

	is := svr.(*server)
	debugInfo := is.opt.DebugInfo
	optDebugInfo := optFunc().(utils.Slice)
	test.Assert(t, len(debugInfo) == len(optDebugInfo))
	test.Assert(t, len(debugInfo) == 3)

	for i := range debugInfo {
		test.Assert(t, debugInfo[i] == optDebugInfo[i])
	}
}

func TestProxyOption(t *testing.T) {
	var opts []Option
	addr, err := net.ResolveTCPAddr("tcp", ":8888")
	test.Assert(t, err == nil, err)
	opts = append(opts, WithServiceAddr(addr))
	opts = append(opts, internal_server.WithProxy(&proxy{}))
	svr := NewServer(opts...)
	time.AfterFunc(100*time.Millisecond, func() {
		err := svr.Stop()
		test.Assert(t, err == nil, err)
	})

	err = svr.RegisterService(mocks.ServiceInfo(), mocks.MyServiceHandler())
	test.Assert(t, err == nil)

	err = svr.Run()
	test.Assert(t, err == nil, err)
	iSvr := svr.(*server)
	test.Assert(t, iSvr.opt.RemoteOpt.Address.String() == "mock", iSvr.opt.RemoteOpt.Address.Network())
}

type proxy struct{}

func (p *proxy) Replace(addr net.Addr) (nAddr net.Addr, err error) {
	if nAddr, err = net.ResolveUnixAddr("unix", "mock"); err == nil {
		return nAddr, err
	}
	return
}

func newMockDiagnosis() *mockDiagnosis {
	return &mockDiagnosis{probes: make(map[diagnosis.ProbeName]diagnosis.ProbeFunc)}
}

type mockDiagnosis struct {
	probes map[diagnosis.ProbeName]diagnosis.ProbeFunc
}

func (m *mockDiagnosis) RegisterProbeFunc(name diagnosis.ProbeName, probeFunc diagnosis.ProbeFunc) {
	m.probes[name] = probeFunc
}

func (m *mockDiagnosis) ProbePairs() map[diagnosis.ProbeName]diagnosis.ProbeFunc {
	return m.probes
}
