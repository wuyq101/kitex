package server

import (
	"fmt"
	"net"

	"github.com/cloudwego/kitex/internal/utils"
	"github.com/cloudwego/kitex/pkg/remote/trans/shmipc"
)

// WithShmIPC to enable shm-ipc mode in proxy scene
func WithShmIPC(udsAddr, fallbackAddr net.Addr) Option {
	return Option{F: func(o *Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithShmIPC(%+v, %+v)", udsAddr, fallbackAddr))

		o.RemoteOpt.TransServerFactory = shmipc.NewTransServerFactory(fallbackAddr)
		o.RemoteOpt.SvrHandlerFactory = shmipc.NewSvrTransHandlerFactory()
	}}
}
