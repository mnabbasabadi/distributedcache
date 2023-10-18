package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitHTTP "github.com/mnabbasbaadi/distributedcache/foundation/http"
	"github.com/mnabbasbaadi/distributedcache/master/service/internal/hash"
	"github.com/mnabbasbaadi/distributedcache/master/service/pkg/app"
	"golang.org/x/exp/slog"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Start server
	serverErrors := make(chan error, 1)

	addr := "0.0.0.0:8080"

	httpServer := setupHTTPServer(5*time.Second, 5*time.Second, 5*time.Second, *logger)

	h := hash.NewHashRing()

	params := app.Params{
		Logger:       logger,
		HTTPRegister: httpServer.Register,
		Hasher:       h,
	}

	env := app.NewEnvironment(ctx, params)

	logger.Info("setup complete, starting server")
	startHTTPServer(httpServer, addr, serverErrors)
	shutdown := listenForShutdown()
	select {
	case err := <-serverErrors:
		logger.Error("error starting server", "err", err)
		os.Exit(1)
	case sig := <-shutdown:
		logger.Error("caught signal", "signal", sig)
		env.Shutdown()
		httpServer.Stop()
	}
}

// ListenForShutdown creates a channel and subscribes to specific signals to trigger a shutdown of the service.
func listenForShutdown() chan os.Signal {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	return shutdown
}

func setupHTTPServer(shutdownTimeout, readTimeout, writeTimeout time.Duration, logger slog.Logger) kitHTTP.Server {
	serverOptions := []kitHTTP.ServerOption{
		kitHTTP.ShutdownTimeout(shutdownTimeout),
		kitHTTP.ReadTimeout(readTimeout),
		kitHTTP.WriteTimeout(writeTimeout),
	}
	return kitHTTP.NewServer(logger, serverOptions...)
}

func startHTTPServer(httpServer kitHTTP.Server, addr string, serverErrors chan<- error) {
	httpServer.Start(addr, serverErrors)
}
