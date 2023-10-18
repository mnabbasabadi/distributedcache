package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Respond(w http.ResponseWriter, data interface{}, statusCode int) error {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}
	if data == nil {
		return fmt.Errorf("data expected for all statuses other than %s", http.StatusText(http.StatusNoContent))
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write(bytes); err != nil {
		return err
	}

	return nil
}

func UnmarshalJSONFromBody(r *http.Request, target interface{}) (err error) {
	defer func() {
		if deferErr := r.Body.Close(); deferErr != nil {
			err = errors.Join(err, deferErr)
		}
	}()

	bodyBytes, bodyErr := io.ReadAll(r.Body)
	if bodyErr != nil {
		return fmt.Errorf("while reading request body: %w", bodyErr)
	}

	if err := json.Unmarshal(bodyBytes, target); err != nil {
		switch err.(type) {
		case *json.SyntaxError:
			return fmt.Errorf("while unmarshalling data from body: %w", err)
		case *json.UnmarshalTypeError:
			return fmt.Errorf("while unmarshalling data from body: %w", err)
		default:
			return fmt.Errorf("while unmarshalling data from body: %w", err)
		}
	}

	return nil
}
