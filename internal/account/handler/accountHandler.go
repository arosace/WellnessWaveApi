package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/arosace/WellnessWaveApi/internal/account/model"
	"github.com/arosace/WellnessWaveApi/internal/account/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/gorilla/mux"
)

// UserHandler handles HTTP requests for user operations.
type AccountHandler struct {
	accountService service.AccountService
}

// NewUserHandler creates a new instance of UserHandler.
func NewAccountHandler(accountService service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// HandleAddUser handles the POST request to add a new user.
func (h *AccountHandler) HandleAddAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var account model.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "wrong_data_type", http.StatusBadRequest)
		return
	}

	if err := account.ValidateModel(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if alreadyExists := h.accountService.CheckAccountExists(ctx, account.Email); alreadyExists {
		http.Error(w, "email_already_in_use", http.StatusConflict)
		return
	}

	if _, err := h.accountService.AddAccount(ctx, account); err != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AccountHandler) HandleGetAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	accounts, err := h.accountService.GetAccounts(ctx)
	if err != nil {
		http.Error(w, "Failed to get accounts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func (h *AccountHandler) HandleAttachAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var attachBody model.AttachAccountBody
	if err := json.NewDecoder(r.Body).Decode(&attachBody); err != nil {
		http.Error(w, "wrong_data_type", http.StatusBadRequest)
		return
	}

	if err := attachBody.ValidateModel(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	account, err := h.accountService.AttachAccount(ctx, attachBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to attach account: %v", err), http.StatusInternalServerError)
		return
	}

	if account == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *AccountHandler) HandleGetAttachedAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	fmt.Println(vars)
	requestParams := utils.GetHTTPVars(r)
	parentId := requestParams["parent_id"]
	if parentId == "" {
		http.Error(w, "parameter parent_id is missing", http.StatusBadRequest)
		return
	}
	if _, err := strconv.ParseInt(parentId, 10, 32); err != nil {
		http.Error(w, "unexpected http parameter", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	attachedAccounts, err := h.accountService.GetAttachedAccounts(ctx, parentId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get accounts attached to account (%s): %v", parentId, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attachedAccounts)
}

func (h *AccountHandler) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	infoType := r.URL.Query().Get("type")
	if infoType == "" || (infoType != "personal" && infoType != "authentication") {
		http.Error(w, fmt.Sprintf("type url parameter is either missing or invalid. Got: %s. Expected either 'personal' or 'authentication'", infoType), http.StatusBadRequest)
		return
	}

	var account model.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		return
	}

	switch infoType {
	case "personal":
		if err := account.ValidateModelForInfoUpdate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case "authentication":
		if err := account.ValidateModelForAuthUpdate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if _, err := h.accountService.UpdateAccount(ctx, account, infoType); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update account: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
