package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitHTTP "github.com/mnabbasbaadi/distributedcache/foundation/http"
	"github.com/mnabbasbaadi/distributedcache/node/service/internal/node"
	"github.com/mnabbasbaadi/distributedcache/node/service/pkg/app"
	"golang.org/x/exp/slog"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("main: started")
	// Define flags
	masterPtr := flag.String("master", "localhost:7171", "address of the master node")

	// Parse the flags
	flag.Parse()

	if *masterPtr == "" {
		logger.Error("main: address and master address cannot be empty")
		os.Exit(1)
	}

	logger.Info("starting server", "master", *masterPtr)
	ctx := context.Background()

	// Start server
	serverErrors := make(chan error, 1)
	logger.Info("initializing server")
	httpServer := setupHTTPServer(5*time.Second, 5*time.Second, 5*time.Second, *logger)
	logger.Info("initializing environment")

	addr := "0.0.0.0:8080"

	params := app.Params{
		Logger:       logger,
		HTTPRegister: httpServer.Register,
		Node:         node.New(*masterPtr),
		MasterAddr:   *masterPtr,
	}

	env, err := app.NewEnvironment(ctx, params)
	if err != nil {
		logger.Error("error initializing environment", "err", err)
		os.Exit(1)
	}
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
