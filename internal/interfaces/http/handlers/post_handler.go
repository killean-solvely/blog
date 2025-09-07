package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"blog/internal/application"
	"blog/internal/interfaces/http/requests"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	postService *application.PostService
}

func NewPostHandler(postService *application.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h PostHandler) Register(mux chi.Router) {
	mux.Route("/posts", func(r chi.Router) {
		// Create post
		r.Post("/", h.CreatePost)

		// Get posts
		r.Get("/", h.GetPosts)

		// Get post
		r.Get("/{id}", h.GetPost)

		// Update post title
		r.Patch("/{id}/title", h.UpdatePostTitle)

		// Update post content
		r.Patch("/{id}/content", h.UpdatePostContent)

		// Archive post
		r.Delete("/{id}", h.ArchivePost)
	})
}

func (h PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("CreatePost: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("CreatePost: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	// Create the post
	post, err := h.postService.CreatePost(req.AuthorID, req.Title, req.Content)
	if err != nil {
		log.Println("CreatePost: failed to create post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the post to the requester
	data, err := json.Marshal(post)
	if err != nil {
		log.Println("CreatePost: failed to marshal post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetPosts()
	if err != nil {
		log.Println("GetPosts: failed to get posts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the posts to the requester
	data, err := json.Marshal(posts)
	if err != nil {
		log.Println("GetPosts: failed to marshal posts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing post id"))
		return
	}

	log.Printf("ID: %s", id)

	post, err := h.postService.GetPost(id)
	if err != nil {
		log.Println("GetPost: failed to get post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the post to the requester
	data, err := json.Marshal(post)
	if err != nil {
		log.Println("GetPost: failed to marshal post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h PostHandler) UpdatePostTitle(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.UpdatePostTitle
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("UpdatePostTitle: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("UpdatePostTitle: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	// Get the id from the URL path
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing post id"))
		return
	}

	// Update the post title
	if err := h.postService.UpdatePostTitle(id, req.Title); err != nil {
		log.Println("UpdatePostTitle: failed to update post title")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h PostHandler) UpdatePostContent(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.UpdatePostContent
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("UpdatePostContent: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("UpdatePostContent: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	// Get the id from the URL path
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing post id"))
		return
	}

	// Update the post content
	if err := h.postService.UpdatePostContent(id, req.Content); err != nil {
		log.Println("UpdatePostContent: failed to update post content")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h PostHandler) ArchivePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing post id"))
		return
	}

	// Archive the post
	if err := h.postService.ArchivePost(id); err != nil {
		log.Println("ArchivePost: failed to archive post")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
