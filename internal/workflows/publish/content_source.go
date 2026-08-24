package publish

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type articleContentSource struct {
	Content string
	BaseDir string
}

var obsidianImageEmbedPattern = regexp.MustCompile(`!\[\[([^\]|]+)(?:\|([^\]]*))?\]\]`)

var articleMarkdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(goldmarkhtml.WithUnsafe()),
)

func loadArticleContentSource(path string) (articleContentSource, error) {
	abs, err := filepath.Abs(strings.TrimSpace(path))
	if err != nil {
		return articleContentSource{}, yxerrors.Usage("invalid article content file", map[string]interface{}{
			"path": path,
		}).WithCategory("file_not_readable")
	}
	if strings.TrimSpace(path) == "" {
		return articleContentSource{}, yxerrors.Usage("article content file must not be empty", nil).
			WithCategory("file_not_readable")
	}
	raw, err := os.ReadFile(abs)
	if err != nil {
		return articleContentSource{}, yxerrors.Usage("failed to read article content file", map[string]interface{}{
			"path":  abs,
			"cause": err.Error(),
		}).WithCategory("file_not_readable").
			WithHint("请确认 Markdown 文件存在且当前用户可读取。")
	}

	var rendered bytes.Buffer
	if err := articleMarkdown.Convert([]byte(normalizeObsidianImageEmbeds(string(raw))), &rendered); err != nil {
		return articleContentSource{}, yxerrors.Usage("failed to render article Markdown", map[string]interface{}{
			"path":  abs,
			"cause": err.Error(),
		}).WithCategory("content_render")
	}
	return articleContentSource{Content: rendered.String(), BaseDir: filepath.Dir(abs)}, nil
}

func normalizeObsidianImageEmbeds(markdown string) string {
	return obsidianImageEmbedPattern.ReplaceAllStringFunc(markdown, func(embed string) string {
		match := obsidianImageEmbedPattern.FindStringSubmatch(embed)
		if len(match) < 2 {
			return embed
		}
		target := strings.TrimSpace(match[1])
		alt := ""
		if len(match) > 2 {
			alt = strings.TrimSpace(match[2])
		}
		if target == "" {
			return embed
		}
		return fmt.Sprintf("![%s](<%s>)", alt, target)
	})
}
