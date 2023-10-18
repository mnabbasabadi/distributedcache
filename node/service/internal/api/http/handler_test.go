package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	cacheAPI "github.com/mnabbasbaadi/distributedcache/node/api/v1"
	"github.com/mnabbasbaadi/distributedcache/node/service/internal/node"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

func TestServer_GetValue(t *testing.T) {

	// Register nodes (this would be dynamic in a real-world scenario)
	n := node.New("localhost:8081")
	// Create a new server instance with the tests cache
	s := server{node: n}

	// Define tests cases
	testCases := map[string]struct {
		key      cacheAPI.Key
		prep     func()
		expected string
		status   int
	}{
		"Valid key": {
			key:      "tests-key",
			expected: "tests-value",
			prep: func() {
				n.Set([]byte("tests-key"), []byte("tests-value"))
			},
			status: http.StatusOK,
		},
		"Invalid key": {
			key:      "invalid-key",
			expected: "",
			prep:     func() {},
			status:   http.StatusNotFound,
		},
	}

	// Run tests cases
	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {

			// Prepare the tests case
			tc.prep()

			// Create a new HTTP request with the tests key
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// Create a new HTTP response recorder
			rec := httptest.NewRecorder()
			// Call the GetValue method with the tests request and response recorder
			s.GetValue(rec, req, tc.key)

			// Check the response status code
			require.Equal(t, tc.status, rec.Code)

			if tc.status != http.StatusOK {
				return
			}
			var response cacheAPI.KeyValue
			err := json.NewDecoder(rec.Body).Decode(&response)
			require.NoError(t, err)

			// Check the response body
			require.Equal(t, tc.expected, response.Value)
			require.Equal(t, tc.key, response.Key)
		})
	}
}

func TestServer_AddKey(t *testing.T) {

	// Register nodes (this would be dynamic in a real-world scenario)
	n := node.New("127.0.0.1:5000")
	// Create a new server instance with the tests cache
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	s := server{logger: logger, node: n}

	// Define tests cases
	testCases := map[string]struct {
		name    string
		reqBody []byte
		status  int
	}{
		"Valid key": {

			reqBody: []byte(`{"key":"tests-key","value":"tests-value"}`),
			status:  http.StatusNoContent,
		},
		"Invalid key": {
			reqBody: []byte(`wrong json`),
			status:  http.StatusBadRequest,
		},
	}

	// Run tests cases
	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(tc.reqBody))
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			s.AddKey(rec, req)

			require.Equal(t, tc.status, rec.Code)

		})
	}
}
