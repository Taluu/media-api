package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Taluu/media-go/pkg/domain/media"
)

type httpError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func jsonError(w http.ResponseWriter, error string, code int) {
	HTTPError := httpError{
		Code:  code,
		Error: error,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	encoder.Encode(HTTPError)
}

func jsonResponse(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func toHttpCode(err error) (code int) {
	switch {
	case err == nil:
		code = http.StatusOK
	case errors.Is(err, media.ErrFileNotFound):
		fallthrough
	case errors.Is(err, media.ErrMediaNotFound):
		code = http.StatusNotFound
	default:
		code = http.StatusInternalServerError
	}

	return
}
