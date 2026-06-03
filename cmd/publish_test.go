package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestPublishCommandSuccessCallsTaskSetAPI(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	if publishBody["publishType"] != "video" {
		t.Fatalf("unexpected publish body: %+v", publishBody)
	}
	if publishBody["publishChannel"] != "cloud" {
		t.Fatalf("expected cloud publishChannel in publish body, got %+v", publishBody)
	}
	if _, ok := publishBody["action"]; ok {
		t.Fatalf("did not expect action in publish body: %+v", publishBody)
	}
	if _, ok := publishBody["clientId"]; ok {
		t.Fatalf("did not expect clientId for cloud publish body: %+v", publishBody)
	}
	if platforms := publishBody["platforms"].([]interface{}); platforms[0] != "抖音" {
		t.Fatalf("expected Chinese platform name in publish body, got %+v", platforms)
	}
}

func TestPublishCommandWithClientIDUsesLocalChannel(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())
	publishChannelFlag = ""
	publishClientID = ""
	t.Cleanup(func() {
		publishChannelFlag = ""
		publishClientID = ""
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath, "client_1"})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	if publishBody["publishChannel"] != "local" || publishBody["clientId"] != "client_1" {
		t.Fatalf("expected local publish metadata in publish body, got %+v", publishBody)
	}
}

func TestPublishCommandMapsPlatformKeyToChineseForAPIRequests(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "douyin", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if platforms := publishBody["platforms"].([]interface{}); platforms[0] != "抖音" {
		t.Fatalf("expected platform key to map to Chinese platform name, got %+v", platforms)
	}
}

func TestPublishCommandConvertsScheduledTimeMillisecondsToSeconds(t *testing.T) {
	withRepoRoot(t)
	payload := validPublishPayload()
	cpf := payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	cpf["scheduledTime"] = float64(1760000000000)
	payloadPath := writePublishPayload(t, payload)

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}

	got := publishBody["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})["scheduledTime"]
	if got != float64(1760000000) {
		t.Fatalf("expected scheduledTime seconds in publish body, got %#v", got)
	}
}

func TestPublishCommandRejectsMultiPlatformArgument(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())

	var accountCalls int
	var publishCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			accountCalls++
		case "/taskSets/v2":
			publishCalls++
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音,知乎", payloadPath})
	if err == nil {
		t.Fatal("expected multi-platform publish error")
	}
	if accountCalls != 0 || publishCalls != 0 {
		t.Fatalf("expected no API calls, got accounts=%d publish=%d", accountCalls, publishCalls)
	}
}

func TestPublishCommandAcceptsFullPublishRequestPayload(t *testing.T) {
	withRepoRoot(t)
	inner := validPublishArgs()
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":      "publish",
		"publishType": "video",
		"platforms":   []interface{}{"抖音"},
		"cover": map[string]interface{}{
			"key":    "cover-key",
			"size":   float64(512),
			"width":  float64(1080),
			"height": float64(1920),
		},
		"coverKey":       "cover-key",
		"publishChannel": "cloud",
		"publishArgs":    inner,
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishBody["publishArgs"].(map[string]interface{})["accountForms"] == nil {
		t.Fatalf("expected publishArgs to contain accountForms directly, got %+v", publishBody)
	}
	if _, ok := publishBody["action"]; ok {
		t.Fatalf("did not expect action to be forwarded to publish API: %+v", publishBody)
	}
	if nested := publishBody["publishArgs"].(map[string]interface{})["publishArgs"]; nested != nil {
		t.Fatalf("did not expect nested publishArgs: %+v", publishBody)
	}
	if publishBody["coverKey"] != "cover-key" {
		t.Fatalf("expected top-level coverKey to be preserved, got %+v", publishBody)
	}
	if publishBody["cover"].(map[string]interface{})["key"] != "cover-key" {
		t.Fatalf("expected top-level cover to be preserved, got %+v", publishBody)
	}
}

func TestPublishCommandAcceptsStandardRequestPayloadShape(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "video",
		"platforms":      []interface{}{"抖音"},
		"coverKey":       "top-cover-key",
		"desc":           "任务描述",
		"clientId":       "local-client",
		"publishChannel": "local",
		"isDraft":        false,
		"isAppContent":   false,
		"publishArgs": map[string]interface{}{
			"video": map[string]interface{}{
				"duration": float64(10),
				"width":    float64(1000),
				"height":   float64(1000),
				"size":     float64(10000000),
				"key":      "video-key",
			},
			"images": []interface{}{
				map[string]interface{}{
					"width":  float64(1000),
					"height": float64(1000),
					"size":   float64(1000000),
					"key":    "image-key",
				},
			},
			"cover": map[string]interface{}{
				"width":  float64(1000),
				"height": float64(1000),
				"size":   float64(1000000),
				"key":    "shared-cover-key",
			},
			"coverKey": "shared-cover-key",
			"accountForms": []interface{}{
				map[string]interface{}{
					"mediaId":           "media_1",
					"platformName":      "抖音",
					"platformAccountId": "acc_001",
					"publishContentId":  "pub_1",
					"fps":               float64(0),
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "标题",
						"description": "<p>视频内容</p>",
						"tags":        []interface{}{"tag1"},
					},
				},
			},
			"content": "正文",
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}

	args := publishBody["publishArgs"].(map[string]interface{})
	form := args["accountForms"].([]interface{})[0].(map[string]interface{})
	if form["video"] == nil {
		t.Fatalf("expected shared publishArgs.video to be available to account form, got %+v", form)
	}
	if form["cover"] == nil || form["coverKey"] != "shared-cover-key" {
		t.Fatalf("expected shared cover fields on account form, got %+v", form)
	}
	if form["mediaId"] != "media_1" || form["platformName"] != "抖音" || form["publishContentId"] != "pub_1" {
		t.Fatalf("expected business fields to be preserved, got %+v", form)
	}
	if publishBody["publishChannel"] != "local" || publishBody["clientId"] != "local-client" {
		t.Fatalf("expected local publish metadata in API body, got %+v", publishBody)
	}
	if publishBody["desc"] != "任务描述" || publishBody["coverKey"] != "top-cover-key" {
		t.Fatalf("expected top-level standard fields to be preserved, got %+v", publishBody)
	}
	if publishBody["isAppContent"] != false || publishBody["isDraft"] != false {
		t.Fatalf("expected standard outer flags to be preserved, got %+v", publishBody)
	}
	if _, ok := publishBody["action"]; ok {
		t.Fatalf("did not expect action to be forwarded to publish API: %+v", publishBody)
	}
}

func TestPublishCommandAcceptsNodeStyleLocalStandardPayloadWithoutDuplicatedAccountResources(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "video",
		"platforms":      []interface{}{"抖音"},
		"coverKey":       "video_cover_key",
		"desc":           "视频发布任务",
		"publishChannel": "local",
		"clientId":       "local-client",
		"publishArgs": map[string]interface{}{
			"video": map[string]interface{}{
				"key":    "video_oss_key",
				"width":  float64(1080),
				"height": float64(1920),
				"size":   float64(52428800),
			},
			"cover": map[string]interface{}{
				"key":    "video_cover_key",
				"width":  float64(720),
				"height": float64(1280),
				"size":   float64(307200),
			},
			"coverKey": "video_cover_key",
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_001",
					"mediaId":           "media_001",
					"platformName":      "抖音",
					"publishContentId":  "publish_content_001",
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "演示视频标题",
						"description": "<p>演示视频简介</p>",
					},
				},
			},
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}

	args := publishBody["publishArgs"].(map[string]interface{})
	form := args["accountForms"].([]interface{})[0].(map[string]interface{})
	if form["video"] == nil || form["cover"] == nil || form["coverKey"] != "video_cover_key" {
		t.Fatalf("expected shared node-style local resources to be copied into account form, got %+v", form)
	}
	if publishBody["publishChannel"] != "local" || publishBody["clientId"] != "local-client" {
		t.Fatalf("expected local publish metadata in node-style API body, got %+v", publishBody)
	}
	if _, ok := publishBody["action"]; ok {
		t.Fatalf("did not expect action to be forwarded to publish API: %+v", publishBody)
	}
}

func TestPublishCommandAutoBuildsOuterEnvelopeFromPublishArgs(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":      "publish",
		"publishType": "video",
		"platforms":   []interface{}{"抖音"},
		"publishArgs": map[string]interface{}{
			"video": map[string]interface{}{
				"duration": float64(10),
				"width":    float64(1000),
				"height":   float64(1000),
				"size":     float64(10000000),
				"key":      "video-key",
			},
			"cover": map[string]interface{}{
				"width":  float64(1000),
				"height": float64(1000),
				"size":   float64(1000000),
				"key":    "shared-cover-key",
			},
			"coverKey": "shared-cover-key",
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_001",
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "视频标题",
						"description": "<p>精彩视频</p>",
					},
				},
			},
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}

	if publishBody["coverKey"] != "shared-cover-key" {
		t.Fatalf("expected top-level coverKey to be synthesized, got %+v", publishBody)
	}
	if publishBody["desc"] != "<p>精彩视频</p>" {
		t.Fatalf("expected top-level desc to be synthesized from contentPublishForm.description, got %+v", publishBody)
	}
	if publishBody["isDraft"] != false {
		t.Fatalf("expected top-level isDraft default false, got %+v", publishBody["isDraft"])
	}
	if publishBody["isAppContent"] != false {
		t.Fatalf("expected top-level isAppContent default false, got %+v", publishBody["isAppContent"])
	}
	cover, _ := publishBody["cover"].(map[string]interface{})
	if cover["key"] != "shared-cover-key" {
		t.Fatalf("expected top-level cover to be synthesized, got %+v", publishBody["cover"])
	}
}

func TestPublishDryRunAutoBuildsOuterEnvelopeFromPublishArgs(t *testing.T) {
	withRepoRoot(t)
	service := publishflow.NewService()
	result, err := service.DryRun(publishflow.ExecuteInput{
		PublishType:   "video",
		PlatformInput: "抖音",
		Payload: map[string]interface{}{
			"action":      "publish",
			"publishType": "video",
			"platforms":   []interface{}{"抖音"},
			"publishArgs": map[string]interface{}{
				"video": map[string]interface{}{
					"duration": float64(10),
					"width":    float64(1000),
					"height":   float64(1000),
					"size":     float64(10000000),
					"key":      "video-key",
				},
				"cover": map[string]interface{}{
					"width":  float64(1000),
					"height": float64(1000),
					"size":   float64(1000000),
					"key":    "shared-cover-key",
				},
				"coverKey": "shared-cover-key",
				"accountForms": []interface{}{
					map[string]interface{}{
						"platformAccountId": "acc_001",
						"contentPublishForm": map[string]interface{}{
							"formType":    "task",
							"title":       "视频标题",
							"description": "<p>精彩视频</p>",
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.PublishBody["coverKey"] != "shared-cover-key" {
		t.Fatalf("expected dry-run body coverKey to be synthesized, got %+v", result.PublishBody)
	}
	if result.PublishBody["desc"] != "<p>精彩视频</p>" {
		t.Fatalf("expected dry-run body desc to be synthesized, got %+v", result.PublishBody)
	}
	if result.PublishBody["isDraft"] != false || result.PublishBody["isAppContent"] != false {
		t.Fatalf("expected dry-run defaults for outer envelope, got %+v", result.PublishBody)
	}
	if result.PublishBody["publishChannel"] != "cloud" {
		t.Fatalf("expected publishChannel in dry-run request body, got %+v", result.PublishBody)
	}
	if _, ok := result.PublishBody["clientId"]; ok {
		t.Fatalf("did not expect clientId for cloud dry-run request body: %+v", result.PublishBody)
	}
	if result.PublishMode != "cloud" {
		t.Fatalf("expected dry-run publish mode metadata to stay cloud, got %q", result.PublishMode)
	}
}

func TestPublishCommandUsesLocalFlagsLikeNodeExample(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())
	publishChannelFlag = "local"
	publishClientID = "flag_client_1"
	t.Cleanup(func() {
		publishChannelFlag = ""
		publishClientID = ""
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishBody["publishChannel"] != "local" || publishBody["clientId"] != "flag_client_1" {
		t.Fatalf("expected flagged local publish metadata in API body, got %+v", publishBody)
	}
}

func TestPublishCommandRejectsLocalWithoutClientID(t *testing.T) {
	withRepoRoot(t)
	configPath := filepath.Join(t.TempDir(), "yxer-config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "video",
		"platforms":      []interface{}{"抖音"},
		"publishChannel": "local",
		"publishArgs":    validPublishArgs(),
	})
	publishChannelFlag = ""
	publishClientID = ""
	t.Cleanup(func() {
		publishChannelFlag = ""
		publishClientID = ""
	})

	var accountCalls int
	var publishCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			accountCalls++
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_001", "name": "账号", "status": 1},
				},
			})
		case "/taskSets/v2":
			publishCalls++
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err == nil {
		t.Fatal("expected local publish to require clientId")
	}
	if !strings.Contains(err.Error(), `clientId is required when publishChannel is "local"`) {
		t.Fatalf("expected local clientId requirement error, got %v", err)
	}
	if publishCalls != 0 {
		t.Fatalf("expected no publish call, got %d", publishCalls)
	}
}

func TestPublishCommandUsesConfiguredLocalClientID(t *testing.T) {
	withRepoRoot(t)
	configPath := filepath.Join(t.TempDir(), "yxer-config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)
	if _, err := config.SaveLocalClientID("configured_client_1"); err != nil {
		t.Fatal(err)
	}
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "video",
		"platforms":      []interface{}{"抖音"},
		"publishChannel": "local",
		"publishArgs":    validPublishArgs(),
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishBody["publishChannel"] != "local" || publishBody["clientId"] != "configured_client_1" {
		t.Fatalf("expected configured local publish metadata in API body, got %+v", publishBody)
	}
}

func TestPublishCommandReturnsStructuredFallbackErrorByDefault(t *testing.T) {
	withRepoRoot(t)
	configPath := filepath.Join(t.TempDir(), "yxer-config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)
	if _, err := config.SaveLocalClientID("configured_client_1"); err != nil {
		t.Fatal(err)
	}
	payloadPath := writePublishPayload(t, validPublishPayload())

	var publishCalls int
	var publishBodies []map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_001", "name": "账号", "status": 1},
				},
			})
		case "/taskSets/v2":
			publishCalls++
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			publishBodies = append(publishBodies, body)
			if publishCalls == 1 {
				http.Error(w, `{"message":"账号代理不存在"}`, http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_local_retry"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err == nil {
		t.Fatal("expected structured fallback error")
	}
	typed, ok := err.(*yxerrors.Error)
	if !ok {
		t.Fatalf("expected structured fallback error, got %T: %v", err, err)
	}
	if typed.Category != "publish_channel_fallback" {
		t.Fatalf("expected publish_channel_fallback category, got %+v", typed)
	}
	if !strings.Contains(typed.NextCommand, "--publish-channel local") {
		t.Fatalf("expected local fallback nextCommand, got %+v", typed)
	}
	if !strings.Contains(typed.Hint, "--auto-fallback-local") {
		t.Fatalf("expected auto fallback hint, got %+v", typed)
	}
	if publishCalls != 1 {
		t.Fatalf("expected single cloud attempt before fallback error, got %d", publishCalls)
	}
	if publishBodies[0]["publishChannel"] != "cloud" {
		t.Fatalf("expected first publish attempt to stay cloud, got %+v", publishBodies[0])
	}
}

func TestPublishCommandAutoFallbacksToLocalWhenFlagEnabled(t *testing.T) {
	withRepoRoot(t)
	configPath := filepath.Join(t.TempDir(), "yxer-config.json")
	t.Setenv("YIXIAOER_CONFIG", configPath)
	if _, err := config.SaveLocalClientID("configured_client_1"); err != nil {
		t.Fatal(err)
	}
	payloadPath := writePublishPayload(t, validPublishPayload())
	publishAutoFallbackLocal = true
	t.Cleanup(func() {
		publishAutoFallbackLocal = false
	})

	var publishCalls int
	var publishBodies []map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_001", "name": "账号", "status": 1},
				},
			})
		case "/taskSets/v2":
			publishCalls++
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			publishBodies = append(publishBodies, body)
			if publishCalls == 1 {
				http.Error(w, `{"message":"账号代理不存在"}`, http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_local_retry"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 2 {
		t.Fatalf("expected automatic local retry, got %d calls", publishCalls)
	}
	if publishBodies[0]["publishChannel"] != "cloud" {
		t.Fatalf("expected first publish attempt to stay cloud, got %+v", publishBodies[0])
	}
	if publishBodies[1]["publishChannel"] != "local" || publishBodies[1]["clientId"] != "configured_client_1" {
		t.Fatalf("expected second publish attempt to switch to local, got %+v", publishBodies[1])
	}
}

func TestPublishCommandSchemaFailureDoesNotCallAPIs(t *testing.T) {
	withRepoRoot(t)
	payload := validPublishPayload()
	cpf := payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	delete(cpf, "title")
	payloadPath := writePublishPayload(t, payload)

	var accountCalls int
	var publishCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			accountCalls++
		case "/taskSets/v2":
			publishCalls++
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err == nil {
		t.Fatal("expected schema validation error")
	}
	if accountCalls != 0 || publishCalls != 0 {
		t.Fatalf("expected no API calls, got accounts=%d publish=%d", accountCalls, publishCalls)
	}
}

func TestPublishCommandRejectsKuaishouImageTextWithMoreThanFourTags(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "imageText",
		"platforms":      []interface{}{"快手"},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_kuaishou_1",
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"description": "<p>快手图文内容</p>",
						"visibleType": float64(0),
						"tags":        []interface{}{"话题1", "话题2", "话题3", "话题4", "话题5"},
						"images": []interface{}{
							map[string]interface{}{
								"key":    "image-key",
								"size":   float64(1024),
								"width":  float64(1080),
								"height": float64(1920),
							},
						},
					},
					"cover": map[string]interface{}{
						"key":    "image-key",
						"size":   float64(1024),
						"width":  float64(1080),
						"height": float64(1920),
					},
					"coverKey": "image-key",
				},
			},
		},
	})

	var accountCalls int
	var publishCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			accountCalls++
		case "/taskSets/v2":
			publishCalls++
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"imageText", "快手", payloadPath})
	if err == nil {
		t.Fatal("expected kuaishou imageText schema validation error")
	}
	if !strings.Contains(err.Error(), "Schema validation failed") {
		t.Fatalf("expected schema validation error, got %v", err)
	}
	if accountCalls != 0 || publishCalls != 0 {
		t.Fatalf("expected no API calls, got accounts=%d publish=%d", accountCalls, publishCalls)
	}
}

func TestPublishCommandPreflightFailureDoesNotCallAPIs(t *testing.T) {
	withRepoRoot(t)
	payload := validPublishPayload()
	form := payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
	form["video"].(map[string]interface{})["key"] = "https://example.com/video.mp4"
	payloadPath := writePublishPayload(t, payload)

	var accountCalls int
	var publishCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			accountCalls++
		case "/taskSets/v2":
			publishCalls++
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err == nil {
		t.Fatal("expected preflight error")
	}
	if accountCalls != 0 || publishCalls != 0 {
		t.Fatalf("expected no API calls, got accounts=%d publish=%d", accountCalls, publishCalls)
	}
}

func TestPublishCommandRejectsInvalidTopLevelCoverInFullPublishRequest(t *testing.T) {
	withRepoRoot(t)
	inner := validPublishArgs()
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":      "publish",
		"publishType": "video",
		"platforms":   []interface{}{"抖音"},
		"cover": map[string]interface{}{
			"key": "https://example.com/cover.jpg",
		},
		"coverKey":    "cover-key",
		"publishArgs": inner,
	})

	var accountCalls int
	var publishCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			accountCalls++
		case "/taskSets/v2":
			publishCalls++
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err == nil {
		t.Fatal("expected top-level cover preflight error")
	}
	if accountCalls != 0 || publishCalls != 0 {
		t.Fatalf("expected no API calls, got accounts=%d publish=%d", accountCalls, publishCalls)
	}
}

func TestPublishCommandOfflineAccountDoesNotPublish(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 0, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err == nil {
		t.Fatal("expected offline account error")
	}
	if publishCalls != 0 {
		t.Fatalf("expected no publish call, got %d", publishCalls)
	}
}

func TestPublishCommandPreservesDistinctImageTextDescriptionAndContentFromPayload(t *testing.T) {
	withRepoRoot(t)
	topicHTML := `<p>今日穿搭分享</p><p><topic text="穿搭">#穿搭</topic><topic text="夏日">#夏日</topic></p>`
	separateContent := "<p>独立 content 字段</p>"
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "imageText",
		"platforms":      []interface{}{"小红书"},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"content": separateContent,
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_xhs_1",
					"images": []interface{}{
						map[string]interface{}{
							"key":    "uploaded/cover.png",
							"size":   float64(512),
							"width":  float64(1080),
							"height": float64(1080),
							"format": "png",
						},
					},
					"cover": map[string]interface{}{
						"key":    "uploaded/cover.png",
						"size":   float64(512),
						"width":  float64(1080),
						"height": float64(1080),
						"format": "png",
					},
					"coverKey": "uploaded/cover.png",
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "夏日穿搭",
						"description": topicHTML,
						"visibleType": float64(0),
						"images": []interface{}{
							map[string]interface{}{
								"key":    "uploaded/cover.png",
								"size":   float64(512),
								"width":  float64(1080),
								"height": float64(1080),
								"format": "png",
							},
						},
					},
				},
			},
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := imageTextPublishTestServer(t, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"imageText", "小红书", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	args := publishBody["publishArgs"].(map[string]interface{})
	if args["content"] != separateContent {
		t.Fatalf("expected publishArgs.content to stay independent, got %+v", args["content"])
	}
	cpf := args["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if cpf["description"] != topicHTML {
		t.Fatalf("expected contentPublishForm.description to keep topic HTML, got %+v", cpf)
	}
}

func TestPublishCommandNormalizesTopicHTMLIntoDescriptionAndContent(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "imageText",
		"platforms":      []interface{}{"抖音"},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_001",
					"images": []interface{}{
						map[string]interface{}{
							"key":    "uploaded/cover.png",
							"size":   float64(512),
							"width":  float64(1080),
							"height": float64(1080),
							"format": "png",
						},
					},
					"cover": map[string]interface{}{
						"key":    "uploaded/cover.png",
						"size":   float64(512),
						"width":  float64(1080),
						"height": float64(1080),
						"format": "png",
					},
					"coverKey": "uploaded/cover.png",
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "夏日穿搭",
						"description": "今日穿搭分享",
						"tags":        []interface{}{"穿搭", "#夏日"},
						"images": []interface{}{
							map[string]interface{}{
								"key":    "uploaded/cover.png",
								"size":   float64(512),
								"width":  float64(1080),
								"height": float64(1080),
								"format": "png",
							},
						},
					},
				},
			},
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"imageText", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}

	expected := `<p>今日穿搭分享</p><p><topic text="穿搭">#穿搭</topic><topic text="夏日">#夏日</topic></p>`
	args := publishBody["publishArgs"].(map[string]interface{})
	if args["content"] != expected {
		t.Fatalf("expected publishArgs.content topic HTML, got %+v", args["content"])
	}
	cpf := args["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if cpf["description"] != expected {
		t.Fatalf("expected contentPublishForm.description topic HTML, got %+v", cpf["description"])
	}
}

func TestPublishCommandNormalizesDouyinShoppingCartStructure(t *testing.T) {
	withRepoRoot(t)
	payload := validPublishPayload()
	cpf := payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	cpf["shoppingCart"] = []interface{}{
		map[string]interface{}{
			"sale_title":   "点击购买",
			"yixiaoerId":   "goods_001",
			"yixiaoerName": "测试商品",
			"raw": map[string]interface{}{
				"gid":        "goods_001",
				"goods_imgs": []interface{}{"https://example.invalid/goods.png"},
			},
		},
	}
	payloadPath := writePublishPayload(t, payload)

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}

	form := publishBody["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
	got := form["contentPublishForm"].(map[string]interface{})
	if _, exists := got["shoppingCart"]; exists {
		t.Fatalf("expected shoppingCart to normalize to shopping_cart, got %+v", got)
	}
	items := got["shopping_cart"].([]interface{})
	item := items[0].(map[string]interface{})
	if item["sale_title"] != "点击购买" {
		t.Fatalf("expected sale_title to stay unchanged, got %+v", item)
	}
	if len(item["images"].([]interface{})) != 1 {
		t.Fatalf("expected images derived from raw, got %+v", item)
	}
	data := item["data"].(map[string]interface{})
	if data["yixiaoerId"] != "goods_001" || data["yixiaoerName"] != "测试商品" {
		t.Fatalf("expected nested data object, got %+v", item)
	}
}

func TestPublishCommandUsesImageTextPublishType(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "imageText",
		"platforms":      []interface{}{"小红书"},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_xhs_1",
					"images": []interface{}{
						map[string]interface{}{
							"key":    "uploaded/cover.png",
							"size":   float64(512),
							"width":  float64(1080),
							"height": float64(1080),
							"format": "png",
						},
					},
					"cover": map[string]interface{}{
						"key":    "uploaded/cover.png",
						"size":   float64(512),
						"width":  float64(1080),
						"height": float64(1080),
						"format": "png",
					},
					"coverKey": "uploaded/cover.png",
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "夏日穿搭",
						"description": "<p>今日穿搭分享</p>",
						"visibleType": float64(0),
						"images": []interface{}{
							map[string]interface{}{
								"key":    "uploaded/cover.png",
								"size":   float64(512),
								"width":  float64(1080),
								"height": float64(1080),
								"format": "png",
							},
						},
					},
				},
			},
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := imageTextPublishTestServer(t, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"imageText", "小红书", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishBody["publishType"] != "imageText" {
		t.Fatalf("expected publishType imageText, got %+v", publishBody["publishType"])
	}
}

func imageTextPublishTestServer(t *testing.T, publishCalls *int, publishBody *map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_xhs_1", "name": "图文账号", "status": 1},
				},
			})
		case "/taskSets/v2":
			*publishCalls++
			if err := json.NewDecoder(r.Body).Decode(publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_1"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
}

func testPNGBytes(t *testing.T) []byte {
	t.Helper()
	return []byte{
		137, 80, 78, 71, 13, 10, 26, 10,
		0, 0, 0, 13, 73, 72, 68, 82,
		0, 0, 0, 1, 0, 0, 0, 1,
		8, 2, 0, 0, 0, 144, 119, 83,
		222, 0, 0, 0, 12, 73, 68, 65,
		84, 120, 156, 99, 248, 15, 4, 0,
		9, 251, 3, 253, 160, 164, 95, 165,
		0, 0, 0, 0, 73, 69, 78, 68,
		174, 66, 96, 130,
	}
}

func publishTestServer(t *testing.T, accountStatus int, publishCalls *int, publishBody *map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			if got := r.URL.Query().Get("platform"); got != "抖音" {
				t.Fatalf("unexpected platform query: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_001", "name": "账号", "status": accountStatus},
				},
			})
		case "/taskSets/v2":
			*publishCalls++
			if err := json.NewDecoder(r.Body).Decode(publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_1"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
}

func validPublishPayload() map[string]interface{} {
	return map[string]interface{}{
		"action":         "publish",
		"publishType":    "video",
		"platforms":      []interface{}{"抖音"},
		"publishChannel": "cloud",
		"publishArgs":    validPublishArgs(),
	}
}

func validPublishArgs() map[string]interface{} {
	return map[string]interface{}{
		"accountForms": []interface{}{
			map[string]interface{}{
				"platformAccountId": "acc_001",
				"video": map[string]interface{}{
					"key":    "video-key",
					"size":   float64(1024),
					"width":  float64(1080),
					"height": float64(1920),
				},
				"cover": map[string]interface{}{
					"key":    "cover-key",
					"size":   float64(512),
					"width":  float64(1080),
					"height": float64(1920),
				},
				"coverKey": "cover-key",
				"contentPublishForm": map[string]interface{}{
					"formType":    "task",
					"title":       "视频标题",
					"description": "视频描述",
				},
			},
		},
	}
}

func writePublishPayload(t *testing.T, payload map[string]interface{}) string {
	t.Helper()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "payload.json")
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func configureAPIKey(t *testing.T, apiKey string) {
	t.Helper()
	configPath := os.Getenv("YIXIAOER_CONFIG")
	if configPath == "" {
		configPath = filepath.Join(t.TempDir(), "yxer-config.json")
		t.Setenv("YIXIAOER_CONFIG", configPath)
	}
	if _, err := config.SaveAPIKey(apiKey); err != nil {
		t.Fatal(err)
	}
}

func useTestAPIBaseURL(t *testing.T, rawURL string) {
	t.Helper()
	t.Cleanup(api.SetBaseURLForTest(rawURL))
}

func testCobraCommand() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	return cmd
}

func withRepoRoot(t *testing.T) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatal(err)
		}
	})
}
