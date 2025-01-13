package palikka

import (
	"context"
	"errors"
	"fmt"
)

const (
	ErrInvalid      = "invalid"
	ErrUnauthorized = "unauthorized"
	ErrForbidden    = "forbidden"
	ErrNotFound     = "not_found"
	ErrInternal     = "internal"
)

type Error struct {
	Cause   error  `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	causeError := ""
	if e.Cause != nil {
		causeError = e.Cause.Error()
	}
	return fmt.Sprintf("palikka error: code=%s, message=%s, cause=%s", e.Code, e.Message, causeError)
}

func ErrorCode(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return ErrInternal
}

func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Message
	}
	return "internal error"
}

func ErrorCause(err error) error {
	var e *Error
	if errors.As(err, &e) && e.Cause != nil {
		return e.Cause
	}
	return err
}

func ErrorCauseRecursive(source error) error {
	var e *Error
	if errors.As(source, &e) {
		if e.Cause != nil {
			return ErrorCauseRecursive(e.Cause)
		}
		return e
	}
	return source
}

func Errorf(code string, message string, args ...interface{}) error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(message, args...),
	}
}

func Errorf2(cause error, code string, message string, args ...interface{}) error {
	return &Error{
		Cause:   cause,
		Code:    code,
		Message: fmt.Sprintf(message, args...),
	}
}

var ReportError = func(ctx context.Context, err error, args ...interface{}) {
	// default no-op
}
