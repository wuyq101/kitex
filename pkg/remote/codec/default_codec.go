/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package codec

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync/atomic"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/remote"
	"github.com/cloudwego/kitex/pkg/remote/codec/perrors"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
	"github.com/cloudwego/kitex/transport"
)

// The byte count of 32 and 16 integer values.
const (
	Size32 = 4
	Size16 = 2
)

const (
	// ThriftV1Magic is the magic code for thrift.VERSION_1
	ThriftV1Magic = 0x80010000
	// ProtobufV1Magic is the magic code for kitex protobuf
	ProtobufV1Magic = 0x90010000

	// MagicMask is bit mask for checking version.
	MagicMask = 0xffff0000
)

var ttHeaderCodec = ttHeader{}
var meshHeaderCodec = meshHeader{}

// NewDefaultCodec creates the default protocol sniffing codec supporting thrift and protobuf.
func NewDefaultCodec() remote.Codec {
	return &defaultCodec{}
}

type defaultCodec struct {
}

// Encode implements the remote.Codec interface, it does complete message encode include header and payload.
func (c *defaultCodec) Encode(ctx context.Context, message remote.Message, out remote.ByteBuffer) error {
	defer func() {
		ri := message.RPCInfo()
		// notice: mallocLen() must exec before flush, or it will be reset
		rpcinfo.AsMutableRPCStats(ri.Stats()).SetSendSize(uint64(out.MallocLen()))
	}()
	var err error
	var totalLenField []byte
	var framedLenField []byte
	var headerLen int
	tp := message.ProtocolInfo().TransProto

	// 1. encode header and return totalLenField if needed
	// totalLenField will be filled after payload encoded
	if tp == transport.TTHeader || tp == transport.TTHeaderFramed {
		if totalLenField, err = ttHeaderCodec.encode(ctx, message, out); err != nil {
			return err
		}
		headerLen = out.MallocLen()
	}
	// 2. malloc framed field if needed
	if tp == transport.TTHeaderFramed || tp == transport.Framed {
		if framedLenField, err = out.Malloc(Size32); err != nil {
			return err
		}
		headerLen += Size32
	}

	// 3. encode payload
	if err := c.encodePayload(ctx, message, out); err != nil {
		return err
	}

	// 4. fill framed field if needed
	if tp == transport.TTHeaderFramed || tp == transport.Framed {
		if framedLenField == nil {
			return perrors.NewProtocolErrorWithMsg("no buffer allocated for the framed length field")
		}
		binary.BigEndian.PutUint32(framedLenField, uint32(out.MallocLen()-headerLen))
	} else if message.ServiceInfo().PayloadCodec == serviceinfo.Protobuf {
		return perrors.NewProtocolErrorWithMsg("protobuf just support 'framed' trans proto")
	}
	// 5. fill totalLen field for header if needed
	if tp == transport.TTHeader || tp == transport.TTHeaderFramed {
		if totalLenField == nil {
			return perrors.NewProtocolErrorWithMsg("no buffer allocated for the header length field")
		}
		binary.BigEndian.PutUint32(totalLenField, uint32(out.MallocLen()-Size32))
	}
	return nil
}

// Decode implements the remote.Codec interface, it does complete message decode include header and payload.
func (c *defaultCodec) Decode(ctx context.Context, message remote.Message, in remote.ByteBuffer) (err error) {
	defer func() {
		ri := message.RPCInfo()
		rpcinfo.AsMutableRPCStats(ri.Stats()).SetRecvSize(uint64(in.ReadLen()))
	}()

	var flagBuf []byte
	if flagBuf, err = in.Peek(2 * Size32); err != nil {
		return perrors.NewProtocolErrorWithErrMsg(err, fmt.Sprintf("default codec read failed: %s", err.Error()))
	}

	if err = checkRPCState(ctx); err != nil {
		// there is one call has finished in retry task, it doesn't need to do decode for this call
		return err
	}
	// 1. decode header
	if IsTTHeader(flagBuf) {
		// TTHeader
		if err = ttHeaderCodec.decode(ctx, message, in); err != nil {
			return err
		}
		if flagBuf, err = in.Peek(2 * Size32); err != nil {
			return perrors.NewProtocolErrorWithErrMsg(err, fmt.Sprintf("ttheader read payload fisrt 8 byte failed: %s", err.Error()))
		}
		if err = checkPayload(flagBuf, message, in, true); err != nil {
			return err
		}
	} else if isMeshHeader(flagBuf) {
		message.Tags()[remote.MeshHeader] = true
		// MeshHeader
		if err = meshHeaderCodec.decode(ctx, message, in); err != nil {
			return err
		}
		if flagBuf, err = in.Peek(2 * Size32); err != nil {
			return perrors.NewProtocolErrorWithErrMsg(err, fmt.Sprintf("meshHeader read payload fisrt 8 byte failed: %s", err.Error()))
		}
		if err = checkPayload(flagBuf, message, in, false); err != nil {
			return err
		}
	} else {
		// no Header
		if err = checkPayload(flagBuf, message, in, false); err != nil {
			return err
		}
	}

	// 2. decode body
	if err := c.decodePayload(ctx, message, in); err != nil {
		return err
	}
	return nil
}

func (c *defaultCodec) Name() string {
	return "default"
}

// Select to use thrift or protobuf according to the protocol.
func (c *defaultCodec) encodePayload(ctx context.Context, message remote.Message, out remote.ByteBuffer) error {
	pCodec, err := remote.GetPayloadCodec(message)
	if err != nil {
		return err
	}
	return pCodec.Marshal(ctx, message, out)
}

func (c *defaultCodec) decodePayload(ctx context.Context, message remote.Message, in remote.ByteBuffer) error {
	pCodec, err := remote.GetPayloadCodec(message)
	if err != nil {
		return err
	}
	return pCodec.Unmarshal(ctx, message, in)
}

/**
 * +------------------------------------------------------------+
 * |                  4Byte                 |       2Byte       |
 * +------------------------------------------------------------+
 * |   			     Length			    	|   HEADER MAGIC    |
 * +------------------------------------------------------------+
 */
func IsTTHeader(flagBuf []byte) bool {
	return binary.BigEndian.Uint32(flagBuf[Size32:])&MagicMask == TTHeaderMagic
}

/**
 * +----------------------------------------+
 * |       2Byte        |       2Byte       |
 * +----------------------------------------+
 * |    HEADER MAGIC    |   HEADER SIZE     |
 * +----------------------------------------+
 */
func isMeshHeader(flagBuf []byte) bool {
	return binary.BigEndian.Uint32(flagBuf[:Size32])&MagicMask == MeshHeaderMagic
}

/**
 * protobuf 默认有Framed
 * +------------------------------------------------------------+
 * |                  4Byte                 |       2Byte       |
 * +------------------------------------------------------------+
 * |   			     Length			    	|   HEADER MAGIC    |
 * +------------------------------------------------------------+
 */
func isProtobufKitex(flagBuf []byte) bool {
	return binary.BigEndian.Uint32(flagBuf[Size32:])&MagicMask == ProtobufV1Magic
}

/**
 * +-------------------+
 * |       2Byte       |
 * +-------------------+
 * |   HEADER MAGIC    |
 * +-------------------
 */
func isThriftBinary(flagBuf []byte) bool {
	return binary.BigEndian.Uint32(flagBuf[:Size32])&MagicMask == ThriftV1Magic
}

/**
 * +------------------------------------------------------------+
 * |                  4Byte                 |       2Byte       |
 * +------------------------------------------------------------+
 * |   			     Length			    	|   HEADER MAGIC    |
 * +------------------------------------------------------------+
 */
func isThriftFramedBinary(flagBuf []byte) bool {
	return binary.BigEndian.Uint32(flagBuf[Size32:])&MagicMask == ThriftV1Magic
}

func checkRPCState(ctx context.Context) error {
	if respOp, ok := ctx.Value(retry.CtxRespOp).(*int32); ok && atomic.LoadInt32(respOp) == retry.OpDone {
		return kerrors.ErrRPCFinish
	}
	return nil
}

func checkPayload(flagBuf []byte, message remote.Message, in remote.ByteBuffer, isTTHeader bool) error {
	var transProto transport.Protocol
	var codecType serviceinfo.PayloadCodec
	if isThriftBinary(flagBuf) {
		codecType = serviceinfo.Thrift
		if isTTHeader {
			transProto = transport.TTHeader
		} else {
			transProto = transport.PurePayload
		}
	} else if isThriftFramedBinary(flagBuf) {
		codecType = serviceinfo.Thrift
		if isTTHeader {
			transProto = transport.TTHeaderFramed
		} else {
			transProto = transport.Framed
		}
		payloadLen := binary.BigEndian.Uint32(flagBuf[:Size32])
		message.SetPayloadLen(int(payloadLen))
		if _, err := in.Next(Size32); err != nil {
			return err
		}
	} else if isProtobufKitex(flagBuf) {
		codecType = serviceinfo.Protobuf
		if isTTHeader {
			transProto = transport.TTHeaderFramed
		} else {
			transProto = transport.Framed
		}
		payloadLen := binary.BigEndian.Uint32(flagBuf[:Size32])
		message.SetPayloadLen(int(payloadLen))
		if _, err := in.Next(Size32); err != nil {
			return err
		}
	} else {
		first4Bytes := binary.BigEndian.Uint32(flagBuf[:Size32])
		second4Bytes := binary.BigEndian.Uint32(flagBuf[Size32:])
		// 0xfff4fffd is the interrupt message of telnet
		err := perrors.NewProtocolErrorWithMsg(fmt.Sprintf("invalid payload (first4Bytes=%#x, second4Bytes=%#x)", first4Bytes, second4Bytes))
		return err
	}
	message.SetProtocolInfo(remote.NewProtocolInfo(transProto, codecType))
	return nil
}
