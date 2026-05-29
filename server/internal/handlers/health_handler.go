package handlers

import (
	"net/http"

	"github.com/thedevrems/tuto_lua/server/internal/httpx"
)

// Health is a liveness probe used by tooling and uptime checks.
func Health(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
