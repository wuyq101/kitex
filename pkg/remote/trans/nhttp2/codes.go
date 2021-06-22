package nhttp2

import (
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/remote/trans/nhttp2/codes"
	"github.com/cloudwego/kitex/pkg/remote/trans/nhttp2/status"
)

var (
	kitexErrConvTab = map[error]codes.Code{
		kerrors.ErrInternalException: codes.Internal,
		kerrors.ErrOverlimit:         codes.ResourceExhausted,
		kerrors.ErrRemoteOrNetwork:   codes.Unavailable,
		kerrors.ErrACL:               codes.PermissionDenied,
	}
	statusCodeConvTab = map[codes.Code]error{
		codes.Internal:          kerrors.ErrInternalException,
		codes.ResourceExhausted: kerrors.ErrOverlimit,
		codes.Unavailable:       kerrors.ErrRemoteOrNetwork,
		codes.PermissionDenied:  kerrors.ErrACL,
	}
)

func convertFromKitexToGrpc(err error) *status.Status {
	if err == nil {
		return status.New(codes.OK, "")
	}
	if c, ok := kitexErrConvTab[err]; ok {
		return status.New(c, err.Error())
	}
	return status.New(codes.Internal, err.Error())
}

func convertErrorFromGrpcToKitex(err error) error {
	s, ok := status.FromError(err)
	if !ok {
		return err
	}
	if s == nil || s.Code() == codes.OK {
		return nil
	}
	if e, ok := statusCodeConvTab[s.Code()]; ok {
		return e.(interface{ WithCause(cause error) error }).WithCause(s.Err())
	}
	return kerrors.ErrInternalException.WithCause(s.Err())
}
