/*
Package cache provides an in-memory cache implementation.
The cache is divided into multiple shards, each with its own lock,
reducing lock contention and allowing more concurrent operations.

Hash Function: Implemented a hash function (fnv32) to distribute keys evenly across shards.

Cache Cleanup: Implemented a lazy expiration mechanism to remove expired items from the cache.
*/
package cache

import (
	"encoding/base64"
	"hash/fnv"
	"sync"
	"time"

	"golang.org/x/text/unicode/norm"
)

const shardCount = 32

type (
	Item struct {
		Value      []byte
		Expiration int64
	}

	CacheShard struct {
		mu    sync.RWMutex
		cache map[string]Item
	}

	InMemoryCache struct {
		shards [shardCount]*CacheShard
	}
)

func NewCache() *InMemoryCache {
	c := &InMemoryCache{}
	for i := 0; i < shardCount; i++ {
		c.shards[i] = &CacheShard{cache: make(map[string]Item)}
	}
	return c
}

func (c *InMemoryCache) getShard(key []byte) *CacheShard {
	return c.shards[uint(fnv32(key))%uint(shardCount)]
}

func (c *InMemoryCache) Set(key, value []byte, ttl time.Duration) {
	strKey := base64.StdEncoding.EncodeToString(key)

	key = normalizeKey(key)
	shard := c.getShard(key)

	shard.mu.Lock()
	shard.cache[strKey] = Item{
		Value:      value,
		Expiration: time.Now().Add(ttl).UnixNano(),
	}
	shard.mu.Unlock()
}

func (c *InMemoryCache) Get(key []byte) ([]byte, bool) {
	strKey := base64.StdEncoding.EncodeToString(key)

	key = normalizeKey(key)
	shard := c.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	item, found := shard.cache[strKey]
	if !found || item.Expiration < time.Now().UnixNano() {
		if found {
			delete(shard.cache, strKey) // Lazy expiration: remove the item if it is expired
		}
		return nil, false
	}
	return item.Value, true
}

func (c *InMemoryCache) Delete(key []byte) {
	strKey := base64.StdEncoding.EncodeToString(key)

	shard := c.getShard(key)
	shard.mu.Lock()
	delete(shard.cache, strKey)
	shard.mu.Unlock()
}

func fnv32(key []byte) uint32 {
	f := fnv.New32()
	_, _ = f.Write(key)
	return f.Sum32()
}

func normalizeKey(key []byte) []byte {
	return norm.NFC.Bytes(key)
}
