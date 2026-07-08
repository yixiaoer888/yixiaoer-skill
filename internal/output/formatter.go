package output

import (
	"encoding/json"
	"io"

	"github.com/yixiaoer/yixiaoer-skill/internal/domain"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func Success(w io.Writer, action string, data interface{}) error {
	return SuccessWithNotice(w, action, data, nil)
}

func SuccessWithNotice(w io.Writer, action string, data interface{}, notice interface{}) error {
	return writeJSON(w, domain.SuccessResponse{
		OK:      true,
		Action:  action,
		Version: domain.SkillVersion,
		Data:    data,
		Notice:  notice,
	})
}

func Error(w io.Writer, err error, context string) int {
	typed := yxerrors.Normalize(err, context)
	if typed == nil {
		return yxerrors.ExitOK
	}
	errType := typed.Type
	code := typed.Code
	message := typed.Message
	details := typed.Details
	hint := typed.Hint
	if hint == "" {
		hint = typed.Suggestion
	}
	retryable := typed.Retryable
	nextCommand := typed.NextCommand
	category := typed.Category
	if category == "" {
		category = errType
	}

	_ = writeJSON(w, domain.ErrorResponse{
		OK:      false,
		Version: domain.SkillVersion,
		Error: domain.ErrorEnvelope{
			Type:        errType,
			Code:        code,
			Category:    category,
			Message:     message,
			Hint:        hint,
			Retryable:   retryable,
			NextCommand: nextCommand,
			Details:     details,
		},
	})
	return typed.ProcessExitCode()
}

func writeJSON(w io.Writer, value interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
