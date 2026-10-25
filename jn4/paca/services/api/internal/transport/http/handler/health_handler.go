// Package handler contains the HTTP request handlers for the API service.
package handler

import (
	"encoding/json"
	"net/http"
)

// HealthHandler serves the /healthz endpoint.
type HealthHandler struct{}

// NewHealthHandler returns a HealthHandler.
func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

// Check responds with a 200 OK and a short status payload.
func (h *HealthHandler) Check(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}
