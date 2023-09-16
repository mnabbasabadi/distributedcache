package receiver

import (
	"net"
	"net/http"
	"net/http/httptest"
)

// CreateTestLocalServer ...
func CreateTestLocalServer(handler http.Handler) *httptest.Server {
	server := httptest.NewUnstartedServer(handler)
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	server.Listener = l
	server.Start()
	return server
}
