package main

import (
	"log"
	"net/http"

	dddmemory "blog/pkg/ddd/memory"

	"blog/internal/application"
	"blog/internal/infrastructure/events"
	"blog/internal/infrastructure/persistence/memory"
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

	commentRepo := memory.NewCommentRepository()
	postRepo := memory.NewPostRepository()
	ratingRepo := memory.NewRatingRepository()
	userRepo := memory.NewUserRepository()

	commentService := application.NewCommentService(
		commentRepo,
		userRepo,
		postRepo,
		eventDispatcher,
	)
	postService := application.NewPostService(postRepo, userRepo, eventDispatcher)
	ratingService := application.NewRatingService(ratingRepo, userRepo, postRepo, eventDispatcher)
	userService := application.NewUserService(userRepo, eventDispatcher)

	newUser, err := userService.CreateUser(
		"a@abc.com",
		"12345678",
		"a",
		[]string{"COMMENTER", "AUTHOR"},
	)
	if err != nil {
		panic(err)
	}

	postService.CreatePost(newUser.ID, "New Post", "Post Content")

	router := httphandler.NewRouter(postService, userService, commentService, ratingService)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
