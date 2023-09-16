// Package http provides utilities for working with HTTP packages.
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/slog"

	"github.com/rs/cors"
)

type (
	// serverOptions contains various configuration options for the server.
	serverOptions struct {
		readTimeout     time.Duration
		writeTimeout    time.Duration
		shutdownTimeout time.Duration
		cors            cors.Options
	}

	// Server is a utility wrapper for the default http.Server.
	Server struct {
		logger     slog.Logger
		server     *http.Server
		mux        *http.ServeMux
		address    string
		options    *serverOptions
		middleware []Middleware
	}

	// ServerOption is used to provide overrides to the Server implementation.
	ServerOption func(options *serverOptions)

	// Registrar is a function that registers handlers using the underlying http.ServeMux.
	Registrar func(register func(mux *http.ServeMux))

	// Middleware is a function that runs code before and/or after another Handler.
	// It is primarily used for applying cross-cutting concerns.
	Middleware func(next http.Handler) http.Handler
)

// defaultServerOptions contains the default options for the server.
var defaultServerOptions = []ServerOption{
	ShutdownTimeout(30 * time.Second),
}

// ShutdownTimeout sets the time to wait for the HTTP server to stop. The default is 30 seconds.
func ShutdownTimeout(timeout time.Duration) ServerOption {
	return func(options *serverOptions) {
		options.shutdownTimeout = timeout
	}
}

// ReadTimeout sets the maximum duration for reading the entire request, including the body.
func ReadTimeout(timeout time.Duration) ServerOption {
	return func(options *serverOptions) {
		options.readTimeout = timeout
	}
}

// WriteTimeout sets the maximum duration before timing out writes of the response.
func WriteTimeout(timeout time.Duration) ServerOption {
	return func(options *serverOptions) {
		options.writeTimeout = timeout
	}
}

// NewServer creates a new instance of the HTTP server wrapper.
func NewServer(logger slog.Logger, options ...ServerOption) Server {
	serverOptions := &serverOptions{}
	// Apply defaults
	for _, defaultOpt := range defaultServerOptions {
		defaultOpt(serverOptions)
	}

	// Apply overrides if any
	for _, opt := range options {
		opt(serverOptions)
	}

	return Server{
		logger:     logger,
		server:     nil,
		mux:        http.NewServeMux(),
		address:    "",
		options:    serverOptions,
		middleware: make([]Middleware, 0),
	}
}

// UseMiddleware adds a middleware to be executed for every request.
func (s *Server) UseMiddleware(mw Middleware) {
	s.middleware = append(s.middleware, mw)
}

// Register allows handlers to be registered with the underlying mux.
func (s *Server) Register(registrar func(mux *http.ServeMux)) {
	registrar(s.mux)
}

// Start initializes a new http.Server and starts the server at the provided address.
func (s *Server) Start(address string, errors chan<- error) {
	s.address = address
	rootHandler := Chain(s.mux)

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      rootHandler,
		ReadTimeout:  s.options.readTimeout,
		WriteTimeout: s.options.writeTimeout,
	}

	go func() {
		s.logger.Info("Starting REST server at", "addr", s.server.Addr)

		errors <- s.server.ListenAndServe()
	}()

}

// Address returns the active address of the HTTP server.
func (s *Server) Address() string {
	return s.address
}

// Stop shuts down and cleans up the HTTP server.
func (s *Server) Stop() {
	if s.server == nil {
		return
	}

	// Create a context for the Stop call.
	ctx, cancel := context.WithTimeout(context.Background(), s.options.shutdownTimeout)
	defer cancel()

	// Ask the listener to shutdown.
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Info(fmt.Sprintf("Graceful shutdown of the REST server did not complete in %s", s.options.shutdownTimeout), "err", err)
		err = s.server.Close()
		if err != nil {
			s.logger.Warn("Error closing REST server connectio", "err", err)
		}
		return
	}
	s.logger.Info("REST server completed graceful shutdown.")
}

// Chain chains the middleware and returns a single handler.
func Chain(base http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) == 0 {
		return base
	}

	for i := range middlewares {
		mw := middlewares[len(middlewares)-i-1]
		if mw == nil {
			continue
		}
		base = mw(base)
	}

	return base
}
