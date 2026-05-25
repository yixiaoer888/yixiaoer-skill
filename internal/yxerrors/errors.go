package yxerrors

import "fmt"

const (
	UsageErr  = "YIXIAOER_USAGE_ERR"
	AuthErr   = "YIXIAOER_AUTH_ERR"
	RemoteErr = "YIXIAOER_REMOTE_ERR"
)

type Error struct {
	Code       string
	Message    string
	Details    interface{}
	Suggestion string
}

func (e *Error) Error() string {
	return e.Message
}

func New(code, message string, details interface{}) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func Usage(message string, details interface{}) *Error {
	return New(UsageErr, message, details)
}

func Auth(message string) *Error {
	return New(AuthErr, message, nil)
}

func Remote(message string, details interface{}) *Error {
	return New(RemoteErr, message, details)
}

func WrapRemote(format string, args ...interface{}) *Error {
	return Remote(fmt.Sprintf(format, args...), nil)
}
