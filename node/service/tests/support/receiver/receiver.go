// Package receiver is the HTTP client for the service.
package receiver

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

type (
	// ResponseMessage ...
	ResponseMessage struct {
		Body []byte `json:"message"`
		Code int    `json:"code"`
	}

	// Receiver ...
	Receiver struct {
		responseCodes []ResponseMessage
		mu            *sync.RWMutex
	}
)

// NewReceiver ...
func NewReceiver(responseCodes []ResponseMessage) *Receiver {
	return &Receiver{
		responseCodes: responseCodes,
		mu:            &sync.RWMutex{},
	}
}

// UpdateResponseCodes ...
func (rec *Receiver) UpdateResponseCodes(codes []ResponseMessage) {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	rec.responseCodes = append(rec.responseCodes, codes...)
}

// Serve ...
func (rec *Receiver) Serve() (string, func()) {

	i := 0
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.mu.Lock()
		defer rec.mu.Unlock()

		if i >= len(rec.responseCodes) {
			w.WriteHeader(504)
			return
		}
		defer func() {
			_ = r.Body.Close()
		}()

		msg := rec.responseCodes[i]
		i++
		w.WriteHeader(msg.Code)
		_, _ = w.Write(msg.Body)
	})

	server := CreateTestLocalServer(handlerFunc)
	address := fmt.Sprintf("http://127.0.0.1:%d", server.Listener.Addr().(*net.TCPAddr).Port)

	return address, func() {
		if server != nil {
			server.Close()
		}
	}

}
