package dto

type LoginRequest struct {
	UniqueIdentifier string `json:"unique_identifier"`
	Password         string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UserInfoResponse struct {
	Id          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	IsActive    bool   `json:"is_active"`
}
