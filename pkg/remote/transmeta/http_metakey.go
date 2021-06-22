package transmeta

// key of http egress meta
const (
	HTTPDestService     = "destination-service"
	HTTPDestCluster     = "destination-cluster"
	HTTPDestIDC         = "destination-idc"
	HTTPDestEnv         = "destination-env"
	HTTPDestAddr        = "destination-addr"
	HTTPDestConnTimeout = "destination-connect-timeout"
	HTTPDestReqTimeout  = "destination-request-timeout"
	HTTPSourceService   = "source-service"
	HTTPSourceCluster   = "source-cluster"
	HTTPSourceIDC       = "source-idc"
	HTTPSourceEnv       = "source-env"
	HTTPRequestHash     = "http-request-hash"
	HTTPLogID           = "x-tt-logid"
	HTTPStress          = "x-tt-stress"
	HTTPFilterType      = "bytemesh-filter-type"
	// supply metadata
	HTTPDestMethod      = "destination-method"
	HTTPSourceMethod    = "source-method"
	HTTPSpanTag         = "span-tag"
	HTTPRawRingHashKey  = "raw-ring-hash-key"
	HTTPLoadBalanceType = "lb-type"
	HTTPMeshFilterType  = "bytemesh-filter-type"
)

// value of http egress meta
const (
	RPCFilter = "rpc"
)
