// Package transmeta .
package transmeta

// Protocol constants.
const (
	MeshTHeaderProtocolVersion = "1.0.0"
	MeshOrTTHeaderProtocol     = "3"
)

// Keys in mesh header.
const (
	MeshVersion uint16 = iota
	TransportType
	LogID
	FromService
	FromCluster
	FromIDC
	ToService
	ToCluster
	ToIDC
	ToMethod
	Env
	DestAddress
	RPCTimeout
	ReadTimeout
	RingHashKey
	DDPTag
	WithMeshHeader
	ConnectTimeout
	SpanContext
	ShortConnection
	FromMethod
	StressTag
	MsgType
	HTTPContentType
	RawRingHashKey
	LBType
)

// key of header transport
const (
	HeaderTransRemoteAddr     = "rip"
	HeaderTransToCluster      = "tc"
	HeaderTransToIDC          = "ti"
	HeaderTransPerfTConnStart = "pcs"
	HeaderTransPerfTConnEnd   = "pce"
	HeaderTransPerfTSendStart = "pss"
	HeaderTransPerfTRecvStart = "prs"
	HeaderTransPerfTRecvEnd   = "pre"
	// the connection peer will shutdown later,so it send back the header to tell client to close the connection.
	HeaderConnectionReadyToReset = "crrst"
	HeaderProcessAtTime          = "K_ProcessAtTime"
)
