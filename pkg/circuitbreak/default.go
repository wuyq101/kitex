package circuitbreak

import (
	"context"
	"errors"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

// ErrorTypeOnServiceLevel determines the error type with a service level criteria.
func ErrorTypeOnServiceLevel(ctx context.Context, request, response interface{}, err error) ErrorType {
	if err != nil {
		var errTypes = map[error]ErrorType{
			kerrors.ErrInternalException: TypeIgnorable,
			kerrors.ErrServiceDiscovery:  TypeIgnorable,
			kerrors.ErrACL:               TypeIgnorable,
			kerrors.ErrLoadbalance:       TypeIgnorable,
			kerrors.ErrCircuitBreak:      TypeFailure,
			kerrors.ErrGetConnection:     TypeFailure,
			kerrors.ErrNoMoreInstance:    TypeFailure,
			kerrors.ErrRemoteOrNetwork:   TypeFailure,
			kerrors.ErrRPCTimeout:        TypeTimeout,
		}
		for e, t := range errTypes {
			if errors.Is(err, e) {
				return t
			}
		}
	}
	return TypeSuccess
}

// ErrorTypeOnInstanceLevel determines the error type with a instance level criteria.
// Basically, it treats only the connection error as failure.
func ErrorTypeOnInstanceLevel(ctx context.Context, request, response interface{}, err error) ErrorType {
	if errors.Is(err, kerrors.ErrGetConnection) {
		return TypeFailure
	}
	return TypeSuccess
}

// FailIfError return TypeFailure if err is not nil, otherwise TypeSuccess.
func FailIfError(ctx context.Context, request, response interface{}, err error) ErrorType {
	if err != nil {
		return TypeFailure
	}
	return TypeSuccess
}

// NoDecoration returns the original err.
func NoDecoration(ctx context.Context, request interface{}, err error) error {
	return err
}
