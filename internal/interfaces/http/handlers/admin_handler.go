package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"blog/internal/application"
	"blog/internal/interfaces/http/middleware"
	"blog/internal/interfaces/http/requests"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	userService    *application.UserService
	postService    *application.PostService
	commentService *application.CommentService
	sessionManager *scs.SessionManager
}

func (h AdminHandler) Register(mux chi.Router) {
	mux.Route("/admin", func(r chi.Router) {
		// Admin authorized routes
		r.Use(middleware.RequireAdminAuth(h.sessionManager, h.userService))

		r.Route("/users", func(r chi.Router) {
			// Set user roles
			r.Post("/{id}/roles", h.SetUserRoles)

			// Update user description
			r.Post("/{id}/description", h.UpdateUserDescription)

			// Update user password
			r.Post("/{id}/password", h.UpdateUserPassword)
		})
	})
}

func (h AdminHandler) SetUserRoles(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.SetUserRolesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("SetUserRoles: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("SetUserRoles: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	userID := chi.URLParam(r, "id")
	if userID == "" {
		log.Println("SetUserRoles: missing user_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the user's roles
	if err := h.userService.SetUserRoles(userID, req.UserRoles); err != nil {
		log.Println("SetUserRoles: failed to set user's roles")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h AdminHandler) UpdateUserDescription(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.UpdateUserDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("UpdateUserDescription: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("UpdateUserDescription: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	userID := chi.URLParam(r, "id")
	if userID == "" {
		log.Println("UpdateUserDescription: missing user_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the user's description
	if err := h.userService.UpdateDescription(userID, req.Description); err != nil {
		log.Println("UpdateUserDescription: failed to update user's description")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h AdminHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.UpdateUserPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("UpdateUserPassword: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("UpdateUserPassword: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	userID := chi.URLParam(r, "id")
	if userID == "" {
		log.Println("UpdateUserPassword: missing user_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the user's password
	if err := h.userService.UpdatePassword(userID, req.Password); err != nil {
		log.Println("UpdateUserPassword: failed to update user's password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
