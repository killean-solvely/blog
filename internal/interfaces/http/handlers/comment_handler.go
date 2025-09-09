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

type CommentHandler struct {
	commentService *application.CommentService
	sessionManager *scs.SessionManager
}

func NewCommentHandler(
	commentService *application.CommentService,
	sessionManager *scs.SessionManager,
) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		sessionManager: sessionManager,
	}
}

func (h CommentHandler) Register(mux chi.Router) {
	mux.Route("/comments", func(r chi.Router) {
		// Public routes
		r.Get("/", h.GetComments)
		r.Get("/{id}", h.GetComment)

		r.Group(func(r chi.Router) {
			// Protected routes
			r.Use(middleware.RequireAuth(h.sessionManager))

			// Edit comment (only owner can edit)
			r.Patch("/{id}", h.EditComment)

			// Archive comment (only owner can archive)
			r.Delete("/{id}", h.ArchiveComment)
		})
	})

	// Nested route for post comments
	mux.Route("/posts/{postId}/comments", func(r chi.Router) {
		// Public routes
		r.Get("/", h.GetCommentsByPost)

		r.Group(func(r chi.Router) {
			// Protected routes
			r.Use(middleware.RequireAuth(h.sessionManager))

			// Create comment on post
			r.Post("/", h.CreateComment)
		})
	})
}

func (h CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.commentService.GetComments()
	if err != nil {
		log.Println("GetComments: failed to get comments")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(comments)
	if err != nil {
		log.Println("GetComments: failed to marshal comments")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing comment id"))
		return
	}

	comment, err := h.commentService.GetComment(id)
	if err != nil {
		log.Println("GetComment: failed to get comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(comment)
	if err != nil {
		log.Println("GetComment: failed to marshal comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h CommentHandler) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postId")
	if postID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing post id"))
		return
	}

	comments, err := h.commentService.GetCommentsByPost(postID)
	if err != nil {
		log.Println("GetCommentsByPost: failed to get comments for post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(comments)
	if err != nil {
		log.Println("GetCommentsByPost: failed to marshal comments")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("CreateComment: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("CreateComment: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	// Get the postID from URL and userID from session
	postID := chi.URLParam(r, "postId")
	if postID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing post id"))
		return
	}

	userID := h.sessionManager.GetString(r.Context(), "user_id")

	// Create the comment
	comment, err := h.commentService.CreateComment(postID, userID, req.Content)
	if err != nil {
		log.Println("CreateComment: failed to create comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the comment
	data, err := json.Marshal(comment)
	if err != nil {
		log.Println("CreateComment: failed to marshal comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (h CommentHandler) EditComment(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.EditCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("EditComment: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("EditComment: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	// Get the comment ID from URL
	commentID := chi.URLParam(r, "id")
	if commentID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing comment id"))
		return
	}

	// Get the userID from the session
	userID := h.sessionManager.GetString(r.Context(), "user_id")

	// Get the comment and make sure that the user editing it is the owner
	comment, err := h.commentService.GetComment(commentID)
	if err != nil {
		log.Println("EditComment: failed to get comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if comment.CommenterID != userID {
		log.Println("EditComment: non-owner attempting to edit comment")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Edit the comment
	if err := h.commentService.EditComment(commentID, req.Content); err != nil {
		log.Println("EditComment: failed to edit comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h CommentHandler) ArchiveComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	if commentID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing comment id"))
		return
	}

	// Get the userID from the session
	userID := h.sessionManager.GetString(r.Context(), "user_id")

	// Get the comment and make sure that the user archiving it is the owner
	comment, err := h.commentService.GetComment(commentID)
	if err != nil {
		log.Println("ArchiveComment: failed to get comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if comment.CommenterID != userID {
		log.Println("ArchiveComment: non-owner attempting to archive comment")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Archive the comment
	if err := h.commentService.ArchiveComment(commentID); err != nil {
		log.Println("ArchiveComment: failed to archive comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}