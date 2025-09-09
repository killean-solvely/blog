package middleware

import (
	"net/http"
	"slices"

	"blog/internal/application"

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

func RequireAdminAuth(
	sessionManager *scs.SessionManager,
	userService *application.UserService,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the user ID from sesson
			userID := sessionManager.GetString(r.Context(), "user_id")
			if userID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Get the user and validate that they've got the admin permission
			user, err := userService.GetUserByID(userID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !slices.Contains(user.UserRoles, "ADMIN") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
