package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"blog/internal/application"
	"blog/internal/interfaces/http/requests"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService    *application.UserService
	sessionManager *scs.SessionManager
}

func NewUserHandler(
	userService *application.UserService,
	sessionManager *scs.SessionManager,
) *UserHandler {
	return &UserHandler{
		userService:    userService,
		sessionManager: sessionManager,
	}
}

func (h UserHandler) Register(mux chi.Router) {
	mux.Route("/users", func(r chi.Router) {
		// Register user
		r.Post("/register", h.RegisterUser)

		// Login user
		r.Post("/login", h.LoginUser)

		// Logout user
		r.Post("/logout", h.LogoutUser)

		// Set user roles
		r.Post("/{id}/roles", h.SetUserRoles)

		// Update user description
		r.Post("/{id}/description", h.UpdateUserDescription)

		// Update user password
		r.Post("/{id}/password", h.UpdateUserPassword)

		// Get user by id
		r.Get("/{id}", h.GetUser)

		// Get users
		r.Get("/", h.GetUsers)
	})
}

func (h UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Decode the request and validate it
	var req requests.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("RegisterUser: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("RegisterUser: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	// Register the user
	user, err := h.userService.CreateUser(
		req.Email,
		req.Password,
		req.Username,
		[]string{"COMMENTER"},
	)
	if err != nil {
		log.Println("RegisterUser: failed to create user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the user ID to the requester
	data, err := json.Marshal(map[string]any{
		"user_id": user.ID,
	})
	if err != nil {
		log.Println("RegisterUser: failed to marshal user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
}

func (h UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
}

func (h UserHandler) SetUserRoles(w http.ResponseWriter, r *http.Request) {
}

func (h UserHandler) UpdateUserDescription(w http.ResponseWriter, r *http.Request) {
}

func (h UserHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
}

func (h UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
}

func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
}
