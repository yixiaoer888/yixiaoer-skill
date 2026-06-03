package api

import (
	"fmt"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

// Publish creates a task set on the gateway and returns a normalized result
// containing the taskSetId. The gateway wraps success as
// {statusCode:0, data:<taskSetId>} (data may also be an object or carry the id
// at the top level), so we unwrap the envelope, confirm a taskSetId is present,
// and surface a remote error when it is missing (e.g. an empty body that the
// backend still answered with HTTP 200).
func (c *Client) Publish(body map[string]interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := c.Post("/taskSets/v2", body, &result); err != nil {
		return nil, err
	}
	taskSetID := extractTaskSetID(result)
	if taskSetID == "" {
		return nil, yxerrors.Remote("publish succeeded over HTTP but the response did not contain a taskSetId", map[string]interface{}{
			"response": result,
		}).WithHint("发布请求已发出但未拿到 taskSetId，请用 yxer records list 确认任务集是否真正创建，再决定是否重试。").
			WithNextCommand("yxer records list")
	}
	out := map[string]interface{}{"taskSetId": taskSetID}
	if data, ok := result["data"]; ok {
		out["response"] = data
	} else if len(result) > 0 {
		out["response"] = result
	}
	return out, nil
}

// extractTaskSetID pulls the task set identifier out of the gateway response,
// tolerating the three shapes seen in practice: a bare string under "data", an
// object under "data" with an id field, and the id at the top level.
func extractTaskSetID(result map[string]interface{}) string {
	if result == nil {
		return ""
	}
	if id := taskSetIDFromValue(result["data"]); id != "" {
		return id
	}
	return taskSetIDFromObject(result)
}

func taskSetIDFromValue(value interface{}) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case map[string]interface{}:
		return taskSetIDFromObject(typed)
	default:
		return ""
	}
}

func taskSetIDFromObject(obj map[string]interface{}) string {
	if obj == nil {
		return ""
	}
	for _, key := range []string{"taskSetId", "taskSetID", "task_set_id", "taskIdentityId", "id"} {
		if value, ok := obj[key]; ok && value != nil {
			text := strings.TrimSpace(fmt.Sprint(value))
			if text != "" && text != "<nil>" {
				return text
			}
		}
	}
	return ""
}
