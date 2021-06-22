// Package acl implements ACL functionality.
package acl

import (
	"context"
	"errors"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

// RejectFunc judges if to reject a request by the given context and request.
// Returns a reason if rejected, otherwise returns nil.
type RejectFunc func(ctx context.Context, request interface{}) (reason error)

// NewACLMiddleware creates a new ACL middleware using the provided reject funcs.
func NewACLMiddleware(rules []RejectFunc) endpoint.Middleware {
	if len(rules) == 0 {
		return endpoint.DummyMiddleware
	}
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request, response interface{}) error {
			for _, r := range rules {
				if e := r(ctx, request); e != nil {
					if !errors.Is(e, kerrors.ErrACL) {
						e = kerrors.ErrACL.WithCause(e)
					}
					return e
				}
			}
			return next(ctx, request, response)
		}
	}
}
