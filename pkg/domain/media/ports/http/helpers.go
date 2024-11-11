package http

import (
	"encoding/json"
	"net/http"
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
