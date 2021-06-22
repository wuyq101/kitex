// Package thrift thrift parser and codec
package thrift

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
)

// MessageReader read from thrift.TProtocol with method
type MessageReader interface {
	Read(ctx context.Context, method string, in thrift.TProtocol) (interface{}, error)
}

// MessageWriter write to thrift.TProtocol
type MessageWriter interface {
	Write(ctx context.Context, out thrift.TProtocol, msg interface{}, requestBase *Base) error
}
