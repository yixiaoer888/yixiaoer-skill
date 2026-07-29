package publish

import (
	"fmt"
	htmlpkg "html"
	"net/url"
	"regexp"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	publishmod "github.com/yixiaoer/yixiaoer-skill/internal/modules/publish"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const articleImageMaterializeBucket = "material-library"

var (
	imgTagPattern     = regexp.MustCompile(`(?is)<img\b[^>]*>`)
	imgSrcAttrPattern = regexp.MustCompile("(?is)(^|\\s)src\\s*=\\s*(\"([^\"]*)\"|'([^']*)'|([^\\s\"'=<>`]+))")
)

type ArticleContentImageMaterialization struct {
	Path   string `json:"path"`
	From   string `json:"from"`
	To     string `json:"to,omitempty"`
	Bucket string `json:"bucket"`
	Key    string `json:"key,omitempty"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func materializeArticleContentImages(apiClient *api.Client, body map[string]interface{}, continueOnError bool) ([]ArticleContentImageMaterialization, error) {
	return rewriteArticleContentImages(body, continueOnError, func(sourceURL string) (string, string, error) {
		uploaded, err := apiClient.Upload(sourceURL, articleImageMaterializeBucket, false)
		if err != nil {
			return "", "", articleImageMaterializeError(sourceURL, err)
		}
		stableURL, err := apiClient.StableURL(articleImageMaterializeBucket, uploaded.Key)
		if err != nil {
			return "", "", articleImageMaterializeError(sourceURL, err)
		}
		return stableURL, uploaded.Key, nil
	})
}

func previewArticleContentImageMaterialization(body map[string]interface{}) []ArticleContentImageMaterialization {
	events, _ := rewriteArticleContentImages(body, true, func(sourceURL string) (string, string, error) {
		return sourceURL, "", nil
	})
	return events
}

func rewriteArticleContentImages(body map[string]interface{}, continueOnError bool, upload func(string) (string, string, error)) ([]ArticleContentImageMaterialization, error) {
	if body == nil || publishmod.NormalizePublishType(stringField(body, "publishType")) != "article" {
		return nil, nil
	}
	publishArgs := objectField(body, "publishArgs")
	if publishArgs == nil {
		return nil, nil
	}
	var events []ArticleContentImageMaterialization
	if err := rewriteStringFieldImages(publishArgs, "content", "publishArgs.content", upload, continueOnError, &events); err != nil {
		return events, err
	}
	if platformForm := weixinAccountArticlePlatformForm(publishArgs); platformForm != nil {
		articles, _ := platformForm["articles"].([]interface{})
		for i, item := range articles {
			article, _ := item.(map[string]interface{})
			if article == nil {
				continue
			}
			path := fmt.Sprintf(`publishArgs.platformForms[微信公众号].articles[%d].content`, i)
			if err := rewriteStringFieldImages(article, "content", path, upload, continueOnError, &events); err != nil {
				return events, err
			}
		}
	}
	return events, nil
}

func rewriteStringFieldImages(container map[string]interface{}, field, path string, upload func(string) (string, string, error), continueOnError bool, events *[]ArticleContentImageMaterialization) error {
	content := stringField(container, field)
	if strings.TrimSpace(content) == "" || !strings.Contains(strings.ToLower(content), "<img") {
		return nil
	}
	rewritten, fieldEvents, err := rewriteHTMLImageSources(content, path, upload, continueOnError)
	if err != nil {
		*events = append(*events, fieldEvents...)
		return err
	}
	if len(fieldEvents) == 0 {
		return nil
	}
	container[field] = rewritten
	*events = append(*events, fieldEvents...)
	return nil
}

func rewriteHTMLImageSources(content, path string, upload func(string) (string, string, error), continueOnError bool) (string, []ArticleContentImageMaterialization, error) {
	uploadedBySource := map[string]ArticleContentImageMaterialization{}
	var ordered []string
	var firstErr error

	rewritten := imgTagPattern.ReplaceAllStringFunc(content, func(tag string) string {
		if firstErr != nil {
			return tag
		}
		match := imgSrcAttrPattern.FindStringSubmatch(tag)
		if len(match) == 0 {
			return tag
		}
		source := strings.TrimSpace(firstNonEmptyString(match[3], match[4], match[5]))
		source = htmlpkg.UnescapeString(source)
		if !shouldMaterializeArticleImageSource(source) {
			return tag
		}
		event, exists := uploadedBySource[source]
		if !exists {
			targetURL, key, err := upload(source)
			if err != nil {
				event = ArticleContentImageMaterialization{Path: path, From: source, Bucket: articleImageMaterializeBucket, Status: "failed", Error: err.Error()}
				uploadedBySource[source] = event
				ordered = append(ordered, source)
				if !continueOnError {
					firstErr = err
				}
				return tag
			}
			status := "materialized"
			if targetURL == source {
				status = "would_materialize"
			}
			event = ArticleContentImageMaterialization{Path: path, From: source, To: targetURL, Bucket: articleImageMaterializeBucket, Key: key, Status: status}
			uploadedBySource[source] = event
			ordered = append(ordered, source)
		}
		if event.To == "" || event.To == source {
			return tag
		}
		return replaceImageSrcAttr(tag, event.To)
	})
	if firstErr != nil {
		return content, eventsInOrder(uploadedBySource, ordered), firstErr
	}
	return rewritten, eventsInOrder(uploadedBySource, ordered), nil
}

func shouldMaterializeArticleImageSource(source string) bool {
	source = strings.TrimSpace(source)
	if !strings.HasPrefix(strings.ToLower(source), "http://") && !strings.HasPrefix(strings.ToLower(source), "https://") {
		return false
	}
	parsed, err := url.Parse(source)
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	if host == "" {
		return false
	}
	if host == "oss-v2.yixiaoer.cn" || host == "yixiaoer-lite-asserts.oss-cn-shanghai.aliyuncs.com" {
		return false
	}
	return true
}

func replaceImageSrcAttr(tag, targetURL string) string {
	return imgSrcAttrPattern.ReplaceAllStringFunc(tag, func(attr string) string {
		prefix := attr[:strings.Index(strings.ToLower(attr), "src")]
		return prefix + `src="` + htmlpkg.EscapeString(targetURL) + `"`
	})
}

func eventsInOrder(events map[string]ArticleContentImageMaterialization, ordered []string) []ArticleContentImageMaterialization {
	out := make([]ArticleContentImageMaterialization, 0, len(ordered))
	for _, source := range ordered {
		out = append(out, events[source])
	}
	return out
}

func articleImageMaterializeError(sourceURL string, cause error) error {
	return yxerrors.Remote("failed to materialize article content image", map[string]interface{}{
		"url":   sourceURL,
		"cause": cause.Error(),
	}).WithCategory("article_content_image_materialization").
		WithHint("文章正文中的外链图片需要先转存为蚁小二可访问的稳定地址；请确认图片 URL 可下载，或先手动上传到 material-library 后替换正文 img src。").
		WithNextCommand("yxer upload --url " + sourceURL + " --bucket material-library")
}
