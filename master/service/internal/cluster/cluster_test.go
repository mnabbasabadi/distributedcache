package cluster

import "testing"

func TestLeaderElection(t *testing.T) {
	m := NewManager("a", []string{"b", "c"})
	if !m.IsLeader() {
		t.Fatalf("expected a to be leader")
	}
	if m.Leader() != "a" {
		t.Fatalf("got leader %s", m.Leader())
	}
}

func TestParsePeers(t *testing.T) {
	peers := ParsePeers("b,c , d")
	if len(peers) != 3 {
		t.Fatalf("expected 3 peers got %d", len(peers))
	}
	if peers[1] != "c" || peers[2] != "d" {
		t.Fatalf("unexpected peers %v", peers)
	}
}
