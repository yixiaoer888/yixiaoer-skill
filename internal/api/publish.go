package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

// DeletePublishedTask deletes the published work for one platform task.
// taskID is the task id returned by `yxer query details`, not a taskSetId.
func (c *Client) DeletePublishedTask(taskID string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	if err := c.Delete(fmt.Sprintf("/tasks/%s/publish", url.PathEscape(taskID)), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Publish creates a task set on the gateway and returns a normalized result
// containing the taskSetId. The gateway wraps success as
// {statusCode:0, data:<taskSetId>} (data may also be an object or carry the id
// at the top level), so we unwrap the envelope, confirm a taskSetId is present,
// and surface a remote error when it is missing (e.g. an empty body that the
// backend still answered with HTTP 200).
func (c *Client) Publish(body map[string]interface{}) (map[string]interface{}, error) {
	// The gateway expects the frontend category selection shape. Keep query,
	// schema validation, and dry-run payloads in their canonical yixiaoer shape;
	// perform this wire-format conversion only at the actual publish boundary.
	wireBody, _ := clonePublishValue(body).(map[string]interface{})
	normalizePublishCategories(wireBody)
	var result map[string]interface{}
	if err := c.Post("/taskSets/v2", wireBody, &result); err != nil {
		return nil, err
	}
	taskSetID := extractTaskSetID(result)
	if taskSetID == "" {
		return nil, yxerrors.Remote("publish succeeded over HTTP but the response did not contain a taskSetId", map[string]interface{}{
			"response": result,
		}).WithHint("发布请求已发出但未拿到 taskSetId，请用 yxer query records 确认任务集是否真正创建，再决定是否重试。").
			WithNextCommand("yxer query records --limit 10 --json")
	}
	out := map[string]interface{}{"taskSetId": taskSetID}
	if data, ok := result["data"]; ok {
		out["response"] = data
	} else if len(result) > 0 {
		out["response"] = result
	}
	return out, nil
}

func clonePublishValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(typed))
		for key, item := range typed {
			out[key] = clonePublishValue(item)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(typed))
		for i, item := range typed {
			out[i] = clonePublishValue(item)
		}
		return out
	default:
		return value
	}
}

func normalizePublishCategories(body map[string]interface{}) {
	if body == nil {
		return
	}
	publishArgs, _ := body["publishArgs"].(map[string]interface{})
	forms, _ := publishArgs["accountForms"].([]interface{})
	for _, item := range forms {
		form, _ := item.(map[string]interface{})
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if category, ok := cpf["category"]; ok {
			cpf["category"] = wrapPublishCategoryValue(category)
		}
	}
}

func wrapPublishCategoryValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case []interface{}:
		for i, item := range typed {
			typed[i] = wrapPublishCategoryValue(item)
		}
		return typed
	case map[string]interface{}:
		if _, hasID := typed["id"]; hasID {
			return typed
		}
		id, idOK := publishStringField(typed, "yixiaoerId")
		name, nameOK := publishStringField(typed, "yixiaoerName")
		if !idOK || !nameOK {
			return typed
		}
		wrapped := map[string]interface{}{
			"id":   id,
			"text": name,
			"raw":  typed,
		}
		if children, ok := typed["child"].([]interface{}); ok && len(children) > 0 {
			wrapped["children"] = children
		}
		return wrapped
	default:
		return value
	}
}

func publishStringField(obj map[string]interface{}, key string) (string, bool) {
	value, ok := obj[key]
	if !ok || value == nil {
		return "", false
	}
	text := strings.TrimSpace(fmt.Sprint(value))
	return text, text != "" && text != "<nil>"
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
