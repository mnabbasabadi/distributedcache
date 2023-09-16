package http

import (
	"net/http"

	HTTPKit "github.com/mnabbasbaadi/distributedcache/foundation/http"
)

func (s server) GetLiveness(w http.ResponseWriter, r *http.Request) {
	if err := HTTPKit.Respond(w, "OK", http.StatusOK); err != nil {
		s.logger.With(err).Error("while responding")
		return
	}
}

func (s server) GetReadiness(w http.ResponseWriter, r *http.Request) {
	if err := HTTPKit.Respond(w, "OK", http.StatusOK); err != nil {
		s.logger.With(err).Error("while responding")
		return
	}
}
