// Package httpx provides shared HTTP request/response helpers used by the
// presenter and middleware packages.
package httpx

import (
	"context"
	"encoding/json"
	"net/http"
)

type requestIDKeyType struct{}

var requestIDKey = requestIDKeyType{}

// WithRequestID returns a copy of ctx carrying the given request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestIDFromContext retrieves the request ID stored by WithRequestID.
// Returns "" if not set.
func RequestIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}

// WriteJSON encodes data as JSON and writes it with the given HTTP status.
// The Content-Type header is always set to application/json.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
