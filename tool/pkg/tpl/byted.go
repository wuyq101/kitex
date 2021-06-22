package tpl

func init() {
	clientOptionTpl = append(clientOptionTpl, bytedClientOptionTpl)
	clientExtendTpl = append(clientExtendTpl, bytedClientExtendTpl)
	invokerOptionTpl = append(invokerOptionTpl, bytedInvokerOptionTpl)
	invokerExtendTpl = append(invokerExtendTpl, bytedInvokerExtendTpl)
	serverOptionTpl = append(serverOptionTpl, bytedServerOptionTpl)
	serverExtendTpl = append(serverExtendTpl, bytedServerExtendTpl)
}

var bytedClientOptionTpl = `
{{if HasFeature .Features "WithByted"}}
config := byted.NewClientConfig()
config.DestService = destService
{{if .HasStreaming}}
config.Transport = byted.GRPC
options = append(options, client.WithTransportProtocol(transport.GRPC))
{{end}}
options = append(options, byted.ClientSuiteWithConfig(serviceInfo(), config))
{{end}}`

var bytedClientExtendTpl = `
{{if HasFeature .Features "WithByted"}}
// NewClientWithBytedConfig creates a client for the service defined in IDL.
func NewClientWithBytedConfig(destService string, config *byted.ClientConfig, opts ...client.Option) (Client, error) {
	if config == nil {
		config = byted.NewClientConfig()
	}
	config.DestService = destService

	var options []client.Option
	options = append(options, client.WithDestService(destService))
	{{if .HasStreaming}}
	config.Transport = byted.GRPC
	options = append(options, client.WithTransportProtocol(transport.GRPC))
	{{end}}
	options = append(options, byted.ClientSuiteWithConfig(serviceInfo(), config))
	options = append(options, opts...)
    kc, err := client.NewClient(serviceInfo(), options...)
    if err != nil {
        return nil, err
	}
	return &k{{.ServiceName}}Client{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClientWithBytedConfig creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClientWithBytedConfig(destService string, config *byted.ClientConfig, opts ...client.Option) Client {
	kc, err := NewClientWithBytedConfig(destService, config, opts...)
    if err != nil {
        panic(err)
    }
    return kc
}
{{end}}`

var bytedInvokerOptionTpl = `
{{if HasFeature .Features "WithByted"}} options = append(options, byted.InvokeSuite(serviceInfo())){{end}}
`

var bytedInvokerExtendTpl = `
{{if HasFeature .Features "WithByted"}}
// NewInvokerWithBytedConfig creates a server.Invoker with the given handler and options.
func NewInvokerWithBytedConfig(handler {{call .ServiceTypeName}}, config *byted.ServerConfig, opts ...server.Option) server.Invoker {
	var options []server.Option
	options = append(options, byted.InvokeSuiteWithConfig(serviceInfo(), config))
	options = append(options, opts...)

    s := server.NewInvoker(options...)
    if err := s.RegisterService(serviceInfo(), handler); err != nil {
            panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}
    return s
}
{{end}}`

var bytedServerOptionTpl = `
{{if HasFeature .Features "WithByted"}} options = append(options, byted.ServerSuite(serviceInfo())){{end}}
`

var bytedServerExtendTpl = `
{{if HasFeature .Features "WithByted"}}
// NewServerWithBytedConfig creates a server.Server with the given handler and options.
func NewServerWithBytedConfig(handler {{call .ServiceTypeName}}, config *byted.ServerConfig, opts ...server.Option) server.Server {
	var options []server.Option
	options = append(options, byted.ServerSuiteWithConfig(serviceInfo(), config))
	options = append(options, opts...)

    svr := server.NewServer(options...)
    if err := svr.RegisterService(serviceInfo(), handler); err != nil {
            panic(err)
    }
    return svr
}
{{end}}`
