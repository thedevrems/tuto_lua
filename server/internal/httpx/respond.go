// Package httpx provides small helpers for JSON request/response handling so
// handlers stay focused on business logic.
package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// maxBody caps request payloads to guard against abuse (1 MiB).
const maxBody = 1 << 20

// JSON writes v as a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// Error writes a uniform error envelope: {"error": "message"}.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

// Decode reads and strictly parses a JSON request body into dst.
// It returns a user-friendly error suitable for a 400 response.
func Decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return decodeError(err)
	}
	if dec.Decode(&struct{}{}) != io.EOF {
		return errors.New("le corps de la requête doit contenir un seul objet JSON")
	}
	return nil
}

// decodeError maps low-level JSON failures to a single readable message.
func decodeError(err error) error {
	var synErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	switch {
	case errors.As(err, &synErr), errors.As(err, &typeErr):
		return errors.New("JSON invalide dans le corps de la requête")
	case errors.Is(err, io.EOF):
		return errors.New("le corps de la requête est vide")
	default:
		return errors.New("impossible de lire le corps de la requête")
	}
}
