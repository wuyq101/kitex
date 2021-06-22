package generic

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"

	gthrift "github.com/cloudwego/kitex/pkg/generic/thrift"
	codecThrift "github.com/cloudwego/kitex/pkg/remote/codec/thrift"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
)

// Service generic service interface
type Service interface {
	// GenericCall handle the generic call
	GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error)
}

// ServiceInfo create a generic ServiceInfo
func ServiceInfo(pcType serviceinfo.PayloadCodec) *serviceinfo.ServiceInfo {
	return newServiceInfo(pcType)
}

func newServiceInfo(pcType serviceinfo.PayloadCodec) *serviceinfo.ServiceInfo {
	serviceName := serviceinfo.GenericService
	handlerType := (*Service)(nil)
	methods := map[string]serviceinfo.MethodInfo{
		serviceinfo.GenericMethod: serviceinfo.NewMethodInfo(callHandler, newGenericServiceCallArgs, newGenericServiceCallResult, false),
	}
	svcInfo := &serviceinfo.ServiceInfo{
		ServiceName:  serviceName,
		HandlerType:  handlerType,
		Methods:      methods,
		PayloadCodec: pcType,
		Extra:        make(map[string]interface{}),
	}
	svcInfo.Extra["generic"] = true
	return svcInfo
}

func callHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*Args)
	realResult := result.(*Result)
	success, err := handler.(Service).GenericCall(ctx, realArg.Method, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}

func newGenericServiceCallArgs() interface{} {
	return &Args{}
}

func newGenericServiceCallResult() interface{} {
	return &Result{}
}

// WithCodec set codec instance for Args or Result
type WithCodec interface {
	SetCodec(codec interface{})
}

// Args generic request
type Args struct {
	Request interface{}
	Method  string
	base    *gthrift.Base
	inner   interface{}
}

var _ codecThrift.MessageReaderWithMethodWithContext = (*Args)(nil)
var _ codecThrift.MessageWriterWithContext = (*Args)(nil)
var _ WithCodec = (*Args)(nil)

// SetCodec ...
func (g *Args) SetCodec(inner interface{}) {
	g.inner = inner
}

func (g *Args) GetOrSetBase() interface{} {
	if g.base == nil {
		g.base = gthrift.NewBase()
	}
	return g.base
}

// Write ...
func (g *Args) Write(ctx context.Context, out thrift.TProtocol) error {
	if w, ok := g.inner.(gthrift.MessageWriter); ok {
		return w.Write(ctx, out, g.Request, g.base)
	}
	return fmt.Errorf("unexpected Args writer type: %T", g.inner)
}

// Read ...
func (g *Args) Read(ctx context.Context, method string, in thrift.TProtocol) error {
	if w, ok := g.inner.(gthrift.MessageReader); ok {
		var err error
		g.Request, err = w.Read(ctx, method, in)
		return err
	}
	return fmt.Errorf("unexpected Args reader type: %T", g.inner)
}

// Result generic response
type Result struct {
	Success interface{}
	inner   interface{}
}

var _ codecThrift.MessageReaderWithMethodWithContext = (*Result)(nil)
var _ codecThrift.MessageWriterWithContext = (*Result)(nil)
var _ WithCodec = (*Result)(nil)

// SetCodec ...
func (r *Result) SetCodec(inner interface{}) {
	r.inner = inner
}

// Write ...
func (r *Result) Write(ctx context.Context, out thrift.TProtocol) error {
	if w, ok := r.inner.(gthrift.MessageWriter); ok {
		return w.Write(ctx, out, r.Success, nil)
	}
	return fmt.Errorf("unexpected Result writer type: %T", r.inner)
}

// Read ...
func (r *Result) Read(ctx context.Context, method string, in thrift.TProtocol) error {
	if w, ok := r.inner.(gthrift.MessageReader); ok {
		var err error
		r.Success, err = w.Read(ctx, method, in)
		return err
	}
	return fmt.Errorf("unexpected Result reader type: %T", r.inner)
}

// GetSuccess ...
func (r *Result) GetSuccess() interface{} {
	if !r.IsSetSuccess() {
		return nil
	}
	return r.Success
}

// SetSuccess ...
func (r *Result) SetSuccess(x interface{}) {
	r.Success = x
}

// IsSetSuccess ...
func (r *Result) IsSetSuccess() bool {
	return r.Success != nil
}
