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

func (h RatingHandler) GetRatingsOnPost(w http.ResponseWriter, r *http.Request) {
	// Get the post id from the path and validate it
	postID := chi.URLParam(r, "post_id")
	if postID == "" {
		log.Println("GetRatingsOnPost: missing post_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the ratings for the post
	ratings, err := h.ratingService.GetRatingsOnPost(postID)
	if err != nil {
		log.Println("GetRatingsOnPost: failed to get ratings on post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the data and return it to the user
	data, err := json.Marshal(ratings)
	if err != nil {
		log.Println("GetRatingsOnPost: failed to marshal ratings")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h RatingHandler) GetRating(w http.ResponseWriter, r *http.Request) {
	// Get the rating id and validate it
	ratingID := chi.URLParam(r, "id")
	if ratingID == "" {
		log.Println("GetRating: missing post_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the rating
	ratings, err := h.ratingService.GetRating(ratingID)
	if err != nil {
		log.Println("GetRating: failed to get rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the data and return it to the user
	data, err := json.Marshal(ratings)
	if err != nil {
		log.Println("GetRating: failed to marshal rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h RatingHandler) CreateRating(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.CreateRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("CreateRating: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("CreateRating: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	// Get the userID making the request
	userID := h.sessionManager.GetString(r.Context(), "user_id")

	// Create the rating
	rating, err := h.ratingService.CreateRating(req.PostID, userID, req.RatingType)
	if err != nil {
		log.Println("CreateRating: failed to create rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the rating to the requester
	data, err := json.Marshal(rating)
	if err != nil {
		log.Println("CreateRating: failed to marshal rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h RatingHandler) ChangeRating(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.ChangeRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("ChangeRating: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("ChangeRating: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	// Get the userID making the request
	userID := h.sessionManager.GetString(r.Context(), "user_id")

	// Validate the rating belongs to the user
	rating, err := h.ratingService.GetRating(req.RatingID)
	if err != nil {
		log.Println("ChangeRating: failed to get rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rating.UserID != userID {
		log.Println("ChangeRating: wrong user attempting to change rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the rating
	if err := h.ratingService.UpdateRating(req.RatingID, req.RatingType); err != nil {
		log.Println("ChangeRating: failed to update rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h RatingHandler) RemoveRating(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.RemoveRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("RemoveRating: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("RemoveRating: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	// Get the userID making the request
	userID := h.sessionManager.GetString(r.Context(), "user_id")

	// Validate the rating belongs to the user
	rating, err := h.ratingService.GetRating(req.RatingID)
	if err != nil {
		log.Println("RemoveRating: failed to get rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rating.UserID != userID {
		log.Println("RemoveRating: wrong user attempting to remove rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Remove the rating
	if err := h.ratingService.RemoveRating(req.RatingID); err != nil {
		log.Println("RemoveRating: failed to remove rating")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
