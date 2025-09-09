package handlers

import (
	"net/http"

	"blog/internal/application"
	"blog/internal/interfaces/http/middleware"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type RatingHandler struct {
	ratingService  *application.RatingService
	sessionManager *scs.SessionManager
}

func NewRatingHandler(
	ratingService *application.RatingService,
	sessionManager *scs.SessionManager,
) *RatingHandler {
	return &RatingHandler{
		ratingService:  ratingService,
		sessionManager: sessionManager,
	}
}

func (h RatingHandler) Register(mux chi.Router) {
	mux.Route("/ratings", func(r chi.Router) {
		// Get ratings on post
		r.Get("/post/{post_id}", nil)

		// Get rating
		r.Get("/{id}", nil)

		r.Group(func(r chi.Router) {
			// Authorized routes
			r.Use(middleware.RequireAuth(h.sessionManager))

			// Create rating
			r.Post("/", nil)

			// Change rating
			r.Patch("/{id}", nil)

			// Remove rating
			r.Delete("/{id}", nil)
		})
	})
}

func (h RatingHandler) CreateRating(w http.ResponseWriter, r *http.Request) {
}

func (h RatingHandler) GetRatingsOnPost(w http.ResponseWriter, r *http.Request) {
}

func (h RatingHandler) GetRating(w http.ResponseWriter, r *http.Request) {
}

func (h RatingHandler) ChangeRating(w http.ResponseWriter, r *http.Request) {
}

func (h RatingHandler) RemoveRating(w http.ResponseWriter, r *http.Request) {
}
