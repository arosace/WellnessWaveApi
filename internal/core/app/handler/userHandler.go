package handler

import (
	"encoding/json"
	"net/http"

	"../../model"
	"../../service"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// HandleAddUser handles the POST request to add a new user.
func (h *UserHandler) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.userService.AddUser(ctx, user); err != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
