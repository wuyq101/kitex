// Package remotecli for remote client
package remotecli

import (
	"context"

	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/trans/nhttp2"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/streaming"
)

// NewStream create a client side stream
func NewStream(ctx context.Context, ri rpcinfo.RPCInfo, handler remote.ClientTransHandler, opt *remote.ClientOption) (streaming.Stream, error) {
	cm := NewConnWrapper(opt.ConnPool, opt.Logger)
	var err error
	for _, shdlr := range opt.StreamingMetaHandlers {
		ctx, err = shdlr.OnConnectStream(ctx)
		if err != nil {
			return nil, err
		}
	}
	rawConn, err := cm.GetConn(ctx, opt.Dialer, ri)
	if err != nil {
		return nil, err
	}
	return nhttp2.NewStream(ctx, opt.SvcInfo, rawConn, handler), nil
}
