package cache

import (
	"sync"
	"time"
)

type (
	Item struct {
		Value      string
		Expiration int64
	}

	Cache struct {
		mu    sync.RWMutex
		cache map[string]Item
	}
)

func NewCache() *Cache {
	c := &Cache{
		cache: make(map[string]Item),
	}
	go c.startCleanupTimer(5 * time.Minute)
	return c
}

func (c *Cache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = Item{
		Value:      value,
		Expiration: time.Now().Add(ttl).UnixNano(),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.cache[key]
	if !found || item.Expiration < time.Now().UnixNano() {
		return "", false
	}
	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

func (c *Cache) startCleanupTimer(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.cleanup()
			}
		}
	}()
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, item := range c.cache {
		if item.Expiration < time.Now().UnixNano() {
			delete(c.cache, key)
		}
	}
}
