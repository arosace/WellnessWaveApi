package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arosace/WellnessWaveApi/internal/account/model"
	"github.com/arosace/WellnessWaveApi/internal/account/service"
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

func (h *AccountHandler) HandleGetAccountsById(c echo.Context) error {
	res := model.AccountResponse{}
	id := c.PathParam("id")

	if id == "" {
		res.Error = "parameter id is missing"
		return apis.NewBadRequestError(res.Error, res)
	}

	if _, err := strconv.ParseInt(id, 10, 32); err != nil {
		res.Error = "unexpected http parameter"
		return apis.NewBadRequestError(res.Error, res)
	}
	account, err := h.accountService.GetAccountById(c, id)
	if err != nil {
		res.Error = fmt.Sprintf("Failed to get account [%s]: %v", id, err)
		return apis.NewApiError(http.StatusInternalServerError, res.Error, res)
	}

	res.Data = account
	return c.JSON(http.StatusOK, res)
}

func (h *AccountHandler) HandleAttachAccount(ctx echo.Context) error {
	var attachBody model.AttachAccountBody

	if err := ctx.Bind(&attachBody); err != nil {
		return apis.NewBadRequestError("wrong_data_type", nil)
	}

	if err := attachBody.ValidateModel(); err != nil {
		return apis.NewBadRequestError(err.Error(), nil)
	}

	account, err := h.accountService.AttachAccount(ctx, attachBody)
	if err != nil {
		return apis.NewApiError(http.StatusBadRequest, fmt.Sprintf("Failed to attach account: %v", err), nil)
	}

	if account == nil {
		return ctx.JSON(http.StatusOK, account)
	} else {
		return ctx.JSON(http.StatusCreated, account)
	}
}

func (h *AccountHandler) HandleGetAttachedAccounts(ctx echo.Context) error {
	res := model.AccountResponse{}
	parentId := ctx.PathParam("parent_id")
	if parentId == "" {
		res.Error = "parameter parent_id is missing"
		return apis.NewBadRequestError(res.Error, nil)
	}
	if _, err := strconv.ParseInt(parentId, 10, 32); err != nil {
		res.Error = "unexpected http parameter"
		return apis.NewBadRequestError(res.Error, nil)
	}

	attachedAccounts, err := h.accountService.GetAttachedAccounts(ctx, parentId)
	if err != nil {
		res.Error = fmt.Sprintf("Failed to get accounts attached to account (%s): %v", parentId, err)
		return apis.NewApiError(http.StatusInternalServerError, res.Error, nil)
	}

	res.Data = attachedAccounts
	return ctx.JSON(http.StatusOK, res)
}

func (h *AccountHandler) HandleUpdateAccount(ctx echo.Context) error {
	infoType := ctx.QueryParam("search")
	if infoType == "" || (infoType != "personal" && infoType != "authentication") {
		return apis.NewBadRequestError(fmt.Sprintf("type url parameter is either missing or invalid. Got: %s. Expected either 'personal' or 'authentication'", infoType), nil)
	}

	var account model.Account
	if err := ctx.Bind(&account); err != nil {
		return apis.NewBadRequestError("Invalid data format", nil)
	}

	switch infoType {
	case "personal":
		if err := account.ValidateModelForInfoUpdate(); err != nil {
			return apis.NewBadRequestError(err.Error(), nil)
		}
	case "authentication":
		if err := account.ValidateModelForAuthUpdate(); err != nil {
			return apis.NewBadRequestError(err.Error(), nil)
		}
	}

	if _, err := h.accountService.UpdateAccount(ctx, account, infoType); err != nil {
		return apis.NewApiError(http.StatusInternalServerError, fmt.Sprintf("Failed to update account: %v", err), nil)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h *AccountHandler) HandleLogIn(ctx echo.Context) error {
	res := model.AccountResponse{}

	var params model.LogInCredentials
	if err := ctx.Bind(&params); err != nil {
		res.Error = "invalid_data_format"
		return apis.NewBadRequestError(res.Error, nil)
	}
	if err := params.ValidateModel(); err != nil {
		res.Error = fmt.Sprintf("invalid_data: %v", err)
		return apis.NewBadRequestError(res.Error, nil)
	}

	authorizedAccount, err := h.accountService.Authorize(ctx, params)
	if err != nil {
		res.Error = "unauthorized"
		return apis.NewUnauthorizedError(res.Error, nil)
	}

	res.Data = authorizedAccount
	return ctx.JSON(http.StatusOK, res)
}
