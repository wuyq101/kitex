package client

import (
	"github.com/apache/thrift/lib/go/thrift"
)

// MockTStruct implements the thrift.TStruct interface.
type MockTStruct struct {
	WriteFunc func(p thrift.TProtocol) (e error)
	ReadFunc  func(p thrift.TProtocol) (e error)
}

// Write implements the thrift.TStruct interface.
func (m MockTStruct) Write(p thrift.TProtocol) (e error) {
	if m.WriteFunc != nil {
		return m.WriteFunc(p)
	}
	return
}

// Read implements the thrift.TStruct interface.
func (m MockTStruct) Read(p thrift.TProtocol) (e error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(p)
	}
	return
}
