package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddNode(t *testing.T) {
	nodes := []node{{address: "127.0.0.1"}, {address: "127.0.0.2"}}
	hr := hashRing{nodes: nodes}

	hr.AddNode("127.0.0.3")

	require.Equal(t, 3, len(hr.nodes))

	require.Equal(t, 768, len(hr.hashes))
}

func TestDeleteNode(t *testing.T) {
	nodes := []node{{address: "127.0.0.1"}, {address: "127.0.0.2"}, {address: "127.0.0.3"}}
	hr := hashRing{nodes: nodes}

	hr.DeleteNode("127.0.0.2")

	require.Equal(t, 2, len(hr.nodes))
	require.Equal(t, 512, len(hr.hashes))
}

func TestGetNode(t *testing.T) {
	nodes := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}
	hr := NewHashRing(nodes...)

	node := hr.GetNode("key0")
	require.Equal(t, "127.0.0.2", node)

	node = hr.GetNode("key1")
	require.Equal(t, "127.0.0.2", node)

	node = hr.GetNode("key2")
	require.Equal(t, "127.0.0.3", node)

	node = hr.GetNode("key3")
	require.Equal(t, "127.0.0.3", node)

	node = hr.GetNode("key4")
	require.Equal(t, "127.0.0.2", node)

	node = hr.GetNode("key5")
	require.Equal(t, "127.0.0.1", node)

	node = hr.GetNode("key6")
	require.Equal(t, "127.0.0.2", node)

	hr.AddNode("127.0.0.4")

	node = hr.GetNode("key0")
	require.Equal(t, "127.0.0.3", node)

	node = hr.GetNode("key1")
	require.Equal(t, "127.0.0.3", node)

	node = hr.GetNode("key2")
	require.Equal(t, "127.0.0.3", node)

	node = hr.GetNode("key3")
	require.Equal(t, "127.0.0.4", node)

	node = hr.GetNode("key4")
	require.Equal(t, "127.0.0.3", node)

	node = hr.GetNode("key5")
	require.Equal(t, "127.0.0.1", node)

	node = hr.GetNode("key6")
	require.Equal(t, "127.0.0.3", node)

	node = hr.GetNode("key95")
	require.Equal(t, "127.0.0.4", node)

}

func TestHashStr(t *testing.T) {

	testCases := map[string]struct {
		input    string
		expected uint32
	}{
		"empty string": {
			input:    "",
			expected: 3820012610,
		},
		"non-empty string": {
			input:    "hello world",
			expected: 3108841401,
		},
		"string with non-ASCII character": {
			input:    "こんにちは",
			expected: 307948255,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := hashStr(tc.input)
			require.Equal(t, tc.expected, actual)
		})
	}

}
