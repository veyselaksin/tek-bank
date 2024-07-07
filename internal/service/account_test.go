package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
	"tek-bank/internal/db/models"
	"tek-bank/internal/dto"
	"tek-bank/internal/i18n"
	"tek-bank/internal/mocks/repository"
	"tek-bank/mocks/converter"
	"tek-bank/mocks/crypto"
	"testing"
)

var mockData = []models.User{
	{
		Id:             "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1b",
		IdentityNumber: 4000000001,
		CustomerNumber: 1000000001,
		Email:          "john.doe@company.com",
		Password:       "password",
		FirstName:      "John",
		LastName:       "Doe",
	},
	{
		Id:             "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1c",
		IdentityNumber: 4000000002,
		CustomerNumber: 1000000002,
		Email:          "jane.doe@company.com",
		Password:       "password",
		FirstName:      "Jane",
		LastName:       "Doe",
	},
}

var mockAccountData = []models.Account{
	{
		Id:            "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1b",
		OwnerId:       "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1b",
		AccountNumber: 1000000001,
		IBAN:          "US1000000001",
		Balance:       0,
		CreatedBy:     "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1b",
		UpdatedBy:     "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1b",
		Owner:         mockData[0],
	},
	{
		Id:            "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1c",
		OwnerId:       "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1c",
		AccountNumber: 1000000002,
		IBAN:          "US1000000002",
		Balance:       0,
		CreatedBy:     "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1c",
		UpdatedBy:     "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1c",
		Owner:         mockData[1],
	},
}

var fiberCtx *fiber.Ctx
var s AccountService

var userRepoMock *repository.MockUserRepository
var accountRepoMock *repository.MockAccountRepository
var transferRepoMock *repository.MockTransferHistoryRepository
var pkgCryptoMock *crypto.MockCrypto
var pkgConverterMock *converter.MockConverter

func setupAccountTest(t *testing.T) func() {
	ct := gomock.NewController(t)
	defer ct.Finish()

	app := fiber.New()
	fiberCtx = app.AcquireCtx(&fasthttp.RequestCtx{})

	// Assign language to fiber context header
	fiberCtx.Request().Header.Set("Accept-Language", "en")

	i18n.InitBundle("./../i18n/languages")

	userRepoMock = repository.NewMockUserRepository(ct)
	accountRepoMock = repository.NewMockAccountRepository(ct)
	transferRepoMock = repository.NewMockTransferHistoryRepository(ct)
	pkgCryptoMock = crypto.NewMockCrypto(ct)
	pkgConverterMock = converter.NewMockConverter(ct)

	s = NewAccountService(accountRepoMock, userRepoMock, transferRepoMock, pkgCryptoMock, pkgConverterMock)
	return func() {
		s = nil
		defer ct.Finish()
	}
}

func TestAccountService_RegisterAccount_Success(t *testing.T) {
	teardown := setupAccountTest(t)
	defer teardown()

	request := dto.RegisterAccountRequest{
		FirstName:      "Peter",
		LastName:       "Parker",
		Email:          "peter.parker@company.com",
		ISOCountryCode: "US",
		IdentityNumber: 4000000003,
		PhoneNumber:    1234567890,
	}

	// Test logic here
	userRepoMock.EXPECT().FindByEmail(request.Email).Return(&models.User{}, errors.New("record not found")).Times(1)
	pkgCryptoMock.EXPECT().RandomPassword().Return("password").Times(1)
	pkgCryptoMock.EXPECT().HashPassword("password").Return("$2a$10$1Q7Z6z1z1z1z1z1z1z1z1u", nil).Times(1)
	pkgCryptoMock.EXPECT().RandomNumber().Return(int64(1000000003)).Times(2)

	user := models.User{
		IdentityNumber: request.IdentityNumber,
		CustomerNumber: 1000000003,
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		Email:          request.Email,
		PhoneNumber:    request.PhoneNumber,
		Password:       "$2a$10$1Q7Z6z1z1z1z1z1z1z1z1u",
	}

	userRepoMock.EXPECT().Create(user).Return(&user, nil).Times(1)
	pkgCryptoMock.EXPECT().RandomIBAN("US").Return("US1000000003").Times(1)

	account := models.Account{
		OwnerId:       user.Id,
		AccountNumber: 1000000003,
		IBAN:          "US1000000003",
		Balance:       0,
		CreatedBy:     user.Id,
		UpdatedBy:     user.Id,
	}

	accountRepoMock.EXPECT().Create(account).Return(&account, nil).Times(1)

	status, err := s.RegisterAccount(fiberCtx, request)
	if err != nil {
		t.Errorf("Error was not expected: %v", err)
	}

	assert.Equal(t, status, fiber.StatusCreated)
}

func TestAccountService_RegisterAccount_AlreadyExists(t *testing.T) {
	teardown := setupAccountTest(t)
	defer teardown()

	request := dto.RegisterAccountRequest{
		FirstName:      "Peter",
		LastName:       "Parker",
		Email:          "peter.parker@company.com",
		ISOCountryCode: "US",
		IdentityNumber: 4000000003,
		PhoneNumber:    1234567890,
	}

	// Test logic here
	userRepoMock.EXPECT().FindByEmail(request.Email).Return(&models.User{}, nil)

	status, err := s.RegisterAccount(fiberCtx, request)
	if err == nil {
		t.Errorf("Error was expected")
	}

	assert.Equal(t, status, fiber.StatusConflict)
}

func TestAccountService_RegisterAccount_UnexpectedError(t *testing.T) {
	teardown := setupAccountTest(t)
	defer teardown()

	request := dto.RegisterAccountRequest{
		FirstName:      "Peter",
		LastName:       "Parker",
		Email:          "peter.parker@company.com",
		ISOCountryCode: "US",
		IdentityNumber: 4000000003,
		PhoneNumber:    1234567890,
	}

	// Test logic here
	userRepoMock.EXPECT().FindByEmail(request.Email).Return(&models.User{}, errors.New("unexpected error"))

	status, err := s.RegisterAccount(fiberCtx, request)
	if err == nil {
		t.Errorf("Error was expected")
	}

	assert.Equal(t, status, fiber.StatusInternalServerError)
}

func TestAccountService_RegisterAccount_ErrorHashingPassword(t *testing.T) {
	teardown := setupAccountTest(t)
	defer teardown()

	request := dto.RegisterAccountRequest{
		FirstName:      "Peter",
		LastName:       "Parker",
		Email:          "peter.parker@company.com",
		ISOCountryCode: "US",
		IdentityNumber: 4000000003,
		PhoneNumber:    1234567890,
	}

	// Test logic here
	userRepoMock.EXPECT().FindByEmail(request.Email).Return(&models.User{}, errors.New("record not found")).Times(1)
	pkgCryptoMock.EXPECT().RandomPassword().Return("password").Times(1)
	pkgCryptoMock.EXPECT().HashPassword("password").Return("", errors.New("error hashing password")).Times(1)

	status, err := s.RegisterAccount(fiberCtx, request)
	if err == nil {
		t.Errorf("Error was expected")
	}

	assert.Equal(t, status, fiber.StatusInternalServerError)
}

func TestAccountService_RegisterAccount_ErrorCreatingUser(t *testing.T) {
	teardown := setupAccountTest(t)
	defer teardown()

	request := dto.RegisterAccountRequest{
		FirstName:      "Peter",
		LastName:       "Parker",
		Email:          "peter.parker@company.com",
		ISOCountryCode: "US",
		IdentityNumber: 4000000003,
		PhoneNumber:    1234567890,
	}

	// Test logic here
	userRepoMock.EXPECT().FindByEmail(request.Email).Return(&models.User{}, errors.New("record not found")).Times(1)
	pkgCryptoMock.EXPECT().RandomPassword().Return("password").Times(1)
	pkgCryptoMock.EXPECT().HashPassword("password").Return("$2a$10$1Q7Z6z1z1z1z1z1z1z1z1u", nil).Times(1)
	pkgCryptoMock.EXPECT().RandomNumber().Return(int64(1000000003)).Times(2)

	user := models.User{
		IdentityNumber: request.IdentityNumber,
		CustomerNumber: 1000000003,
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		Email:          request.Email,
		PhoneNumber:    request.PhoneNumber,
		Password:       "$2a$10$1Q7Z6z1z1z1z1z1z1z1z1u",
	}

	userRepoMock.EXPECT().Create(user).Return(nil, errors.New("error creating user")).Times(1)

	status, err := s.RegisterAccount(fiberCtx, request)
	if err == nil {
		t.Errorf("Error was expected")
	}

	assert.Equal(t, status, fiber.StatusInternalServerError)
}

func TestAccountService_RegisterAccount_ErrorCreatingAccount(t *testing.T) {
	teardown := setupAccountTest(t)
	defer teardown()

	request := dto.RegisterAccountRequest{
		FirstName:      "Peter",
		LastName:       "Parker",
		Email:          "peter.parker@company.com",
		ISOCountryCode: "US",
		IdentityNumber: 4000000003,
		PhoneNumber:    1234567890,
	}

	// Test logic here
	userRepoMock.EXPECT().FindByEmail(request.Email).Return(&models.User{}, errors.New("record not found")).Times(1)
	pkgCryptoMock.EXPECT().RandomPassword().Return("password").Times(1)
	pkgCryptoMock.EXPECT().HashPassword("password").Return("$2a$10$1Q7Z6z1z1z1z1z1z1z1z1u", nil).Times(1)
	pkgCryptoMock.EXPECT().RandomNumber().Return(int64(1000000003)).Times(2)

	user := models.User{
		IdentityNumber: request.IdentityNumber,
		CustomerNumber: 1000000003,
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		Email:          request.Email,
		PhoneNumber:    request.PhoneNumber,
		Password:       "$2a$10$1Q7Z6z1z1z1z1z1z1z1z1u",
	}

	userRepoMock.EXPECT().Create(user).Return(&user, nil).Times(1)
	pkgCryptoMock.EXPECT().RandomIBAN("US").Return("US1000000003").Times(1)

	account := models.Account{
		OwnerId:       user.Id,
		AccountNumber: 1000000003,
		IBAN:          "US1000000003",
		Balance:       0,
		CreatedBy:     user.Id,
		UpdatedBy:     user.Id,
	}

	accountRepoMock.EXPECT().Create(account).Return(nil, errors.New("error creating account")).Times(1)

	status, err := s.RegisterAccount(fiberCtx, request)
	if err == nil {
		t.Errorf("Error was expected")
	}

	assert.Equal(t, status, fiber.StatusInternalServerError)
}

func TestAccountService_CreateNewAccount_Success(t *testing.T) {
	teardown := setupAccountTest(t)
	defer teardown()

	request := dto.CreateNewAccountRequest{
		UserId:         "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1b",
		ISOCountryCode: "US",
	}

	// Test logic here
	userRepoMock.EXPECT().FindByID(request.UserId).Return(&mockData[0], nil).Times(1)
	pkgCryptoMock.EXPECT().RandomIBAN(request.ISOCountryCode).Return("US1000000001").Times(1)
	pkgCryptoMock.EXPECT().RandomNumber().Return(int64(1000000001)).Times(2)

	account := models.Account{
		OwnerId:       request.UserId,
		AccountNumber: 1000000001,
		IBAN:          "US1000000001",
		Balance:       0,
		CreatedBy:     request.UserId,
		UpdatedBy:     request.UserId,
	}

	accountRepoMock.EXPECT().Create(account).Return(&account, nil).Times(1)

	response, status, err := s.CreateNewAccount(fiberCtx, request)
	if err != nil {
		t.Errorf("Error was not expected: %v", err)
	}

	assert.Equal(t, status, fiber.StatusCreated)
	assert.Equal(t, response.AccountNumber, account.AccountNumber)
	assert.Equal(t, response.IBAN, account.IBAN)
	assert.Equal(t, response.Balance, account.Balance)
	assert.Equal(t, response.UserId, account.OwnerId)
}

func TestAccountService_CreateNewAccount_UserNotFound(t *testing.T) {
	teardown := setupAccountTest(t)
	defer teardown()

	request := dto.CreateNewAccountRequest{
		UserId:         "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1b",
		ISOCountryCode: "US",
	}

	// Test logic here
	userRepoMock.EXPECT().FindByID(request.UserId).Return(nil, errors.New("record not found")).Times(1)

	response, status, err := s.CreateNewAccount(fiberCtx, request)
	if err == nil {
		t.Errorf("Error was expected")
	}

	assert.Equal(t, status, fiber.StatusNotFound)
	assert.Nil(t, response)
}

func TestAccountService_CreateNewAccount_UnexpectedError(t *testing.T) {
	teardown := setupAccountTest(t)
	defer teardown()

	request := dto.CreateNewAccountRequest{
		UserId:         "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1b",
		ISOCountryCode: "US",
	}

	// Test logic here
	userRepoMock.EXPECT().FindByID(request.UserId).Return(nil, errors.New("unexpected error")).Times(1)

	response, status, err := s.CreateNewAccount(fiberCtx, request)
	if err == nil {
		t.Errorf("Error was expected")
	}

	assert.Equal(t, status, fiber.StatusInternalServerError)
	assert.Nil(t, response)
}

func TestAccountService_AddMoney_Success(t *testing.T) {
	teardown := setupAccountTest(t)
	defer teardown()

	request := dto.AddMoneyRequest{
		Amount:        100,
		AccountNumber: 1000000001,
	}

	// Test logic here
	accountRepoMock.EXPECT().FindByAccountNumber(request.AccountNumber).Return(&mockAccountData[0], nil).Times(1)
	accountRepoMock.EXPECT().UpdateBalance(mockAccountData[0].Balance+request.Amount, "e7e1b1b0-7f46-4b6d-8b0d-3b6f1b4f1b1b").Return(nil).Times(1)
	accountRepoMock.EXPECT().FindByAccountNumber(request.AccountNumber).Return(&models.Account{
		Id:            mockAccountData[0].Id,
		OwnerId:       mockAccountData[0].OwnerId,
		AccountNumber: mockAccountData[0].AccountNumber,
		IBAN:          mockAccountData[0].IBAN,
		Balance:       mockAccountData[0].Balance + request.Amount,
		CreatedBy:     mockAccountData[0].CreatedBy,
		UpdatedBy:     mockAccountData[0].UpdatedBy,
		Owner:         mockData[0],
	}, nil).Times(1)

	response, status, err := s.AddMoney(fiberCtx, request)
	if err != nil {
		t.Errorf("Error was not expected: %v", err)
	}

	assert.Equal(t, status, fiber.StatusOK)
	assert.Equal(t, response.Balance, mockAccountData[0].Balance+request.Amount)
	assert.Equal(t, response.CustomerNumber, mockData[0].CustomerNumber)
}
