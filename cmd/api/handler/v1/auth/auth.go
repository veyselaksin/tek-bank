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
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags auth
// @Accept application/json
// @Produce application/json
// @Param registerRequest body dto.RegisterRequest true "Register Request"
// @Success 200 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *authHandler) Register(ctx *fiber.Ctx) error {
	var request dto.RegisterRequest
	if err := ctx.BodyParser(&request); err != nil {
		err = errors.New(i18n.CreateMsg(ctx, messages.PasswordsDoNotMatch))
		log.Error(err.Error())
		return cresponse.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	status, err := h.authService.Register(ctx, request)
	if err != nil {
		return cresponse.ErrorResponse(ctx, status, err.Error())
	}

	return nil
}

// Login godoc
// @Summary Login a user
// @Description Login a user
// @Tags auth
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

	response, status, err := h.authService.Login(ctx, request)
	if err != nil {
		return cresponse.ErrorResponse(ctx, status, err.Error())
	}

	return cresponse.SuccessResponse(ctx, status, response)
}
