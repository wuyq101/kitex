package rpcinfo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

func TestChainHandlers(t *testing.T) {
	var order []int
	m1 := &rpcinfo.SynthesizedMetaHandler{
		Writer: func(ctx context.Context, request, response, protocol interface{}) (context.Context, error) {
			order = append(order, 1)
			return ctx, nil
		},
		Reader: func(ctx context.Context, request, response, protocol interface{}) (context.Context, error) {
			order = append(order, -1)
			return ctx, nil
		},
	}
	m2 := &rpcinfo.SynthesizedMetaHandler{
		Writer: func(ctx context.Context, request, response, protocol interface{}) (context.Context, error) {
			order = append(order, 2)
			return ctx, nil
		},
		Reader: func(ctx context.Context, request, response, protocol interface{}) (context.Context, error) {
			order = append(order, -2)
			return ctx, nil
		},
	}

	m3 := rpcinfo.ChainHandlers([]rpcinfo.MetaHandler{m1, m2})
	m3.WriteMeta(nil, nil, nil, nil)
	m3.ReadMeta(nil, nil, nil, nil)
	test.Assert(t, len(order) == 4)
	test.Assert(t, fmt.Sprint(order) == "[1 2 -1 -2]")
}
