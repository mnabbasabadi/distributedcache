package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCacheSet(t *testing.T) {
	c := NewCache()
	c.Set("key", "value", 5*time.Minute)
	val, ok := c.Get("key")
	require.Equal(t, "value", val)
	require.True(t, ok)

}

func TestCacheGet(t *testing.T) {
	c := NewCache()
	c.Set("key", "value", 5*time.Minute)
	val, ok := c.Get("key")
	require.Equal(t, "value", val)
	require.True(t, ok)
	val, ok = c.Get("nonexistent")
	require.Equal(t, "", val)
	require.False(t, ok)
}

func TestCacheDelete(t *testing.T) {
	c := NewCache()
	c.Set("key", "value", 5*time.Minute)
	c.Delete("key")
	val, ok := c.Get("key")
	require.Equal(t, "", val)
	require.False(t, ok)
}

func TestCacheCleanup(t *testing.T) {
	c := NewCache()
	c.Set("key1", "value1", 2*time.Second)
	c.Set("key2", "value2", 1*time.Second)

	time.Sleep(3 * time.Second)
	val, ok := c.Get("key1")
	require.Equal(t, "", val)
	require.False(t, ok)

	val, ok = c.Get("key2")
	require.Equal(t, "", val)
	require.False(t, ok)
}

func TestCacheStartCleanupTimer(t *testing.T) {
	c := NewCache()
	//c.StartCleanupTimer(1 * time.Second)
	c.Set("key1", "value1", 5*time.Second)
	c.Set("key2", "value2", 1*time.Second)
	time.Sleep(2 * time.Second)
	val, ok := c.Get("key1")
	require.Equal(t, "value1", val)
	require.True(t, ok)
	val, ok = c.Get("key2")
	require.Equal(t, "", val)
	require.False(t, ok)
}
