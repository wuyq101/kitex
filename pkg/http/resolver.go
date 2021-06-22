// Package http is used to implement RPC over HTTP.
package http

import (
	"net"
	"net/url"
	"strconv"
)

// Resolver resolves url to address.
type Resolver interface {
	Resolve(string) (string, error)
}

type defaultResolver struct{}

// NewDefaultResolver creates a default resolver.
func NewDefaultResolver() Resolver {
	return &defaultResolver{}
}

// Resolve implements the Resolver interface.
func (p *defaultResolver) Resolve(URL string) (string, error) {
	pu, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	host := pu.Hostname()
	port := pu.Port()
	if port == "" {
		port = "443"
		if pu.Scheme == "http" {
			port = "80"
		}
	}
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return "", err
	}
	return net.JoinHostPort(addr.IP.String(), strconv.Itoa(addr.Port)), nil
}
