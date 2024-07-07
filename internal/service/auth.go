package service

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"net/http"
	"os"
	"tek-bank/cmd/api/middleware/authware"
	"tek-bank/internal/db/repository"
	"tek-bank/internal/dto"
	"tek-bank/internal/i18n"
	"tek-bank/internal/i18n/messages"
	"tek-bank/pkg/crypto"
)

type AuthService interface {
	Login(ctx *fiber.Ctx, request dto.LoginRequest) (*dto.LoginResponse, int, error)
	GetUserInfo(ctx *fiber.Ctx) (*dto.UserInfoResponse, int, error)
}

type authService struct {
	userRepository repository.UserRepository
	pkgCrypto      crypto.Crypto
}

func NewAuthService(
	userRepository repository.UserRepository,
	pkgCrypto crypto.Crypto,
) AuthService {
	return &authService{
		userRepository: userRepository,
		pkgCrypto:      pkgCrypto,
	}
}

func (s *authService) Login(ctx *fiber.Ctx, request dto.LoginRequest) (*dto.LoginResponse, int, error) {
	user, err := s.userRepository.FindByUniqueIdentifier(request.UniqueIdentifier)
	if err != nil && err.Error() == "record not found" {
		log.Error("User not found.")
		return nil, fiber.StatusNotFound, errors.New(i18n.CreateMsg(ctx, messages.UserNotFound))
	}
	if err != nil {
		log.Error("Error finding a authware: ", err)
		return nil, fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	// Check the password
	if !s.pkgCrypto.CheckPasswordHash(request.Password, user.Password) {
		log.Error("Password is incorrect.")
		return nil, fiber.StatusUnauthorized, errors.New(i18n.CreateMsg(ctx, messages.PasswordIncorrect))
	}

	jwtPayload := authware.JWTClaimsPayload{
		ID:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: fmt.Sprintf("+%d", user.PhoneNumber),
		Email:       user.Email,
	}

	// Generate JWT Token
	token, err := authware.GenerateJwtToken(jwtPayload, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return nil, http.StatusBadGateway, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	response := &dto.LoginResponse{
		Token: token,
	}

	return response, fiber.StatusOK, nil
}

func (s *authService) GetUserInfo(ctx *fiber.Ctx) (*dto.UserInfoResponse, int, error) {
	currentUser, err := authware.GetCurrentUser(ctx)
	if err != nil {
		log.Error(err)
		return nil, fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	user, err := s.userRepository.FindByID(currentUser.Id)
	if err != nil {
		log.Error("Error finding a authware: ", err)
		return nil, fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	response := &dto.UserInfoResponse{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: fmt.Sprintf("+%d", user.PhoneNumber),
		IsActive:    user.IsActive,
	}

	return response, fiber.StatusOK, nil
}
