package middleware

import (
	"net/http"

	"github.com/Paca-AI/api/internal/apierr"
	"github.com/Paca-AI/api/internal/transport/http/presenter"
)

// RequireFreshPassword rejects any request whose access token carries
// MustChangePassword=true. Apply this middleware after Authn on every route
// except PATCH /users/me/password.
func RequireFreshPassword() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ok bool
			r, ok = EnforceFreshPassword(w, r)
			if !ok {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// EnforceFreshPassword checks MustChangePassword without advancing the handler chain.
// Returns the (possibly updated) request and whether to continue.
func EnforceFreshPassword(w http.ResponseWriter, r *http.Request) (*http.Request, bool) {
	claims := ClaimsFrom(r)
	if claims != nil && claims.MustChangePassword {
		presenter.Error(w, r, apierr.New(
			apierr.CodePasswordChangeRequired,
			"you must change your password before accessing this resource",
		))
		return r, false
	}
	return r, true
}
