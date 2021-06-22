package netpoll

import (
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/netpoll"
)

// NewDialer returns the default netpoll dialer.
func NewDialer() remote.Dialer {
	return netpoll.NewDialer()
}
