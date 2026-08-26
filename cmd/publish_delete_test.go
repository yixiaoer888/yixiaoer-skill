package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPublishDeleteDryRunOutputsRequest(t *testing.T) {
	var out bytes.Buffer
	cmd := newPublishDeleteCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"task_1", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "publish.delete.dry-run" {
		t.Fatalf("unexpected action: %#v", response)
	}
	data := response["data"].(map[string]interface{})
	if data["dryRun"] != true || data["taskId"] != "task_1" {
		t.Fatalf("unexpected dry-run metadata: %#v", data)
	}
	request := data["request"].(map[string]interface{})
	if request["method"] != "DELETE" || request["path"] != "/tasks/task_1/publish" {
		t.Fatalf("unexpected delete request: %#v", request)
	}
}

func TestPublishDeleteCallsExpectedEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/tasks/task_1/publish" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"statusCode": 0})
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	var out bytes.Buffer
	cmd := newPublishDeleteCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"task_1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["ok"] != true || response["action"] != "publish.delete" {
		t.Fatalf("unexpected delete envelope: %#v", response)
	}
	data := response["data"].(map[string]interface{})
	if data["taskId"] != "task_1" {
		t.Fatalf("unexpected delete data: %#v", data)
	}
	remote := data["response"].(map[string]interface{})
	if remote["statusCode"] != float64(0) {
		t.Fatalf("unexpected delete response: %#v", remote)
	}
}

func TestPublishDeleteRejectsEmptyTaskID(t *testing.T) {
	err := newPublishDeleteCmd().RunE(testCobraCommand(), []string{" "})
	if err == nil {
		t.Fatal("expected empty task id error")
	}
	if !strings.Contains(err.Error(), "published task id must not be empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPublishedWorkPreviewShowsWorkMetadataWithoutTaskID(t *testing.T) {
	preview, err := publishedWorkPreview("set_1", map[string]interface{}{
		"tasks": []interface{}{map[string]interface{}{
			"id":                  "task_1",
			"platformName":        "视频号",
			"platformAccountName": "流星运气逆天",
			"publishType":         "video",
			"title":               "春日片段",
			"coverUrl":            "https://example.com/cover.jpg",
			"openUrl":             "https://example.com/work",
			"stageStatus":         "success",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	public := publicPublishedWorkPreview(preview)
	works := public["works"].([]interface{})
	if len(works) != 1 {
		t.Fatalf("unexpected preview works: %#v", public)
	}
	work := works[0].(map[string]interface{})
	for field, want := range map[string]interface{}{
		"index":       1,
		"platform":    "视频号",
		"accountName": "流星运气逆天",
		"publishType": "video",
		"title":       "春日片段",
		"status":      "发布成功",
	} {
		if work[field] != want {
			t.Fatalf("work.%s = %#v, want %#v", field, work[field], want)
		}
	}
	if _, exists := work["taskId"]; exists {
		t.Fatalf("preview must not expose taskId: %#v", work)
	}
}

func TestSelectPublishedWorkUsesOneBasedIndexAndRetainsInternalTaskID(t *testing.T) {
	preview, err := publishedWorkPreview("set_1", map[string]interface{}{
		"tasks": []interface{}{
			map[string]interface{}{"id": "task_1", "platformName": "视频号", "stageStatus": "success"},
			map[string]interface{}{"id": "task_2", "platformName": "抖音", "stageStatus": "success"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	work, err := selectPublishedWork(preview, 2)
	if err != nil {
		t.Fatal(err)
	}
	if work["taskId"] != "task_2" || work["platform"] != "抖音" {
		t.Fatalf("unexpected selected work: %#v", work)
	}
}

func TestSelectPublishedWorkRejectsOutOfRangeIndex(t *testing.T) {
	preview, err := publishedWorkPreview("set_1", map[string]interface{}{
		"tasks": []interface{}{map[string]interface{}{"id": "task_1", "platformName": "视频号"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := selectPublishedWork(preview, 2); err == nil {
		t.Fatal("expected out-of-range selection error")
	}
}

func TestPublishDeletePreviewOutputsSelectableWorksWithoutTaskID(t *testing.T) {
	server := publishedWorkDetailsServer(t, "set_1")
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	var out bytes.Buffer
	cmd := newPublishDeletePreviewCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"set_1"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "publish.delete.preview" {
		t.Fatalf("unexpected preview envelope: %#v", response)
	}
	works := response["data"].(map[string]interface{})["works"].([]interface{})
	work := works[0].(map[string]interface{})
	if work["platform"] != "视频号" || work["accountName"] != "流星运气逆天" || work["status"] != "发布成功" {
		t.Fatalf("unexpected preview work: %#v", work)
	}
	if _, exists := work["taskId"]; exists {
		t.Fatalf("preview must not expose taskId: %#v", work)
	}
}

func TestPublishDeleteFromRecordDryRunOutputsSelectedWorkWithoutTaskID(t *testing.T) {
	server := publishedWorkDetailsServer(t, "set_1")
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	var out bytes.Buffer
	cmd := newPublishDeleteFromRecordCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"set_1", "--index", "1", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "publish.delete.from-record.dry-run" {
		t.Fatalf("unexpected dry-run envelope: %#v", response)
	}
	data := response["data"].(map[string]interface{})
	work := data["work"].(map[string]interface{})
	if work["index"] != float64(1) || work["platform"] != "视频号" {
		t.Fatalf("unexpected selected work: %#v", work)
	}
	if _, exists := work["taskId"]; exists {
		t.Fatalf("dry-run must not expose taskId: %#v", work)
	}
	request := data["request"].(map[string]interface{})
	if request["method"] != "DELETE" || request["path"] != "/tasks/task_1/publish" {
		t.Fatalf("unexpected delete request: %#v", request)
	}
}

func TestPublishDeleteFromRecordDoesNotExposeTaskIDForDeletedWork(t *testing.T) {
	server := publishedWorkDetailsServerWithStatus(t, "set_1", "deleted")
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := ExecuteWithIO([]string{"publish", "delete", "from-record", "set_1", "--index", "1", "--dry-run"}, &stdout, &stderr); code == 0 {
		t.Fatal("expected deleted-work selection to fail")
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout on error, got %q", stdout.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stderr.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	work := response["error"].(map[string]interface{})["details"].(map[string]interface{})["work"].(map[string]interface{})
	if _, exists := work["taskId"]; exists {
		t.Fatalf("error must not expose taskId: %#v", work)
	}
}

func publishedWorkDetailsServer(t *testing.T, taskSetID string) *httptest.Server {
	return publishedWorkDetailsServerWithStatus(t, taskSetID, "success")
}

func publishedWorkDetailsServerWithStatus(t *testing.T, taskSetID, taskStatus string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/v2/taskSets/"+taskSetID+"/tasks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"tasks": []map[string]interface{}{
					{
						"id":                  "task_1",
						"platformName":        "视频号",
						"platformAccountName": "流星运气逆天",
						"publishType":         "video",
						"title":               "春日片段",
						"coverUrl":            "https://example.com/cover.jpg",
						"stageStatus":         "success",
						"taskStatus":          taskStatus,
					},
				},
			},
		})
	}))
}
