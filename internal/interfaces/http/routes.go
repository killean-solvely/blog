package http

import (
	"net/http"
	"time"

	"blog/internal/application"
	"blog/internal/interfaces/http/handlers"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	postService *application.PostService,
	userService *application.UserService,
	commentService *application.CommentService,
	ratingService *application.RatingService,
) *chi.Mux {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(sessionManager.LoadAndSave)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		postHandler := handlers.NewPostHandler(postService, sessionManager)
		postHandler.Register(r)
	})

	return r
}
