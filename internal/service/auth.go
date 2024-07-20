package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"tek-bank/cmd/api/middleware/authware"
	"tek-bank/internal/db/repository"
	"tek-bank/internal/dto"
	"tek-bank/internal/i18n/messages"
	"tek-bank/pkg/crypto"
)

type AuthService interface {
	Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error)
	GetUserInfo(ctx context.Context) (*dto.UserInfoResponse, error)
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

func (s *authService) Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepository.FindByUniqueIdentifier(request.UniqueIdentifier)
	if err != nil && err.Error() == "record not found" {
		return nil, errors.New(messages.UserNotFound)
	}
	if err != nil {
		return nil, errors.New(messages.UnexpectedError)
	}

	// Check the password
	if !s.pkgCrypto.CheckPasswordHash(request.Password, user.Password) {
		return nil, errors.New(messages.PasswordIncorrect)
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
		return nil, errors.New(messages.UnexpectedError)
	}

	response := &dto.LoginResponse{
		Token: token,
	}

	return response, nil
}

func (s *authService) GetUserInfo(ctx context.Context) (*dto.UserInfoResponse, error) {
	currentUser, err := authware.GetCurrentUser(ctx)
	if err != nil {
		return nil, errors.New(messages.Unauthorized)
	}

	user, err := s.userRepository.FindByID(currentUser.Id)
	if err != nil {
		return nil, errors.New(messages.UnexpectedError)
	}

	response := &dto.UserInfoResponse{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: fmt.Sprintf("+%d", user.PhoneNumber),
		IsActive:    user.IsActive,
	}

	return response, nil
}
