package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/arosace/WellnessWaveApi/internal/account/model"
	"github.com/arosace/WellnessWaveApi/internal/account/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
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
func (h *AccountHandler) HandleAddAccount(c echo.Context) error {
	res := model.AccountResponse{}

	var account model.Account
	if err := c.Bind(&account); err != nil {
		return apis.NewBadRequestError("wrong_data_type", err)
	}

	if err := account.ValidateModel(); err != nil {
		res.Error = fmt.Sprintf("inavlid_data_format: %v", err)
		return apis.NewBadRequestError(res.Error, res)
	}

	if alreadyExists := h.accountService.CheckAccountExists(c, account.Email); alreadyExists {
		res.Error = "email_already_in_use"
		return apis.NewApiError(http.StatusConflict, res.Error, res)
	}

	if _, err := h.accountService.AddAccount(c, account); err != nil {
		res.Error = fmt.Sprintf("failed_to_add_user: %v", err)
		return apis.NewApiError(http.StatusInternalServerError, res.Error, res)
	}

	return c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) HandleGetAccounts(c echo.Context) error {
	res := model.AccountResponse{}

	accounts, err := h.accountService.GetAccounts(c)
	if err != nil {
		res.Error = fmt.Sprintf("Failed to get accounts: %v", err)
		return c.JSON(http.StatusInternalServerError, res)
	}

	res.Data = accounts
	return c.JSON(http.StatusOK, res)
}

func (h *AccountHandler) HandleGetAccountsById(w http.ResponseWriter, r *http.Request) {
	res := model.AccountResponse{}
	ctx := r.Context()
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	requestParams := utils.GetHTTPVars(r)
	id := requestParams["id"]
	if id == "" {
		res.Error = "parameter id is missing"
		utils.FormatResponse(w, res, http.StatusBadRequest)
		return
	}

	if _, err := strconv.ParseInt(id, 10, 32); err != nil {
		res.Error = "unexpected http parameter"
		utils.FormatResponse(w, res, http.StatusBadRequest)
		return
	}
	account, err := h.accountService.GetAccountById(ctx, id)
	if err != nil {
		res.Error = fmt.Sprintf("Failed to get account [%s]: %v", id, err)
		utils.FormatResponse(w, res, http.StatusInternalServerError)
		return
	}

	res.Data = account
	utils.FormatResponse(w, res, http.StatusOK)
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
	res := model.AccountResponse{}
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	requestParams := utils.GetHTTPVars(r)
	parentId := requestParams["parent_id"]
	if parentId == "" {
		res.Error = "parameter parent_id is missing"
		utils.FormatResponse(w, res, http.StatusBadRequest)
		return
	}
	if _, err := strconv.ParseInt(parentId, 10, 32); err != nil {
		res.Error = "unexpected http parameter"
		utils.FormatResponse(w, res, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	attachedAccounts, err := h.accountService.GetAttachedAccounts(ctx, parentId)
	if err != nil {
		res.Error = fmt.Sprintf("Failed to get accounts attached to account (%s): %v", parentId, err)
		utils.FormatResponse(w, res, http.StatusInternalServerError)
		return
	}

	res.Data = attachedAccounts
	utils.FormatResponse(w, res, http.StatusOK)
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

func (h *AccountHandler) HandleLogIn(w http.ResponseWriter, r *http.Request) {
	res := model.AccountResponse{}
	ctx := r.Context()

	if r.Method != http.MethodPost {
		res.Error = fmt.Sprintf("%s_method_not_allowed", r.Method)
		utils.FormatResponse(w, res, http.StatusMethodNotAllowed)
		return
	}
	var params model.LogInCredentials
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		res.Error = "invalid_data_format"
		utils.FormatResponse(w, res, http.StatusBadRequest)
		return
	}
	if err := params.ValidateModel(); err != nil {
		res.Error = fmt.Sprintf("invalid_data: %v", err)
		utils.FormatResponse(w, res, http.StatusBadRequest)
		return
	}

	authorizedAccount, err := h.accountService.Authorize(ctx, params)
	if err != nil {
		res.Error = "unauthorized"
		utils.FormatResponse(w, res, http.StatusUnauthorized)
		return
	}

	res.Data = authorizedAccount
	utils.FormatResponse(w, res, http.StatusOK)
}
