package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"net/http"
	"os"
	"tek-bank/cmd/api/middleware/authware"
	"tek-bank/internal/db/models"
	"tek-bank/internal/db/repository"
	"tek-bank/internal/dto"
	"tek-bank/internal/i18n"
	"tek-bank/internal/i18n/messages"
	"tek-bank/pkg/crypto"
)

type AuthService interface {
	Register(ctx *fiber.Ctx, request dto.RegisterRequest) (int, error)
	Login(ctx *fiber.Ctx, request dto.LoginRequest) (dto.LoginResponse, int, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

// Register is a function to register a new authware
func (s *authService) Register(ctx *fiber.Ctx, request dto.RegisterRequest) (int, error) {
	// Check the passwords match
	if request.Password != request.ConfirmPassword {
		log.Error("Passwords do not match.")
		return fiber.StatusBadRequest, errors.New(i18n.CreateMsg(ctx, messages.PasswordsDoNotMatch))
	}

	// Check if the authware already exists
	_, err := s.userRepository.FindByEmail(request.Email)
	if err == nil {
		log.Error("User already exists.")
		return fiber.StatusConflict, errors.New(i18n.CreateMsg(ctx, messages.UserAlreadyExists))
	}

	if err.Error() != "record not found" {
		log.Error("Error finding a authware: ", err)
		return fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	// Hash the password
	hashedPassword, err := crypto.HashPassword(request.Password)
	if err != nil {
		log.Error("Error hashing the password: ", err)
		return fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	// Register a new authware
	user := models.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  hashedPassword,
	}

	err = s.userRepository.Create(user)
	if err != nil {
		log.Error("Error creating a new authware: ", err)
		return fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	return fiber.StatusOK, nil
}

// Login is a function to login a authware
func (s *authService) Login(ctx *fiber.Ctx, request dto.LoginRequest) (dto.LoginResponse, int, error) {
	user, err := s.userRepository.FindByEmail(request.UniqueIdentifier)
	if err != nil && err.Error() == "record not found" {
		log.Error("User not found.")
		return dto.LoginResponse{}, fiber.StatusNotFound, errors.New(i18n.CreateMsg(ctx, messages.UserNotFound))
	}
	if err != nil {
		log.Error("Error finding a authware: ", err)
		return dto.LoginResponse{}, fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	// Check the password
	if !crypto.CheckPasswordHash(request.Password, user.Password) {
		log.Error("Password is incorrect.")
		return dto.LoginResponse{}, fiber.StatusUnauthorized, errors.New(i18n.CreateMsg(ctx, messages.PasswordIncorrect))
	}

	jwtPayload := authware.JWTClaimsPayload{
		ID:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
	}

	// Generate JWT Token
	token, err := authware.GenerateJwtToken(jwtPayload, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return dto.LoginResponse{}, http.StatusBadGateway, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	response := dto.LoginResponse{
		Token: token,
	}

	return response, fiber.StatusOK, nil
}
