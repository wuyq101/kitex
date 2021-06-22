package remote

import (
	"context"
	"net"
	"time"

	internal_stats "github.com/cloudwego/kitex/internal/stats"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
)

// Option is used to pack the inbound and outbound handlers.
type Option struct {
	Outbounds []OutboundHandler

	Inbounds []InboundHandler

	StreamingMetaHandlers []StreamingMetaHandler
}

// ServerOption contains option that is used to init the remote server.
type ServerOption struct {
	SvcInfo *serviceinfo.ServiceInfo

	TransServerFactory TransServerFactory

	SvrHandlerFactory ServerTransHandlerFactory

	Codec Codec

	PayloadCodec PayloadCodec

	Address net.Addr

	// Duration that server waits for to allow any existing connection to be closed gracefully.
	ExitWaitTime time.Duration

	// Duration that server waits for after error occurs during connection accepting.
	AcceptFailedDelayTime time.Duration

	// Duration that the accepted connection waits for to read or write data.
	MaxConnectionIdleTime time.Duration

	ReadWriteTimeout time.Duration

	InitRPCInfoFunc func(context.Context, net.Addr) (rpcinfo.RPCInfo, context.Context)

	TracerCtl *internal_stats.Controller

	Logger klog.FormatLogger

	Option
}

// ClientOption is used to init the remote client.
type ClientOption struct {
	SvcInfo *serviceinfo.ServiceInfo

	CliHandlerFactory ClientTransHandlerFactory

	Codec Codec

	PayloadCodec PayloadCodec

	ConnPool ConnPool

	Dialer Dialer

	Logger klog.FormatLogger

	Option

	EnableConnPoolReporter bool
}
