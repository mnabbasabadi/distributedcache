// Package integration is the integration test suite of the service.

//go:build integration
// +build integration

package integration

import (
	"context"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	cacheAPI "github.com/mnabbasbaadi/distributedcache/node/api/v1"
	"github.com/mnabbasbaadi/distributedcache/node/service/internal/node"
	"github.com/mnabbasbaadi/distributedcache/node/service/pkg/app"
	"github.com/mnabbasbaadi/distributedcache/node/service/tests/support/client"
	"github.com/mnabbasbaadi/distributedcache/node/service/tests/support/receiver"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slog"
)

const (
	host = "localhost"
	port = "8080"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	ctx := context.TODO()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	httpServer := setupHTTPServer(5*time.Second, 5*time.Second, 5*time.Second, *logger)
	defer httpServer.Stop()

	addr := net.JoinHostPort(host, port)
	n := node.New(addr)

	server := receiver.NewReceiver([]receiver.ResponseMessage{
		{
			Body: nil,
			Code: 204,
		},
	})
	url, closer := server.Serve()
	defer closer()

	params := app.Params{
		Logger:       logger,
		HTTPRegister: httpServer.Register,
		Node:         n,
		MasterAddr:   url,
	}

	env, err := app.NewEnvironment(ctx, params)
	require.NoError(t, err)
	defer env.Shutdown()

	httpServer.Start(addr, nil)

	testClient, err := client.NewGradingAPITestClient(addr, cacheAPI.WithHTTPClient(http.DefaultClient))
	require.NoError(t, err)

	suites := map[string]suite.TestingSuite{
		"E2E": &E2ETestSuite{
			client: testClient,
			node:   n,
			rec:    server,
		},
	}
	for _, s := range suites {
		suite.Run(t, s)
	}
}
