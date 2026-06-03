package yxerrors

import "fmt"

const (
	UsageErr  = "YIXIAOER_USAGE_ERR"
	AuthErr   = "YIXIAOER_AUTH_ERR"
	RemoteErr = "YIXIAOER_REMOTE_ERR"

	ValidationType = "validation_error"
	AuthType       = "auth_error"
	RemoteType     = "remote_error"
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
