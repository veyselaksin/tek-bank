package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"tek-bank/internal/db/models"
	"tek-bank/internal/db/repository"
	"tek-bank/internal/dto"
	"tek-bank/internal/i18n"
	"tek-bank/internal/i18n/messages"
	"tek-bank/pkg/converter"
	"tek-bank/pkg/crypto"
	"tek-bank/pkg/enum"
	"tek-bank/pkg/gomailer"
)

type AccountService interface {
	RegisterAccount(ctx *fiber.Ctx, request dto.RegisterAccountRequest) (int, error)
	CreateNewAccount(ctx *fiber.Ctx, request dto.CreateNewAccountRequest) (*dto.CreateNewAccountResponse, int, error)
	AddMoney(ctx *fiber.Ctx, request dto.AddMoneyRequest) (*dto.AddMoneyResponse, int, error)
	TransferMoney(ctx *fiber.Ctx, request dto.TransferMoneyRequest) (int, error)
	TransferApproval(ctx *fiber.Ctx, token string) (int, error)

	WithTx(trxHandle *gorm.DB) AccountService
}

type accountService struct {
	accountRepository         repository.AccountRepository
	userRepository            repository.UserRepository
	transferHistoryRepository repository.TransferHistoryRepository
}

func NewAccountService(
	accountRepository repository.AccountRepository,
	userRepository repository.UserRepository,
	transferHistoryRepository repository.TransferHistoryRepository,
) AccountService {
	return &accountService{
		accountRepository:         accountRepository,
		userRepository:            userRepository,
		transferHistoryRepository: transferHistoryRepository,
	}
}

func (s *accountService) WithTx(trxHandle *gorm.DB) AccountService {
	s.accountRepository = s.accountRepository.WithTx(trxHandle)
	s.userRepository = s.userRepository.WithTx(trxHandle)
	return s

}

func (s *accountService) RegisterAccount(ctx *fiber.Ctx, request dto.RegisterAccountRequest) (int, error) {

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

	randomPassword := crypto.RandomPassword()

	// Hash the password
	hashedPassword, err := crypto.HashPassword(randomPassword)
	if err != nil {
		log.Error("Error hashing the password: ", err)
		return fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	// Register a new authware
	user := models.User{
		IdentityNumber: request.IdentityNumber,
		CustomerNumber: crypto.RandomNumber(),
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		Email:          request.Email,
		PhoneNumber:    request.PhoneNumber,
		Password:       hashedPassword,
	}

	createdUser, err := s.userRepository.Create(user)
	if err != nil {
		log.Error("Error creating a new authware: ", err)
		return fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	// Create a new account for the user
	account := models.Account{
		IBAN:          crypto.RandomIBAN(request.ISOCountryCode),
		OwnerId:       createdUser.Id,
		AccountNumber: crypto.RandomNumber(),
		Balance:       0,
		CreatedBy:     createdUser.Id,
		UpdatedBy:     createdUser.Id,
	}

	_, err = s.accountRepository.Create(account)
	if err != nil {
		log.Error("Error creating a new account: ", err)
		return fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	go func() {
		// Send the password to the user's email
		content := gomailer.Content{
			Subject: "TEK Bank - First Password",
			Body:    "Welcome to TEK Bank! Your first password is: " + randomPassword + ". Please change your password after you login.",
			To:      []string{user.Email},
		}

		err = gomailer.SendMail(content)
		if err != nil {
			log.Error("Error sending the email: ", err)
		}
	}()

	return fiber.StatusCreated, nil
}

// CreateNewAccount creates a new account for the registered user
func (s *accountService) CreateNewAccount(ctx *fiber.Ctx, request dto.CreateNewAccountRequest) (*dto.CreateNewAccountResponse, int, error) {
	// Check if the user exists

	// Create a random IBAN for the user
	iban := crypto.RandomIBAN(request.ISOCountryCode)

	// Create a new account for the user
	account := models.Account{
		IBAN:          iban,
		OwnerId:       request.UserId,
		Balance:       0,
		AccountNumber: crypto.RandomNumber(),
		CreatedBy:     request.UserId,
		UpdatedBy:     request.UserId,
	}

	createdAccount, err := s.accountRepository.Create(account)
	if err != nil {
		log.Error("Error creating a new account: ", err)
		return nil, fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	fmt.Println(createdAccount)

	response := &dto.CreateNewAccountResponse{
		UserId:        createdAccount.OwnerId,
		IBAN:          createdAccount.IBAN,
		AccountNumber: crypto.RandomNumber(),
		FirstName:     createdAccount.Owner.FirstName,
		LastName:      createdAccount.Owner.LastName,
		Balance:       createdAccount.Balance,
		IsActive:      createdAccount.IsActive,
	}

	return response, fiber.StatusCreated, nil
}

func (s *accountService) AddMoney(ctx *fiber.Ctx, request dto.AddMoneyRequest) (*dto.AddMoneyResponse, int, error) {
	// Check if the account exists
	account, err := s.accountRepository.FindByAccountNumber(request.AccountNumber)
	if err != nil {
		log.Error("Error finding the account: ", err)
		return nil, fiber.StatusNotFound, errors.New(i18n.CreateMsg(ctx, messages.AccountNotFound))
	}

	// Add money to the account
	err = s.accountRepository.UpdateBalance(account.Balance+request.Amount, account.Id)
	if err != nil {
		log.Error("Error updating the account: ", err)
		return nil, fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	updatedAccount, err := s.accountRepository.FindByAccountNumber(request.AccountNumber)
	if err != nil {
		log.Error("Error finding the account: ", err)
		return nil, fiber.StatusNotFound, errors.New(i18n.CreateMsg(ctx, messages.AccountNotFound))
	}

	response := &dto.AddMoneyResponse{
		CustomerNumber: updatedAccount.Owner.CustomerNumber,
		Balance:        updatedAccount.Balance,
	}

	return response, fiber.StatusOK, nil
}

func (s *accountService) TransferMoney(ctx *fiber.Ctx, request dto.TransferMoneyRequest) (int, error) {
	// Check if the sender account exists
	senderAccount, err := s.accountRepository.FindByAccountNumber(request.FromAccountNumber)
	if err != nil {
		log.Error("Error finding the sender account: ", err)
		return fiber.StatusNotFound, errors.New(i18n.CreateMsg(ctx, messages.AccountNotFound))
	}

	// Check if the sender account has enough balance
	totalAmount := request.Amount + enum.TransferFee
	if senderAccount.Balance < totalAmount {
		return fiber.StatusUnprocessableEntity, errors.New(i18n.CreateMsg(ctx, messages.InSufficientBalance))
	}

	// Create a token for the transaction approval
	token, err := crypto.GenerateToken(32)
	if err != nil {
		log.Error(err)
		return fiber.StatusBadGateway, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
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

	content, err := converter.Stos(value)
	if err != nil {
		log.Error(err)
		return fiber.StatusBadGateway, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	// One hour expiration time
	err = s.accountRepository.SetToken(context.Background(), token, *content)
	if err != nil {
		log.Error(err)
		return fiber.StatusBadGateway, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	transferApprovalLink := fmt.Sprintf("%s://%s/v1/account/transfer-approval?token=%s", ctx.Protocol(), ctx.Hostname(), token)
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

	go func() {
		// Send the password to the user's email
		content := gomailer.Content{
			Subject: "TEK Bank - Transfer Approval",
			Body:    body,
			To:      []string{senderAccount.Owner.Email},
		}

		err = gomailer.SendMail(content)
		if err != nil {
			log.Error("Error sending the email: ", err)
		}
	}()

	return fiber.StatusOK, nil
}

// TransferApproval approves the transaction
func (s *accountService) TransferApproval(ctx *fiber.Ctx, token string) (int, error) {
	// Get the token from Redis
	value, err := s.accountRepository.GetToken(context.Background(), token)
	if err != nil {
		log.Error(err)
		return fiber.StatusBadGateway, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	var content struct {
		FromAccountNumber int64
		ToAccountNumber   int64
		Amount            float64
		TransactionFee    float64
		Note              string
		Token             string
	}

	err = converter.Stom(*value, &content)
	if err != nil {
		log.Error(err)
		return fiber.StatusBadGateway, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	// Check if the sender account exists
	senderAccount, err := s.accountRepository.FindByAccountNumber(content.FromAccountNumber)
	if err != nil {
		log.Error("Error finding the sender account: ", err)
		return fiber.StatusNotFound, errors.New(i18n.CreateMsg(ctx, messages.AccountNotFound))
	}

	// Check if the receiver account exists
	receiverAccount, err := s.accountRepository.FindByAccountNumber(content.ToAccountNumber)
	if err != nil {
		log.Error("Error finding the receiver account: ", err)
		return fiber.StatusNotFound, errors.New(i18n.CreateMsg(ctx, messages.AccountNotFound))
	}

	// Check if the sender account has enough balance
	totalAmount := content.Amount + content.TransactionFee
	if senderAccount.Balance < totalAmount {
		return fiber.StatusUnprocessableEntity, errors.New(i18n.CreateMsg(ctx, messages.InSufficientBalance))
	}

	// Deduct the amount from the sender account
	err = s.accountRepository.UpdateBalance(senderAccount.Balance-totalAmount, senderAccount.Id)
	if err != nil {
		log.Error("Error updating the sender account: ", err)
		return fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	// Add the amount to the receiver account
	err = s.accountRepository.UpdateBalance(receiverAccount.Balance+content.Amount, receiverAccount.Id)
	if err != nil {
		log.Error("Error updating the receiver account: ", err)
		return fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
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
		log.Error("Error creating a new transfer history: ", err)
		return fiber.StatusInternalServerError, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

	// Delete the token from Redis
	err = s.accountRepository.DeleteToken(context.Background(), token)
	if err != nil {
		log.Error(err)
		return fiber.StatusBadGateway, errors.New(i18n.CreateMsg(ctx, messages.UnexpectedError))
	}

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
			log.Error("Error sending the email: ", err)
		}

		contentReceiver := gomailer.Content{
			Subject: "TEK Bank - Transfer Approval",
			Body:    "You have received a new transfer.",
			To:      []string{receiverAccount.Owner.Email},
		}

		err = gomailer.SendMail(contentReceiver)
		if err != nil {
			log.Error("Error sending the email: ", err)
		}
	}()

	return fiber.StatusOK, nil
}
