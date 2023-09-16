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

	cacheAPI "github.com/mnabbasbaadi/distributedcache/master/api/v1"
	"github.com/mnabbasbaadi/distributedcache/master/service/internal/hash"
	"github.com/mnabbasbaadi/distributedcache/master/service/pkg/app"
	"github.com/mnabbasbaadi/distributedcache/master/service/tests/support/client"

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
	params := app.Params{
		Logger:       logger,
		HTTPRegister: httpServer.Register,
		Hasher:       hash.NewHashRing(),
	}

	env := app.NewEnvironment(ctx, params)
	defer env.Shutdown()

	httpServer.Start(addr, nil)

	testClient, err := client.NewGradingAPITestClient(addr, cacheAPI.WithHTTPClient(http.DefaultClient))
	require.NoError(t, err)

	suites := map[string]suite.TestingSuite{
		"E2E": &E2ETestSuite{
			client: testClient,
		},
	}
	for _, s := range suites {
		suite.Run(t, s)
	}
}
