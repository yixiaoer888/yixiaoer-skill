package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yixiaoer/yixiaoer-skill/internal/domain"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func Success(w io.Writer, action string, data interface{}) error {
	return writeJSON(w, domain.SuccessResponse{
		Success: true,
		Action:  action,
		Version: domain.SkillVersion,
		Data:    data,
	})
}

func Error(w io.Writer, err error, context string) {
	code := yxerrors.RemoteErr
	message := fmt.Sprintf("Failed to %s", context)
	var details interface{} = err.Error()
	suggestion := "请依次检查: 1. 技能版本号是否一致; 2. 请求参数是否符合 DTO 规范; 3. 查阅 docs/troubleshooting-guide.md。"

	if typed, ok := err.(*yxerrors.Error); ok {
		code = typed.Code
		message = typed.Message
		if typed.Details != nil {
			details = typed.Details
		}
		if typed.Suggestion != "" {
			suggestion = typed.Suggestion
		}
	}

	_ = writeJSON(w, domain.ErrorResponse{
		Success:    false,
		ErrorCode:  code,
		Message:    message,
		Details:    details,
		Suggestion: suggestion,
	})
}

func writeJSON(w io.Writer, value interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
