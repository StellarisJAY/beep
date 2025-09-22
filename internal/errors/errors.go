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
	ErrInternalServer ErrorCode = 500
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

func (e *ServiceError) WithDetails(details any) *ServiceError {
	e.Details = details
	return e
}

func NewBadRequestError(message string) *ServiceError {
	return &ServiceError{
		Code:    ErrBadRequest,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *ServiceError {
	return &ServiceError{
		Code:    ErrUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *ServiceError {
	return &ServiceError{
		Code:    ErrForbidden,
		Message: message,
	}
}

func NewNotFoundError(message string) *ServiceError {
	return &ServiceError{
		Code:    ErrNotFound,
		Message: message,
	}
}

func NewInternalServerError(message string) *ServiceError {
	if message == "" {
		message = "服务器内部错误"
	}
	return &ServiceError{
		Code:    ErrInternalServer,
		Message: message,
	}
}

func AsServiceError(err error) (*ServiceError, bool) {
	var e *ServiceError
	ok := errors.As(err, &e)
	return e, ok
}
