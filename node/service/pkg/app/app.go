// Package app is the entry point of the service binary.
package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	cacheAPI "github.com/mnabbasbaadi/distributedcache/node/service/internal/api/http"

	kitHTTP "github.com/mnabbasbaadi/distributedcache/foundation/http"
	registerAPI "github.com/mnabbasbaadi/distributedcache/master/api/v1"
	"github.com/mnabbasbaadi/distributedcache/node/service/internal/node"

	"golang.org/x/exp/slog"
)

// Environment ...
type Environment struct {
	logger *slog.Logger

	HTTPRegister kitHTTP.Registrar
	node         node.Node

	masterAddr string
}

// Params ...
type Params struct {
	//metrics *Metrics
	Logger       *slog.Logger
	HTTPRegister kitHTTP.Registrar
	Node         node.Node

	MasterAddr string
}

// NewEnvironment ...
func NewEnvironment(ctx context.Context, params Params) (*Environment, error) {
	e := &Environment{
		logger:       params.Logger,
		HTTPRegister: params.HTTPRegister,
		node:         params.Node,
		masterAddr:   params.MasterAddr,
	}
	err := e.Setup(ctx)
	return e, err
}

// Setup ...
func (e *Environment) Setup(_ context.Context) error {

	cacheHandler := cacheAPI.NewHandler(e.logger, e.node)
	e.logger.Info("registering node", "addr", e.masterAddr)
	client, err := registerAPI.NewClient(e.masterAddr)
	if err != nil {
		e.logger.With(err).Error("while creating client")
		return err
	}

	operation := func() error {
		rsp, err := client.RegisterNode(context.Background())
		if err != nil {
			e.logger.Error("while registering node", "err", err)
			return err
		}
		if rsp.StatusCode != http.StatusNoContent {
			e.logger.Error("while registering node", "status", rsp.StatusCode)
			return errors.New("status code not 204")
		}
		return nil
	}

	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), 1000))
	if err != nil {
		e.logger.Error("while registering node")
		return err
	}

	e.HTTPRegister(func(mux *http.ServeMux) {
		var mw []kitHTTP.Middleware
		h := kitHTTP.Chain(cacheHandler, mw...)
		mux.Handle("/", h)
	})
	return nil
}

// Shutdown ...
func (e *Environment) Shutdown() {
	client, err := registerAPI.NewClient(e.masterAddr)
	if err != nil {
		e.logger.With(err).Error("while creating client")
		return
	}
	_, err = client.UnRegisterNode(context.Background(), e.masterAddr)
	if err != nil {
		e.logger.With(err).Error("while unregistering node")
		return
	}

}
