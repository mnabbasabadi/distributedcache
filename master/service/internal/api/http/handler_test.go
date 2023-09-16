package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mnabbasbaadi/distributedcache/master/service/internal/hash"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

func TestServer_RegisterNode(t *testing.T) {
	// Create a new server instance with a mock logger and shard handler
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sh := hash.NewHashRing("localhost:8080", "localhost:8081", "localhost:8082")
	s := server{
		logger: logger,
		sh:     sh,
	}

	// Define tests cases
	testCases := map[string]struct {
		hasError bool
		status   int
	}{
		"Valid payload": {
			status: http.StatusNoContent,
		},
	}

	// Run tests cases
	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			// Create a new HTTP request with the tests payload
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set("Content-Type", "application/json")

			// Create a new HTTP response recorder
			rec := httptest.NewRecorder()

			// Call the RegisterNode method with the tests request and response recorder
			s.RegisterNode(rec, req)

			// Check the response status code
			require.Equal(t, tc.status, rec.Code)

		})
	}
}

func TestServer_UnRegisterNode(t *testing.T) {
	// Create a new server instance with a mock logger and shard handler
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sh := hash.NewHashRing("localhost:8080", "localhost:8081", "localhost:8082")
	s := server{
		logger: logger,
		sh:     sh,
	}

	// Define tests cases
	testCases := map[string]struct {
		prep   func()
		addr   string
		status int
	}{
		"Valid payload": {
			prep: func() {
				sh.AddNode("http://localhost:8080")
			},
			addr:   "http://localhost:8080",
			status: http.StatusNoContent,
		},
		"Invalid payload": {
			addr:   "http://localhost:8081",
			prep:   func() {},
			status: http.StatusNoContent,
		},
	}

	// Run tests cases
	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			// Create a new HTTP request with the tests payload
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set("Content-Type", "application/json")

			// Create a new HTTP response recorder
			rec := httptest.NewRecorder()
			tc.prep()
			// Call the RegisterNode method with the tests request and response recorder
			s.UnRegisterNode(rec, req, "http://localhost:8080")

			// Check the response status code
			require.Equal(t, tc.status, rec.Code)

		})
	}
}
