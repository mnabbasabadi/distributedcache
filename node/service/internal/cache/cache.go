package cache

import (
	"sync"
	"time"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
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
	key = normalizeKey(key)
	c.cache[key] = Item{
		Value:      value,
		Expiration: time.Now().Add(ttl).UnixNano(),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	if !isValidUTF8(key) {
		return "", false
	}

	key = normalizeKey(key)

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
			<-ticker.C
			c.cleanup()

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

func isValidUTF8(s string) bool {
	return utf8.ValidString(s)
}

func normalizeKey(key string) string {
	return norm.NFC.String(key)
}
