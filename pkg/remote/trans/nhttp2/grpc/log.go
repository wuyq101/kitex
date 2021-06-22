package grpc

import (
	"github.com/cloudwego/kitex/pkg/klog"
)

func infof(format string, args ...interface{}) {
	klog.DefaultLogger().Infof(format, args...)
}

func warningf(format string, args ...interface{}) {
	klog.DefaultLogger().Warnf(format, args...)
}

func errorf(format string, args ...interface{}) {
	klog.DefaultLogger().Errorf(format, args...)
}
