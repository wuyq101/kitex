package invoke

import (
	"testing"

	"github.com/cloudwego/kitex/internal/test"
	"github.com/cloudwego/kitex/internal/utils"
	"github.com/cloudwego/kitex/pkg/remote/trans/invoke"
)

func TestNewMessage(t *testing.T) {
	lAddr := utils.NewNetAddr("tcp", "lAddr")
	rAddr := utils.NewNetAddr("tcp", "rAddr")
	msg := NewMessage(lAddr, rAddr)
	_, ok := msg.(invoke.Message)
	test.Assert(t, ok)
	test.Assert(t, msg.LocalAddr().String() == lAddr.String())
	test.Assert(t, msg.RemoteAddr().String() == rAddr.String())
}
