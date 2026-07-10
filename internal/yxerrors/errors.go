package yxerrors

import (
	"errors"
	"fmt"
)

const (
	UsageErr    = "YIXIAOER_USAGE_ERR"
	AuthErr     = "YIXIAOER_AUTH_ERR"
	RemoteErr   = "YIXIAOER_REMOTE_ERR"
	InternalErr = "YIXIAOER_INTERNAL_ERR"

	ValidationType = "validation_error"
	AuthType       = "auth_error"
	RemoteType     = "remote_error"
	InternalType   = "internal_error"
)

const (
	ExitOK         = 0
	ExitRemote     = 1
	ExitValidation = 2
	ExitAuth       = 3
	ExitInternal   = 5
)

type Error struct {
	Type        string
	Code        string
	Category    string
	Message     string
	Details     interface{}
	Hint        string
	Retryable   bool
	NextCommand string
	Suggestion  string
	ExitCode    int
}

func (e *Error) Error() string {
	return e.Message
}

func New(errType, code, message string, details interface{}) *Error {
	return &Error{
		Type:     errType,
		Code:     code,
		Category: errType,
		Message:  message,
		Details:  details,
		ExitCode: defaultExitCode(errType),
	}
}

func Usage(message string, details interface{}) *Error {
	return New(ValidationType, UsageErr, message, details).
		WithRetryable(true)
}

func Auth(message string) *Error {
	return New(AuthType, AuthErr, message, nil)
}

func Remote(message string, details interface{}) *Error {
	return New(RemoteType, RemoteErr, message, details)
}

func Internal(message string, details interface{}) *Error {
	return New(InternalType, InternalErr, message, details)
}

func WrapRemote(format string, args ...interface{}) *Error {
	return Remote(fmt.Sprintf(format, args...), nil)
}

func (e *Error) WithHint(hint string) *Error {
	e.Hint = hint
	e.Suggestion = hint
	return e
}

func (e *Error) WithRetryable(retryable bool) *Error {
	e.Retryable = retryable
	return e
}

func (e *Error) WithNextCommand(nextCommand string) *Error {
	e.NextCommand = nextCommand
	return e
}

func (e *Error) WithCategory(category string) *Error {
	e.Category = category
	return e
}

func (e *Error) WithExitCode(exitCode int) *Error {
	e.ExitCode = exitCode
	return e
}

func (e *Error) ProcessExitCode() int {
	if e == nil {
		return ExitOK
	}
	if e.ExitCode != 0 {
		return e.ExitCode
	}
	return defaultExitCode(e.Type)
}

func Normalize(err error, context string) *Error {
	if err == nil {
		return nil
	}
	var typed *Error
	if errors.As(err, &typed) {
		if typed.ExitCode == 0 {
			typed.ExitCode = defaultExitCode(typed.Type)
		}
		return typed
	}
	message := "Internal command error"
	if context != "" {
		message = fmt.Sprintf("Failed to %s", context)
	}
	return Internal(message, err.Error()).
		WithHint("请查看错误 details；如果是命令参数问题，请运行对应命令的 --help。")
}

func defaultExitCode(errType string) int {
	switch errType {
	case ValidationType:
		return ExitValidation
	case AuthType:
		return ExitAuth
	case RemoteType:
		return ExitRemote
	case InternalType:
		return ExitInternal
	default:
		return ExitRemote
	}
}
