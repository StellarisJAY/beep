package errors

import (
	"errors"
	"fmt"
)

// ErrorCode 系统错误码
type ErrorCode int

// 常用错误码
const (
	ErrBadRequest     ErrorCode = 400
	ErrUnauthorized   ErrorCode = 401
	ErrForbidden      ErrorCode = 403
	ErrNotFound       ErrorCode = 404
	ErrConflict       ErrorCode = 409
	ErrInternalServer ErrorCode = 500
	ErrLoginFailed    ErrorCode = 401001
)

// ServiceError 自定义错误类型
type ServiceError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("error code: %d, error message: %s", e.Code, e.Message)
}

func NewError(code ErrorCode, msg string, details any) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: msg,
		Details: details,
	}
}

func (e *ServiceError) WithDetails(details any) *ServiceError {
	e.Details = details
	return e
}

func NewBadRequestError(message string, details any) *ServiceError {
	if message == "" {
		message = "bad request"
	}
	return &ServiceError{
		Code:    ErrBadRequest,
		Message: message,
		Details: details,
	}
}

func NewUnauthorizedError(message string, details any) *ServiceError {
	if message == "" {
		message = "unauthorized"
	}
	return &ServiceError{
		Code:    ErrUnauthorized,
		Message: message,
		Details: details,
	}
}

func NewForbiddenError(message string, details any) *ServiceError {
	if message == "" {
		message = "forbidden"
	}
	return &ServiceError{
		Code:    ErrForbidden,
		Message: message,
		Details: details,
	}
}

func NewNotFoundError(message string, details any) *ServiceError {
	if message == "" {
		message = "not found"
	}
	return &ServiceError{
		Code:    ErrNotFound,
		Message: message,
		Details: details,
	}
}

func NewConflictError(message string, details any) *ServiceError {
	if message == "" {
		message = "conflict"
	}
	return &ServiceError{
		Code:    ErrConflict,
		Message: message,
		Details: details,
	}
}

func NewInternalServerError(message string, detail any) *ServiceError {
	if message == "" {
		message = "internal server error"
	}
	return &ServiceError{
		Code:    ErrInternalServer,
		Message: message,
		Details: detail,
	}
}

func AsServiceError(err error) (*ServiceError, bool) {
	var e *ServiceError
	ok := errors.As(err, &e)
	return e, ok
}
