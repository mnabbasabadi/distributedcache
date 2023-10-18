//go:build integration
// +build integration

package integration

import (
	"time"

	kitHTTP "github.com/mnabbasbaadi/distributedcache/foundation/http"
	"golang.org/x/exp/slog"
)

func setupHTTPServer(shutdownTimeout, readTimeout, writeTimeout time.Duration, logger slog.Logger) kitHTTP.Server {
	serverOptions := []kitHTTP.ServerOption{
		kitHTTP.ShutdownTimeout(shutdownTimeout),
		kitHTTP.ReadTimeout(readTimeout),
		kitHTTP.WriteTimeout(writeTimeout),
	}
	return kitHTTP.NewServer(logger, serverOptions...)
}
