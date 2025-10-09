package errx

import (
	"errors"
	"net/http"
)

var (
	ErrUserNotFound        = NewNotFoundError("User not found")
	ErrEmailAlreadyExists  = NewConflictError("Email already exists")
	ErrInvalidCredentials  = NewUnauthorizedError("Invalid credentials")
	ErrMissingAuthorizationHeader = NewUnauthorizedError("Missing bearer token")
	ErrInvalidAuthorizationHeader = NewUnauthorizedError("Invalid authorization header")
	ErrInvalidBearerToken  = NewUnauthorizedError("Invalid bearer token")
	ErrInvalidUserIDFormat  = NewUnauthorizedError("Invalid user ID format in token")
	ErrUnauthorized        = NewUnauthorizedError("Unauthorized")
	ErrDatabaseError       = NewInternalServerError("Database error")
	ErrRedisError          = NewInternalServerError("Redis error")
	ErrInternalServer      = NewInternalServerError("Internal server error")
	ErrTransactionNotFound = NewNotFoundError("Transaction not found")
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}