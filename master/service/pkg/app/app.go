// Package app is the entry point of the service binary.
package app

import (
	"context"
	"net/http"

	kitHTTP "github.com/mnabbasbaadi/distributedcache/foundation/http"
	registryAPI "github.com/mnabbasbaadi/distributedcache/master/service/internal/api/http"
	"github.com/mnabbasbaadi/distributedcache/master/service/internal/hash"

	"golang.org/x/exp/slog"
)

// Environment ...
type Environment struct {
	logger *slog.Logger

	HTTPRegister kitHTTP.Registrar

	hasher hash.Hasher
}

// Params ...
type Params struct {
	//metrics *Metrics
	Logger       *slog.Logger
	HTTPRegister kitHTTP.Registrar
	Hasher       hash.Hasher
}

// NewEnvironment ...
func NewEnvironment(ctx context.Context, params Params) *Environment {
	e := &Environment{
		logger:       params.Logger,
		HTTPRegister: params.HTTPRegister,
		hasher:       params.Hasher,
	}
	e.Setup(ctx)
	return e
}

// Setup ...
func (e *Environment) Setup(_ context.Context) {

	handler := registryAPI.NewHandler(e.logger, e.hasher)

	e.HTTPRegister(func(mux *http.ServeMux) {
		h := kitHTTP.Chain(handler)
		mux.Handle("/", h)
	})

}

// Shutdown ...
func (e *Environment) Shutdown() {
	//e.metrics.Close()
}
