package registry

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mnabbasbaadi/distributedcache/foundation/hash"
	"golang.org/x/exp/slog"
)

var _ registerAPI.ServerInterface = new(server)
var (
	errKeyNotFound = errors.New("key not found")
)

type (
	server struct {
		logger *slog.Logger
		sh     hash.Hasher
	}
)

func (s server) RegisterNode(w http.ResponseWriter, r *http.Request) {
	var payload registerAPI.Node
	if err := rootHTTP.UnmarshalJSONFromBody(r, &payload); err != nil {
		s.logger.With(err).Error("while unmarshalling request body")
		s.respondError(w, err, http.StatusBadRequest)
		return
	}
	_ = s.sh.AddNode(payload.Address)
	s.respond(w, nil, http.StatusNoContent)
}

func (s server) UnRegisterNode(w http.ResponseWriter, r *http.Request, address registerAPI.Address) {
	b := s.sh.DeleteNode(address)
	if !b {
		s.respondError(w, errKeyNotFound, http.StatusNotFound)
		return
	}
	s.respond(w, nil, http.StatusNoContent)
}

func (s server) respond(w http.ResponseWriter, data any, statusCode int) {
	err := rootHTTP.Respond(w, data, statusCode)
	if err != nil {
		s.logger.With(err).Error("while responding")
		return
	}
}
func (s server) respondError(w http.ResponseWriter, err error, code int) {
	if err := rootHTTP.Respond(w, map[string]string{"error": err.Error()}, code); err != nil {
		s.logger.With(err).Error("error responding to request")

	}
}

// NewHandler returns a new http.Handler that implements the ServerInterface
func NewHandler(logger *slog.Logger, sh hash.Hasher) http.Handler {
	s := server{
		logger: logger,
		sh:     sh,
	}

	options := registerAPI.ChiServerOptions{
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			logger.With(err).Error("error")
			s.respondError(w, err, http.StatusBadRequest)
		},
	}

	return registerAPI.HandlerWithOptions(s, options)

}
