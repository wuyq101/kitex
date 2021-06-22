package generic

import (
	"context"

	"github.com/cloudwego/kitex/pkg/remote"
)

// JSONRequest alias of string
type JSONRequest = string

type jsonThriftCodec struct {
}

func (c *jsonThriftCodec) UpdateIDL(main string, includes map[string]string) error {
	return nil
}

func (c *jsonThriftCodec) Marshal(ctx context.Context, message remote.Message, out remote.ByteBuffer) error {
	// ...
	return nil
}

func (c *jsonThriftCodec) Unmarshal(ctx context.Context, message remote.Message, in remote.ByteBuffer) error {
	// ...
	return nil
}

func (c *jsonThriftCodec) Name() string {
	return "JSONThrift"
}
