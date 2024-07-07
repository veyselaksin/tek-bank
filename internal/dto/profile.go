package dto

type AccountItem struct {
	Id            string  `json:"id"`
	AccountNumber int64   `json:"account_number"`
	IBAN          string  `json:"iban"`
	Balance       float64 `json:"balance"`
}

type GetProfileResponse struct {
	Id          string        `json:"id"`
	FirstName   string        `json:"first_name"`
	LastName    string        `json:"last_name"`
	Email       string        `json:"email"`
	PhoneNumber string        `json:"phone_number"`
	AccountList []AccountItem `json:"account_list"`
}

type GetTransferHistoryRequest struct {
	AccountNumber int64 `json:"-"`
}

type GetTransferHistoryResponse struct {
	Id     string  `json:"id"`
	From   int64   `json:"from"`
	To     int64   `json:"to"`
	Note   string  `json:"note"`
	Amount float64 `json:"amount"`
}
