package service

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tek-bank/internal/db/models"
	"tek-bank/internal/db/repository"
	"tek-bank/internal/dto"
	"tek-bank/internal/i18n/messages"
	"tek-bank/pkg/converter"
	"tek-bank/pkg/crypto"
	"tek-bank/pkg/enum"
	"tek-bank/pkg/gomailer"
	"time"
)

type AccountService interface {
	RegisterAccount(ctx context.Context, request dto.RegisterAccountRequest) error
	CreateNewAccount(ctx context.Context, request dto.CreateNewAccountRequest) (*dto.CreateNewAccountResponse, error)
	AddMoney(ctx context.Context, request dto.AddMoneyRequest) (*dto.AddMoneyResponse, error)
	TransferMoney(ctx context.Context, request dto.TransferMoneyRequest) error
	TransferApproval(ctx context.Context, token string) error

	WithTx(trxHandle *gorm.DB) AccountService
}

type accountService struct {
	accountRepository         repository.AccountRepository
	userRepository            repository.UserRepository
	transferHistoryRepository repository.TransferHistoryRepository
	pkgCrypto                 crypto.Crypto
	pkgConverter              converter.Converter
}

func NewAccountService(
	accountRepository repository.AccountRepository,
	userRepository repository.UserRepository,
	transferHistoryRepository repository.TransferHistoryRepository,
	pkgCrypto crypto.Crypto,
	pkgConverter converter.Converter,
) AccountService {
	return &accountService{
		accountRepository:         accountRepository,
		userRepository:            userRepository,
		transferHistoryRepository: transferHistoryRepository,
		pkgCrypto:                 pkgCrypto,
		pkgConverter:              pkgConverter,
	}
}

func (s *accountService) WithTx(trxHandle *gorm.DB) AccountService {
	s.accountRepository = s.accountRepository.WithTx(trxHandle)
	s.userRepository = s.userRepository.WithTx(trxHandle)
	return s

}

func (s *accountService) RegisterAccount(ctx context.Context, request dto.RegisterAccountRequest) error {

	// Check if the authware already exists
	_, err := s.userRepository.FindByEmail(request.Email)
	if err == nil {
		return errors.New(messages.UserAlreadyExists)
	}

	if err.Error() != "record not found" {
		return errors.New(messages.UnexpectedError)
	}

	randomPassword := s.pkgCrypto.RandomPassword()

	// Hash the password
	hashedPassword, err := s.pkgCrypto.HashPassword(randomPassword)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	// Register a new authware
	user := models.User{
		IdentityNumber: request.IdentityNumber,
		CustomerNumber: s.pkgCrypto.RandomNumber(),
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		Email:          request.Email,
		PhoneNumber:    request.PhoneNumber,
		Password:       hashedPassword,
	}

	createdUser, err := s.userRepository.Create(user)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	// Create a new account for the user
	account := models.Account{
		IBAN:          s.pkgCrypto.RandomIBAN(request.ISOCountryCode),
		OwnerId:       createdUser.Id,
		AccountNumber: s.pkgCrypto.RandomNumber(),
		Balance:       0,
		CreatedBy:     createdUser.Id,
		UpdatedBy:     createdUser.Id,
	}

	_, err = s.accountRepository.Create(account)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	errCh := make(chan error)

	go func() {
		// Send the password to the user's email
		content := gomailer.Content{
			Subject: "TEK Bank - First Password",
			Body:    "Welcome to TEK Bank! Your first password is: " + randomPassword + ". Please change your password after you login.",
			To:      []string{user.Email},
		}

		err = gomailer.SendMail(content)
		if err != nil {
			errCh <- err
		}
	}()

	// Use a select statement to wait for an error or a timeout
	select {
	case err := <-errCh:
		if err != nil {
			return errors.New(messages.UnexpectedError)
		}
	case <-time.After(2 * time.Second):
		break
	}

	// Close the error channel
	close(errCh)

	return nil
}

// CreateNewAccount creates a new account for the registered user
func (s *accountService) CreateNewAccount(ctx context.Context, request dto.CreateNewAccountRequest) (*dto.CreateNewAccountResponse, error) {
	// Check if the user exists
	_, err := s.userRepository.FindByID(request.UserId)
	if err != nil && err.Error() == "record not found" {
		return nil, errors.New(messages.UserNotFound)
	}

	if err != nil {
		return nil, errors.New(messages.UnexpectedError)
	}

	// Create a random IBAN for the user
	iban := s.pkgCrypto.RandomIBAN(request.ISOCountryCode)

	// Create a new account for the user
	account := models.Account{
		IBAN:          iban,
		OwnerId:       request.UserId,
		Balance:       0,
		AccountNumber: s.pkgCrypto.RandomNumber(),
		CreatedBy:     request.UserId,
		UpdatedBy:     request.UserId,
	}

	createdAccount, err := s.accountRepository.Create(account)
	if err != nil {
		return nil, errors.New(messages.UnexpectedError)
	}

	response := &dto.CreateNewAccountResponse{
		UserId:        createdAccount.OwnerId,
		IBAN:          createdAccount.IBAN,
		AccountNumber: s.pkgCrypto.RandomNumber(),
		FirstName:     createdAccount.Owner.FirstName,
		LastName:      createdAccount.Owner.LastName,
		Balance:       createdAccount.Balance,
		IsActive:      createdAccount.IsActive,
	}

	return response, nil
}

func (s *accountService) AddMoney(ctx context.Context, request dto.AddMoneyRequest) (*dto.AddMoneyResponse, error) {
	// Check if the account exists
	account, err := s.accountRepository.FindByAccountNumber(request.AccountNumber)
	if err != nil {
		return nil, errors.New(messages.AccountNotFound)
	}

	// Add money to the account
	err = s.accountRepository.UpdateBalance(account.Balance+request.Amount, account.Id)
	if err != nil {
		return nil, errors.New(messages.UnexpectedError)
	}

	updatedAccount, err := s.accountRepository.FindByAccountNumber(request.AccountNumber)
	if err != nil {
		return nil, errors.New(messages.AccountNotFound)
	}

	response := &dto.AddMoneyResponse{
		CustomerNumber: updatedAccount.Owner.CustomerNumber,
		Balance:        updatedAccount.Balance,
	}

	return response, nil
}

func (s *accountService) TransferMoney(ctx context.Context, request dto.TransferMoneyRequest) error {
	// Check if the sender account exists
	senderAccount, err := s.accountRepository.FindByAccountNumber(request.FromAccountNumber)
	if err != nil {
		return errors.New(messages.AccountNotFound)
	}

	// Check if the sender account has enough balance
	totalAmount := request.Amount + enum.TransferFee
	if senderAccount.Balance < totalAmount {
		return errors.New(messages.InSufficientBalance)
	}

	// Create a token for the transaction approval
	token, err := s.pkgCrypto.GenerateToken(32)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	// Save token to Redis
	value := struct {
		FromAccountNumber int64
		ToAccountNumber   int64
		Amount            float64
		TransactionFee    float64
		Note              string
		Token             string
	}{
		FromAccountNumber: request.FromAccountNumber,
		ToAccountNumber:   request.ToAccountNumber,
		Amount:            request.Amount,
		Note:              request.Note,
		TransactionFee:    enum.TransferFee,
		Token:             token,
	}

	content, err := s.pkgConverter.Stos(value)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	// One hour expiration time
	err = s.accountRepository.SetToken(context.Background(), token, *content)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	transferApprovalLink := fmt.Sprintf("http://localhost/v1/account/transfer-approval?token=%s", token)
	var body string = `
			<body>
				<p>Your Account Number: <strong>` + fmt.Sprint(request.FromAccountNumber) + `</strong></p>
				<p>Receiver Account Number: <strong>` + fmt.Sprint(request.ToAccountNumber) + `</strong></p>
				<p>Amount: <strong>` + fmt.Sprint(request.Amount) + `</strong></p>
				<p>Fee: <strong>` + fmt.Sprint(enum.TransferFee) + `</strong></p>
				<p>You have a new transfer request. Please click the link below to approve the transaction.</p>
				<p><a href="` + transferApprovalLink + `">` + transferApprovalLink + `</a></p>
				<p>If you did not request a transfer, please ignore this email.</p>
				<br>
				<p>Best Regards,</p>
			</body>
	`

	errCh := make(chan error)

	go func() {
		// Send the password to the user's email
		content := gomailer.Content{
			Subject: "TEK Bank - Transfer Approval",
			Body:    body,
			To:      []string{senderAccount.Owner.Email},
		}

		err = gomailer.SendMail(content)
		if err != nil {
			errCh <- err
		}
	}()

	// Use a select statement to wait for an error or a timeout
	select {
	case err := <-errCh:
		if err != nil {
			return errors.New(messages.UnexpectedError)
		}
	case <-time.After(2 * time.Second):
		break
	}

	// Close the error channel
	close(errCh)

	return nil
}

// TransferApproval approves the transaction
func (s *accountService) TransferApproval(ctx context.Context, token string) error {
	// Get the token from Redis
	value, err := s.accountRepository.GetToken(context.Background(), token)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	var content struct {
		FromAccountNumber int64
		ToAccountNumber   int64
		Amount            float64
		TransactionFee    float64
		Note              string
		Token             string
	}

	err = s.pkgConverter.Stom(*value, &content)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	// Check if the sender account exists
	senderAccount, err := s.accountRepository.FindByAccountNumber(content.FromAccountNumber)
	if err != nil {
		return errors.New(messages.AccountNotFound)
	}

	// Check if the receiver account exists
	receiverAccount, err := s.accountRepository.FindByAccountNumber(content.ToAccountNumber)
	if err != nil {
		return errors.New(messages.AccountNotFound)
	}

	// Check if the sender account has enough balance
	totalAmount := content.Amount + content.TransactionFee
	if senderAccount.Balance < totalAmount {
		return errors.New(messages.InSufficientBalance)
	}

	// Deduct the amount from the sender account
	err = s.accountRepository.UpdateBalance(senderAccount.Balance-totalAmount, senderAccount.Id)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	// Add the amount to the receiver account
	err = s.accountRepository.UpdateBalance(receiverAccount.Balance+content.Amount, receiverAccount.Id)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	// Multiple insertions to transfer history table
	var transferHistories []models.TransferHistory
	transferHistories = append(transferHistories, models.TransferHistory{
		From:      content.FromAccountNumber,
		To:        content.ToAccountNumber,
		Amount:    content.Amount,
		Note:      content.Note,
		CreatedBy: senderAccount.OwnerId,
		UpdatedBy: senderAccount.OwnerId,
	})

	transferHistories = append(transferHistories, models.TransferHistory{
		From:      content.FromAccountNumber,
		To:        content.ToAccountNumber,
		Amount:    content.TransactionFee,
		Note:      "Transaction Fee",
		IsFee:     true,
		CreatedBy: senderAccount.OwnerId,
		UpdatedBy: senderAccount.OwnerId,
	})

	err = s.transferHistoryRepository.Create(transferHistories)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	// Delete the token from Redis
	err = s.accountRepository.DeleteToken(context.Background(), token)
	if err != nil {
		return errors.New(messages.UnexpectedError)
	}

	errCh := make(chan error)

	// Send an email to the sender and receiver
	go func() {
		// Send the password to the user's email
		contentSender := gomailer.Content{
			Subject: "TEK Bank - Transfer Approval",
			Body:    "Your transfer has been successfully completed.",
			To:      []string{senderAccount.Owner.Email},
		}

		err = gomailer.SendMail(contentSender)
		if err != nil {
			errCh <- err
		}

		contentReceiver := gomailer.Content{
			Subject: "TEK Bank - Transfer Approval",
			Body:    "You have received a new transfer.",
			To:      []string{receiverAccount.Owner.Email},
		}

		err = gomailer.SendMail(contentReceiver)
		if err != nil {
			errCh <- err
		}
	}()

	// Use a select statement to wait for an error or a timeout
	select {
	case err := <-errCh:
		if err != nil {
			return errors.New(messages.UnexpectedError)
		}
	case <-time.After(2 * time.Second):
		break
	}

	// Close the error channel
	close(errCh)

	return nil
}
