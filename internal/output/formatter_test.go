package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestErrorIncludesMachineFields(t *testing.T) {
	var out bytes.Buffer

	Error(&out, yxerrors.Usage("bad payload", []string{"missing title"}).
		WithCategory("validation").
		WithHint("fill title").
		WithNextCommand("yxer validate xhs imageText payload.json"), "run command")

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	errorObj := response["error"].(map[string]interface{})
	if errorObj["code"] != yxerrors.UsageErr {
		t.Fatalf("unexpected code: %#v", errorObj)
	}
	if errorObj["category"] != "validation" {
		t.Fatalf("unexpected category: %#v", errorObj)
	}
	if errorObj["nextCommand"] != "yxer validate xhs imageText payload.json" {
		t.Fatalf("unexpected nextCommand: %#v", errorObj)
	}
}
