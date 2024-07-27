package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func encode(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

type ResponseError struct {
	Status string `json:"status"`
	Err    string `json:"error"`
}

func errorJSON(w http.ResponseWriter, status int, err string) {
	encode(w, status, &ResponseError{Status: http.StatusText(status), Err: err})
}
