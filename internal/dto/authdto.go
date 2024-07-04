package dto

type RegisterRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginRequest struct {
	UniqueIdentifier string `json:"unique_identifier"`
	Password         string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
