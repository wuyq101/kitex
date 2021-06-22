package rpcinfo

import (
	"context"

	"github.com/cloudwego/kitex/internal"
)

type ctxRPCInfoKeyType struct{}

var ctxRPCInfoKey ctxRPCInfoKeyType

// NewCtxWithRPCInfo creates a new context with the RPCInfo given.
func NewCtxWithRPCInfo(ctx context.Context, ri RPCInfo) context.Context {
	if ri != nil {
		return context.WithValue(ctx, ctxRPCInfoKey, ri)
	}
	return ctx
}

// GetRPCInfo gets RPCInfo from ctx.
// Returns nil if not found.
func GetRPCInfo(ctx context.Context) RPCInfo {
	if ri, ok := ctx.Value(ctxRPCInfoKey).(RPCInfo); ok {
		return ri
	}
	return nil
}

// PutRPCInfo recycles the RPCInfo. This function is for internal use only.
func PutRPCInfo(ri RPCInfo) {
	if ri == nil {
		return
	}
	if r, ok := ri.From().(internal.Reusable); ok {
		r.Recycle()
	}
	if r, ok := ri.To().(internal.Reusable); ok {
		r.Recycle()
	}
	if r, ok := ri.Invocation().(internal.Reusable); ok {
		r.Recycle()
	}
	if r, ok := ri.Config().(internal.Reusable); ok {
		r.Recycle()
	}
	if r, ok := ri.Stats().(internal.Reusable); ok {
		r.Recycle()
	}
	if r, ok := ri.(internal.Reusable); ok {
		r.Recycle()
	}
}
