package draft

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/docx"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

func TestRequireHorizontalCoverRejectsPortraitAndSquareImages(t *testing.T) {
	for _, tc := range []struct {
		name          string
		width, height int
	}{
		{name: "portrait", width: 901, height: 1600},
		{name: "square", width: 1000, height: 1000},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := RequireHorizontalCover("cover.jpg", tc.width, tc.height)
			if err == nil {
				t.Fatal("expected non-horizontal cover to be rejected")
			}
			var typed *yxerrors.Error
			if !errors.As(err, &typed) {
				t.Fatalf("expected structured error, got %T", err)
			}
			if typed.Category != "cover_orientation" {
				t.Fatalf("expected cover_orientation category, got %+v", typed)
			}
		})
	}
}

func TestBuildWeixinArticlePayloadDefaultsToNonOriginal(t *testing.T) {
	plan := importPlan{
		Input: ImportInput{
			AccountID: "acc_1",
			DraftID:   "draft-1",
			Author:    "蚁小二",
		},
		Document: docx.Document{HTML: "<p>正文</p>"},
		Title:    "文章标题",
		Digest:   "文章摘要",
	}
	cover := uploadResultResource(api.UploadResult{Key: "cover-key", Bucket: "cloud-publish", Width: 900, Height: 383}, "https://oss-v2.yixiaoer.cn/cover.jpg")
	payload := buildWeixinArticlePayload(plan, cover)
	article := payload["publishArgs"].(map[string]interface{})["platformForms"].(map[string]interface{})[weixinArticlePlatform].(map[string]interface{})["articles"].([]interface{})[0].(map[string]interface{})
	if article["type"] != float64(0) {
		t.Fatalf("expected non-original article by default, got %#v", article["type"])
	}
	if article["author"] != "蚁小二" {
		t.Fatalf("expected author to be retained, got %#v", article["authorName"])
	}
	if _, exists := article["authorName"]; exists {
		t.Fatalf("did not expect ordinary author to populate original author, got %#v", article["authorName"])
	}
	platformForm := payload["publishArgs"].(map[string]interface{})["platformForms"].(map[string]interface{})[weixinArticlePlatform].(map[string]interface{})
	if platformForm["sex"] != float64(0) {
		t.Fatalf("expected WeChat article platform form to include sex=0, got %#v", platformForm["sex"])
	}
	if article["cover"].(map[string]interface{})["pathOrUrl"] != "https://oss-v2.yixiaoer.cn/cover.jpg" {
		t.Fatalf("expected stable cover URL in article cover, got %#v", article["cover"])
	}
	if article["cover"].(map[string]interface{})["url"] != "https://oss-v2.yixiaoer.cn/cover.jpg" {
		t.Fatalf("expected editor cover URL in article cover, got %#v", article["cover"])
	}

	plan.Input.DraftID = "draft-1"
	payload = buildWeixinArticlePayload(plan, cover)
	payload["taskSetId"] = plan.Input.DraftID
	if payload["taskSetId"] != "draft-1" {
		t.Fatalf("expected existing draft ID to be preserved, got %#v", payload["taskSetId"])
	}
}

func TestBuildWeixinArticlePayloadSeparatesOriginalAuthor(t *testing.T) {
	plan := importPlan{
		Input: ImportInput{
			Author:         "普通作者",
			Original:       true,
			OriginalAuthor: "原创作者",
		},
		Document: docx.Document{HTML: "<p>正文</p>"},
		Title:    "文章标题",
	}
	article := buildWeixinArticlePayload(plan, map[string]interface{}{})["publishArgs"].(map[string]interface{})["platformForms"].(map[string]interface{})[weixinArticlePlatform].(map[string]interface{})["articles"].([]interface{})[0].(map[string]interface{})
	if article["author"] != "普通作者" || article["authorName"] != "原创作者" {
		t.Fatalf("expected ordinary and original author fields to remain separate, got %#v", article)
	}
}

func TestImportDOCXMaterializesCoverAndBodyImages(t *testing.T) {
	coverPath := writeHorizontalPNG(t, 900, 383)
	docxPath := writeImportTestDOCX(t)
	var draftBody map[string]interface{}
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/storages/") && strings.HasSuffix(r.URL.Path, "/stable-url"):
			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			bucket := parts[1]
			key := strings.TrimPrefix(r.URL.Query().Get("fileKey"), bucket+"/")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": "https://oss-v2.yixiaoer.cn/" + bucket + "/" + key})
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/storages/") && strings.HasSuffix(r.URL.Path, "/upload-url"):
			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
			bucket := parts[1]
			fileKey := r.URL.Query().Get("fileKey")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{
				"serviceUrl": server.URL + "/oss/" + bucket + "/" + fileKey,
				"key":        bucket + "/" + fileKey,
			}})
		case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/oss/"):
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodPut && r.URL.Path == "/taskSets/drafts":
			if err := json.NewDecoder(r.Body).Decode(&draftBody); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"id": "draft-1"}})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	}))
	defer server.Close()

	result, err := NewService(app.New(config.Config{APIKey: "test-key", APIURL: server.URL})).ImportDOCX(ImportInput{
		DocxPath:  docxPath,
		CoverPath: coverPath,
		AccountID: "account-1",
		DraftID:   "draft-existing",
		Author:    "蚁小二",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result["draft"].(map[string]interface{})["data"].(map[string]interface{})["id"] != "draft-1" {
		t.Fatalf("unexpected draft response: %+v", result)
	}
	cover := draftBody["cover"].(map[string]interface{})
	if cover["key"] == nil || cover["pathOrUrl"] == nil || cover["url"] == nil {
		t.Fatalf("expected cover key and stable URL, got %+v", cover)
	}
	if draftBody["taskSetId"] != "draft-existing" {
		t.Fatalf("expected existing draft ID in update request, got %#v", draftBody["taskSetId"])
	}
	args := draftBody["publishArgs"].(map[string]interface{})
	platformForms := args["platformForms"].(map[string]interface{})
	form := platformForms[weixinArticlePlatform].(map[string]interface{})
	article := form["articles"].([]interface{})[0].(map[string]interface{})
	content := article["content"].(string)
	if !strings.Contains(content, `src="https://oss-v2.yixiaoer.cn/material-library/`) {
		t.Fatalf("expected body image to use a stable material-library URL, got %s", content)
	}
	if strings.Contains(content, `src="image-001.jpeg"`) {
		t.Fatalf("expected local DOCX image reference to be replaced, got %s", content)
	}
}

func writeHorizontalPNG(t *testing.T, width, height int) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cover.png")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{R: 20, G: 120, B: 200, A: 255})
		}
	}
	if err := png.Encode(file, img); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeImportTestDOCX(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "article.docx")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	archive := zip.NewWriter(file)
	entries := map[string]string{
		"word/document.xml":            `<?xml version="1.0"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><w:body><w:p><w:r><w:t>文章标题</w:t></w:r></w:p><w:p><w:r><w:t>正文</w:t></w:r><w:r><w:drawing><wp:inline><a:blip r:embed="rId1"/></wp:inline></w:drawing></w:r></w:p></w:body></w:document>`,
		"word/_rels/document.xml.rels": `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/image1.jpeg"/></Relationships>`,
	}
	for name, body := range entries {
		writer, writerErr := archive.Create(name)
		if writerErr != nil {
			t.Fatal(writerErr)
		}
		if _, writerErr = writer.Write([]byte(body)); writerErr != nil {
			t.Fatal(writerErr)
		}
	}
	writer, err := archive.Create("word/media/image1.jpeg")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("body-image")); err != nil {
		t.Fatal(err)
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}
