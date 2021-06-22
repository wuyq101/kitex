package proxy

import (
	"context"
	"net"

	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/remote"
)

// Config contains basic components used in service discovery process.
type Config struct {
	ServiceName  string
	Resolver     discovery.Resolver
	Balancer     loadbalance.Loadbalancer
	Pool         remote.ConnPool
	FixedTargets string // A comma separated list of host ports that user specify to use.
}

// ForwardProxy manages the service discovery, load balance and connection pooling processes.
type ForwardProxy interface {
	// Configure is provided to initialize the proxy.
	Configure(*Config) error

	// ResolveProxyInstance set instance for remote endpoint.
	ResolveProxyInstance(ctx context.Context) error
}

// BackwardProxy replaces the listen address with another one.
type BackwardProxy interface {
	Replace(net.Addr) (net.Addr, error)
}
