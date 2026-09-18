package draft

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/docx"
	publishflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const (
	weixinArticlePlatform       = "微信公众号"
	articleCoverBucket          = "cloud-publish"
	articleBodyImageBucket      = "material-library"
	maxWeixinArticleTitleRunes  = 64
	maxWeixinArticleDigestRunes = 129
	maxWeixinArticleAuthorRunes = 8
)

// ImportInput describes the DOCX-to-WeChat-article import contract exposed by
// the CLI. It mirrors the fields in the Yixiaoer article editor while keeping
// the actual resource upload and stable-URL resolution in this workflow.
type ImportInput struct {
	DocxPath       string
	CoverPath      string
	AccountID      string
	DraftID        string
	Title          string
	Digest         string
	Author         string
	OriginalAuthor string
	Original       bool
}

type importPlan struct {
	Input       ImportInput
	Document    docx.Document
	Title       string
	Digest      string
	CoverMeta   api.UploadResult
	CoverSource string
}

// PreviewImport parses the DOCX and inspects the cover without uploading or
// saving anything. It is the dry-run contract for `yxer draft import`.
func PreviewImport(input ImportInput) (map[string]interface{}, error) {
	plan, err := prepareImportPlan(input)
	if err != nil {
		return nil, err
	}
	defer plan.Document.Cleanup()

	return map[string]interface{}{
		"platform":  weixinArticlePlatform,
		"type":      "article",
		"accountId": plan.Input.AccountID,
		"source": map[string]interface{}{
			"docx":  plan.Document.SourcePath,
			"cover": plan.CoverSource,
		},
		"article": map[string]interface{}{
			"title":          plan.Title,
			"digest":         plan.Digest,
			"author":         plan.Input.Author,
			"originalAuthor": plan.Input.OriginalAuthor,
			"original":       plan.Input.Original,
			"paragraphCount": len(plan.Document.Paragraphs),
			"embeddedImages": imageNames(plan.Document.Images),
		},
		"cover": map[string]interface{}{
			"bucket":      articleCoverBucket,
			"source":      plan.CoverSource,
			"contentType": plan.CoverMeta.ContentType,
			"size":        plan.CoverMeta.Size,
			"width":       plan.CoverMeta.Width,
			"height":      plan.CoverMeta.Height,
			"orientation": "horizontal",
		},
		"contentImages": map[string]interface{}{
			"bucket": articleBodyImageBucket,
			"count":  len(plan.Document.Images),
			"mode":   "upload-and-rewrite-to-stable-url",
		},
		"draft": map[string]interface{}{
			"isDraft": true,
			"pubType": 0,
			"draftId": input.DraftID,
		},
	}, nil
}

// ImportDOCX uploads the horizontal cover, materializes every embedded DOCX
// image to a stable Yixiaoer URL, and saves the normalized article as an
// internal Yixiaoer draft.
func (s Service) ImportDOCX(input ImportInput) (map[string]interface{}, error) {
	plan, err := prepareImportPlan(input)
	if err != nil {
		return nil, err
	}
	defer plan.Document.Cleanup()
	if s.rt == nil || s.rt.Client == nil {
		return nil, yxerrors.Internal("draft DOCX import requires a configured API client", nil).
			WithCategory("docx_import")
	}

	coverUpload, err := s.rt.Client.Upload(plan.CoverSource, articleCoverBucket, true)
	if err != nil {
		return nil, err
	}
	coverURL, err := s.rt.Client.StableURL(articleCoverBucket, coverUpload.Key)
	if err != nil {
		return nil, yxerrors.Remote("failed to resolve uploaded article cover URL", map[string]interface{}{
			"bucket": articleCoverBucket,
			"key":    coverUpload.Key,
			"cause":  err.Error(),
		}).WithCategory("article_cover_stable_url").
			WithHint("封面已上传但无法取得可预览的稳定地址；请重试，或先运行 yxer upload 后检查 stable-url。")
	}

	cover := uploadResultResource(coverUpload, coverURL)
	payload := buildWeixinArticlePayload(plan, cover)
	if plan.Input.DraftID != "" {
		// The web editor sends taskSetId when saving an existing draft. Keeping
		// it at the envelope level lets this importer repair the current draft
		// without creating a second copy.
		payload["taskSetId"] = plan.Input.DraftID
	}
	body := publishflow.BuildDraftBody(payload)
	contentImages, err := publishflow.MaterializeArticleContentImages(s.rt.Client, body, plan.Document.BaseDir, false)
	if err != nil {
		return nil, err
	}
	draftResult, err := s.rt.Client.SaveDraft(body)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"platform":      weixinArticlePlatform,
		"type":          "article",
		"accountId":     plan.Input.AccountID,
		"title":         plan.Title,
		"cover":         cover,
		"contentImages": contentImages,
		"draft":         draftResult,
	}, nil
}

func prepareImportPlan(input ImportInput) (importPlan, error) {
	input.DocxPath = strings.TrimSpace(input.DocxPath)
	input.CoverPath = strings.TrimSpace(input.CoverPath)
	input.AccountID = strings.TrimSpace(input.AccountID)
	input.DraftID = strings.TrimSpace(input.DraftID)
	input.Title = strings.TrimSpace(input.Title)
	input.Digest = strings.TrimSpace(input.Digest)
	input.Author = strings.TrimSpace(input.Author)
	input.OriginalAuthor = strings.TrimSpace(input.OriginalAuthor)
	if input.DocxPath == "" {
		return importPlan{}, importUsage("DOCX import requires an article file", map[string]interface{}{"field": "docx"}).
			WithNextCommand("yxer draft import <article.docx> --cover <horizontal-cover.jpg> --account-id <online-account-id> --dry-run")
	}
	if input.CoverPath == "" {
		return importPlan{}, importUsage("DOCX import requires a horizontal cover file", map[string]interface{}{"field": "cover"}).
			WithNextCommand("yxer draft import <article.docx> --cover <horizontal-cover.jpg> --account-id <online-account-id> --dry-run")
	}
	if input.AccountID == "" {
		return importPlan{}, importUsage("DOCX import requires a WeChat account ID", map[string]interface{}{"field": "account-id"}).
			WithHint("账号 ID 必须来自 yxer accounts list 微信公众号 --status 1 --json，不要手写账号对象。").
			WithNextCommand("yxer accounts list 微信公众号 --status 1 --json")
	}
	if err := validateTextLength(input.Title, maxWeixinArticleTitleRunes, "title"); err != nil {
		return importPlan{}, err
	}
	if err := validateTextLength(input.Author, maxWeixinArticleAuthorRunes, "author"); err != nil {
		return importPlan{}, err
	}
	if err := validateTextLength(input.OriginalAuthor, maxWeixinArticleAuthorRunes, "original-author"); err != nil {
		return importPlan{}, err
	}
	if input.Original && input.OriginalAuthor == "" {
		return importPlan{}, importUsage("original articles require an original author", map[string]interface{}{"field": "original-author"}).
			WithHint("普通作者和原创作者是两个字段；请用 --author 填普通作者，用 --original-author 填原创作者。")
	}
	if err := validateTextLength(input.Digest, maxWeixinArticleDigestRunes, "digest"); err != nil {
		return importPlan{}, err
	}

	document, err := docx.Parse(input.DocxPath)
	if err != nil {
		return importPlan{}, err
	}
	cleanupDocument := true
	defer func() {
		if cleanupDocument {
			_ = document.Cleanup()
		}
	}()

	title := input.Title
	if title == "" {
		title = inferTitle(document.Paragraphs)
	}
	if title == "" {
		return importPlan{}, importUsage("DOCX does not contain an article title", map[string]interface{}{"path": document.SourcePath}).
			WithHint("请用 --title 指定微信公众号文章标题；标题最多 64 个字符。")
	}
	if err := validateTextLength(title, maxWeixinArticleTitleRunes, "title"); err != nil {
		return importPlan{}, err
	}
	digest := input.Digest
	if digest == "" {
		digest = inferDigest(document.Paragraphs, title)
	}

	coverMeta, _, err := api.InspectUpload(input.CoverPath, true)
	if err != nil {
		return importPlan{}, err
	}
	if err := RequireHorizontalCover(input.CoverPath, coverMeta.Width, coverMeta.Height); err != nil {
		return importPlan{}, err
	}

	cleanupDocument = false
	return importPlan{
		Input:       input,
		Document:    document,
		Title:       title,
		Digest:      truncateRunes(digest, maxWeixinArticleDigestRunes),
		CoverMeta:   coverMeta,
		CoverSource: input.CoverPath,
	}, nil
}

// RequireHorizontalCover is shared by import and payload preflight. A square
// or portrait image is intentionally rejected; the importer does not crop it
// silently because that could remove the subject selected by the user.
func RequireHorizontalCover(source string, width, height int) error {
	if width <= 0 || height <= 0 {
		return importUsage("article cover dimensions could not be detected", map[string]interface{}{
			"source": source,
			"width":  width,
			"height": height,
		}).WithCategory("cover_dimensions").
			WithHint("请传入可读取的横版图片；建议使用至少 900×383、宽度大于高度的封面。")
	}
	if width <= height {
		return importUsage("微信公众号文章封面必须是横版图片", map[string]interface{}{
			"source":      source,
			"width":       width,
			"height":      height,
			"orientation": "portrait-or-square",
		}).WithCategory("cover_orientation").
			WithHint("请改传横版封面，要求宽度大于高度；建议按公众号常用比例 900×383 准备，CLI 不会把竖版图强行裁剪。")
	}
	return nil
}

func buildWeixinArticlePayload(plan importPlan, cover map[string]interface{}) map[string]interface{} {
	article := map[string]interface{}{
		"title":   plan.Title,
		"content": plan.Document.HTML,
		"cover":   cover,
		"type":    float64(0),
	}
	if plan.Input.Author != "" {
		article["author"] = plan.Input.Author
	}
	if plan.Input.OriginalAuthor != "" {
		article["authorName"] = plan.Input.OriginalAuthor
	}
	if plan.Digest != "" {
		article["digest"] = plan.Digest
	}
	if plan.Input.Original {
		article["type"] = float64(1)
	}
	return map[string]interface{}{
		"action":      "publish",
		"publishType": "article",
		"platforms":   []interface{}{weixinArticlePlatform},
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": plan.Input.AccountID,
				},
			},
			"platformForms": map[string]interface{}{
				weixinArticlePlatform: map[string]interface{}{
					"articles":          []interface{}{article},
					"notifySubscribers": float64(0),
					"pubType":           float64(0),
					"sex":               float64(0),
				},
			},
		},
	}
}

func uploadResultResource(uploaded api.UploadResult, stableURL string) map[string]interface{} {
	resource := map[string]interface{}{
		"key":         uploaded.Key,
		"bucket":      uploaded.Bucket,
		"contentType": uploaded.ContentType,
		"size":        uploaded.Size,
		"format":      uploaded.Format,
		"raw":         map[string]interface{}{},
	}
	if uploaded.Width > 0 {
		resource["width"] = uploaded.Width
	}
	if uploaded.Height > 0 {
		resource["height"] = uploaded.Height
	}
	if stableURL != "" {
		// The standard upload resource contract uses pathOrUrl, while the
		// Yixiaoer WeChat article editor restores an article cover as an
		// ImageFormItem and reads its preview from cover.url. Keep both forms
		// so the saved draft remains valid for the CLI contract and editable in
		// the web editor.
		resource["pathOrUrl"] = stableURL
		resource["url"] = stableURL
	}
	return resource
}

func imageNames(images []docx.Image) []string {
	result := make([]string, 0, len(images))
	for _, image := range images {
		result = append(result, image.Name)
	}
	return result
}

func inferTitle(paragraphs []docx.Paragraph) string {
	for _, paragraph := range paragraphs {
		text := strings.TrimSpace(paragraph.Text)
		if text == "" || paragraph.HasImage || looksLikeLabel(text) {
			continue
		}
		return text
	}
	return ""
}

func inferDigest(paragraphs []docx.Paragraph, title string) string {
	title = strings.TrimSpace(title)
	seenTitle := false
	for _, paragraph := range paragraphs {
		text := strings.TrimSpace(paragraph.Text)
		if text == "" || paragraph.HasImage || looksLikeLabel(text) {
			continue
		}
		if !seenTitle && text == title {
			seenTitle = true
			continue
		}
		if text != title {
			return truncateRunes(text, maxWeixinArticleDigestRunes)
		}
	}
	return truncateRunes(title, maxWeixinArticleDigestRunes)
}

func looksLikeLabel(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || !strings.Contains(value, "/") {
		return false
	}
	letters := 0
	upper := 0
	for _, r := range value {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			letters++
			if r >= 'A' && r <= 'Z' {
				upper++
			}
		}
	}
	return letters > 0 && upper == letters
}

func validateTextLength(value string, max int, field string) error {
	if value == "" || utf8.RuneCountInString(value) <= max {
		return nil
	}
	return importUsage(fmt.Sprintf("article %s exceeds the WeChat limit", field), map[string]interface{}{
		"field": field,
		"limit": max,
		"runes": utf8.RuneCountInString(value),
	}).WithCategory("article_field_length").
		WithHint(fmt.Sprintf("请缩短 %s；微信公众号文章该字段最多 %d 个字符。", field, max))
}

func truncateRunes(value string, max int) string {
	value = strings.TrimSpace(value)
	if max <= 0 || utf8.RuneCountInString(value) <= max {
		return value
	}
	runes := []rune(value)
	return strings.TrimSpace(string(runes[:max]))
}

func importUsage(message string, details interface{}) *yxerrors.Error {
	return yxerrors.Usage(message, details).
		WithCategory("docx_import").
		WithHint("微信公众号 DOCX 导入会保留正文图片，但封面必须在导入前提供横版图片。")
}
