package account

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"strconv"
	"tek-bank/cmd/api/middleware/transaction"
	"tek-bank/internal/dto"
	"tek-bank/internal/i18n"
	"tek-bank/internal/i18n/messages"
	"tek-bank/internal/service"
	"tek-bank/pkg/cresponse"
)

type AccountHandler interface {
	RegisterAccount(ctx *fiber.Ctx) error
	CreateNewAccount(ctx *fiber.Ctx) error
	AddMoney(ctx *fiber.Ctx) error
	TransferMoney(ctx *fiber.Ctx) error
	TransferApproval(ctx *fiber.Ctx) error
}

type accountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler(accountService service.AccountService) AccountHandler {
	return &accountHandler{
		accountService: accountService,
	}
}

// RegisterAccount godoc
// @Summary Register a new account and user
// @Description It should be used for users who will create an account for the first time, because when creating a user account, one user must also be created.
// @Description The user password will be sent via e-mail.
// @Tags Account
// @Accept application/json
// @Produce application/json
// @Param registerRequest body dto.RegisterAccountRequest true "Register Request"
// @Success 200 {object} map[string]interface{}
// @Router /account/register [post]
func (h *accountHandler) RegisterAccount(ctx *fiber.Ctx) error {
	var request dto.RegisterAccountRequest
	if err := ctx.BodyParser(&request); err != nil {
		log.Error(err.Error())
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, i18n.CreateMsg(ctx, messages.BadRequest))
	}

	// Database transaction
	tx, err := transaction.GetDbTx(ctx)
	if err != nil {
		log.Error(err)
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, i18n.CreateMsg(ctx, messages.TransactionFailed))
	}

	err = h.accountService.WithTx(tx).RegisterAccount(ctx.Context(), request)
	if err != nil {
		var status int = fiber.StatusInternalServerError
		if err.Error() == messages.UserAlreadyExists {
			status = fiber.StatusConflict
		}
		return cresponse.ErrorResponse(ctx, status, i18n.CreateMsg(ctx, err.Error()))
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusOK, nil, i18n.CreateMsg(ctx, messages.AccountCreated))
}

// CreateAccount godoc
// @Summary Create a new account for the registered user
// @Description Create a new account for the registered user, the user must be registered before creating an account.
// @Description If you want to create an account for a user who has not registered yet, you should use the register endpoint.
// @Tags Account
// @Accept application/json
// @Produce application/json
// @Param createAccountRequest body dto.CreateNewAccountRequest true "Create Account Request"
// @Success 200 {object} map[string]interface{}
// @Router /account/create [post]
func (h *accountHandler) CreateNewAccount(ctx *fiber.Ctx) error {
	var request dto.CreateNewAccountRequest
	if err := ctx.BodyParser(&request); err != nil {
		err = errors.New(i18n.CreateMsg(ctx, messages.InvalidCreateAccountRequest))
		log.Error(err.Error())
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	response, err := h.accountService.CreateNewAccount(ctx.Context(), request)
	if err != nil {
		var status int = fiber.StatusInternalServerError
		if err.Error() == messages.UserNotFound {
			status = fiber.StatusNotFound
		}
		return cresponse.ErrorResponse(ctx, status, err.Error())
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusOK, response)
}

// AddMoney godoc
// @Summary Add money to the account
// @Description Add money to the account by providing the account number and the amount to be added.
// @Description You can imagine this as a deposit operation. Like depositing money using an ATM.
// @Tags Account
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Param accountNumber path int true "Account Number"
// @Param addMoneyRequest body dto.AddMoneyRequest true "Add Money Request"
// @Success 200 {object} dto.AddMoneyResponse
// @Router /account/add-money/{accountNumber} [put]
func (h *accountHandler) AddMoney(ctx *fiber.Ctx) error {
	accountNumber, err := strconv.Atoi(ctx.Params("accountNumber"))
	if err != nil {
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, i18n.CreateMsg(ctx, messages.BadRequest))
	}

	var request dto.AddMoneyRequest
	if err := ctx.BodyParser(&request); err != nil {
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, i18n.CreateMsg(ctx, messages.BadRequest))
	}

	request.AccountNumber = int64(accountNumber)

	// Database transaction
	tx, err := transaction.GetDbTx(ctx)
	if err != nil {
		log.Error(err)
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, i18n.CreateMsg(ctx, messages.TransactionFailed))
	}

	response, err := h.accountService.WithTx(tx).AddMoney(ctx.Context(), request)
	if err != nil {
		var status int = fiber.StatusInternalServerError
		if err.Error() == messages.AccountNotFound {
			log.Error(err)
			status = fiber.StatusNotFound
		}

		return cresponse.ErrorResponse(ctx, status, i18n.CreateMsg(ctx, err.Error()))
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusOK, response)
}

// TransferMoney godoc
// @Summary Transfer money between accounts
// @Description Transfer money between accounts by providing the account numbers and the amount to be transferred.
// @Tags Account
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Param transferMoneyRequest body dto.TransferMoneyRequest true "Transfer Money Request"
// @Success 200 {object} map[string]interface{}
// @Router /account/transfer [post]
func (h *accountHandler) TransferMoney(ctx *fiber.Ctx) error {
	var request dto.TransferMoneyRequest
	if err := ctx.BodyParser(&request); err != nil {
		log.Error(err.Error())
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, messages.BadRequest)
	}

	// Database transaction
	tx, err := transaction.GetDbTx(ctx)
	if err != nil {
		log.Error(err)
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, messages.TransactionFailed)
	}

	err = h.accountService.WithTx(tx).TransferMoney(ctx.Context(), request)
	if err != nil {
		var status int = fiber.StatusInternalServerError
		if err.Error() == messages.AccountNotFound {
			status = fiber.StatusNotFound
		} else if err.Error() == messages.InSufficientBalance {
			status = fiber.StatusBadRequest
		}
		log.Error(err.Error())
		return cresponse.ErrorResponse(ctx, status, i18n.CreateMsg(ctx, err.Error()))
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusOK, nil)
}

// TransferApproval godoc
// @Summary Approve the transfer
// @Description Approve the transfer by providing the transfer token.
// @Tags Account
// @Accept application/json
// @Produce application/json
// @Param token query string true "Token"
// @Success 200 {object} map[string]interface{}
// @Router /account/transfer-approval [get]
func (h *accountHandler) TransferApproval(ctx *fiber.Ctx) error {
	token := ctx.Query("token")

	// Database transaction
	tx, err := transaction.GetDbTx(ctx)
	if err != nil {
		log.Error(err)
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, i18n.CreateMsg(ctx, messages.TransactionFailed))
	}

	err = h.accountService.WithTx(tx).TransferApproval(ctx.Context(), token)
	if err != nil {
		var status int = fiber.StatusInternalServerError
		if err.Error() == messages.AccountNotFound {
			status = fiber.StatusNotFound
		} else if err.Error() == messages.InSufficientBalance {
			status = fiber.StatusBadRequest
		}
		log.Error(err.Error())
		return cresponse.ErrorResponse(ctx, status, err.Error())
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusOK, nil, i18n.CreateMsg(ctx, messages.TransferApproved))
}
