package netpollmux

import (
	"sync"

	"github.com/cloudwego/kitex/pkg/remote"
)

// EventHandler is used to handle events
type EventHandler interface {
	Recv(bufReader remote.ByteBuffer, err error) error
}

// A concurrent safe <seqID,EventHandler> map, which fork from code.byted.org/rpc/gauss
// To avoid lock bottlenecks this map is dived to several (SHARD_COUNT) map shards.
type sharedMap struct {
	size   int32
	shared []*shared
}

// A "thread" safe string to anything map.
type shared struct {
	msgs map[int32]EventHandler
	sync.Mutex
}

// Creates a new concurrent map.
func newSharedMap(size int32) *sharedMap {
	m := &sharedMap{
		size:   size,
		shared: make([]*shared, size),
	}
	for i := range m.shared {
		m.shared[i] = &shared{
			msgs: make(map[int32]EventHandler),
		}
	}
	return m
}

// getShard returns shard under given seq id
func (m *sharedMap) getShard(seqID int32) *shared {
	return m.shared[seqID%m.size]
}

// store stores msg under given seq id.
func (m *sharedMap) store(seqID int32, msg EventHandler) {
	if seqID == 0 {
		return
	}
	// Get map shard.
	shard := m.getShard(seqID)
	shard.Lock()
	shard.msgs[seqID] = msg
	shard.Unlock()
}

// load loads the msg under the seq id.
func (m *sharedMap) load(seqID int32) (msg EventHandler, ok bool) {
	if seqID == 0 {
		return nil, false
	}
	shard := m.getShard(seqID)
	shard.Lock()
	msg, ok = shard.msgs[seqID]
	shard.Unlock()
	return msg, ok
}

// delete deletes the msg under the given seq id.
func (m *sharedMap) delete(seqID int32) {
	if seqID == 0 {
		return
	}
	shard := m.getShard(seqID)
	shard.Lock()
	delete(shard.msgs, seqID)
	shard.Unlock()
}

// rangeMap iterates over the map.
func (m *sharedMap) rangeMap(fn func(seqID int32, msg EventHandler)) {
	for _, shard := range m.shared {
		shard.Lock()
		for k, v := range shard.msgs {
			fn(k, v)
		}
		shard.Unlock()
	}
}
