// Package genericclient ...
package genericclient

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	internal_client "github.com/cloudwego/kitex/internal/client"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
)

// NewClient create a generic client
func NewClient(psm string, g generic.Generic, opts ...client.Option) (Client, error) {
	svcInfo := generic.ServiceInfo(g.PayloadCodecType())
	return NewClientWithServiceInfo(psm, g, svcInfo, opts...)
}

// NewClientWithServiceInfo create a generic client with serviceInfo
func NewClientWithServiceInfo(psm string, g generic.Generic, svcInfo *serviceinfo.ServiceInfo, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, internal_client.WithGeneric(g))
	options = append(options, client.WithDestService(psm))
	options = append(options, opts...)

	kc, err := client.NewClient(svcInfo, options...)
	if err != nil {
		return nil, err
	}
	return &genericServiceClient{
		kClient: kc,
		g:       g,
	}, nil
}

// Client generic client
type Client interface {
	// GenericCall generic call
	GenericCall(ctx context.Context, method string, request interface{}, callOptions ...callopt.Option) (response interface{}, err error)
}

type genericServiceClient struct {
	kClient client.Client
	g       generic.Generic
}

func (gc *genericServiceClient) GenericCall(ctx context.Context, method string, request interface{}, callOptions ...callopt.Option) (response interface{}, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	var _args generic.Args
	_args.Method = method
	_args.Request = request
	mt, err := gc.g.GetMethod(request, method)
	if err != nil {
		return nil, err
	}
	if mt.Oneway {
		return nil, gc.kClient.Call(ctx, mt.Name, &_args, nil)
	}
	var _result generic.Result
	if err = gc.kClient.Call(ctx, mt.Name, &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
