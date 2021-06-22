package client

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudwego/kitex/internal/client"
	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/event"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/rpcinfo/remoteinfo"
)

var (
	instance404 = []discovery.Instance{
		discovery.NewInstance("tcp", "localhost:404", 1, make(map[string]string)),
	}
	instance505 = discovery.NewInstance("tcp", "localhost:505", 1, make(map[string]string))
	resolver404 = &discovery.SynthesizedResolver{
		ResolveFunc: func(ctx context.Context, key string) (discovery.Result, error) {
			return discovery.Result{
				Cacheable: true,
				CacheKey:  "test",
				Instances: instance404,
			}, nil
		},
		NameFunc: func() string { return "middlewares_test" },
	}

	ctx = func() context.Context {
		ctx := context.Background()
		ctx = context.WithValue(ctx, endpoint.CtxEventBusKey, event.NewEventBus())
		ctx = context.WithValue(ctx, endpoint.CtxEventQueueKey, event.NewQueue(10))
		return ctx
	}()
)

func TestResolverMWNoResolver(t *testing.T) {
	var invoked bool
	opts := &client.Options{}
	mw := newResolveMWBuilder(opts)(ctx)
	ep := func(ctx context.Context, request, response interface{}) error {
		invoked = true
		return nil
	}

	ri := rpcinfo.NewRPCInfo(nil,
		remoteinfo.NewRemoteInfo(&rpcinfo.EndpointBasicInfo{}, ""),
		rpcinfo.NewInvocation("", ""),
		nil, rpcinfo.NewRPCStats())

	ctx := rpcinfo.NewCtxWithRPCInfo(context.Background(), ri)
	req := new(MockTStruct)
	res := new(MockTStruct)
	err := mw(ep)(ctx, req, res)
	test.Assert(t, err != nil)
	test.Assert(t, errors.Is(err, kerrors.ErrNoResolver))
	test.Assert(t, !invoked)
}

func TestResolverMW(t *testing.T) {
	var invoked bool
	opts := &client.Options{Resolver: resolver404}
	mw := newResolveMWBuilder(opts)(ctx)
	ep := func(ctx context.Context, request, response interface{}) error {
		invoked = true
		return nil
	}

	to := remoteinfo.NewRemoteInfo(&rpcinfo.EndpointBasicInfo{}, "")
	ri := rpcinfo.NewRPCInfo(nil, to, rpcinfo.NewInvocation("", ""), nil, rpcinfo.NewRPCStats())

	ctx := rpcinfo.NewCtxWithRPCInfo(context.Background(), ri)
	req := new(MockTStruct)
	res := new(MockTStruct)
	err := mw(ep)(ctx, req, res)
	test.Assert(t, err == nil)
	test.Assert(t, invoked)
	test.Assert(t, to.GetInstance() == instance404[0])
}

func TestResolverMWOutOfInstance(t *testing.T) {
	resolver := &discovery.SynthesizedResolver{
		ResolveFunc: func(ctx context.Context, key string) (discovery.Result, error) {
			return discovery.Result{}, nil
		},
		NameFunc: func() string { return t.Name() },
	}
	var invoked bool
	opts := &client.Options{Resolver: resolver}
	mw := newResolveMWBuilder(opts)(ctx)
	ep := func(ctx context.Context, request, response interface{}) error {
		invoked = true
		return nil
	}

	to := remoteinfo.NewRemoteInfo(&rpcinfo.EndpointBasicInfo{}, "")
	ri := rpcinfo.NewRPCInfo(nil, to, rpcinfo.NewInvocation("", ""), nil, rpcinfo.NewRPCStats())

	ctx := rpcinfo.NewCtxWithRPCInfo(context.Background(), ri)
	req := new(MockTStruct)
	res := new(MockTStruct)

	err := mw(ep)(ctx, req, res)
	test.Assert(t, err != nil, err)
	test.Assert(t, errors.Is(err, kerrors.ErrNoMoreInstance))
	test.Assert(t, to.GetInstance() == nil)
	test.Assert(t, !invoked)
}

func BenchmarkResolverMW(b *testing.B) {
	opts := &client.Options{Resolver: resolver404}
	mw := newResolveMWBuilder(opts)(ctx)
	ep := func(ctx context.Context, request, response interface{}) error { return nil }
	ri := rpcinfo.NewRPCInfo(nil, nil, rpcinfo.NewInvocation("", ""), nil, rpcinfo.NewRPCStats())

	ctx := rpcinfo.NewCtxWithRPCInfo(context.Background(), ri)
	req := new(MockTStruct)
	res := new(MockTStruct)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mw(ep)(ctx, req, res)
	}
}

func BenchmarkResolverMWParallel(b *testing.B) {
	opts := &client.Options{Resolver: resolver404}

	mw := newResolveMWBuilder(opts)(ctx)
	ep := func(ctx context.Context, request, response interface{}) error { return nil }
	ri := rpcinfo.NewRPCInfo(nil, nil, rpcinfo.NewInvocation("", ""), nil, rpcinfo.NewRPCStats())

	ctx := rpcinfo.NewCtxWithRPCInfo(context.Background(), ri)
	req := new(MockTStruct)
	res := new(MockTStruct)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mw(ep)(ctx, req, res)
		}
	})
}
