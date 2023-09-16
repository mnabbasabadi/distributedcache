package node

import (
	"time"

	"github.com/mnabbasbaadi/distributedcache/node/service/internal/cache"
)

var defaultExpiration = 30 * time.Minute

type (
	node struct {
		IpAddr string
		Cache  *cache.Cache
	}
	Node interface {
		Get(string) (string, bool)
		Set(string, string)
	}
)

func (n node) Get(key string) (string, bool) {
	return n.Cache.Get(key)
}
func (n node) Set(key, value string) {
	n.Cache.Set(key, value, defaultExpiration)
}

func New(ipAddr string) Node {
	return node{
		IpAddr: ipAddr,
		Cache:  cache.NewCache(),
	}
}
