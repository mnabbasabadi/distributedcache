package hash

import (
	"crypto/sha256"
	"sort"
	"strconv"
)

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

func (hr *hashRing) rehash() {
	hr.hashes = []uint32{}
	for i, node := range hr.nodes {
		for j := 0; j < 100; j++ {
			hr.hashes = append(hr.hashes, hashStr(node.address+strconv.Itoa(i)+strconv.Itoa(j)))
		}
	}
	sort.Slice(hr.hashes, func(i, j int) bool { return hr.hashes[i] < hr.hashes[j] })
}

func (hr *hashRing) GetNode(key string) string {
	hash := hashStr(key)
	idx := sort.Search(len(hr.hashes), func(i int) bool { return hr.hashes[i] >= hash })
	if idx == len(hr.hashes) {
		idx = 0
	}
	return hr.nodes[idx/100].address
}

func hashStr(key string) uint32 {
	hash := sha256.New()
	hash.Write([]byte(key))
	hashBytes := hash.Sum(nil)
	return uint32(hashBytes[0])<<24 | uint32(hashBytes[1])<<16 | uint32(hashBytes[2])<<8 | uint32(hashBytes[3])
}
