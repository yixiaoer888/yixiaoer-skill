package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	publishmod "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

func TestPublishCommandFindsTargetAccountOnSecondPage(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())

	var publishCalls int
	var publishBody map[string]interface{}
	accountRequests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			accountRequests++
			switch r.URL.Query().Get("page") {
			case "1":
				items := make([]map[string]interface{}, 0, 50)
				for i := 0; i < 50; i++ {
					items = append(items, map[string]interface{}{
						"platformAccountId": "acc_other_" + string(rune('a'+i)),
						"name":              "账号",
						"status":            1,
					})
				}
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"data": map[string]interface{}{
						"list":      items,
						"page":      1,
						"size":      50,
						"totalPage": 2,
						"totalSize": 51,
					},
				})
			case "2":
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"data": map[string]interface{}{
						"list": []map[string]interface{}{
							{"platformAccountId": "acc_001", "name": "目标账号", "status": 1},
						},
						"page":      2,
						"size":      50,
						"totalPage": 2,
						"totalSize": 51,
					},
				})
			default:
				t.Fatalf("unexpected page query: %s", r.URL.Query().Get("page"))
			}
		case "/taskSets/v2":
			publishCalls++
			if err := json.NewDecoder(r.Body).Decode(&publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_1"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if accountRequests != 2 {
		t.Fatalf("expected two account page requests, got %d", accountRequests)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
}

func TestPublishCommandRejectsPositionalClientIDWithoutLocalChannel(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath, "client_1"})
	if err == nil {
		t.Fatal("expected positional clientId to require explicit local channel")
	}
	if !strings.Contains(err.Error(), "positional clientId requires local publish channel") {
		t.Fatalf("unexpected error: %v", err)
	}
	if publishCalls != 0 {
		t.Fatalf("expected no publish call, got %d", publishCalls)
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "douyin", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if platforms := publishBody["platforms"].([]interface{}); platforms[0] != "抖音" {
		t.Fatalf("expected platform key to map to Chinese platform name, got %+v", platforms)
	}
}

func TestPublishCommandPreservesScheduledTimeMilliseconds(t *testing.T) {
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}

	got := publishBody["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})["scheduledTime"]
	if got != float64(1760000000000) {
		t.Fatalf("expected scheduledTime milliseconds in publish body, got %#v", got)
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音,知乎", payloadPath})
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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
				"key":      "video_oss_key",
				"width":    float64(1080),
				"height":   float64(1920),
				"size":     float64(52428800),
				"duration": float64(30),
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

func TestPublishCommandAutoExtractsResourceMetadataFromLocalSourceFields(t *testing.T) {
	withRepoRoot(t)
	imagePath := filepath.Join(t.TempDir(), "cover.png")
	if err := os.WriteFile(imagePath, testPNGBytes(t), 0o644); err != nil {
		t.Fatal(err)
	}
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "imageText",
		"platforms":      []interface{}{"抖音"},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_001",
					"cover": map[string]interface{}{
						"key":    "uploaded/cover.png",
						"source": imagePath,
					},
					"coverKey": "uploaded/cover.png",
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "夏日穿搭",
						"description": "今日穿搭分享",
						"images": []interface{}{
							map[string]interface{}{
								"key":    "uploaded/cover.png",
								"source": imagePath,
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"imageText", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}

	form := publishBody["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
	cover := form["cover"].(map[string]interface{})
	if cover["width"] != float64(1) || cover["height"] != float64(1) {
		t.Fatalf("expected cover metadata to be auto-extracted, got %+v", cover)
	}
	if _, exists := cover["source"]; exists {
		t.Fatalf("expected source helper field to be stripped before publish, got %+v", cover)
	}
	imageItem := form["images"].([]interface{})[0].(map[string]interface{})
	if imageItem["width"] != float64(1) || imageItem["height"] != float64(1) {
		t.Fatalf("expected image metadata to be auto-extracted, got %+v", imageItem)
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

func TestPublishCommandRejectsInstagramVideoKeyWithChineseCharacters(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "video",
		"platforms":      []interface{}{"Instagram"},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_instagram_1",
					"video": map[string]interface{}{
						"key":      "yfb/test/t-68db/飞书20250424-172618.mp4",
						"size":     float64(1024),
						"width":    float64(1080),
						"height":   float64(1920),
						"duration": float64(30),
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
						"description": "instagram reel",
					},
				},
			},
		},
	})

	var publishCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_instagram_1", "name": "Instagram账号", "status": 1},
				},
			})
		case "/taskSets/v2":
			publishCalls++
			t.Fatal("publish API should not be called when instagram media key contains Chinese characters")
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "Instagram", payloadPath})
	if err == nil {
		t.Fatal("expected instagram media key validation error")
	}
	if !strings.Contains(err.Error(), "ASCII-only uploaded video key") {
		t.Fatalf("unexpected error: %v", err)
	}
	if publishCalls != 0 {
		t.Fatalf("expected no publish call, got %d", publishCalls)
	}
}

func TestPublishDryRunAutoBuildsOuterEnvelopeFromPublishArgs(t *testing.T) {
	withRepoRoot(t)
	configureEmptyConfig(t)
	service := publishflow.NewService(testRuntime(t))
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
	if result.InferredFields["desc"].SourcePath != "publishArgs.accountForms[0].contentPublishForm.description" {
		t.Fatalf("expected dry-run inferred desc source, got %+v", result.InferredFields)
	}
	if result.InferredFields["coverKey"].SourcePath != "publishArgs.coverKey" {
		t.Fatalf("expected dry-run inferred coverKey source, got %+v", result.InferredFields)
	}
}

func TestPublishDryRunChecksCloudProxyWhenAPIKeyConfigured(t *testing.T) {
	withRepoRoot(t)
	payload := validPublishPayload()
	payload["platforms"] = []interface{}{"视频号"}
	form := payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
	form["platformAccountId"] = "acc_shipinhao_1"
	cpf := form["contentPublishForm"].(map[string]interface{})
	cpf["createType"] = float64(2)
	cpf["pubType"] = float64(1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_shipinhao_1", "name": "视频号账号", "status": 1},
				},
			})
		case "/taskSets/v2":
			t.Fatal("dry-run should not call publish API")
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	service := publishflow.NewService(testRuntime(t))
	_, err := service.DryRun(publishflow.ExecuteInput{
		PublishType:   "video",
		PlatformInput: "视频号",
		Payload:       payload,
	})
	if err == nil || !strings.Contains(err.Error(), "Cloud publish preflight failed") {
		t.Fatalf("expected cloud proxy preflight error, got %v", err)
	}
}

func TestPublishDryRunMarksPlatformDraftSeparatelyFromYixiaoerDraft(t *testing.T) {
	withRepoRoot(t)
	configureEmptyConfig(t)
	service := publishflow.NewService(testRuntime(t))
	result, err := service.DryRun(publishflow.ExecuteInput{
		PublishType:   "imageText",
		PlatformInput: "shipinhao",
		Payload: map[string]interface{}{
			"action":      "publish",
			"publishType": "imageText",
			"platforms":   []interface{}{"视频号"},
			"publishArgs": map[string]interface{}{
				"accountForms": []interface{}{
					map[string]interface{}{
						"platformAccountId": "acc_001",
						"coverKey":          "uploaded/cover.png",
						"cover": map[string]interface{}{
							"key":    "uploaded/cover.png",
							"size":   float64(512),
							"width":  float64(1080),
							"height": float64(1440),
							"format": "png",
						},
						"contentPublishForm": map[string]interface{}{
							"formType": "task",
							"pubType":  float64(0),
							"title":    "标题",
							"images": []interface{}{
								map[string]interface{}{
									"key":    "uploaded/cover.png",
									"size":   float64(512),
									"width":  float64(1080),
									"height": float64(1440),
									"format": "png",
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.PlatformDraft {
		t.Fatalf("expected platformDraft=true when pubType=0, got %#v", result)
	}
	if result.YixiaoerDraft {
		t.Fatalf("expected yixiaoerDraft=false for platform draft publish, got %#v", result)
	}
}

func TestPublishCommandUsesLocalFlagsLikeNodeExample(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, validPublishPayload())

	var publishCalls int
	var publishBody map[string]interface{}
	server := publishTestServer(t, 1, &publishCalls, &publishBody)
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	cmd := newPublishCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"video", "抖音", payloadPath, "--publish-channel", "local", "--client-id", "flag_client_1"})
	if err := cmd.Execute(); err != nil {
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

	cmd := newPublishCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"video", "抖音", payloadPath, "--auto-fallback-local"})
	if err := cmd.Execute(); err != nil {
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"imageText", "快手", payloadPath})
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"imageText", "小红书", payloadPath})
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

func TestPublishCommandCopiesImageTextCoverFieldsFromContentPublishFormToAccountForm(t *testing.T) {
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
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "夏日穿搭",
						"description": "<p>今日穿搭分享</p>",
						"cover": map[string]interface{}{
							"key":    "uploaded/cover.png",
							"size":   float64(512),
							"width":  float64(1080),
							"height": float64(1080),
							"format": "png",
						},
						"coverKey": "uploaded/cover.png",
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"imageText", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}

	form := publishBody["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
	if form["coverKey"] != "uploaded/cover.png" {
		t.Fatalf("expected account form coverKey copied from contentPublishForm, got %+v", form)
	}
	cover, _ := form["cover"].(map[string]interface{})
	if cover["key"] != "uploaded/cover.png" {
		t.Fatalf("expected account form cover copied from contentPublishForm, got %+v", form)
	}
}

func TestPublishCommandNormalizesDescriptionTopics(t *testing.T) {
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
						"description": "今日穿搭分享 #穿搭 #夏日",
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"imageText", "抖音", payloadPath})
	if err != nil {
		t.Fatal(err)
	}

	expected := `<p>今日穿搭分享</p><p><topic text="穿搭">#穿搭</topic><topic text="夏日">#夏日</topic></p>`
	args := publishBody["publishArgs"].(map[string]interface{})
	cpf := args["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if cpf["description"] != expected {
		t.Fatalf("expected contentPublishForm.description topic HTML, got %+v", cpf["description"])
	}
	if _, exists := args["content"]; exists {
		t.Fatalf("did not expect description topic rule to synthesize publishArgs.content, got %+v", args["content"])
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "抖音", payloadPath})
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

func TestPublishDryRunReportsDynamicFieldNormalizations(t *testing.T) {
	withRepoRoot(t)
	configureEmptyConfig(t)
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			if got := r.URL.Query().Get("platform"); got != "抖音" {
				t.Fatalf("unexpected platform query: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_001", "name": "账号", "status": 1},
				},
			})
		case "/platform-accounts/acc_001/entitlements":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"shopping_cart": true},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	service := publishflow.NewService(testRuntime(t))
	result, err := service.DryRun(publishflow.ExecuteInput{
		PublishType:   "video",
		PlatformInput: "抖音",
		Payload:       payload,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.RemoteChecks {
		t.Fatal("expected shopping-cart dry-run to perform a remote entitlement check")
	}
	for _, action := range []string{"rename_field", "wrap_data", "derive_images"} {
		if !hasNormalizationAction(result.Normalizations, action) {
			t.Fatalf("expected dry-run normalization action %q, got %+v", action, result.Normalizations)
		}
	}
}

func hasNormalizationAction(events []publishmod.NormalizationEvent, action string) bool {
	for _, event := range events {
		if event.Action == action {
			return true
		}
	}
	return false
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"imageText", "小红书", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishBody["publishType"] != "imageText" {
		t.Fatalf("expected publishType imageText, got %+v", publishBody["publishType"])
	}
}

func TestPublishCommandKeepsArticleContentOnlyUnderPublishArgs(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "article",
		"platforms":      []interface{}{"知乎"},
		"desc":           "文章任务描述",
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"content": "<p>文章正文</p>",
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_zhihu_1",
					"cover": map[string]interface{}{
						"key":    "cover-key",
						"size":   float64(512),
						"width":  float64(1080),
						"height": float64(1080),
					},
					"coverKey": "cover-key",
					"contentPublishForm": map[string]interface{}{
						"formType": "task",
						"title":    "知乎文章标题示例一",
						"pubType":  float64(1),
					},
				},
			},
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_zhihu_1", "name": "文章账号", "status": 1},
				},
			})
		case "/taskSets/v2":
			publishCalls++
			if err := json.NewDecoder(r.Body).Decode(&publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_article_1"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := newPublishCmd().RunE(testCobraCommand(), []string{"article", "知乎", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	args := publishBody["publishArgs"].(map[string]interface{})
	if args["content"] != "<p>文章正文</p>" {
		t.Fatalf("expected article content under publishArgs, got %#v", args["content"])
	}
	cpf := args["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if _, exists := cpf["content"]; exists {
		t.Fatalf("did not expect article content inside contentPublishForm publish body, got %+v", cpf)
	}
	if publishBody["desc"] != "文章任务描述" {
		t.Fatalf("expected top-level desc to be preserved, got %#v", publishBody["desc"])
	}
}

func TestPublishCommandMaterializesArticleContentImageURLs(t *testing.T) {
	withRepoRoot(t)
	imageBytes := testPNGBytes(t)
	sourcePath := "/external/photo.png"
	stableURL := "https://oss-v2.yixiaoer.cn/materialized/photo.png"
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "article",
		"platforms":      []interface{}{"知乎"},
		"desc":           "文章任务描述",
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"content": `<p>文章正文配图<img src="SOURCE_URL" alt="配图"></p>`,
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_zhihu_1",
					"cover": map[string]interface{}{
						"key":    "cover-key",
						"size":   float64(512),
						"width":  float64(1080),
						"height": float64(1080),
					},
					"coverKey": "cover-key",
					"contentPublishForm": map[string]interface{}{
						"formType": "task",
						"title":    "知乎文章标题示例一",
						"pubType":  float64(1),
					},
				},
			},
		},
	})

	var uploadURLCalls int
	var stableURLCalls int
	var publishBody map[string]interface{}
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{{"platformAccountId": "acc_zhihu_1", "name": "文章账号", "status": 1}},
			})
		case sourcePath:
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write(imageBytes)
		case "/storages/material-library/upload-url":
			uploadURLCalls++
			if got := r.URL.Query().Get("fileKey"); got != "photo.png" {
				t.Fatalf("unexpected fileKey: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"serviceUrl": server.URL + "/oss/photo.png",
					"key":        "uploaded/photo.png",
				},
			})
		case "/oss/photo.png":
			if r.Method != http.MethodPut {
				t.Fatalf("unexpected upload method: %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
		case "/storages/material-library/stable-url":
			stableURLCalls++
			if got := r.URL.Query().Get("fileKey"); got != "uploaded/photo.png" {
				t.Fatalf("unexpected stable-url fileKey: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": stableURL})
		case "/taskSets/v2":
			if err := json.NewDecoder(r.Body).Decode(&publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"taskSetId": "task_set_article_1"}})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	replacePayloadFileText(t, payloadPath, "SOURCE_URL", server.URL+sourcePath)
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := newPublishCmd().RunE(testCobraCommand(), []string{"article", "知乎", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if uploadURLCalls != 1 || stableURLCalls != 1 {
		t.Fatalf("expected one upload and one stable-url call, got upload=%d stable=%d", uploadURLCalls, stableURLCalls)
	}
	content := publishBody["publishArgs"].(map[string]interface{})["content"].(string)
	if !strings.Contains(content, `src="`+stableURL+`"`) {
		t.Fatalf("expected article content img src to be materialized, got %s", content)
	}
	if strings.Contains(content, server.URL+sourcePath) {
		t.Fatalf("expected original external image URL to be replaced, got %s", content)
	}
}

func TestPublishCommandAutoCompressesShipinhaoCoverBeforePublish(t *testing.T) {
	withRepoRoot(t)
	coverPath := filepath.Join(t.TempDir(), "cover.png")
	if err := os.WriteFile(coverPath, noisyPNGBytes(t, 1400, 1000), 0o644); err != nil {
		t.Fatal(err)
	}
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "video",
		"platforms":      []interface{}{"视频号"},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_shipinhao_1",
					"cover": map[string]interface{}{
						"key":    "cover-key",
						"size":   float64(700 * 1024),
						"width":  float64(1080),
						"height": float64(1920),
						"source": coverPath,
					},
					"coverKey": "cover-key",
					"video": map[string]interface{}{
						"key":      "video-key",
						"size":     float64(1024),
						"width":    float64(1080),
						"height":   float64(1920),
						"duration": float64(30),
					},
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "视频号标题",
						"description": "视频号正文",
						"createType":  float64(2),
						"pubType":     float64(1),
					},
				},
			},
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	var uploadBody []byte
	var uploadContentType string
	var requestedFileKey string
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			if got := r.URL.Query().Get("platform"); got != "视频号" {
				t.Fatalf("unexpected platform query: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{{"platformAccountId": "acc_shipinhao_1", "name": "视频号账号", "status": 1, "proxyId": "proxy_1"}},
			})
		case "/storages/cloud-publish/upload-url":
			requestedFileKey = r.URL.Query().Get("fileKey")
			if got := r.URL.Query().Get("contentType"); got != "image/jpeg" {
				t.Fatalf("unexpected compressed content type query: %s", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"serviceUrl": server.URL + "/oss/cover.jpg",
					"key":        "uploaded/cover.jpg",
				},
			})
		case "/oss/cover.jpg":
			uploadContentType = r.Header.Get("Content-Type")
			var err error
			uploadBody, err = io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			_ = r.Body.Close()
			w.WriteHeader(http.StatusOK)
		case "/taskSets/v2":
			publishCalls++
			if err := json.NewDecoder(r.Body).Decode(&publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_shipinhao_1"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "视频号", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if requestedFileKey != "cover.jpg" {
		t.Fatalf("expected compressed cover upload to use jpg key, got %q", requestedFileKey)
	}
	if uploadContentType != "image/jpeg" {
		t.Fatalf("expected compressed upload content type image/jpeg, got %q", uploadContentType)
	}
	if len(uploadBody) > 512*1024 {
		t.Fatalf("expected compressed upload body <= 512KB, got %d", len(uploadBody))
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	args := publishBody["publishArgs"].(map[string]interface{})
	form := args["accountForms"].([]interface{})[0].(map[string]interface{})
	cover := form["cover"].(map[string]interface{})
	if cover["key"] != "uploaded/cover.jpg" {
		t.Fatalf("expected published cover to use compressed upload key, got %+v", cover)
	}
	if cover["size"] != float64(len(uploadBody)) {
		t.Fatalf("expected published cover size to match compressed upload body, got %+v", cover)
	}
	if _, exists := cover["source"]; exists {
		t.Fatalf("did not expect source helper to reach publish body, got %+v", cover)
	}
	if form["coverKey"] != "uploaded/cover.jpg" {
		t.Fatalf("expected coverKey to follow compressed upload key, got %+v", form["coverKey"])
	}
	if publishBody["coverKey"] != "uploaded/cover.jpg" {
		t.Fatalf("expected top-level coverKey to follow compressed upload key, got %+v", publishBody["coverKey"])
	}
}

func TestPublishCommandStopsWhenArticleContentImageMaterializationFails(t *testing.T) {
	withRepoRoot(t)
	sourcePath := "/external/blocked.png"
	payloadPath := writePublishPayload(t, articlePayloadWithContentImage("SOURCE_URL"))

	var publishCalls int
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{{"platformAccountId": "acc_zhihu_1", "name": "文章账号", "status": 1}},
			})
		case sourcePath:
			http.Error(w, "forbidden", http.StatusForbidden)
		case "/storages/proxy-url":
			http.Error(w, "proxy failed", http.StatusBadGateway)
		case "/taskSets/v2":
			publishCalls++
			t.Fatal("publish must not continue without explicit confirmation")
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	replacePayloadFileText(t, payloadPath, "SOURCE_URL", server.URL+sourcePath)
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := runPublish(testCobraCommand(), []string{"article", "知乎", payloadPath}, publishOptions{})
	if err == nil {
		t.Fatal("expected materialization confirmation error")
	}
	var typed *yxerrors.Error
	if !errors.As(err, &typed) {
		t.Fatalf("expected structured error, got %T: %v", err, err)
	}
	if typed.Category != "article_content_image_materialization_confirmation" {
		t.Fatalf("unexpected error category: %+v", typed)
	}
	if !strings.Contains(typed.Hint, "如确认可以保留原图地址继续发布") {
		t.Fatalf("expected confirmation hint, got %q", typed.Hint)
	}
	if !strings.Contains(typed.NextCommand, "--continue-on-content-image-error") {
		t.Fatalf("expected next command to include confirmation flag, got %q", typed.NextCommand)
	}
	if publishCalls != 0 {
		t.Fatalf("expected publish not to be called, got %d", publishCalls)
	}
}

func TestPublishCommandContinuesWhenArticleContentImageMaterializationFailureIsConfirmed(t *testing.T) {
	withRepoRoot(t)
	sourcePath := "/external/blocked.png"
	payloadPath := writePublishPayload(t, articlePayloadWithContentImage("SOURCE_URL"))

	var publishBody map[string]interface{}
	var uploadURLCalls int
	var publishCalls int
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{{"platformAccountId": "acc_zhihu_1", "name": "文章账号", "status": 1}},
			})
		case sourcePath:
			http.Error(w, "forbidden", http.StatusForbidden)
		case "/storages/proxy-url":
			http.Error(w, "proxy failed", http.StatusBadGateway)
		case "/storages/material-library/upload-url":
			uploadURLCalls++
			t.Fatal("upload-url should not be requested when source and proxy downloads fail")
		case "/taskSets/v2":
			publishCalls++
			if err := json.NewDecoder(r.Body).Decode(&publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"taskSetId": "task_set_article_1"}})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	originalURL := server.URL + sourcePath
	replacePayloadFileText(t, payloadPath, "SOURCE_URL", originalURL)
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := runPublish(testCobraCommand(), []string{"article", "知乎", payloadPath}, publishOptions{ContinueOnContentImageError: true})
	if err != nil {
		t.Fatal(err)
	}
	if uploadURLCalls != 0 || publishCalls != 1 {
		t.Fatalf("expected publish to continue without upload-url, got upload=%d publish=%d", uploadURLCalls, publishCalls)
	}
	content := publishBody["publishArgs"].(map[string]interface{})["content"].(string)
	if !strings.Contains(content, originalURL) {
		t.Fatalf("expected original image URL to remain when continuing after failure, got %s", content)
	}
}

func TestPublishDryRunReportsArticleContentImageMaterializationWithoutUploading(t *testing.T) {
	withRepoRoot(t)
	externalURL := "https://example.invalid/photo.png"
	payload := map[string]interface{}{
		"action":         "publish",
		"publishType":    "article",
		"platforms":      []interface{}{"知乎"},
		"desc":           "文章任务描述",
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"content": `<p>文章正文配图<img src="` + externalURL + `" alt="配图"></p>`,
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_zhihu_1",
					"cover": map[string]interface{}{
						"key":    "cover-key",
						"size":   float64(512),
						"width":  float64(1080),
						"height": float64(1080),
					},
					"coverKey": "cover-key",
					"contentPublishForm": map[string]interface{}{
						"formType": "task",
						"title":    "知乎文章标题示例一",
						"pubType":  float64(1),
					},
				},
			},
		},
	}

	uploadURLCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{{"platformAccountId": "acc_zhihu_1", "name": "文章账号", "status": 1}},
			})
		case "/storages/material-library/upload-url":
			uploadURLCalls++
			t.Fatalf("dry-run must not request upload url")
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	service := publishflow.NewService(testRuntime(t))
	result, err := service.DryRun(publishflow.ExecuteInput{
		PublishType:   "article",
		PlatformInput: "知乎",
		Payload:       payload,
	})
	if err != nil {
		t.Fatal(err)
	}
	if uploadURLCalls != 0 {
		t.Fatalf("expected dry-run not to upload, got %d upload calls", uploadURLCalls)
	}
	if len(result.ContentImages) != 1 {
		t.Fatalf("expected one content image materialization preview, got %+v", result.ContentImages)
	}
	if result.ContentImages[0].From != externalURL || result.ContentImages[0].Status != "would_materialize" {
		t.Fatalf("unexpected dry-run materialization event: %+v", result.ContentImages[0])
	}
}

func TestPublishCommandSupportsWeixinAccountArticlePlatformForms(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":      "publish",
		"publishType": "article",
		"platforms":   []interface{}{"微信公众号"},
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_weixin_1",
					"platformName":      "微信公众号",
				},
			},
			"platformForms": map[string]interface{}{
				"微信公众号": map[string]interface{}{
					"articles": []interface{}{
						map[string]interface{}{
							"title":   "公众号文章标题",
							"content": "<p>公众号文章正文</p>",
							"digest":  "公众号摘要",
							"type":    float64(1),
							"cover": map[string]interface{}{
								"key": "wx-cover-key",
								"raw": map[string]interface{}{},
							},
						},
					},
					"notifySubscribers": float64(1),
					"pubType":           float64(1),
				},
			},
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_weixin_1", "name": "公众号账号", "status": 1},
				},
			})
		case "/taskSets/v2":
			publishCalls++
			if err := json.NewDecoder(r.Body).Decode(&publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_weixin_article_1"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := newPublishCmd().RunE(testCobraCommand(), []string{"article", "微信公众号", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	args := publishBody["publishArgs"].(map[string]interface{})
	wxForm := args["platformForms"].(map[string]interface{})["微信公众号"].(map[string]interface{})
	if wxForm["pubType"] != float64(1) {
		t.Fatalf("expected weixin platform form to be forwarded, got %#v", wxForm)
	}
	if publishBody["coverKey"] != "wx-cover-key" {
		t.Fatalf("expected top-level coverKey synthesized from weixin article cover, got %#v", publishBody["coverKey"])
	}
	if publishBody["desc"] != "公众号摘要" {
		t.Fatalf("expected top-level desc synthesized from weixin article digest, got %#v", publishBody["desc"])
	}
}

func TestPublishCommandWeixinAccountImageTextDryRunAppliesDefaultsAndFirstCover(t *testing.T) {
	withRepoRoot(t)
	configureEmptyConfig(t)
	t.Setenv("YIXIAOER_API_KEY", "")
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "imageText",
		"platforms":      []interface{}{"WeiXinGongZhongHao"},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_weixin_image_text_1",
					"images": []interface{}{
						map[string]interface{}{"key": "wx-image-1", "size": float64(1024), "width": float64(1080), "height": float64(1440)},
						map[string]interface{}{"key": "wx-image-2", "size": float64(1024), "width": float64(1080), "height": float64(1440)},
					},
					"contentPublishForm": map[string]interface{}{
						"formType": "task",
						"title":    "公众号图文标题",
						"desc":     "<p>公众号图文描述</p>",
					},
				},
			},
		},
	})

	var out bytes.Buffer
	cmd := newPublishCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"imageText", "WeiXinGongZhongHao", payloadPath, "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "publish.dry-run" {
		t.Fatalf("expected dry-run action, got %#v", response["action"])
	}
	data := response["data"].(map[string]interface{})
	request := data["request"].(map[string]interface{})
	if request["desc"] != "<p>公众号图文描述</p>" {
		t.Fatalf("expected top-level desc to be inferred from the WeChat imageText form, got %#v", request["desc"])
	}
	args := request["publishArgs"].(map[string]interface{})
	form := args["accountForms"].([]interface{})[0].(map[string]interface{})
	cpf := form["contentPublishForm"].(map[string]interface{})
	for field, want := range map[string]interface{}{
		"notifySubscribers": float64(0),
		"sex":               float64(0),
		"needOpenComment":   float64(0),
		"statement":         float64(0),
		"disableRecommend":  float64(0),
		"pubType":           float64(1),
	} {
		if got := cpf[field]; got != want {
			t.Fatalf("dry-run account content form %s = %#v, want %#v", field, got, want)
		}
	}
	if _, exists := args["platformForms"]; exists {
		t.Fatalf("did not expect a platform form in the dry-run request: %#v", args)
	}
	if form["coverKey"] != "wx-image-1" {
		t.Fatalf("expected first image coverKey in dry-run request, got %#v", form["coverKey"])
	}
	if form["cover"].(map[string]interface{})["key"] != "wx-image-1" {
		t.Fatalf("expected first image cover in dry-run request, got %#v", form["cover"])
	}
}

func TestPublishCommandAcceptsBaijiahaoImageTextPayload(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "imageText",
		"platforms":      []interface{}{"百家号"},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_bjh_1",
					"cover": map[string]interface{}{
						"key":    "cover-key",
						"size":   float64(512),
						"width":  float64(1080),
						"height": float64(1080),
					},
					"coverKey": "cover-key",
					"contentPublishForm": map[string]interface{}{
						"formType":      "task",
						"title":         "百家号图文标题",
						"description":   "<p>百家号图文内容</p>",
						"pubType":       float64(0),
						"declaration":   float64(0),
						"scheduledTime": float64(1760000000000),
						"images": []interface{}{
							map[string]interface{}{
								"key":    "image-key",
								"size":   float64(512),
								"width":  float64(1080),
								"height": float64(1080),
								"format": "jpg",
							},
						},
					},
				},
			},
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_bjh_1", "name": "百家号图文账号", "status": 1},
				},
			})
		case "/taskSets/v2":
			publishCalls++
			if err := json.NewDecoder(r.Body).Decode(&publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_bjh_image_text_1"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := newPublishCmd().RunE(testCobraCommand(), []string{"imageText", "百家号", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	cpf := publishBody["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if cpf["pubType"] != float64(0) || cpf["declaration"] != float64(0) {
		t.Fatalf("expected baijiahao imageText fields to survive publish normalization, got %+v", cpf)
	}
	if cpf["scheduledTime"] != float64(1760000000000) {
		t.Fatalf("expected scheduledTime to remain in milliseconds, got %+v", cpf["scheduledTime"])
	}
}

func TestPublishCommandAcceptsFirstImageCoverImageTextPayloadWithoutExternalCover(t *testing.T) {
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
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "小红书图文标题",
						"description": "小红书图文内容",
						"visibleType": float64(0),
						"images": []interface{}{
							map[string]interface{}{
								"key":    "first-image-key",
								"size":   float64(512),
								"width":  float64(1080),
								"height": float64(1440),
								"format": "jpg",
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

	err := newPublishCmd().RunE(testCobraCommand(), []string{"imageText", "小红书", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	form := publishBody["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
	if form["coverKey"] != "first-image-key" {
		t.Fatalf("expected internal coverKey derived from first image, got %+v", form)
	}
	if cover := form["cover"].(map[string]interface{}); cover["key"] != "first-image-key" {
		t.Fatalf("expected internal cover derived from first image, got %+v", form)
	}
}

func TestPublishCommandAcceptsSouhuhaoVideoPayload(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "video",
		"platforms":      []interface{}{"搜狐号"},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"video": map[string]interface{}{
				"key":      "video-key",
				"size":     float64(1024),
				"width":    float64(1080),
				"height":   float64(1920),
				"duration": float64(30),
			},
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_sh_1",
					"cover": map[string]interface{}{
						"key":    "cover-key",
						"size":   float64(512),
						"width":  float64(1080),
						"height": float64(1080),
					},
					"coverKey": "cover-key",
					"contentPublishForm": map[string]interface{}{
						"formType":    "task",
						"title":       "搜狐号视频标题示例",
						"description": "这是搜狐号视频描述内容。",
						"tags":        []interface{}{"科技"},
						"declaration": float64(2),
						"pubType":     float64(1),
						"category": []interface{}{
							map[string]interface{}{
								"id":   "1",
								"text": "科技",
								"raw":  map[string]interface{}{"id": "1"},
							},
						},
					},
				},
			},
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/platform-accounts/acc_sh_1/categories":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"dataList": []map[string]interface{}{
						{
							"yixiaoerId":   "1",
							"yixiaoerName": "科技",
							"raw":          map[string]interface{}{"id": 1, "name": "科技"},
						},
					},
				},
			})
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_sh_1", "name": "搜狐号视频账号", "status": 1},
				},
			})
		case "/taskSets/v2":
			publishCalls++
			if err := json.NewDecoder(r.Body).Decode(&publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_souhuhao_video_1"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := newPublishCmd().RunE(testCobraCommand(), []string{"video", "搜狐号", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	cpf := publishBody["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if cpf["pubType"] != float64(1) || cpf["declaration"] != float64(2) {
		t.Fatalf("expected souhuhao video fields to survive publish normalization, got %+v", cpf)
	}
	category := cpf["category"].([]interface{})[0].(map[string]interface{})
	canonical := category["raw"].(map[string]interface{})
	platformRaw := canonical["raw"].(map[string]interface{})
	if category["id"] != "1" || canonical["yixiaoerId"] != "1" || platformRaw["id"] != float64(1) {
		t.Fatalf("expected nested Sohu category wire payload, got %#v", category)
	}
}

func TestPublishCommandAcceptsToutiaohaoArticleExtendedFields(t *testing.T) {
	withRepoRoot(t)
	payloadPath := writePublishPayload(t, map[string]interface{}{
		"action":         "publish",
		"publishType":    "article",
		"platforms":      []interface{}{"头条号"},
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"content": "<p>文章正文</p>",
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_tt_1",
					"cover": map[string]interface{}{
						"key":    "cover-key",
						"size":   float64(512),
						"width":  float64(1080),
						"height": float64(1080),
					},
					"coverKey": "cover-key",
					"contentPublishForm": map[string]interface{}{
						"formType":      "task",
						"title":         "头条号文章标题",
						"pubType":       float64(0),
						"isFirst":       true,
						"advertisement": float64(3),
						"declaration":   float64(3),
						"scheduledTime": float64(1760000000000),
						"location": map[string]interface{}{
							"yixiaoerId":   "loc_1",
							"yixiaoerName": "上海",
							"raw":          map[string]interface{}{"id": "loc_1"},
						},
					},
				},
			},
		},
	})

	var publishCalls int
	var publishBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/platform/accounts":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_tt_1", "name": "头条号文章账号", "status": 1},
				},
			})
		case "/taskSets/v2":
			publishCalls++
			if err := json.NewDecoder(r.Body).Decode(&publishBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"taskSetId": "task_set_toutiaohao_article_1"},
			})
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()
	configureAPIKey(t, "test-key")
	useTestAPIBaseURL(t, server.URL)

	err := newPublishCmd().RunE(testCobraCommand(), []string{"article", "头条号", payloadPath})
	if err != nil {
		t.Fatal(err)
	}
	if publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", publishCalls)
	}
	cpf := publishBody["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if cpf["advertisement"] != float64(3) || cpf["declaration"] != float64(3) || cpf["isFirst"] != true {
		t.Fatalf("expected toutiaohao article fields to survive publish normalization, got %+v", cpf)
	}
	if cpf["scheduledTime"] != float64(1760000000000) {
		t.Fatalf("expected scheduledTime to remain in milliseconds, got %+v", cpf["scheduledTime"])
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
					{"platformAccountId": "acc_001", "name": "抖音图文账号", "status": 1},
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

func noisyPNGBytes(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	seed := uint32(1)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			seed = seed*1664525 + 1013904223
			r := uint8(seed >> 24)
			seed = seed*1664525 + 1013904223
			g := uint8(seed >> 24)
			seed = seed*1664525 + 1013904223
			b := uint8(seed >> 24)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
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
		case "/platform-accounts/acc_001/entitlements":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"shopping_cart": true},
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
					"key":      "video-key",
					"size":     float64(1024),
					"width":    float64(1080),
					"height":   float64(1920),
					"duration": float64(30),
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

func articlePayloadWithContentImage(imageURL string) map[string]interface{} {
	return map[string]interface{}{
		"action":         "publish",
		"publishType":    "article",
		"platforms":      []interface{}{"知乎"},
		"desc":           "文章任务描述",
		"publishChannel": "cloud",
		"publishArgs": map[string]interface{}{
			"content": `<p>文章正文配图<img src="` + imageURL + `" alt="配图"></p>`,
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_zhihu_1",
					"cover": map[string]interface{}{
						"key":    "cover-key",
						"size":   float64(512),
						"width":  float64(1080),
						"height": float64(1080),
					},
					"coverKey": "cover-key",
					"contentPublishForm": map[string]interface{}{
						"formType": "task",
						"title":    "知乎文章标题示例一",
						"pubType":  float64(1),
					},
				},
			},
		},
	}
}

func replacePayloadFileText(t *testing.T, path, old, new string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.ReplaceAll(string(raw), old, new)
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
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

func configureEmptyConfig(t *testing.T) {
	t.Helper()
	t.Setenv("YIXIAOER_CONFIG", filepath.Join(t.TempDir(), "yxer-config.json"))
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

func testRuntime(t *testing.T) *app.Runtime {
	t.Helper()
	rt, err := app.Load()
	if err != nil {
		t.Fatal(err)
	}
	return rt
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
