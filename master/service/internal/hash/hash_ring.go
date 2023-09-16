/*

a hash ring data structure that can be used to distribute keys evenly across a set of nodes.
The hash ring is implemented using a circular array of hash values,
where each node is assigned a range of hash values based on its index in the array.
When a key is added to the hash ring, it is hashed and the corresponding node is determined based on its hash value.
The hash ring is rehashed whenever a node is added or removed to ensure that the distribution of keys is evenly across all nodes.


*/

package hash

import (
	"crypto/sha256"
	"fmt"
	"sort"
)

const hashSpace = 1 << 8

type (
	Hasher interface {
		AddNode(addr string)
		GetNode(key string) string
		DeleteNode(addr string)
	}
	hashRing struct {
		nodes  []node
		hashes []uint32
	}
	node struct {
		address string
	}
)

func NewHashRing(addrs ...string) Hasher {
	nodes := make([]node, len(addrs))
	for i, addr := range addrs {
		nodes[i] = node{address: addr}
	}
	hr := &hashRing{nodes: nodes}
	hr.rehash()
	return hr
}

func (hr *hashRing) AddNode(addr string) {
	hr.nodes = append(hr.nodes, node{address: addr})
	hr.rehash()
}

func (hr *hashRing) DeleteNode(address string) {
	for i, node := range hr.nodes {
		if node.address == address {
			hr.nodes = append(hr.nodes[:i], hr.nodes[i+1:]...)
			break
		}
	}
	hr.rehash()
}

// rehash pre-allocates the hashes array and fills it with the hash values of all nodes.
func (hr *hashRing) rehash() {
	hr.hashes = []uint32{}
	for i, node := range hr.nodes {
		for j := 0; j < hashSpace; j++ {
			hr.hashes = append(hr.hashes, hashStr(fmt.Sprintf("%s-%d", node.address, i*hashSpace+j)))
		}
	}
	sort.Slice(hr.hashes, func(i, j int) bool { return hr.hashes[i] < hr.hashes[j] })
}

// GetNode returns the address of the node that should handle the given key.
func (hr *hashRing) GetNode(key string) string {
	hash := hashStr(key)
	idx := sort.Search(len(hr.hashes), func(i int) bool { return hr.hashes[i] >= hash })
	if idx == len(hr.hashes) {
		idx = 0
	}
	return hr.nodes[idx/hashSpace].address
}

func hashStr(key string) uint32 {
	hash := sha256.New()
	hash.Write([]byte(key))
	hashBytes := hash.Sum(nil)
	// convert first 4 bytes to uint32 it helps to reduce the number of collisions
	// and distribute keys more evenly across nodes
	return uint32(hashBytes[0])<<24 | uint32(hashBytes[1])<<16 | uint32(hashBytes[2])<<8 | uint32(hashBytes[3])
}
