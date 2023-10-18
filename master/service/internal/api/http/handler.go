package http

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpFramework "github.com/mnabbasbaadi/distributedcache/foundation/http"
	registerAPI "github.com/mnabbasbaadi/distributedcache/master/api/v1"
	"github.com/mnabbasbaadi/distributedcache/master/service/internal/hash"
	cacheAPI "github.com/mnabbasbaadi/distributedcache/node/api/v1"
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

func (s server) AddKey(w http.ResponseWriter, r *http.Request) {

	var keyValue cacheAPI.KeyValue
	if err := httpFramework.UnmarshalJSONFromBody(r, &keyValue); err != nil {
		s.logger.With(err).Error("while unmarshalling request body")
		s.respondError(w, err, http.StatusBadRequest)
		return
	}
	if keyValue.Key == "" {
		s.respondError(w, errKeyNotFound, http.StatusBadRequest)
		return
	}
	s.logger.Debug("adding key", "key", keyValue.Key, "value", keyValue.Value)

	url := s.sh.GetNode(keyValue.Key)

	s.logger.Debug("url", "url", url)
	if url == "" {
		s.respondError(w, errKeyNotFound, http.StatusNotFound)
		return
	}
	// TODO: use client from cache
	client, err := cacheAPI.NewClientWithResponses(fmt.Sprintf("http://%s", url))
	if err != nil {
		s.logger.With(err).Error("while creating client")
		s.respondError(w, err, http.StatusInternalServerError)
		return
	}
	resp, err := client.AddKeyWithResponse(r.Context(), keyValue)
	if err != nil {
		s.logger.With(err).Error("while adding key")
		s.respondError(w, err, http.StatusInternalServerError)
		return
	}

	s.respond(w, resp.Body, http.StatusOK)
}

func (s server) GetValue(w http.ResponseWriter, r *http.Request, key registerAPI.Key) {
	url := s.sh.GetNode(key)
	if url == "" {
		s.respondError(w, errKeyNotFound, http.StatusNotFound)
		return
	}
	// TODO: use client from cache
	client, err := cacheAPI.NewClientWithResponses(fmt.Sprintf("http://%s", url))
	if err != nil {
		s.logger.With(err).Error("while creating client")
		s.respondError(w, err, http.StatusInternalServerError)
		return
	}
	resp, err := client.GetValueWithResponse(r.Context(), key)
	if err != nil {
		s.logger.With(err).Error("while getting value")
		s.respondError(w, err, http.StatusInternalServerError)
		return
	}
	if resp.JSON200 == nil {
		s.respondError(w, errKeyNotFound, http.StatusNotFound)
		return
	}
	s.respond(w, *resp.JSON200, http.StatusOK)
}

func (s server) RegisterNode(w http.ResponseWriter, r *http.Request) {
	remoteAddr := r.RemoteAddr
	s.logger.Info("registering node", "addr", remoteAddr)
	host, err := extractHost(w, remoteAddr, s)
	if err != nil {
		return
	}

	s.sh.AddNode(host)
	s.respond(w, nil, http.StatusNoContent)
}

func (s server) UnRegisterNode(w http.ResponseWriter, r *http.Request, _ registerAPI.Address) {
	remoteAddr := r.RemoteAddr
	s.logger.Info("registering node", "addr", remoteAddr)
	host, err := extractHost(w, remoteAddr, s)
	if err != nil {
		return
	}
	s.sh.DeleteNode(host)

	s.respond(w, nil, http.StatusNoContent)
}

func extractHost(w http.ResponseWriter, remoteAddr string, s server) (string, error) {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		s.logger.With(err).Error("while splitting host and port")
		s.respondError(w, err, http.StatusInternalServerError)
		return "", err
	}
	// TODO: confgiure port
	addr := fmt.Sprintf("%s:8080", host)
	return addr, nil
}

func (s server) respond(w http.ResponseWriter, data any, statusCode int) {
	err := httpFramework.Respond(w, data, statusCode)
	if err != nil {
		s.logger.With(err).Error("while responding")
		return
	}
}
func (s server) respondError(w http.ResponseWriter, err error, code int) {
	if err := httpFramework.Respond(w, map[string]string{"error": err.Error()}, code); err != nil {
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
