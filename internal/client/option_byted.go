package client

import (
	"fmt"
	"net"

	"github.com/cloudwego/kitex/internal/utils"
	"github.com/cloudwego/kitex/pkg/remote/trans/shmipc"
)

// WithShmIPC to enable shm-ipc mode in proxy scene
func WithShmIPC(concurrency int, udsAddr, fallbackAddr net.Addr) Option {
	return Option{F: func(o *Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithShmIPC(%d, %+v, %+v)", concurrency, udsAddr, fallbackAddr))

		o.RemoteOpt.CliHandlerFactory = shmipc.NewCliTransHandlerFactory()
	}}
}
