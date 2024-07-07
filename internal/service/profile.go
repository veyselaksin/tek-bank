package service

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"tek-bank/cmd/api/middleware/authware"
	"tek-bank/internal/db/repository"
	"tek-bank/internal/dto"
	"tek-bank/internal/i18n"
	"tek-bank/internal/i18n/messages"
)

type ProfileService interface {
	MyProfile(ctx *fiber.Ctx) (*dto.GetProfileResponse, int, error)
	MyTransferHistory(ctx *fiber.Ctx, accountNumber int64) ([]dto.GetTransferHistoryResponse, int, error)
}

type profileService struct {
	accountRepository  repository.AccountRepository
	transferRepository repository.TransferHistoryRepository
	userRepository     repository.UserRepository
}

func NewProfileService(
	accountRepository repository.AccountRepository,
	transferRepository repository.TransferHistoryRepository,
	userRepository repository.UserRepository,
) ProfileService {
	return &profileService{
		accountRepository:  accountRepository,
		transferRepository: transferRepository,
		userRepository:     userRepository,
	}
}

func (s *profileService) MyProfile(ctx *fiber.Ctx) (*dto.GetProfileResponse, int, error) {
	currentUser, err := authware.GetCurrentUser(ctx)
	if err != nil {
		log.Error(err)
		return nil, fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	accounts, err := s.accountRepository.FindByOwnerId(currentUser.Id)
	if err != nil {
		return nil, fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.AccountNotFound))
	}

	user, err := s.userRepository.FindByID(currentUser.Id)
	if err != nil {
		log.Error("Error finding a authware: ", err)
		return nil, fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	var accountItems []dto.AccountItem
	for _, account := range accounts {
		accountItems = append(accountItems, dto.AccountItem{
			Id:            account.Id,
			AccountNumber: account.AccountNumber,
			IBAN:          account.IBAN,
			Balance:       account.Balance,
		})
	}

	response := dto.GetProfileResponse{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: fmt.Sprintf("+%d", user.PhoneNumber),
		Email:       user.Email,
		AccountList: accountItems,
	}

	return &response, fiber.StatusOK, nil
}

func (s *profileService) MyTransferHistory(ctx *fiber.Ctx, accountNumber int64) ([]dto.GetTransferHistoryResponse, int, error) {
	currentUser := ctx.Locals("user").(authware.CurrentUser)

	account, err := s.accountRepository.FindByAccountNumber(accountNumber)
	if err != nil {
		log.Error(err)
		return nil, fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.AccountNotFound))
	}

	if account.OwnerId != currentUser.Id {
		log.Error("Unauthorized access.")
		return nil, fiber.StatusUnauthorized, errors.New(i18n.CreateMsg(ctx, messages.Unauthorized))
	}

	transferHistory, err := s.transferRepository.FetchByAccountNumber(accountNumber)
	if err != nil {
		log.Error(err)
		return nil, fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.AccountNotFound))
	}

	var response []dto.GetTransferHistoryResponse
	for _, transfer := range transferHistory {

		if transfer.To == accountNumber && transfer.IsFee {
			continue
		}

		if transfer.From == accountNumber {
			transfer.Amount = -transfer.Amount
		}
		response = append(response, dto.GetTransferHistoryResponse{
			Id:     transfer.Id,
			From:   transfer.From,
			To:     transfer.To,
			Note:   transfer.Note,
			Amount: transfer.Amount,
		})
	}

	return response, fiber.StatusOK, nil
}
