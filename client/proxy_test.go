package client

import (
	"context"
	"testing"

	"github.com/cloudwego/kitex/internal/client"
	"github.com/cloudwego/kitex/internal/mocks"
	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/proxy"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/rpcinfo/remoteinfo"
)

type mockForwardProxy struct {
	ConfigureFunc            func(*proxy.Config) error
	ResolveProxyInstanceFunc func(ctx context.Context) error
}

func (m *mockForwardProxy) Configure(cfg *proxy.Config) error {
	if m.ConfigureFunc != nil {
		return m.ConfigureFunc(cfg)
	}
	return nil
}
func (m *mockForwardProxy) ResolveProxyInstance(ctx context.Context) error {
	if m.ResolveProxyInstanceFunc != nil {
		return m.ResolveProxyInstanceFunc(ctx)
	}
	return nil
}

func TestForwardProxy(t *testing.T) {
	fp := &mockForwardProxy{
		ConfigureFunc: func(cfg *proxy.Config) error {
			test.Assert(t, cfg.Resolver == nil, cfg.Resolver)
			test.Assert(t, cfg.Balancer == nil, cfg.Balancer)
			return nil
		},
		ResolveProxyInstanceFunc: func(ctx context.Context) error {
			ri := rpcinfo.GetRPCInfo(ctx)
			re := remoteinfo.AsRemoteInfo(ri.To())
			in := re.GetInstance()
			test.Assert(t, in == nil, in)
			re.SetInstance(instance505)
			return nil
		},
	}

	var opts []Option
	opts = append(opts, client.WithTransHandlerFactory(mocks.NewMockCliTransHandlerFactory(hdlr)))
	opts = append(opts, client.WithDialer(dialer))
	opts = append(opts, WithDestService("psm"))
	opts = append(opts, client.WithProxy(fp))

	svcInfo := mocks.ServiceInfo()
	cli, err := NewClient(svcInfo, opts...)
	test.Assert(t, err == nil)

	mtd := mocks.MockMethod
	ctx := context.Background()
	req := new(MockTStruct)
	res := new(MockTStruct)

	err = cli.Call(ctx, mtd, req, res)
	test.Assert(t, err == nil, err)
}

func TestProxyWithResolver(t *testing.T) {
	fp := &mockForwardProxy{
		ConfigureFunc: func(cfg *proxy.Config) error {
			test.Assert(t, cfg.Resolver == resolver404)
			test.Assert(t, cfg.Balancer == nil, cfg.Balancer)
			return nil
		},
		ResolveProxyInstanceFunc: func(ctx context.Context) error {
			ri := rpcinfo.GetRPCInfo(ctx)
			re := remoteinfo.AsRemoteInfo(ri.To())
			in := re.GetInstance()
			test.Assert(t, in == instance404[0])
			re.SetInstance(instance505)
			return nil
		},
	}

	var opts []Option
	opts = append(opts, client.WithTransHandlerFactory(mocks.NewMockCliTransHandlerFactory(hdlr)))
	opts = append(opts, client.WithDialer(dialer))
	opts = append(opts, WithDestService("psm"))
	opts = append(opts, client.WithProxy(fp))
	opts = append(opts, WithResolver(resolver404))

	svcInfo := mocks.ServiceInfo()
	cli, err := NewClient(svcInfo, opts...)
	test.Assert(t, err == nil)

	mtd := mocks.MockMethod
	ctx := context.Background()
	req := new(MockTStruct)
	res := new(MockTStruct)

	err = cli.Call(ctx, mtd, req, res)
	test.Assert(t, err == nil, err)
}

func TestProxyWithBalancer(t *testing.T) {
	lb := &loadbalance.SynthesizedLoadbalancer{
		GetPickerFunc: func(entry discovery.Result) loadbalance.Picker {
			return &loadbalance.SynthesizedPicker{
				NextFunc: func(ctx context.Context, request interface{}) discovery.Instance {
					if len(entry.Instances) > 0 {
						return entry.Instances[0]
					}
					return nil
				},
			}
		},
	}
	fp := &mockForwardProxy{
		ConfigureFunc: func(cfg *proxy.Config) error {
			test.Assert(t, cfg.Resolver == nil, cfg.Resolver)
			test.Assert(t, cfg.Balancer == lb)
			return nil
		},
		ResolveProxyInstanceFunc: func(ctx context.Context) error {
			ri := rpcinfo.GetRPCInfo(ctx)
			re := remoteinfo.AsRemoteInfo(ri.To())
			in := re.GetInstance()
			test.Assert(t, in == nil)
			re.SetInstance(instance505)
			return nil
		},
	}

	var opts []Option
	opts = append(opts, client.WithTransHandlerFactory(mocks.NewMockCliTransHandlerFactory(hdlr)))
	opts = append(opts, client.WithDialer(dialer))
	opts = append(opts, WithDestService("psm"))
	opts = append(opts, client.WithProxy(fp))
	opts = append(opts, WithLoadBalancer(lb))

	svcInfo := mocks.ServiceInfo()
	cli, err := NewClient(svcInfo, opts...)
	test.Assert(t, err == nil)

	mtd := mocks.MockMethod
	ctx := context.Background()
	req := new(MockTStruct)
	res := new(MockTStruct)

	err = cli.Call(ctx, mtd, req, res)
	test.Assert(t, err == nil)
}

func TestProxyWithResolverAndBalancer(t *testing.T) {
	lb := &loadbalance.SynthesizedLoadbalancer{
		GetPickerFunc: func(entry discovery.Result) loadbalance.Picker {
			return &loadbalance.SynthesizedPicker{
				NextFunc: func(ctx context.Context, request interface{}) discovery.Instance {
					return instance505
				},
			}
		},
	}
	fp := &mockForwardProxy{
		ConfigureFunc: func(cfg *proxy.Config) error {
			test.Assert(t, cfg.Resolver == resolver404)
			test.Assert(t, cfg.Balancer == lb)
			return nil
		},
		ResolveProxyInstanceFunc: func(ctx context.Context) error {
			ri := rpcinfo.GetRPCInfo(ctx)
			re := remoteinfo.AsRemoteInfo(ri.To())
			in := re.GetInstance()
			test.Assert(t, in == instance505)
			return nil
		},
	}

	var opts []Option
	opts = append(opts, client.WithTransHandlerFactory(mocks.NewMockCliTransHandlerFactory(hdlr)))
	opts = append(opts, client.WithDialer(dialer))
	opts = append(opts, WithDestService("psm"))
	opts = append(opts, client.WithProxy(fp))
	opts = append(opts, WithResolver(resolver404))
	opts = append(opts, WithLoadBalancer(lb))

	svcInfo := mocks.ServiceInfo()
	cli, err := NewClient(svcInfo, opts...)
	test.Assert(t, err == nil)

	mtd := mocks.MockMethod
	ctx := context.Background()
	req := new(MockTStruct)
	res := new(MockTStruct)

	err = cli.Call(ctx, mtd, req, res)
	test.Assert(t, err == nil)
}
