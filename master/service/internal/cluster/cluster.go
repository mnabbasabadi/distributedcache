// Package cluster provides a minimal leader election mechanism used by the
// master service to determine which node should act as the leader.
package cluster

import (
	"sort"
	"strings"
)

// Manager handles simple leader election between master nodes.
type Manager struct {
	self   string
	peers  []string
	leader string
}

// NewManager returns a Manager configured with the local node address and peer list.
func NewManager(self string, peers []string) *Manager {
	all := append([]string{self}, peers...)
	sort.Strings(all)
	c := &Manager{self: self, peers: peers, leader: all[0]}
	return c
}

// Leader returns the elected leader's address.
func (c *Manager) Leader() string { return c.leader }

// IsLeader reports whether this node is the leader.
func (c *Manager) IsLeader() bool { return c.self == c.leader }

// ParsePeers reads a comma separated list of addresses from the provided
// environment variable value. Empty strings are ignored.
func ParsePeers(env string) []string {
	if env == "" {
		return nil
	}
	parts := strings.Split(env, ",")
	var peers []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			peers = append(peers, trimmed)
		}
	}
	return peers
}
