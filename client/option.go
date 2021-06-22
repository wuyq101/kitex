package client

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/cloudwego/kitex/internal/client"
	"github.com/cloudwego/kitex/internal/utils"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/http"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/loadbalance/lbcache"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/trans/netpollmux"
	"github.com/cloudwego/kitex/pkg/remote/trans/nhttp2"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/stats"
	"github.com/cloudwego/kitex/transport"
)

// Option is the only way to config client.
type Option = client.Option

// A Suite is a collection of Options. It is useful to assemble multiple associated
// Options as a single one to keep the order or presence in a desired manner.
type Suite interface {
	Options() []Option
}

// WithTransportProtocol sets the transport protocol for client.
func WithTransportProtocol(tp transport.Protocol) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		tpName := tp.String()
		if tpName == transport.Unknown {
			panic(fmt.Errorf("WithTransportProtocol: invalid '%v'", tp))
		}
		if tp == transport.GRPC {
			o.RemoteOpt.ConnPool = nhttp2.NewConnPool()
			o.RemoteOpt.CliHandlerFactory = nhttp2.NewCliTransHandlerFactory()
		}
		di.Push(fmt.Sprintf("WithTransportProtocol(%s)", tpName))
		rpcinfo.AsMutableRPCConfig(o.Configs).SetTransportProtocol(tp)
	}}
}

// WithSuite adds a option suite for client.
func WithSuite(suite Suite) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		var nested struct {
			Suite   string
			Options utils.Slice
		}
		nested.Suite = fmt.Sprintf("%T(%+v)", suite, suite)

		for _, op := range suite.Options() {
			op.F(o, &nested.Options)
		}
		di.Push(nested)
	}}
}

// WithMiddleware adds middleware for client to handle request.
func WithMiddleware(mw endpoint.Middleware) Option {
	mwb := func(ctx context.Context) endpoint.Middleware {
		return mw
	}
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithMiddleware(%+v)", utils.GetFuncName(mw)))
		o.MWBs = append(o.MWBs, mwb)
	}}
}

// WithMiddlewareBuilder adds middleware that depend on context for client to handle request
func WithMiddlewareBuilder(mwb endpoint.MiddlewareBuilder) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithMiddlewareBuilder(%+v)", utils.GetFuncName(mwb)))
		o.MWBs = append(o.MWBs, mwb)
	}}
}

// WithInstanceMW adds middleware for client to handle request after service discovery and loadbalance process.
func WithInstanceMW(mw endpoint.Middleware) Option {
	imwb := func(ctx context.Context) endpoint.Middleware {
		return mw
	}
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithInstanceMW(%+v)", utils.GetFuncName(mw)))
		o.IMWBs = append(o.IMWBs, imwb)
	}}
}

// WithDestService specifies the name of target service.
func WithDestService(svr string) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		o.Once.OnceOrPanic()
		di.Push(fmt.Sprintf("WithDestService(%s)", svr))
		o.Svr.ServiceName = svr
	}}
}

// WithHostPorts specifies the target instance addresses when doing service discovery.
// It overwrites the results from the Resolver.
func WithHostPorts(hostports ...string) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithHostPorts(%v)", hostports))

		var ins []discovery.Instance
		for _, hp := range hostports {
			if _, err := net.ResolveTCPAddr("tcp", hp); err == nil {
				ins = append(ins, discovery.NewInstance("tcp", hp, discovery.DefaultWeight, nil))
				continue
			}
			if _, err := net.ResolveUnixAddr("unix", hp); err == nil {
				ins = append(ins, discovery.NewInstance("unix", hp, discovery.DefaultWeight, nil))
				continue
			}
			panic(fmt.Errorf("WithHostPorts: invalid '%s'", hp))
		}

		if len(ins) == 0 {
			panic("WithHostPorts() requires at least one argument")
		}

		o.Targets = strings.Join(hostports, ",")
		o.Resolver = &discovery.SynthesizedResolver{
			ResolveFunc: func(ctx context.Context, key string) (discovery.Result, error) {
				return discovery.Result{
					Cacheable: true,
					CacheKey:  "fixed",
					Instances: ins,
				}, nil
			},
			NameFunc: func() string { return o.Targets },
		}
	}}
}

// WithResolver provides the Resolver for kitex client.
func WithResolver(r discovery.Resolver) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithResolver(%T)", r))

		o.Resolver = r
	}}
}

// WithHTTPResolver specifies resolver for url (which specified by WithURL).
func WithHTTPResolver(r http.Resolver) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithHTTPResolver(%T)", r))

		o.HTTPResolver = r
	}}
}

// WithLongConnection enables long connection with kitex's built-in pooling implementation.
func WithLongConnection(cfg connpool.IdleConfig) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithLongConnection(%+v)", cfg))

		if o.PoolCfg == nil {
			o.PoolCfg = connpool.CheckPoolConfig(cfg)
		}
	}}
}

// WithMuxConnection specifies the transport type to be mux.
// IMPORTANT: this option is not stable, and will be changed or removed in the future!!!
// We don't promise compatibility for this option in future versions!!!
func WithMuxConnection(connNum int) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push("WithMuxConnection")

		o.RemoteOpt.ConnPool = netpollmux.NewMuxConnPool(connNum)
		client.WithTransHandlerFactory(netpollmux.NewCliTransHandlerFactory()).F(o, di)
		rpcinfo.AsMutableRPCConfig(o.Configs).SetTransportProtocol(transport.TTHeader)
	}}
}

// WithLogger sets the Logger for kitex client.
func WithLogger(logger klog.FormatLogger) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithLogger(%T)", logger))

		o.Logger = logger
	}}
}

// WithLoadBalancer sets the loadbalancer for client.
func WithLoadBalancer(lb loadbalance.Loadbalancer, opts ...*lbcache.Options) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithLoadBalancer(%+v, %+v)", lb, opts))
		o.Balancer = lb
		if len(opts) > 0 {
			o.BalancerCacheOpt = opts[0]
		}
	}}
}

func lockTargetTag(opt, key, val string) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("%s(%s)", opt, val))

		o.Svr.Tags[key] = val
		o.Locks.Tags[key] = struct{}{}
	}}
}

// WithRPCTimeout specifies the RPC timeout.
func WithRPCTimeout(d time.Duration) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithRPCTimeout(%dms)", d.Milliseconds()))

		rpcinfo.AsMutableRPCConfig(o.Configs).SetRPCTimeout(d)
		o.Locks.Bits |= rpcinfo.BitRPCTimeout
	}}
}

// WithConnectTimeout specifies the connection timeout.
func WithConnectTimeout(d time.Duration) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithConnectTimeout(%dms)", d.Milliseconds()))

		rpcinfo.AsMutableRPCConfig(o.Configs).SetConnectTimeout(d)
		o.Locks.Bits |= rpcinfo.BitConnectTimeout
	}}
}

// WithCluster sets the target cluster for service discovery.
func WithCluster(cluster string) Option {
	return lockTargetTag("WithCluster", "cluster", cluster)
}

// WithIDC sets the target IDC for service discovery.
func WithIDC(idc string) Option {
	return lockTargetTag("WithIDC", "idc", idc)
}

// WithTracer adds a tracer to client.
func WithTracer(c stats.Tracer) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithTracer(%T{%+v})", c, c))

		o.TracerCtl.Append(c)
	}}
}

// WithStatsLevel sets the stats level for client.
func WithStatsLevel(level stats.Level) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithStatsLevel(%+v)", level))

		o.StatsLevel = level
	}}
}

// WithCodec to set a codec that handle other protocols which not support by kitex
func WithCodec(c remote.Codec) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithCodec(%+v)", c))

		o.RemoteOpt.Codec = c
	}}
}

// WithPayloadCodec to set a payloadCodec that handle other payload which not support by kitex
func WithPayloadCodec(c remote.PayloadCodec) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithPayloadCodec(%+v)", c))

		o.RemoteOpt.PayloadCodec = c
	}}
}

// WithConnReporterEnabled to enable reporting connection pool stats.
func WithConnReporterEnabled() Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push("WithConnReporterEnabled()")

		o.RemoteOpt.EnableConnPoolReporter = true
	}}
}

// WithFailureRetry sets the failure retry policy for client.
func WithFailureRetry(p *retry.FailurePolicy) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		if p == nil {
			return
		}
		di.Push(fmt.Sprintf("WithFailureRetry(%+v)", *p))
		if o.RetryPolicy == nil {
			o.RetryPolicy = &retry.Policy{}
		}
		if o.RetryPolicy.BackupPolicy != nil {
			panic("BackupPolicy has been setup, cannot support Failure Retry at same time")
		}
		o.RetryPolicy.FailurePolicy = p
		o.RetryPolicy.Enable = true
		o.RetryPolicy.Type = retry.FailureType
	}}
}

// WithBackupRequest sets the backup request policy for client.
func WithBackupRequest(p *retry.BackupPolicy) Option {
	return Option{F: func(o *client.Options, di *utils.Slice) {
		if p == nil {
			return
		}
		di.Push(fmt.Sprintf("WithBackupRequest(%+v)", *p))
		if o.RetryPolicy == nil {
			o.RetryPolicy = &retry.Policy{}
		}
		if o.RetryPolicy.FailurePolicy != nil {
			panic("Failure Retry has been setup, cannot support Backup Request at same time")
		}
		o.RetryPolicy.BackupPolicy = p
		o.RetryPolicy.Enable = true
		o.RetryPolicy.Type = retry.BackupType
	}}
}

// WithMetaHandler adds a MetaHandler.
func WithMetaHandler(h remote.MetaHandler) client.Option {
	return client.Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithMetaHandler(%T)", h))

		o.MetaHandlers = append(o.MetaHandlers, h)
	}}
}
