package auth

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"tek-bank/internal/dto"
	"tek-bank/internal/i18n"
	"tek-bank/internal/i18n/messages"
	"tek-bank/internal/service"
	"tek-bank/pkg/cresponse"
)

type AuthHandler interface {
	Login(ctx *fiber.Ctx) error
	GetUserInfo(ctx *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

// Login godoc
// @Summary Login a user
// @Description You can login with your identity number or customer number. If you are a new user, you can register with the /account/register endpoint.
// @Description If you registered before, your password will be sent to your e-mail address.
// @Tags Auth
// @Accept application/json
// @Produce application/json
// @Param loginRequest body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.LoginResponse
// @Router /auth/login [post]
func (h *authHandler) Login(ctx *fiber.Ctx) error {
	var request dto.LoginRequest
	if err := ctx.BodyParser(&request); err != nil {
		err = errors.New(i18n.CreateMsg(ctx, messages.InvalidLoginCredentials))
		log.Error(err.Error())
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	response, err := h.authService.Login(ctx.Context(), request)
	if err != nil {
		var status int = fiber.StatusInternalServerError
		if err.Error() == messages.UserNotFound || err.Error() == messages.PasswordIncorrect {
			status = fiber.StatusUnauthorized
		}
		log.Error(err.Error())
		return cresponse.ErrorResponse(ctx, status, i18n.CreateMsg(ctx, err.Error()))
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusOK, response)
}

// GetUserInfo godoc
// @Summary Get user info
// @Description You can use this endpoint to get user information.
// @Tags Auth
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer <token>"
// @Success 200 {object} dto.UserInfoResponse
// @Router /auth/user-info [get]
func (h *authHandler) GetUserInfo(ctx *fiber.Ctx) error {
	response, err := h.authService.GetUserInfo(ctx.Context())
	if err != nil {
		var status int = fiber.StatusInternalServerError
		if err.Error() == messages.UserNotFound {
			status = fiber.StatusNotFound
		}
		log.Error(err.Error())
		return cresponse.ErrorResponse(ctx, status, i18n.CreateMsg(ctx, err.Error()))
	}

	return cresponse.SuccessResponse(ctx, fiber.StatusOK, response)
}
