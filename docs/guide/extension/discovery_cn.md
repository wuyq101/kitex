# 服务发现扩展

框架自身不带服务发现，[kitex-contrib](https://github.com/kitex-contrib) 中提供了 dns 服务发现扩展。

用户如果需要更换其他的服务发现，例如 ETCD，用户可以根据需求实现 `Resolver ` 接口，client 通过 `WithResolver` Option 来注入。

```go
// Result contains the result of service discovery process.
// Cacheable tells whether the instance list can/should be cached.
// When Cacheable is true, CacheKey can be used to map the instance list in cache.
type Result struct {
	  Cacheable bool
	  CacheKey  string
	  Instances []Instance
}

// Change contains the difference between the current discovery result and the previous one.
// It is designed for providing detail information when dispatching an event for service
// discovery result change.
// Since the loadbalancer may rely on caching the result of resolver to improve performance,
// the resolver implementation should dispatch an event when result changes.
type Change struct {
	  Result  Result
	  Added   []Instance
	  Updated []Instance
	  Removed []Instance
}

// Resolver resolves the target endpoint into a list of Instance.
type Resolver interface {
	  Resolve(ctx context.Context, target rpcinfo.EndpointInfo) (Result, error)
}
```

