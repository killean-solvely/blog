package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func RequireAuth(sessionManager *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the user ID from sesson
			userID := sessionManager.GetString(r.Context(), "user_id")
			if userID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// User is authenticated, continue
			next.ServeHTTP(w, r)
		})
	}
}
