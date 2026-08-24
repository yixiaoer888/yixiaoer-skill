package publish

import (
	"fmt"
	htmlpkg "html"
	"net/url"
	"os"
	"path/filepath"
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

func materializeArticleContentImages(apiClient *api.Client, body map[string]interface{}, baseDir string, continueOnError bool) ([]ArticleContentImageMaterialization, error) {
	return rewriteArticleContentImagesWithBaseDir(body, baseDir, continueOnError, func(sourceURL string) (string, string, error) {
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

func previewArticleContentImageMaterialization(body map[string]interface{}, baseDir string) []ArticleContentImageMaterialization {
	events, _ := rewriteArticleContentImagesWithBaseDir(body, baseDir, true, func(sourceURL string) (string, string, error) {
		return sourceURL, "", nil
	})
	return events
}

func rewriteArticleContentImages(body map[string]interface{}, continueOnError bool, upload func(string) (string, string, error)) ([]ArticleContentImageMaterialization, error) {
	return rewriteArticleContentImagesWithBaseDir(body, "", continueOnError, upload)
}

func rewriteArticleContentImagesWithBaseDir(body map[string]interface{}, baseDir string, continueOnError bool, upload func(string) (string, string, error)) ([]ArticleContentImageMaterialization, error) {
	if body == nil || publishmod.NormalizePublishType(stringField(body, "publishType")) != "article" {
		return nil, nil
	}
	publishArgs := objectField(body, "publishArgs")
	if publishArgs == nil {
		return nil, nil
	}
	var events []ArticleContentImageMaterialization
	if err := rewriteStringFieldImagesWithBaseDir(publishArgs, "content", "publishArgs.content", baseDir, upload, continueOnError, &events); err != nil {
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
			if err := rewriteStringFieldImagesWithBaseDir(article, "content", path, baseDir, upload, continueOnError, &events); err != nil {
				return events, err
			}
		}
	}
	return events, nil
}

func rewriteStringFieldImages(container map[string]interface{}, field, path string, upload func(string) (string, string, error), continueOnError bool, events *[]ArticleContentImageMaterialization) error {
	return rewriteStringFieldImagesWithBaseDir(container, field, path, "", upload, continueOnError, events)
}

func rewriteStringFieldImagesWithBaseDir(container map[string]interface{}, field, path, baseDir string, upload func(string) (string, string, error), continueOnError bool, events *[]ArticleContentImageMaterialization) error {
	content := stringField(container, field)
	if strings.TrimSpace(content) == "" || !strings.Contains(strings.ToLower(content), "<img") {
		return nil
	}
	rewritten, fieldEvents, err := rewriteHTMLImageSourcesWithBaseDir(content, path, baseDir, upload, continueOnError)
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
	return rewriteHTMLImageSourcesWithBaseDir(content, path, "", upload, continueOnError)
}

func rewriteHTMLImageSourcesWithBaseDir(content, path, baseDir string, upload func(string) (string, string, error), continueOnError bool) (string, []ArticleContentImageMaterialization, error) {
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
		resolvedSource, shouldMaterialize, resolveErr := resolveArticleImageSource(source, baseDir)
		if resolveErr != nil {
			event := ArticleContentImageMaterialization{Path: path, From: source, Bucket: articleImageMaterializeBucket, Status: "failed", Error: resolveErr.Error()}
			uploadedBySource[source] = event
			ordered = append(ordered, source)
			if !continueOnError {
				firstErr = resolveErr
			}
			return tag
		}
		if !shouldMaterialize {
			return tag
		}
		event, exists := uploadedBySource[source]
		if !exists {
			targetURL, key, err := upload(resolvedSource)
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
			if targetURL == source || targetURL == resolvedSource {
				status = "would_materialize"
				targetURL = source
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

func resolveArticleImageSource(source, baseDir string) (string, bool, error) {
	source = strings.TrimSpace(source)
	if source == "" {
		return source, false, nil
	}
	lower := strings.ToLower(source)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		if !shouldMaterializeArticleImageSource(source) {
			return source, false, nil
		}
		return source, true, nil
	}
	if strings.HasPrefix(lower, "data:") || strings.HasPrefix(lower, "cid:") {
		return source, false, nil
	}

	localSource := source
	if strings.HasPrefix(lower, "file://") {
		parsed, err := url.Parse(source)
		if err != nil {
			return "", true, yxerrors.Usage("invalid article content image file URL", map[string]interface{}{
				"source": source,
			}).WithCategory("article_content_image_source")
		}
		localSource = parsed.Path
		if parsed.Host != "" && !strings.EqualFold(parsed.Host, "localhost") {
			localSource = `\\` + parsed.Host + localSource
		}
		localSource, err = url.PathUnescape(localSource)
		if err != nil {
			return "", true, yxerrors.Usage("invalid encoded article content image file URL", map[string]interface{}{
				"source": source,
			}).WithCategory("article_content_image_source")
		}
		if len(localSource) >= 3 && localSource[0] == '/' && localSource[2] == ':' {
			localSource = localSource[1:]
		}
	}
	decodedLocalSource, err := url.PathUnescape(localSource)
	if err != nil {
		return "", true, yxerrors.Usage("invalid encoded article content image path", map[string]interface{}{
			"source": source,
		}).WithCategory("article_content_image_source")
	}
	localSource = decodedLocalSource

	if !filepath.IsAbs(localSource) {
		if strings.TrimSpace(baseDir) == "" {
			return "", true, yxerrors.Usage("relative article content image requires a content file", map[string]interface{}{
				"source": source,
			}).WithCategory("article_content_image_source").
				WithHint("请通过 --content-file 提供 Markdown 文件，以便按文档目录解析图片路径。")
		}
		localSource = filepath.Join(baseDir, filepath.FromSlash(localSource))
	}
	localSource = filepath.Clean(localSource)
	info, err := os.Stat(localSource)
	if err != nil {
		return "", true, yxerrors.Usage("article content image file not found", map[string]interface{}{
			"source": source,
			"path":   localSource,
		}).WithCategory("article_content_image_source").
			WithHint("请检查 Markdown 图片路径，路径应相对于 Markdown 文件所在目录。")
	}
	if info.IsDir() {
		return "", true, yxerrors.Usage("article content image path is a directory", map[string]interface{}{
			"source": source,
			"path":   localSource,
		}).WithCategory("article_content_image_source")
	}
	return localSource, true, nil
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
