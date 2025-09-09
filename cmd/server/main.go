package main

import (
	"log"
	"net/http"

	dddmemory "blog/pkg/ddd/memory"

	"blog/internal/application"
	"blog/internal/infrastructure/events"
	"blog/internal/infrastructure/persistence/sqlite"
	httphandler "blog/internal/interfaces/http"
)

func main() {
	eventDispatcher := dddmemory.NewInMemoryEventDispatcher(nil)

	commentEventHandler := events.NewCommentEventHandler()
	postEventHandler := events.NewPostEventHandler()
	ratingEventHandler := events.NewRatingEventHandler()
	userEventHandler := events.NewUserEventHandler()

	commentEventHandler.Register(eventDispatcher)
	postEventHandler.Register(eventDispatcher)
	ratingEventHandler.Register(eventDispatcher)
	userEventHandler.Register(eventDispatcher)

	db, err := sqlite.NewDB()
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic("failed db ping")
	}

	commentRepo := sqlite.NewCommentRepository(db.DB)
	postRepo := sqlite.NewPostRepository(db.DB)
	ratingRepo := sqlite.NewRatingRepository(db.DB)
	userRepo := sqlite.NewUserRepository(db.DB)

	commentService := application.NewCommentService(
		commentRepo,
		userRepo,
		postRepo,
		eventDispatcher,
	)
	postService := application.NewPostService(postRepo, userRepo, eventDispatcher)
	ratingService := application.NewRatingService(ratingRepo, userRepo, postRepo, eventDispatcher)
	userService := application.NewUserService(userRepo, eventDispatcher)

	router := httphandler.NewRouter(postService, userService, commentService, ratingService)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
