package output

import (
	"encoding/json"
	"fmt"
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

func Error(w io.Writer, err error, context string) {
	errType := yxerrors.RemoteType
	code := yxerrors.RemoteErr
	message := fmt.Sprintf("Failed to %s", context)
	var details interface{} = err.Error()
	hint := "请依次检查: 1. 技能版本号是否一致; 2. 请求参数是否符合 DTO 规范; 3. 查阅 skills/yixiaoer/references/troubleshooting-guide.md。"
	retryable := false
	nextCommand := ""

	if typed, ok := err.(*yxerrors.Error); ok {
		if typed.Type != "" {
			errType = typed.Type
		}
		code = typed.Code
		message = typed.Message
		if typed.Details != nil {
			details = typed.Details
		}
		if typed.Hint != "" {
			hint = typed.Hint
		} else if typed.Suggestion != "" {
			hint = typed.Suggestion
		}
		retryable = typed.Retryable
		nextCommand = typed.NextCommand
	}

	_ = writeJSON(w, domain.ErrorResponse{
		OK:      false,
		Version: domain.SkillVersion,
		Error: domain.ErrorEnvelope{
			Type:        errType,
			Code:        code,
			Message:     message,
			Hint:        hint,
			Retryable:   retryable,
			NextCommand: nextCommand,
			Details:     details,
		},
	})
}

func writeJSON(w io.Writer, value interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
