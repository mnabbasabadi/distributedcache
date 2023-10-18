package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCacheSet(t *testing.T) {
	c := NewCache()
	c.Set([]byte("key"), []byte("value"), 5*time.Minute)
	val, ok := c.Get([]byte("key"))
	require.EqualValues(t, "value", val)
	require.True(t, ok)

}

func TestCacheGet(t *testing.T) {
	c := NewCache()
	c.Set([]byte("key"), []byte("value"), 5*time.Minute)
	val, ok := c.Get([]byte("key"))
	require.EqualValues(t, "value", val)
	require.True(t, ok)
	val, ok = c.Get([]byte("nonexistent"))
	require.Nil(t, val)
	require.False(t, ok)
}

func TestCacheDelete(t *testing.T) {
	c := NewCache()
	c.Set([]byte("key"), []byte("value"), 5*time.Minute)
	c.Delete([]byte("key"))
	val, ok := c.Get([]byte("key"))
	require.Nil(t, val)
	require.False(t, ok)
}

func TestCacheCleanup(t *testing.T) {
	c := NewCache()
	c.Set([]byte("key1"), []byte("value1"), 2*time.Second)
	c.Set([]byte("key2"), []byte("value2"), 1*time.Second)

	time.Sleep(3 * time.Second)
	val, ok := c.Get([]byte("key1"))
	require.Nil(t, val)
	require.False(t, ok)

	val, ok = c.Get([]byte("key2"))
	require.Nil(t, val)
	require.False(t, ok)
}

func TestCacheStartCleanup(t *testing.T) {
	c := NewCache()
	//c.StartCleanupTimer(1 * time.Second)
	c.Set([]byte("key1"), []byte("value1"), 5*time.Second)
	c.Set([]byte("key2"), []byte("value2"), 1*time.Second)
	time.Sleep(2 * time.Second)
	val, ok := c.Get([]byte("key1"))
	require.EqualValues(t, "value1", val)
	require.True(t, ok)
	val, ok = c.Get([]byte("key2"))
	require.Nil(t, val)
	require.False(t, ok)
}

func TestCache_normalizeKey(t *testing.T) {
	require.EqualValues(t, "hello", normalizeKey([]byte("hello")))
	require.EqualValues(t, "café", normalizeKey([]byte("café")))
	require.EqualValues(t, "café", normalizeKey([]byte("cafe\u0301")))
}

func TestCache_fnv32(t *testing.T) {
	require.Equal(t, uint32(3069866343), fnv32([]byte("hello")))
	require.Equal(t, uint32(2609808943), fnv32([]byte("world")))
}
