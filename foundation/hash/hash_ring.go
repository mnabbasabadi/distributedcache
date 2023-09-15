package hash

import (
	"hash/crc32"
	"sort"
	"sync"

	"github.com/mnabbasbaadi/distributedcache/foundation/cache"
)

type (
	hashRing []uint32

	Hasher interface {
		AddNode(addr string) bool
		GetNode(key string) Node
		DeleteNode(addr string) bool
	}

	consistentHash struct {
		Nodes     map[uint32]Node
		IsPresent map[string]bool
		Circle    hashRing
		sync.RWMutex
	}
)

func (hr hashRing) Len() int           { return len(hr) }
func (hr hashRing) Less(i, j int) bool { return hr[i] < hr[j] }
func (hr hashRing) Swap(i, j int)      { hr[i], hr[j] = hr[j], hr[i] }

func NewConsistentHash() Hasher {
	return &consistentHash{
		Nodes:     make(map[uint32]Node),
		IsPresent: make(map[string]bool),
		Circle:    hashRing{},
	}
}

func (ch *consistentHash) AddNode(addr string) bool {
	ch.Lock()
	defer ch.Unlock()
	if _, ok := ch.IsPresent[addr]; ok {
		return false
	}
	node := &node{
		IpAddr: addr,
		Cache:  cache.NewCache(),
	}

	ch.Nodes[ch.hashStr(addr)] = node
	ch.IsPresent[addr] = true
	ch.sortHashRing()
	return true
}

func (ch *consistentHash) GetNode(key string) Node {
	ch.RLock()
	defer ch.RUnlock()
	if len(ch.Circle) == 0 {
		return node{}
	}
	hash := ch.hashStr(key)
	i := sort.Search(len(ch.Circle), func(i int) bool { return ch.Circle[i] >= hash })
	if i == len(ch.Circle) {
		i = 0
	}
	return ch.Nodes[ch.Circle[i]]
}

func (ch *consistentHash) sortHashRing() {
	ch.Circle = hashRing{}
	for k := range ch.Nodes {
		ch.Circle = append(ch.Circle, k)
	}
	sort.Sort(ch.Circle)
}

func (ch *consistentHash) hashStr(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (ch *consistentHash) DeleteNode(addr string) bool {
	ch.Lock()
	defer ch.Unlock()
	if _, ok := ch.IsPresent[addr]; !ok {
		return false
	}
	delete(ch.IsPresent, addr)
	delete(ch.Nodes, ch.hashStr(addr))
	ch.sortHashRing()
	return true
}
