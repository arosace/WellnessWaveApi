package handler

import (
	"encoding/json"
	"net/http"

	"../../model"
	"../../service"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	accountService service.AccountService
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(accountService service.AccountService) *UserHandler {
	return &UserHandler{accountService: accountService}
}

// HandleAddUser handles the POST request to add a new user.
func (h *UserHandler) HandleAddAccount(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.accountService.AddAccount(ctx, user); err != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
