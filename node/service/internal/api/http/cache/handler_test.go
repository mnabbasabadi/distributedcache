package cache

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServer_GetValue(t *testing.T) {

	// Register nodes (this would be dynamic in a real-world scenario)
	node := hash.NewNode("localhost:8081")
	// Create a new server instance with the test cache
	s := server{node: node}

	// Define test cases
	testCases := map[string]struct {
		key      cacheAPI.Key
		prep     func()
		expected string
		status   int
	}{
		"Valid key": {
			key:      "test-key",
			expected: "test-value",
			prep: func() {
				node.Set("test-key", "test-value")
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

	// Run test cases
	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {

			// Prepare the test case
			tc.prep()

			// Create a new HTTP request with the test key
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// Create a new HTTP response recorder
			rec := httptest.NewRecorder()
			// Call the GetValue method with the test request and response recorder
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
	node := hash.NewNode("127.0.0.1:5000")
	// Create a new server instance with the test cache
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	s := server{logger: logger, node: node}

	// Define test cases
	testCases := map[string]struct {
		name    string
		reqBody []byte
		status  int
	}{
		"Valid key": {

			reqBody: []byte(`{"key":"test-key","value":"test-value"}`),
			status:  http.StatusNoContent,
		},
		"Invalid key": {
			reqBody: []byte(`wrong json`),
			status:  http.StatusBadRequest,
		},
	}

	// Run test cases
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
