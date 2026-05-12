package errors

import "fmt"

type ErrorCode string

const (
	ValidationError ErrorCode = "VALIDATION_ERROR"
	AuthError       ErrorCode = "AUTH_ERROR"
	NotFoundError   ErrorCode = "NOT_FOUND"
	ConflictError   ErrorCode = "CONFLICT"
	InternalError   ErrorCode = "INTERNAL_ERROR"
	ForbiddenError  ErrorCode = "FORBIDDEN"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Details error
}

func (e *AppError) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewValidationError(message string, details error) *AppError {
	return &AppError{
		Code:    ValidationError,
		Message: message,
		Details: details,
	}
}

func NewAuthError(message string, details error) *AppError {
	return &AppError{
		Code:    AuthError,
		Message: message,
		Details: details,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    NotFoundError,
		Message: message,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    ConflictError,
		Message: message,
	}
}

func NewInternalError(message string, details error) *AppError {
	return &AppError{
		Code:    InternalError,
		Message: message,
		Details: details,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    ForbiddenError,
		Message: message,
	}
}
