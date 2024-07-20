package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"strconv"
	"tek-bank/internal/i18n"
	"tek-bank/internal/i18n/messages"
	"tek-bank/internal/service"
	"tek-bank/pkg/cresponse"
)

type ProfileHandler interface {
	MyProfile(ctx *fiber.Ctx) error
	MyTransferHistory(ctx *fiber.Ctx) error
}

type profileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) ProfileHandler {
	return &profileHandler{
		profileService: profileService,
	}
}

// MyProfile godoc
// @Summary Get user profile
// @Description This endpoint is used to get the user's profile information. It returns the user's account information and user information.
// @Description The user's account information includes the account number, balance, and account type.
// @Description The user information includes the user's name, surname, and e-mail address.
// @Description NOTE! The user's account information is returned as an array because the user can have more than one account.
// @Tags Profile
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} dto.GetProfileResponse
// @Router /profile [get]
func (h *profileHandler) MyProfile(ctx *fiber.Ctx) error {
	response, err := h.profileService.MyProfile(ctx.Context())
	if err != nil {
		var status int = fiber.StatusInternalServerError
		if err.Error() == messages.AccountNotFound {
			status = fiber.StatusNotFound
		}
		log.Error(err.Error())
		return cresponse.ErrorResponse(ctx, status, i18n.CreateMsg(ctx, err.Error()))
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusOK, response)
}

// MyTransferHistory godoc
// @Summary Get user transfer history
// @Description You can see only your own transfer history with this endpoint.
// @Description You can filter the transfer history by account number.
// @Description The account number is required for this endpoint.
// @Tags Profile
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Param accountNumber query string true "Account Number"
// @Success 200 {object} []dto.GetTransferHistoryResponse
// @Router /profile/transfer-history [get]
func (h *profileHandler) MyTransferHistory(ctx *fiber.Ctx) error {

	accountNumber, err := strconv.Atoi(ctx.Query("accountNumber"))
	if err != nil {
		log.Error(err.Error())
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, i18n.CreateMsg(ctx, messages.BadRequest))
	}

	response, err := h.profileService.MyTransferHistory(ctx.Context(), int64(accountNumber))
	if err != nil {
		var status int = fiber.StatusInternalServerError
		if err.Error() == messages.AccountNotFound {
			status = fiber.StatusNotFound
		}
		return cresponse.ErrorResponse(ctx, status, i18n.CreateMsg(ctx, err.Error()))
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusOK, response)
}
