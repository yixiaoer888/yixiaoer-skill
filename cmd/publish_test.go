package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

func TestPublishCommandSuccessCallsTaskSetAPI(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	if publishBody["publishType"] != "video" || publishBody["publishChannel"] != "cloud" {
		t.Fatalf("unexpected publish body: %+v", publishBody)
	}
	if _, ok := publishBody["action"]; ok {
		t.Fatalf("did not expect action in publish body: %+v", publishBody)
	}
	if _, ok := publishBody["clientId"]; ok {
		t.Fatalf("did not expect clientId in cloud publish body: %+v", publishBody)
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
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath, "client_1"})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	if publishBody["publishChannel"] != "local" || publishBody["clientId"] != "client_1" {
		t.Fatalf("unexpected local publish body: %+v", publishBody)
	}
}

func TestPublishCommandMapsPlatformKeyToChineseForAPIRequests(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

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
	cpf := payload["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	cpf["scheduledTime"] = float64(1760000000000)
	payloadPath := writePublishPayload(t, payload)

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

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
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

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
	inner := validPublishPayload()
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "video",
		"platforms":      []interface{}{"抖音"},
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
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

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
					"mediaId":          "media_1",
					"platformName":     "抖音",
					"platformAccountId": "acc_001",
					"publishContentId": "pub_1",
					"fps":              float64(0),
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
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

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
		t.Fatalf("expected local publish config to be preserved, got %+v", publishBody)
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
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

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
		t.Fatalf("expected local publish config to be preserved, got %+v", publishBody)
	}
	if _, ok := publishBody["action"]; ok {
		t.Fatalf("did not expect action to be forwarded to publish API: %+v", publishBody)
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
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishBody["publishChannel"] != "local" || publishBody["clientId"] != "flag_client_1" {
		t.Fatalf("expected local publish config from flags, got %+v", publishBody)
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
		"publishArgs":    validPublishPayload(),
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
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

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
		"publishArgs":    validPublishPayload(),
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishBody["publishChannel"] != "local" || publishBody["clientId"] != "configured_client_1" {
		t.Fatalf("expected local publish config from saved config, got %+v", publishBody)
	}
}

func TestPublishCommandSchemaFailureDoesNotCallAPIs(t *testing.T) {
	withRepoRoot(t)
	payload := validPublishPayload()
	cpf := payload["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
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
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err == nil {
		t.Fatal("expected schema validation error")
	}
	if accountCalls != 0 || publishCalls != 0 {
		t.Fatalf("expected no API calls, got accounts=%d publish=%d", accountCalls, publishCalls)
	}
}

func TestPublishCommandPreflightFailureDoesNotCallAPIs(t *testing.T) {
	withRepoRoot(t)
	payload := validPublishPayload()
	form := payload["accountForms"].([]interface{})[0].(map[string]interface{})
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
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

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
	inner := validPublishPayload()
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
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

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
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err == nil {
		t.Fatal("expected offline account error")
	}
	if publishCalls != 0 {
		t.Fatalf("expected no publish call, got %d", publishCalls)
	}
}

func TestPublishCommandBuildsImageTextPayloadFromFlags(t *testing.T) {
	withRepoRoot(t)
	imagePath := filepath.Join(t.TempDir(), "cover.png")
	if err := os.WriteFile(imagePath, testPNGBytes(t), 0o644); err != nil {
		t.Fatal(err)
	}
	publishAccount = "图文账号"
	publishTitle = "图文标题"
	publishDescription = "图文描述"
	publishImages = []string{imagePath}
	publishCoverPath = ""
	publishContent = ""
	publishVideoKey = ""
	publishVisibleType = 0
	t.Cleanup(resetPublishFlagsForTest)

	var publishCalls int
	var publishBody map[string]interface{}
	server := imageTextPublishTestServer(t, &publishCalls, &publishBody)
	defer server.Close()
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"image-text", "小红书"})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	args := publishBody["publishArgs"].(map[string]interface{})
	form := args["accountForms"].([]interface{})[0].(map[string]interface{})
	if form["platformAccountId"] != "acc_xhs_1" {
		t.Fatalf("unexpected account form: %+v", form)
	}
	cpf := form["contentPublishForm"].(map[string]interface{})
	if cpf["description"] != "图文描述" || cpf["title"] != "图文标题" {
		t.Fatalf("unexpected contentPublishForm: %+v", cpf)
	}
}

func TestPublishCommandBuildsArticlePayloadFromFlags(t *testing.T) {
	withRepoRoot(t)
	coverPath := filepath.Join(t.TempDir(), "cover.png")
	if err := os.WriteFile(coverPath, testPNGBytes(t), 0o644); err != nil {
		t.Fatal(err)
	}
	contentPath := filepath.Join(t.TempDir(), "article.html")
	if err := os.WriteFile(contentPath, []byte("<p>这是一篇用于测试的文章正文内容</p>"), 0o644); err != nil {
		t.Fatal(err)
	}

	publishAccount = "知乎账号"
	publishTitle = "这是一个足够长的知乎文章标题"
	publishContent = "@" + contentPath
	publishCoverPath = coverPath
	publishDescription = ""
	publishImages = nil
	publishVideoKey = ""
	publishVisibleType = -1
	t.Cleanup(resetPublishFlagsForTest)

	var publishCalls int
	var publishBody map[string]interface{}
	server := articlePublishTestServer(t, &publishCalls, &publishBody)
	defer server.Close()
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"article", "知乎"})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	args := publishBody["publishArgs"].(map[string]interface{})
	form := args["accountForms"].([]interface{})[0].(map[string]interface{})
	if form["platformAccountId"] != "acc_zhihu_1" {
		t.Fatalf("unexpected account form: %+v", form)
	}
	cpf := form["contentPublishForm"].(map[string]interface{})
	if cpf["title"] != "这是一个足够长的知乎文章标题" {
		t.Fatalf("unexpected contentPublishForm title: %+v", cpf)
	}
	if cpf["pubType"] != float64(1) {
		t.Fatalf("expected pubType=1, got %+v", cpf)
	}
	if !strings.Contains(cpf["content"].(string), "测试的文章正文内容") {
		t.Fatalf("expected content read from file, got %+v", cpf["content"])
	}
}

func TestPublishCommandBuildsVideoPayloadFromFlags(t *testing.T) {
	withRepoRoot(t)
	videoPath := filepath.Join(t.TempDir(), "clip.mp4")
	createTestVideo(t, videoPath)
	coverPath := filepath.Join(t.TempDir(), "cover.png")
	if err := os.WriteFile(coverPath, testPNGBytes(t), 0o644); err != nil {
		t.Fatal(err)
	}

	publishAccount = "视频账号"
	publishTitle = "视频标题"
	publishDescription = "视频描述"
	publishVideoPath = videoPath
	publishCoverPath = coverPath
	publishContent = ""
	publishImages = nil
	publishVideoKey = ""
	publishVisibleType = -1
	t.Cleanup(resetPublishFlagsForTest)

	var publishCalls int
	var publishBody map[string]interface{}
	server := videoPublishTestServer(t, &publishCalls, &publishBody)
	defer server.Close()
	t.Setenv("YIXIAOER_API_KEY", "test-key")
	t.Setenv("YIXIAOER_API_URL", server.URL)

	err := publishCmd.RunE(testCobraCommand(), []string{"video", "抖音"})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	args := publishBody["publishArgs"].(map[string]interface{})
	form := args["accountForms"].([]interface{})[0].(map[string]interface{})
	video := form["video"].(map[string]interface{})
	if video["key"] != "uploaded/clip.mp4" {
		t.Fatalf("unexpected video payload: %+v", video)
	}
	if video["width"] == nil || video["height"] == nil || video["duration"] == nil {
		t.Fatalf("expected probed video metadata, got %+v", video)
	}
}

func resetPublishFlagsForTest() {
	publishChannelFlag = ""
	publishClientID = ""
	publishAccount = ""
	publishTitle = ""
	publishDescription = ""
	publishContent = ""
	publishImages = nil
	publishVideoPath = ""
	publishVideoKey = ""
	publishCoverPath = ""
	publishVisibleType = -1
}

func imageTextPublishTestServer(t *testing.T, publishCalls *int, publishBody *map[string]interface{}) *httptest.Server {
	t.Helper()
	imageBytes := testPNGBytes(t)
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_xhs_1", "name": "图文账号", "status": 1},
				},
			})
		case "/storages/cloud-publish/upload-url":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"serviceUrl": server.URL + "/oss/cover.png",
					"key":        "uploaded/cover.png",
				},
			})
		case "/oss/cover.png":
			w.WriteHeader(http.StatusOK)
		case "/taskSets/v2":
			*publishCalls++
			if err := json.NewDecoder(r.Body).Decode(publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_1"},
			})
		default:
			if strings.HasSuffix(r.URL.Path, ".png") {
				w.Header().Set("Content-Type", "image/png")
				_, _ = w.Write(imageBytes)
				return
			}
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	return server
}

func articlePublishTestServer(t *testing.T, publishCalls *int, publishBody *map[string]interface{}) *httptest.Server {
	t.Helper()
	imageBytes := testPNGBytes(t)
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_zhihu_1", "name": "知乎账号", "status": 1},
				},
			})
		case "/storages/cloud-publish/upload-url":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"serviceUrl": server.URL + "/oss/cover.png",
					"key":        "uploaded/cover.png",
				},
			})
		case "/oss/cover.png":
			w.WriteHeader(http.StatusOK)
		case "/taskSets/v2":
			*publishCalls++
			if err := json.NewDecoder(r.Body).Decode(publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_article"},
			})
		default:
			if strings.HasSuffix(r.URL.Path, ".png") {
				w.Header().Set("Content-Type", "image/png")
				_, _ = w.Write(imageBytes)
				return
			}
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	return server
}

func videoPublishTestServer(t *testing.T, publishCalls *int, publishBody *map[string]interface{}) *httptest.Server {
	t.Helper()
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_video_1", "name": "视频账号", "status": 1},
				},
			})
		case "/storages/cloud-publish/upload-url":
			fileKey := r.URL.Query().Get("fileKey")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"serviceUrl": server.URL + "/oss/" + fileKey,
					"key":        "uploaded/" + fileKey,
				},
			})
		case "/oss/clip.mp4", "/oss/cover.png":
			w.WriteHeader(http.StatusOK)
		case "/taskSets/v2":
			*publishCalls++
			if err := json.NewDecoder(r.Body).Decode(publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_video"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	return server
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

func createTestVideo(t *testing.T, outputPath string) {
	t.Helper()
	ffmpegPath := resolveCommandPath(t, "ffmpeg")
	if ffmpegPath == "" {
		t.Skip("ffmpeg not available for video flags test")
	}
	cmd := exec.Command(ffmpegPath, "-y", "-f", "lavfi", "-i", "color=c=black:s=16x16:d=1", "-pix_fmt", "yuv420p", outputPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Skipf("skipping video flags test because ffmpeg could not run: %v, output=%s", err, string(out))
	}
}

func resolveCommandPath(t *testing.T, name string) string {
	t.Helper()
	if path, err := exec.LookPath(name); err == nil {
		return path
	}
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "(Get-Command "+name+" | Select-Object -ExpandProperty Source)")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
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
