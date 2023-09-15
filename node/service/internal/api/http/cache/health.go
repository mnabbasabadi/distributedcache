package cache

import (
	"net/http"
)

func (s server) GetLiveness(w http.ResponseWriter, r *http.Request) {
	if err := rootHTTP.Respond(w, "OK", http.StatusOK); err != nil {
		s.logger.With(err).Error("while responding")
		return
	}
}

func (s server) GetReadiness(w http.ResponseWriter, r *http.Request) {
	if err := rootHTTP.Respond(w, "OK", http.StatusOK); err != nil {
		s.logger.With(err).Error("while responding")
		return
	}
}
