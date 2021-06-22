package bound

import (
	"context"
	"net"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/limiter"
	"github.com/cloudwego/kitex/pkg/remote"
)

// NewServerLimiterHandler creates a new server limiter handler.
func NewServerLimiterHandler(conLimit limiter.ConcurrencyLimiter, qpsLimit limiter.RateLimiter, reporter limiter.LimitReporter) remote.InboundHandler {
	return &serverLimiterHandler{conLimit, qpsLimit, reporter}
}

type serverLimiterHandler struct {
	conLimit limiter.ConcurrencyLimiter
	qpsLimit limiter.RateLimiter
	reporter limiter.LimitReporter
}

// OnActive implements the remote.InboundHandler interface.
func (l *serverLimiterHandler) OnActive(ctx context.Context, conn net.Conn) (context.Context, error) {
	if l.conLimit.Acquire() {
		return ctx, nil
	}
	if l.reporter != nil {
		l.reporter.ConnOverloadReport()
	}
	return ctx, kerrors.ErrConnOverLimit
}

// OnRead implements the remote.InboundHandler interface.
func (l *serverLimiterHandler) OnRead(ctx context.Context, conn net.Conn) (context.Context, error) {
	if l.qpsLimit.Acquire() {
		return ctx, nil
	}
	if l.reporter != nil {
		l.reporter.QPSOverloadReport()
	}
	return ctx, kerrors.ErrQPSOverLimit
}

// OnInactive implements the remote.InboundHandler interface.
func (l *serverLimiterHandler) OnInactive(ctx context.Context, conn net.Conn) context.Context {
	l.conLimit.Release()
	return ctx
}

// OnMessage implements the remote.InboundHandler interface.
func (l *serverLimiterHandler) OnMessage(ctx context.Context, args, result remote.Message) (context.Context, error) {
	return ctx, nil
}
