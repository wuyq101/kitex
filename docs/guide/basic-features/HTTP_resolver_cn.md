# HTTPResolver

## 指定 URL 进行调用

在进行调用时，可以通过`callopt.WithURL`指定，通过该 option 指定的URL，会经过默认的 DNS resolver 解析后拿到 host 和 port，此时其等效于`callopt.WithHostPort`。

```go
import "github.com/cloudwego/kitex/client/callopt"
...
url := callopt.WithURL("http://myserverdomain.com:8888")
resp, err := cli.Echo(context.Background(), req, url)
if err != nil {
	log.Fatal(err)
}
```

## 自定义 DNS resolver

此外也可以自定义 DNS resolver

resolver定义如下(pkg/http)：

```go
type Resolver interface {
	Resolve(string) (string, error)
}
```

参数为 URL，返回值为访问的 server 的 "host:port"。

通过`client.WithHTTPResolver`指定用于 DNS 解析的resolver。

```go
import "github.com/cloudwego/kitex/client/callopt"
...
dr := client.WithHTTPResolver(myResolver)
cli, err := echo.NewClient("echo", dr)
if err != nil {
	log.Fatal(err)
}
```

# Server SDK化

SDK化（invoker）允许用户将 KiteX server 当作一个本地 SDK 调用。

调用通过 `message` 完成，初始化 `message` 需要 `local` 和 `remote` 两个 `net.Addr` ，分别表示本地地址和远端（客户端）地址（此处的地址主要用于日志监控），初始化后通过 `SetRequestBytes(buf []byte) error` 设置请求的二进制数据。最后调用 `invoker` 的 `Call` 方法即可完成调用。调用完成后可通过 `message` 的 `GetResponseBytes() ([]byte, error)` 获取响应的二进制数据。

```go
import (
		...
    "github.com/cloudwego/kitex/sdk/message"
  	...
)

func main() {
    var reqPayload, respPayload []byte
    var local, remote net.Addr
    ...
    // init local/remote
    ...
    ivk := echo.NewInvoker(new(EchoImpl))
    msg := message.NewMessage(local, remote)
    // 装载payload
    msg.SetRequestBytes(reqPayload)
    // 发起调用
    err := ivk.Call(msg)
    if err != nil {
        ...
    }
    respPayload, err = msg.GetResponseBytes()
    if err != nil {
        ...
    }
}
```