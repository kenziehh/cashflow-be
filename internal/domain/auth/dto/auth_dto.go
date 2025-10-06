package dto

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string      `json:"access_token"`
	User  UserProfile `json:"user"`
}

type UserProfile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}