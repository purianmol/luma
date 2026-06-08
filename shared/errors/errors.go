package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func Wrap(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, message)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, message)
}

func UnprocessableEntity(message string) *AppError {
	return New(http.StatusUnprocessableEntity, message)
}

func TooManyRequests(message string) *AppError {
	return New(http.StatusTooManyRequests, message)
}

func Internal(message string, err error) *AppError {
	return Wrap(http.StatusInternalServerError, message, err)
}

func HTTPCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return http.StatusInternalServerError
}

func ClientMessage(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return "an unexpected error occurred"
}

func IsNotFound(err error) bool {
	return HTTPCode(err) == http.StatusNotFound
}

func IsConflict(err error) bool {
	return HTTPCode(err) == http.StatusConflict
}

func IsUnAuthorized(err error) bool {
	return HTTPCode(err) == http.StatusUnauthorized
}

var Is = errors.Is

var As = errors.As
