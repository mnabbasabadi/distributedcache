package node

import (
	"time"

	"github.com/mnabbasbaadi/distributedcache/node/service/internal/cache"
)

var defaultExpiration = 30 * time.Minute

type (
	node struct {
		IpAddr string
		Cache  *cache.InMemoryCache
	}
	Node interface {
		Get([]byte) ([]byte, bool)
		Set([]byte, []byte)
	}
)

func (n node) Get(key []byte) ([]byte, bool) {
	return n.Cache.Get(key)
}
func (n node) Set(key, value []byte) {
	n.Cache.Set(key, value, defaultExpiration)
}

func New(ipAddr string) Node {
	return node{
		IpAddr: ipAddr,
		Cache:  cache.NewCache(),
	}
}
