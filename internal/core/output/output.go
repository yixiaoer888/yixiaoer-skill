package output

import (
	"io"

	base "github.com/yixiaoer/yixiaoer-skill/internal/output"
)

func Success(w io.Writer, action string, data interface{}) error {
	return base.Success(w, action, data)
}

func SuccessWithNotice(w io.Writer, action string, data interface{}, notice interface{}) error {
	return base.SuccessWithNotice(w, action, data, notice)
}

func Error(w io.Writer, err error, context string) {
	base.Error(w, err, context)
}
