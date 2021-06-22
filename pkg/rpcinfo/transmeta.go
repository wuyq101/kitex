package rpcinfo

import "context"

// MetaHandler reads or writes metadata throught certain protocol.
// The protocol will be a thrift.TProtocol usually.
type MetaHandler interface {
	WriteMeta(ctx context.Context, request, response, protocol interface{}) (context.Context, error)
	ReadMeta(ctx context.Context, request, response, protocol interface{}) (context.Context, error)
}

// SynthesizedMetaHandler wraps read/write functions to be a MetaHandler.
type SynthesizedMetaHandler struct {
	Writer func(ctx context.Context, request, response, protocol interface{}) (context.Context, error)
	Reader func(ctx context.Context, request, response, protocol interface{}) (context.Context, error)
}

// WriteMeta implements the MetaHandler interface.
func (sm *SynthesizedMetaHandler) WriteMeta(ctx context.Context, request, response, protocol interface{}) (context.Context, error) {
	if sm.Writer != nil {
		return sm.Writer(ctx, request, response, protocol)
	}
	return ctx, nil
}

// ReadMeta implements the MetaHandler interface.
func (sm *SynthesizedMetaHandler) ReadMeta(ctx context.Context, request, response, protocol interface{}) (context.Context, error) {
	if sm.Reader != nil {
		return sm.Reader(ctx, request, response, protocol)
	}
	return ctx, nil
}

// ChainHandlers combines multiple MetaHandlers as one.
func ChainHandlers(hs []MetaHandler) MetaHandler {
	return &chained{hs: hs}
}

type chained struct {
	hs []MetaHandler
}

// WriteMeta implements the MetaHandler interface.
func (c *chained) WriteMeta(ctx context.Context, request, response, protocol interface{}) (context.Context, error) {
	var err error
	for _, h := range c.hs {
		ctx, err = h.WriteMeta(ctx, request, response, protocol)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

// ReadMeta implements the MetaHandler interface.
func (c *chained) ReadMeta(ctx context.Context, request, response, protocol interface{}) (context.Context, error) {
	var err error
	for _, h := range c.hs {
		ctx, err = h.ReadMeta(ctx, request, response, protocol)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}
