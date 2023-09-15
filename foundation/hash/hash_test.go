package hash

import (
	"testing"

	"github.com/mnabbasbaadi/distributedcache/foundation/cache"
	"github.com/stretchr/testify/require"
)

func TestConsistentHashAddNode(t *testing.T) {
	ch := &consistentHash{
		Nodes:     make(map[uint32]Node),
		IsPresent: make(map[string]bool),
		Circle:    hashRing{},
	}
	node1 := node{IpAddr: "127.0.0.1:8080", Cache: cache.NewCache()}
	node2 := node{IpAddr: "127.0.0.1:8081", Cache: cache.NewCache()}
	require.True(t, ch.AddNode(node1.IpAddr))
	require.True(t, ch.AddNode(node2.IpAddr))
	require.Equal(t, 2, len(ch.Nodes))
	require.True(t, ch.IsPresent[node1.IpAddr])
	require.True(t, ch.IsPresent[node2.IpAddr])
}

func TestConsistentHashAddNodeDuplicate(t *testing.T) {
	ch := &consistentHash{
		Nodes:     make(map[uint32]Node),
		IsPresent: make(map[string]bool),
		Circle:    hashRing{},
	}
	node1 := node{IpAddr: "127.0.0.1:8080", Cache: cache.NewCache()}
	node2 := node{IpAddr: "127.0.0.1:8081", Cache: cache.NewCache()}
	require.True(t, ch.AddNode(node1.IpAddr))
	require.True(t, ch.AddNode(node2.IpAddr))
	require.Equal(t, 2, len(ch.Nodes))
	require.False(t, ch.AddNode(node1.IpAddr))
	require.Equal(t, 2, len(ch.Nodes))
}

func TestConsistentHashGetNode(t *testing.T) {
	ch := &consistentHash{
		Nodes:     make(map[uint32]Node),
		IsPresent: make(map[string]bool),
		Circle:    hashRing{},
	}
	node1 := node{IpAddr: "127.0.0.1:8080", Cache: cache.NewCache()}
	node2 := node{IpAddr: "127.0.0.1:8081", Cache: cache.NewCache()}
	ch.AddNode(node1.IpAddr)
	ch.AddNode(node2.IpAddr)
	require.Equal(t, &node1, ch.GetNode("key2"))
	require.Equal(t, &node2, ch.GetNode("key1"))
	require.Equal(t, &node2, ch.GetNode("key4"))
	require.Equal(t, &node1, ch.GetNode("key3311"))

}

func TestConsistentHashGetNodeEmpty(t *testing.T) {
	ch := NewConsistentHash()
	require.Equal(t, node{}, ch.GetNode("key1"))
}
