package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/Paca-AI/api/internal/apierr"
	"github.com/Paca-AI/api/internal/transport/http/presenter"
)

// BindJSON decodes the JSON body into dst and writes a 400 on failure.
// Use this in handlers to keep binding logic concise.
func BindJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		presenter.Error(w, r, apierr.New(apierr.CodeBadRequest, err.Error()))
		return false
	}
	return true
}
