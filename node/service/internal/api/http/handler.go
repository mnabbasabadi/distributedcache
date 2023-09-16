package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	HTTPKit "github.com/mnabbasbaadi/distributedcache/foundation/http"
	cacheAPI "github.com/mnabbasbaadi/distributedcache/node/api/v1"
	"github.com/mnabbasbaadi/distributedcache/node/service/internal/node"
	"golang.org/x/exp/slog"
)

var _ cacheAPI.ServerInterface = new(server)

var (
	errKeyNotFound = errors.New("key not found")
)

type (
	server struct {
		logger *slog.Logger
		node   node.Node
	}
)

// NewHandler returns a new http.Handler that implements the ServerInterface
func NewHandler(logger *slog.Logger, node node.Node) http.Handler {
	s := server{
		logger: logger,
		node:   node,
	}

	options := cacheAPI.ChiServerOptions{
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.With(err).Error("error")
			s.respondError(w, err, http.StatusBadRequest)
		},
	}

	return cacheAPI.HandlerWithOptions(s, options)

}
func (s server) GetValue(w http.ResponseWriter, _ *http.Request, key cacheAPI.Key) {
	get, b := s.node.Get(key)
	if !b {
		s.respondError(w, errKeyNotFound, http.StatusNotFound)
		return
	}

	s.respond(w, cacheAPI.KeyValue{
		Key:   key,
		Value: get,
	}, http.StatusOK)
}

func (s server) AddKey(w http.ResponseWriter, r *http.Request) {
	var keyValue cacheAPI.KeyValue
	if err := HTTPKit.UnmarshalJSONFromBody(r, &keyValue); err != nil {
		s.logger.With(err).Error("while unmarshalling request body")
		s.respondError(w, err, http.StatusBadRequest)
		return
	}
	s.node.Set(keyValue.Key, keyValue.Value)
	s.respond(w, nil, http.StatusNoContent)
}

func (s server) respond(w http.ResponseWriter, data any, statusCode int) {
	err := HTTPKit.Respond(w, data, statusCode)
	if err != nil {
		s.logger.With(err).Error("while responding")
		return
	}
}
func (s server) respondError(w http.ResponseWriter, err error, code int) {
	if err := HTTPKit.Respond(w, map[string]string{"error": err.Error()}, code); err != nil {
		s.logger.With(err).Error("error responding to request")

	}
}
