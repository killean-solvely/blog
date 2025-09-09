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
	// Register user
	mux.Post("/register", h.RegisterUser)

	// Login user
	mux.Post("/login", h.LoginUser)

	// Logout user
	mux.Post("/logout", h.LogoutUser)

	mux.Route("/users", func(r chi.Router) {
		// Get user by id
		r.Get("/{id}", h.GetUser)

		// Get users
		r.Get("/", h.GetUsers)

		r.Group(func(r chi.Router) {
			// Protected routes
			r.Use(middleware.RequireAuth(h.sessionManager))

			// Update user description
			r.Post("/description", h.UpdateUserDescription)

			// Update user password
			r.Post("/password", h.UpdateUserPassword)
		})
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
	// Decode the request and validate it
	var req requests.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("LoginUser: failed to decode request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errors := req.Validate(); errors != nil {
		log.Println("LoginUser: invalid request data")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errors.Error()))
		return
	}

	// Validate the password and get the user if valid
	user, err := h.userService.ValidatePassword(req.Username, req.Password)
	if err != nil {
		log.Println("LoginUser: failed to get user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect username or password"))
		return
	}

	// Add the user ID to the session
	h.sessionManager.Put(r.Context(), "user_id", user.ID)

	w.WriteHeader(http.StatusOK)
}

func (h UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionManager.Destroy(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		log.Println("GetUser: failed to get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Println("GetUser: failed to marshal user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetAllUsers()
	if err != nil {
		log.Println("GetUsers: failed to get users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Println("GetUsers: failed to marshal users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (h UserHandler) SetUserRoles(w http.ResponseWriter, r *http.Request) {
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

	// Get the userID making the request
	userID := h.sessionManager.GetString(r.Context(), "user_id")

	// Update the user's roles
	if err := h.userService.SetUserRoles(userID, req.UserRoles); err != nil {
		log.Println("SetUserRoles: failed to set user's roles")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h UserHandler) UpdateUserDescription(w http.ResponseWriter, r *http.Request) {
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

	// Get the userID making the request
	userID := h.sessionManager.GetString(r.Context(), "user_id")

	// Update the user's description
	if err := h.userService.UpdateDescription(userID, req.Description); err != nil {
		log.Println("UpdateUserDescription: failed to update user's description")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h UserHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
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

	// Get the userID making the request
	userID := h.sessionManager.GetString(r.Context(), "user_id")

	// Update the user's password
	if err := h.userService.UpdatePassword(userID, req.Password); err != nil {
		log.Println("UpdateUserPassword: failed to update user's password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
