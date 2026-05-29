package handlers

import (
	"net/http"

	"github.com/thedevrems/tuto_lua/server/internal/httpx"
)

// bind decodes the JSON request body into dst. On failure it writes a 400 with
// the reason and returns false, so handlers can `if !bind(...) { return }`.
func bind(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := httpx.Decode(w, r, dst); err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}
