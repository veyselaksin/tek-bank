package service

import (
	"context"
	"errors"
	"fmt"
	"tek-bank/cmd/api/middleware/authware"
	"tek-bank/internal/db/repository"
	"tek-bank/internal/dto"
	"tek-bank/internal/i18n/messages"
)

type ProfileService interface {
	MyProfile(ctx context.Context) (*dto.GetProfileResponse, error)
	MyTransferHistory(ctx context.Context, accountNumber int64) ([]dto.GetTransferHistoryResponse, error)
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

func (s *profileService) MyProfile(ctx context.Context) (*dto.GetProfileResponse, error) {
	currentUser, err := authware.GetCurrentUser(ctx)
	if err != nil {
		return nil, errors.New(messages.UnexpectedError)
	}

	accounts, err := s.accountRepository.FindByOwnerId(currentUser.Id)
	if err != nil {
		return nil, errors.New(messages.AccountNotFound)
	}

	user, err := s.userRepository.FindByID(currentUser.Id)
	if err != nil {
		return nil, errors.New(messages.UnexpectedError)
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

	return &response, nil
}

func (s *profileService) MyTransferHistory(ctx context.Context, accountNumber int64) ([]dto.GetTransferHistoryResponse, error) {
	currentUser, err := authware.GetCurrentUser(ctx)
	if err != nil {
		return nil, errors.New(messages.Unauthorized)
	}

	account, err := s.accountRepository.FindByAccountNumber(accountNumber)
	if err != nil {
		return nil, errors.New(messages.AccountNotFound)
	}

	if account.OwnerId != currentUser.Id {
		return nil, errors.New(messages.Unauthorized)
	}

	transferHistory, err := s.transferRepository.FetchByAccountNumber(accountNumber)
	if err != nil {
		return nil, errors.New(messages.AccountNotFound)
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

	return response, nil
}
