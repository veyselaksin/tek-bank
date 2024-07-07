package dto

type RegisterAccountRequest struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	ISOCountryCode string `json:"iso_country_code"`
	IdentityNumber int64  `json:"identity_number"`
	PhoneNumber    uint64 `json:"phone_number"`
}

type CreateNewAccountRequest struct {
	UserId         string `json:"user_id"`
	ISOCountryCode string `json:"iso_country_code"`
}

type CreateNewAccountResponse struct {
	UserId        string  `json:"user_id"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	AccountNumber int64   `json:"account_number"`
	IBAN          string  `json:"iban"`
	Balance       float64 `json:"balance"`
	IsActive      bool    `json:"is_active"`
}

type AddMoneyRequest struct {
	AccountNumber int64   `json:"-"`
	Amount        float64 `json:"amount"`
}

type AddMoneyResponse struct {
	CustomerNumber int64   `json:"customer_number"`
	Balance        float64 `json:"balance"`
}

type TransferMoneyRequest struct {
	Note              string  `json:"note"`
	Amount            float64 `json:"amount"`
	FromAccountNumber int64   `json:"from_account_number"`
	ToAccountNumber   int64   `json:"to_account_number"`
}
