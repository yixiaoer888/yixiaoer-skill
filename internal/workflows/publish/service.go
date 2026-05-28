package publish

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yixiaoer/yixiaoer-skill/internal/core/client"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/config"
	"github.com/yixiaoer/yixiaoer-skill/internal/core/media"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/core/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type ExecuteInput struct {
	PublishType        string
	PlatformInput      string
	Payload            map[string]interface{}
	PositionalClientID string
	FlagChannel        string
	FlagClientID       string
}

type BuildInput struct {
	PublishType   string
	PlatformInput string
	Account       string
	Title         string
	Description   string
	Content       string
	Images        []string
	VideoPath     string
	VideoKey      string
	CoverPath     string
	VisibleType   int
}

type Service struct{}

func NewService() Service {
	return Service{}
}

func formatImageTextDescription(description string, tags []string) string {
	description = strings.TrimSpace(description)
	if len(tags) == 0 {
		return description
	}
	if strings.Contains(description, "<topic") {
		return description
	}
	topicParts := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		text := strings.TrimPrefix(tag, "#")
		if text == "" {
			continue
		}
		topicParts = append(topicParts, fmt.Sprintf(`<topic text="%s">#%s</topic>`, text, text))
	}
	if len(topicParts) == 0 {
		return description
	}
	return fmt.Sprintf("<p>%s</p><p>%s</p>", description, strings.Join(topicParts, ""))
}

func stringSlice(values []interface{}) []string {
	if len(values) == 0 {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		text := strings.TrimSpace(fmt.Sprint(value))
		if text == "" || text == "<nil>" {
			continue
		}
		result = append(result, text)
	}
	return result
}

func (Service) Execute(input ExecuteInput) (map[string]interface{}, error) {
	input.PublishType = NormalizePublishType(input.PublishType)
	platform, err := SinglePlatform(input.PlatformInput)
	if err != nil {
		return nil, err
	}
	platforms := []string{platform}
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	resolvedPayload := cloneMap(input.Payload)
	channel, clientID, err := ResolvePublishMode(cfg, resolvedPayload, input.PositionalClientID, input.FlagChannel, input.FlagClientID)
	if err != nil {
		return nil, err
	}
	resolvedPayload["publishChannel"] = channel
	if clientID != "" {
		resolvedPayload["clientId"] = clientID
	} else {
		delete(resolvedPayload, "clientId")
	}
	NormalizeStandardPublishArgs(ExtractPublishArgs(resolvedPayload))
	publishArgs := ExtractPublishArgs(resolvedPayload)

	validator := schema.NewValidator(cfg.SchemaDir)
	for _, platform := range platforms {
		result := validator.Validate(platform, input.PublishType, resolvedPayload)
		if !result.Valid {
			return nil, yxerrors.Usage("Schema validation failed", result.Errors).
				WithHint("请根据对应平台 schema 修正 payload 字段后重试。").
				WithNextCommand(fmt.Sprintf("yxer validate %s %s <payload.json>", platform, input.PublishType))
		}
	}
	preflight := Preflight(input.PublishType, platforms, resolvedPayload)
	if len(preflight.Errors) > 0 {
		return nil, yxerrors.Usage("Publish preflight failed", preflight.Errors).
			WithHint("请先完成资源上传、账号校验，并确保发布参数中不包含外部 URL。")
	}

	apiClient := client.New(cfg)
	if err := AssertAccountsOnline(apiClient, platforms, preflight.AccountIDs); err != nil {
		return nil, err
	}

	body := BuildPublishBody(resolvedPayload, publishArgs, input.PublishType, platforms)
	body["publishChannel"] = channel
	if clientID != "" {
		body["clientId"] = clientID
	} else {
		delete(body, "clientId")
	}
	return apiClient.Publish(body)
}

func (Service) BuildPayload(input BuildInput) (map[string]interface{}, error) {
	input.PublishType = NormalizePublishType(input.PublishType)
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	apiClient := client.New(cfg)
	platform, err := SinglePlatform(input.PlatformInput)
	if err != nil {
		return nil, err
	}
	accountID, err := resolveAccountID(apiClient, platform, input.Account)
	if err != nil {
		return nil, err
	}

	form := map[string]interface{}{
		"formType": "task",
	}
	accountForm := map[string]interface{}{
		"platformAccountId":  accountID,
		"contentPublishForm": form,
	}
	payload := map[string]interface{}{
		"publishChannel": "cloud",
	}

	switch input.PublishType {
	case "imageText":
		if strings.TrimSpace(input.Description) == "" {
			return nil, yxerrors.Usage("imageText publish requires description", nil).
				WithHint("请传入 --description，图文发布正文不能为空。")
		}
		tagValues, _ := form["tags"].([]interface{})
		finalDescription := formatImageTextDescription(input.Description, stringSlice(tagValues))
		form["description"] = finalDescription
		if strings.TrimSpace(input.Title) != "" {
			form["title"] = input.Title
		}
		if input.VisibleType >= 0 {
			form["visibleType"] = float64(input.VisibleType)
		}
		images, err := uploadImages(apiClient, input.Images)
		if err != nil {
			return nil, err
		}
		if len(images) == 0 {
			return nil, yxerrors.Usage("imageText publish requires at least one image", nil).
				WithHint("请至少传入一个 --image 本地文件路径或 URL。")
		}
		form["images"] = images
		firstImage, _ := images[0].(map[string]interface{})
		accountForm["cover"] = imageCover(firstImage)
		accountForm["coverKey"] = fmt.Sprint(firstImage["key"])
		payload["coverKey"] = accountForm["coverKey"]
		payload["desc"] = input.Description
	case "article":
		if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Content) == "" {
			return nil, yxerrors.Usage("article publish requires title and content", nil).
				WithHint("请同时传入 --title 和 --content，文章正文支持 @文件路径。")
		}
		if strings.TrimSpace(input.CoverPath) == "" {
			return nil, yxerrors.Usage("article publish requires cover", nil).
				WithHint("请传入 --cover，当前 article flags 模式要求明确提供封面。")
		}
		content, err := resolveContent(input.Content)
		if err != nil {
			return nil, err
		}
		form["title"] = input.Title
		form["content"] = content
		if usesVisibleType(platform) {
			form["visibleType"] = float64(defaultVisibleType(input.VisibleType))
		} else {
			form["pubType"] = float64(1)
		}
		cover, err := uploadImage(apiClient, input.CoverPath)
		if err != nil {
			return nil, err
		}
		form["covers"] = []interface{}{cover}
		accountForm["cover"] = cover
		accountForm["coverKey"] = fmt.Sprint(cover["key"])
		payload["coverKey"] = accountForm["coverKey"]
		payload["desc"] = input.Title
	case "video":
		if strings.TrimSpace(input.Title) == "" || strings.TrimSpace(input.Description) == "" {
			return nil, yxerrors.Usage("video publish requires title and description", nil).
				WithHint("请传入 --title 和 --description。")
		}
		if strings.TrimSpace(input.CoverPath) == "" {
			return nil, yxerrors.Usage("video publish requires cover", nil).
				WithHint("请传入 --cover，当前 video flags 模式要求明确提供封面。")
		}
		var video map[string]interface{}
		if strings.TrimSpace(input.VideoPath) != "" {
			video, err = uploadVideo(apiClient, input.VideoPath)
			if err != nil {
				return nil, err
			}
		} else if strings.TrimSpace(input.VideoKey) != "" {
			video = map[string]interface{}{"key": input.VideoKey}
		} else {
			return nil, yxerrors.Usage("video publish requires --video or --video-key", nil).
				WithHint("请传入本地视频文件路径，或先上传视频后提供 --video-key。")
		}
		form["title"] = input.Title
		form["description"] = input.Description
		if usesVisibleType(platform) {
			form["visibleType"] = float64(defaultVisibleType(input.VisibleType))
		}
		accountForm["video"] = video
		cover, err := uploadImage(apiClient, input.CoverPath)
		if err != nil {
			return nil, err
		}
		accountForm["cover"] = cover
		accountForm["coverKey"] = fmt.Sprint(cover["key"])
		payload["coverKey"] = accountForm["coverKey"]
		payload["desc"] = input.Description
	default:
		return nil, yxerrors.Usage("flags mode does not support publish type", input.PublishType).
			WithHint("目前仅支持 video、imageText、article 三种发布类型。")
	}

	payload["accountForms"] = []interface{}{accountForm}
	return payload, nil
}

func usesVisibleType(platform string) bool {
	switch platform {
	case "抖音":
		return true
	default:
		return false
	}
}

func defaultVisibleType(value int) int {
	if value >= 0 {
		return value
	}
	return 0
}

func cloneMap(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return map[string]interface{}{}
	}
	dst := make(map[string]interface{}, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func BuildPublishBody(payload, publishArgs map[string]interface{}, publishType string, platforms []string) map[string]interface{} {
	body := map[string]interface{}{
		"publishType":    publishType,
		"platforms":      platforms,
		"publishArgs":    publishArgs,
		"publishChannel": "cloud",
	}
	if _, ok := payload["publishArgs"].(map[string]interface{}); ok {
		for key, value := range payload {
			if key == "action" {
				continue
			}
			if key == "publishArgs" {
				body[key] = publishArgs
				continue
			}
			body[key] = value
		}
		body["publishType"] = publishType
		body["platforms"] = platforms
		if _, ok := body["publishChannel"]; !ok {
			body["publishChannel"] = "cloud"
		}
		return body
	}
	copyOptionalPublishFields(body, payload)
	return body
}

func ResolvePublishMode(cfg config.Config, payload map[string]interface{}, positionalClientID, flagChannel, flagClientID string) (string, string, error) {
	channel := "cloud"
	clientID := ""
	if value, ok := payload["publishChannel"]; ok {
		channel = strings.TrimSpace(fmt.Sprint(value))
	}
	if value, ok := payload["clientId"]; ok {
		clientID = strings.TrimSpace(fmt.Sprint(value))
	}
	if strings.TrimSpace(positionalClientID) != "" {
		channel = "local"
		clientID = strings.TrimSpace(positionalClientID)
	}
	if strings.TrimSpace(flagChannel) != "" {
		channel = strings.TrimSpace(flagChannel)
	}
	if strings.TrimSpace(flagClientID) != "" {
		clientID = strings.TrimSpace(flagClientID)
	}
	switch channel {
	case "", "cloud":
		channel = "cloud"
		clientID = ""
	case "local":
		if clientID == "" {
			clientID = strings.TrimSpace(cfg.LocalClientID)
		}
		if clientID == "" {
			return "", "", yxerrors.Usage(`clientId is required when publishChannel is "local"`, []string{
				`Run: yxer config set-local-client-id <clientId>`,
				`Or pass a fourth positional argument: yxer publish video <platform> payload.json <clientId>`,
				`Or pass flags: yxer publish video <platform> payload.json --publish-channel local --client-id <clientId>`,
			}).
				WithHint("本机发布必须指定 clientId，可通过配置或命令参数提供。").
				WithNextCommand("yxer config set-local-client-id <clientId>")
		}
	default:
		return "", "", yxerrors.Usage(`publishChannel must be "cloud" or "local"`, []string{
			fmt.Sprintf("got %q", channel),
		}).
			WithHint(`publishChannel 仅支持 "cloud" 或 "local"。`)
	}
	return channel, clientID, nil
}

func copyOptionalPublishFields(dst, src map[string]interface{}) {
	for _, field := range []string{"cover", "coverKey", "taskSetId", "desc", "isDraft", "isAppContent"} {
		if value, ok := src[field]; ok {
			dst[field] = value
		}
	}
	if channel, ok := src["publishChannel"]; ok {
		dst["publishChannel"] = channel
	}
	if clientID, ok := src["clientId"]; ok {
		dst["clientId"] = clientID
	}
}

func SplitPlatforms(value string) []string {
	var platforms []string
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			platforms = append(platforms, platformutil.ChineseName(item))
		}
	}
	return platforms
}

func SinglePlatform(value string) (string, error) {
	platforms := SplitPlatforms(value)
	if len(platforms) != 1 {
		return "", yxerrors.Usage("publish supports exactly one platform per command", []string{
			`Use Agent orchestration for multi-platform publishing: run "yxer accounts", "yxer schema get", "yxer validate", and "yxer publish" once per platform.`,
			`Example: yxer publish image-text xhs xhs-payload.json; then yxer publish image-text zhihu zhihu-payload.json`,
		}).
			WithHint("单次 publish 命令只支持一个平台，请拆分成多次调用。")
	}
	return platforms[0], nil
}

func AssertAccountsOnline(apiClient *client.Client, platforms []string, accountIDs []string) error {
	wanted := map[string]bool{}
	for _, id := range accountIDs {
		wanted[id] = true
	}
	found := map[string]map[string]interface{}{}
	for _, platform := range platforms {
		accounts, err := apiClient.Accounts(platform)
		if err != nil {
			return err
		}
		for _, account := range accounts {
			id := client.AccountID(account)
			if wanted[id] {
				found[id] = account
			}
		}
	}
	var errors []string
	for id := range wanted {
		account, ok := found[id]
		if !ok {
			errors = append(errors, "account "+id+": not found in target platform account list")
			continue
		}
		if status := client.AccountStatus(account); status != 1 {
			errors = append(errors, fmt.Sprintf("account %s: status=%d; publish requires status=1", id, status))
		}
	}
	if len(errors) > 0 {
		return yxerrors.Usage("Account preflight failed", errors).
			WithHint("请先运行账号查询，确认目标账号存在且状态为在线。").
			WithNextCommand("yxer accounts <platform>")
	}
	return nil
}

func resolveAccountID(apiClient *client.Client, platform, selector string) (string, error) {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return "", yxerrors.Usage("account is required in flags mode", nil).
			WithHint("请传入 --account，值可以是账号 ID、名称或昵称。").
			WithNextCommand(fmt.Sprintf("yxer accounts %s", platform))
	}
	accounts, err := apiClient.Accounts(platform)
	if err != nil {
		return "", err
	}
	var exact []map[string]interface{}
	var fuzzy []map[string]interface{}
	for _, account := range accounts {
		id := strings.TrimSpace(client.AccountID(account))
		name := accountDisplayName(account)
		if id == selector || name == selector {
			exact = append(exact, account)
			continue
		}
		if strings.Contains(name, selector) {
			fuzzy = append(fuzzy, account)
		}
	}
	candidates := exact
	if len(candidates) == 0 {
		candidates = fuzzy
	}
	if len(candidates) == 1 {
		if client.AccountStatus(candidates[0]) != 1 {
			return "", yxerrors.Usage("selected account is offline", selector).
				WithHint("目标账号当前不在线，请先检查账号状态。").
				WithNextCommand(fmt.Sprintf("yxer accounts %s", platform))
		}
		return client.AccountID(candidates[0]), nil
	}
	if len(candidates) == 0 {
		return "", yxerrors.Usage("account not found", selector).
			WithHint("未找到匹配账号，请先查询平台账号列表。").
			WithNextCommand(fmt.Sprintf("yxer accounts %s", platform))
	}
	return "", yxerrors.Usage("account selector is ambiguous", selector).
		WithHint("匹配到多个账号，请改用更精确的名称或直接传账号 ID。").
		WithNextCommand(fmt.Sprintf("yxer accounts %s --name %s", platform, selector))
}

func accountDisplayName(account map[string]interface{}) string {
	for _, key := range []string{"platformAccountName", "name", "nickname", "remark", "platformAuthorId"} {
		if value := strings.TrimSpace(fmt.Sprint(account[key])); value != "" && value != "<nil>" {
			return value
		}
	}
	return ""
}

func uploadImages(apiClient *client.Client, paths []string) ([]interface{}, error) {
	images := make([]interface{}, 0, len(paths))
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		image, err := uploadImage(apiClient, path)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func uploadImage(apiClient *client.Client, path string) (map[string]interface{}, error) {
	result, err := apiClient.Upload(path, "cloud-publish")
	if err != nil {
		return nil, err
	}
	format := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	if format == "" {
		format = result.Format
	}
	return map[string]interface{}{
		"key":    result.Key,
		"size":   float64(result.Size),
		"width":  float64(result.Width),
		"height": float64(result.Height),
		"format": format,
	}, nil
}

func uploadVideo(apiClient *client.Client, path string) (map[string]interface{}, error) {
	result, err := apiClient.Upload(path, "cloud-publish")
	if err != nil {
		return nil, err
	}
	meta, err := media.ProbeVideo(path)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"key":      result.Key,
		"size":     float64(result.Size),
		"width":    float64(meta.Width),
		"height":   float64(meta.Height),
		"duration": meta.Duration,
		"format":   meta.Format,
	}, nil
}

func imageCover(image interface{}) map[string]interface{} {
	typed, _ := image.(map[string]interface{})
	if typed == nil {
		return nil
	}
	return map[string]interface{}{
		"key":    typed["key"],
		"size":   typed["size"],
		"width":  typed["width"],
		"height": typed["height"],
		"format": typed["format"],
	}
}

func resolveContent(content string) (string, error) {
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "@") {
		return content, nil
	}
	raw, err := os.ReadFile(strings.TrimPrefix(content, "@"))
	if err != nil {
		return "", yxerrors.Usage("article content file could not be read", []string{
			err.Error(),
		}).WithHint("请检查 --content @文件路径 是否存在且可读。")
	}
	return string(raw), nil
}
