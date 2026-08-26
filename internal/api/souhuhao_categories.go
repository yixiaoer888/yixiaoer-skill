package api

import (
	"fmt"
	"strings"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

const (
	sohuhaoCategoryErrorCategory = "sohuhao_category"
	sohuhaoCategoryQueryFailed   = "sohuhao_category_query_failed"
	sohuhaoCategoryDataInvalid   = "sohuhao_category_data_invalid"
	sohuhaoCategoryNotFound      = "sohuhao_category_not_found"
	sohuhaoCategoryAmbiguous     = "sohuhao_category_ambiguous"
	sohuhaoCategoryInvalid       = "sohuhao_category_invalid"
)

type sohuhaoCategoryNode struct {
	id        string
	name      string
	canonical map[string]interface{}
	children  []*sohuhaoCategoryNode
}

type sohuhaoCategoryCandidate struct {
	id   string
	name string
}

type sohuhaoCategoryResolveError struct {
	code string
}

func (e *sohuhaoCategoryResolveError) Error() string {
	return e.code
}

// CategoryPathView turns the normalized category tree returned by Categories
// into a machine-readable list of complete root-to-leaf paths. It accepts both
// the usual []interface{} tree and the {dataList: [...]} result shape.
func CategoryPathView(value interface{}) (interface{}, error) {
	items, ok := categoryItems(value)
	if !ok {
		return nil, yxerrors.New(yxerrors.ValidationType, sohuhaoCategoryDataInvalid, "搜狐号视频分类数据无效", map[string]interface{}{
			"cause":          "category result must be an array or an object containing dataList",
			"availablePaths": []interface{}{},
		}).WithCategory(sohuhaoCategoryErrorCategory).
			WithHint("请重新执行 yxer query categories <account_id> --type video --json 获取有效分类数据。").
			WithNextCommand("yxer query categories <account_id> --type video --json")
	}

	roots, err := parseSohuhaoCategoryTree(items)
	if err != nil {
		return nil, yxerrors.New(yxerrors.ValidationType, sohuhaoCategoryDataInvalid, "搜狐号视频分类数据无效", map[string]interface{}{
			"cause":          err.Error(),
			"availablePaths": []interface{}{},
		}).WithCategory(sohuhaoCategoryErrorCategory).
			WithHint("请重新执行 yxer query categories <account_id> --type video --json 检查分类返回结构。").
			WithNextCommand("yxer query categories <account_id> --type video --json")
	}

	paths := buildSohuhaoCategoryPaths(roots)
	categories := make([]interface{}, 0, len(roots))
	for _, root := range roots {
		categories = append(categories, clonePublishValue(root.canonical))
	}
	return map[string]interface{}{
		"categories": categories,
		"paths":      paths,
	}, nil
}

// NormalizeSohuhaoVideoPayload resolves the current Sohu video categories for
// every target account and rewrites category into the wire shape expected by
// the gateway. The input is never mutated.
func (c *Client) NormalizeSohuhaoVideoPayload(payload map[string]interface{}) (map[string]interface{}, error) {
	return c.normalizeSohuhaoVideoPayload(payload, false)
}

// NormalizeSohuhaoVideoPayloadForPlatform is the publish-workflow entry point.
// The command platform is authoritative even when a malformed payload omits or
// mislabels its top-level platforms field.
func (c *Client) NormalizeSohuhaoVideoPayloadForPlatform(payload map[string]interface{}, platform string) (map[string]interface{}, error) {
	return c.normalizeSohuhaoVideoPayload(payload, platformutil.CanonicalKey(platform) == "souhuhao")
}

func (c *Client) normalizeSohuhaoVideoPayload(payload map[string]interface{}, forceSohuhao bool) (map[string]interface{}, error) {
	cloned, _ := clonePublishValue(payload).(map[string]interface{})
	if cloned == nil || normalizeCategoryPublishType(stringField(cloned, "publishType")) != "video" {
		return cloned, nil
	}
	if !isSohuhaoVideoPayload(cloned) && !forceSohuhao {
		return cloned, nil
	}

	publishArgs, ok := cloned["publishArgs"].(map[string]interface{})
	if !ok || publishArgs == nil {
		return cloned, nil
	}
	forms, ok := categoryCandidateItems(publishArgs["accountForms"])
	if !ok || len(forms) == 0 {
		return cloned, nil
	}

	for index, item := range forms {
		form, ok := item.(map[string]interface{})
		if !ok || form == nil {
			return nil, newSohuhaoCategoryError(sohuhaoCategoryInvalid, "搜狐号视频分类参数无效", "", []interface{}{}, []interface{}{}, "account form must be an object", false)
		}
		accountID := firstCategoryString(form, "platformAccountId", "account_id")
		if accountID == "" {
			return nil, newSohuhaoCategoryError(sohuhaoCategoryInvalid, "搜狐号视频分类参数无效", "", []interface{}{}, []interface{}{}, fmt.Sprintf("accountForms[%d] is missing platformAccountId", index), false)
		}
		cpf, cpfOK := form["contentPublishForm"].(map[string]interface{})
		requestedValue := interface{}(nil)
		requestedExists := false
		if cpfOK && cpf != nil {
			requestedValue, requestedExists = cpf["category"]
		}
		requested, requestedOK := categoryCandidateItems(requestedValue)
		requestedDetails := requestedCategoryDetails(requestedValue)

		categoryResult, err := c.Categories(accountID, "video")
		if err != nil {
			return nil, newSohuhaoCategoryError(sohuhaoCategoryQueryFailed, "搜狐号视频分类查询失败", accountID, requestedDetails, []interface{}{}, err.Error(), true)
		}
		items, ok := categoryItems(categoryResult)
		if !ok || len(items) == 0 {
			cause := "category result is empty or does not contain dataList"
			return nil, newSohuhaoCategoryError(sohuhaoCategoryDataInvalid, "搜狐号视频分类数据无效", accountID, []interface{}{}, []interface{}{}, cause, false)
		}
		roots, err := parseSohuhaoCategoryTree(items)
		if err != nil || len(roots) == 0 {
			if err == nil {
				err = fmt.Errorf("category tree is empty")
			}
			return nil, newSohuhaoCategoryError(sohuhaoCategoryDataInvalid, "搜狐号视频分类数据无效", accountID, []interface{}{}, []interface{}{}, err.Error(), false)
		}
		availablePaths := buildSohuhaoCategoryPaths(roots)

		if !cpfOK || cpf == nil {
			return nil, newSohuhaoCategoryError(sohuhaoCategoryInvalid, "搜狐号视频分类参数无效", accountID, requestedDetails, availablePaths, "contentPublishForm must be an object", false)
		}
		if !requestedExists || !requestedOK || len(requested) == 0 {
			return nil, newSohuhaoCategoryError(sohuhaoCategoryInvalid, "搜狐号视频分类参数无效", accountID, requestedDetails, availablePaths, "category must be a non-empty array", false)
		}

		path, resolveErr := resolveSohuhaoCategoryPath(roots, requested)
		if resolveErr != nil {
			code := resolveErr.code
			message := "搜狐号视频分类参数无效"
			retryable := false
			switch code {
			case sohuhaoCategoryNotFound:
				message = "搜狐号视频分类不存在"
			case sohuhaoCategoryAmbiguous:
				message = "搜狐号视频分类名称不唯一"
			}
			return nil, newSohuhaoCategoryError(code, message, accountID, requested, availablePaths, "", retryable)
		}

		wireCategories := make([]interface{}, 0, len(path))
		for _, node := range path {
			wireCategories = append(wireCategories, sohuhaoWireCategory(node))
		}
		cpf["category"] = wireCategories
	}

	return cloned, nil
}

func isSohuhaoVideoPayload(payload map[string]interface{}) bool {
	if normalizeCategoryPublishType(stringField(payload, "publishType")) != "video" {
		return false
	}
	platforms, ok := categoryCandidateItems(payload["platforms"])
	if !ok {
		if values, ok := payload["platforms"].([]string); ok {
			platforms = make([]interface{}, 0, len(values))
			for _, value := range values {
				platforms = append(platforms, value)
			}
		} else {
			return false
		}
	}
	for _, item := range platforms {
		if platformutil.CanonicalKey(fmt.Sprint(item)) == "souhuhao" {
			return true
		}
	}
	return false
}

func normalizeCategoryPublishType(value string) string {
	switch strings.TrimSpace(value) {
	case "video", "视频":
		return "video"
	default:
		return strings.TrimSpace(value)
	}
}

func categoryItems(value interface{}) ([]interface{}, bool) {
	switch typed := value.(type) {
	case []interface{}:
		return typed, true
	case []map[string]interface{}:
		items := make([]interface{}, 0, len(typed))
		for _, item := range typed {
			items = append(items, item)
		}
		return items, true
	case map[string]interface{}:
		items, ok := typed["dataList"]
		if !ok {
			return nil, false
		}
		list, ok := categoryCandidateItems(items)
		return list, ok
	default:
		return nil, false
	}
}

func parseSohuhaoCategoryTree(items []interface{}) ([]*sohuhaoCategoryNode, error) {
	roots := make([]*sohuhaoCategoryNode, 0, len(items))
	for index, item := range items {
		node, err := parseSohuhaoCategoryNode(item)
		if err != nil {
			return nil, fmt.Errorf("categories[%d]: %w", index, err)
		}
		roots = append(roots, node)
	}
	return roots, nil
}

func parseSohuhaoCategoryNode(value interface{}) (*sohuhaoCategoryNode, error) {
	obj, ok := value.(map[string]interface{})
	if !ok || obj == nil {
		return nil, fmt.Errorf("category node must be an object")
	}
	id := categoryScalarString(obj["yixiaoerId"])
	name := categoryScalarString(obj["yixiaoerName"])
	if id == "" || name == "" {
		return nil, fmt.Errorf("category node requires yixiaoerId and yixiaoerName")
	}
	raw, ok := obj["raw"].(map[string]interface{})
	if !ok || raw == nil {
		return nil, fmt.Errorf("category %q raw must be an object", id)
	}
	rawID := categoryScalarString(raw["id"])
	if rawID == "" {
		return nil, fmt.Errorf("category %q raw.id is required", id)
	}
	if rawID != id {
		return nil, fmt.Errorf("category %q raw.id conflicts with yixiaoerId %q", rawID, id)
	}

	node := &sohuhaoCategoryNode{
		id:        id,
		name:      name,
		canonical: clonePublishValue(obj).(map[string]interface{}),
	}
	children, key, present, err := categoryChildren(obj)
	if err != nil {
		return nil, fmt.Errorf("category %q: %w", id, err)
	}
	if !present {
		return node, nil
	}
	parsedChildren, err := parseSohuhaoCategoryTree(children)
	if err != nil {
		return nil, fmt.Errorf("category %q children: %w", id, err)
	}
	node.children = parsedChildren
	canonicalChildren := make([]interface{}, 0, len(parsedChildren))
	for _, child := range parsedChildren {
		canonicalChildren = append(canonicalChildren, clonePublishValue(child.canonical))
	}
	node.canonical[key] = canonicalChildren
	return node, nil
}

func categoryChildren(obj map[string]interface{}) ([]interface{}, string, bool, error) {
	for _, key := range []string{"child", "children"} {
		value, exists := obj[key]
		if !exists || value == nil {
			continue
		}
		children, ok := categoryCandidateItems(value)
		if !ok {
			return nil, "", false, fmt.Errorf("%s must be an array", key)
		}
		return children, key, true, nil
	}
	return nil, "", false, nil
}

func buildSohuhaoCategoryPaths(roots []*sohuhaoCategoryNode) []interface{} {
	paths := make([]interface{}, 0)
	for _, root := range roots {
		if len(root.children) == 0 {
			paths = append(paths, sohuhaoCategoryPathView([]*sohuhaoCategoryNode{root}))
			continue
		}
		visitCategoryChildren(root, []*sohuhaoCategoryNode{root}, &paths)
	}
	return paths
}

func visitCategoryChildren(node *sohuhaoCategoryNode, path []*sohuhaoCategoryNode, paths *[]interface{}) {
	for _, child := range node.children {
		childPath := append(append([]*sohuhaoCategoryNode(nil), path...), child)
		if len(child.children) == 0 {
			*paths = append(*paths, sohuhaoCategoryPathView(childPath))
			continue
		}
		visitCategoryChildren(child, childPath, paths)
	}
}

func sohuhaoCategoryPathView(path []*sohuhaoCategoryNode) map[string]interface{} {
	names := make([]string, 0, len(path))
	nodes := make([]interface{}, 0, len(path))
	categories := make([]interface{}, 0, len(path))
	for _, node := range path {
		names = append(names, node.name)
		nodes = append(nodes, map[string]interface{}{"id": node.id, "name": node.name})
		categories = append(categories, clonePublishValue(node.canonical))
	}
	view := map[string]interface{}{
		"path":     strings.Join(names, " > "),
		"nodes":    nodes,
		"category": categories,
	}
	if len(path) == 2 {
		view["parentId"] = path[0].id
		view["parentName"] = path[0].name
		view["childId"] = path[1].id
		view["childName"] = path[1].name
	}
	return view
}

func sohuhaoWireCategory(node *sohuhaoCategoryNode) map[string]interface{} {
	item := map[string]interface{}{
		"id":   node.id,
		"text": node.name,
		"raw":  clonePublishValue(node.canonical),
	}
	if len(node.children) > 0 {
		children := make([]interface{}, 0, len(node.children))
		for _, child := range node.children {
			children = append(children, sohuhaoWireCategory(child))
		}
		item["children"] = children
	}
	return item
}

func categoryCandidateItems(value interface{}) ([]interface{}, bool) {
	switch typed := value.(type) {
	case []interface{}:
		return typed, true
	case []map[string]interface{}:
		items := make([]interface{}, 0, len(typed))
		for _, item := range typed {
			items = append(items, item)
		}
		return items, true
	default:
		return nil, false
	}
}

func requestedCategoryDetails(value interface{}) []interface{} {
	items, ok := categoryCandidateItems(value)
	if !ok || items == nil {
		return []interface{}{}
	}
	return clonePublishValue(items).([]interface{})
}

func resolveSohuhaoCategoryPath(roots []*sohuhaoCategoryNode, requested []interface{}) ([]*sohuhaoCategoryNode, *sohuhaoCategoryResolveError) {
	candidates := make([]sohuhaoCategoryCandidate, 0, len(requested))
	for _, item := range requested {
		candidate, err := parseSohuhaoCategoryCandidate(item)
		if err != nil {
			return nil, &sohuhaoCategoryResolveError{code: sohuhaoCategoryInvalid}
		}
		candidates = append(candidates, candidate)
	}
	last := candidates[len(candidates)-1]
	paths := buildSohuhaoNodePaths(roots)
	matched := make([][]*sohuhaoCategoryNode, 0)
	for _, path := range paths {
		leaf := path[len(path)-1]
		if last.id != "" {
			if leaf.id == last.id {
				matched = append(matched, path)
			}
		} else if last.name != "" && leaf.name == last.name {
			matched = append(matched, path)
		}
	}
	if len(matched) == 0 {
		return nil, &sohuhaoCategoryResolveError{code: sohuhaoCategoryNotFound}
	}
	if len(matched) > 1 {
		return nil, &sohuhaoCategoryResolveError{code: sohuhaoCategoryAmbiguous}
	}
	path := matched[0]
	if last.name != "" && path[len(path)-1].name != last.name {
		return nil, &sohuhaoCategoryResolveError{code: sohuhaoCategoryInvalid}
	}
	if len(candidates) > 1 {
		if len(candidates) != len(path) {
			return nil, &sohuhaoCategoryResolveError{code: sohuhaoCategoryInvalid}
		}
		for index, candidate := range candidates {
			if candidate.id != "" && candidate.id != path[index].id {
				return nil, &sohuhaoCategoryResolveError{code: sohuhaoCategoryInvalid}
			}
			if candidate.name != "" && candidate.name != path[index].name {
				return nil, &sohuhaoCategoryResolveError{code: sohuhaoCategoryInvalid}
			}
		}
	}
	return path, nil
}

func buildSohuhaoNodePaths(roots []*sohuhaoCategoryNode) [][]*sohuhaoCategoryNode {
	paths := make([][]*sohuhaoCategoryNode, 0)
	for _, root := range roots {
		if len(root.children) == 0 {
			paths = append(paths, []*sohuhaoCategoryNode{root})
			continue
		}
		visitSohuhaoNodePaths(root, []*sohuhaoCategoryNode{root}, &paths)
	}
	return paths
}

func visitSohuhaoNodePaths(node *sohuhaoCategoryNode, path []*sohuhaoCategoryNode, paths *[][]*sohuhaoCategoryNode) {
	for _, child := range node.children {
		childPath := append(append([]*sohuhaoCategoryNode(nil), path...), child)
		if len(child.children) == 0 {
			*paths = append(*paths, childPath)
			continue
		}
		visitSohuhaoNodePaths(child, childPath, paths)
	}
}

func parseSohuhaoCategoryCandidate(value interface{}) (sohuhaoCategoryCandidate, error) {
	obj, ok := value.(map[string]interface{})
	if !ok || obj == nil {
		return sohuhaoCategoryCandidate{}, fmt.Errorf("category candidate must be an object")
	}
	var candidate sohuhaoCategoryCandidate
	addID := func(value interface{}) error {
		if value == nil {
			return nil
		}
		text, ok := categoryScalar(value)
		if !ok {
			return fmt.Errorf("category id must be scalar")
		}
		if text == "" {
			return nil
		}
		if candidate.id != "" && candidate.id != text {
			return fmt.Errorf("category ids conflict")
		}
		candidate.id = text
		return nil
	}
	addName := func(value interface{}) error {
		if value == nil {
			return nil
		}
		text, ok := categoryScalar(value)
		if !ok {
			return fmt.Errorf("category name must be scalar")
		}
		if text == "" {
			return nil
		}
		if candidate.name != "" && candidate.name != text {
			return fmt.Errorf("category names conflict")
		}
		candidate.name = text
		return nil
	}
	for _, key := range []string{"id", "value", "yixiaoerId"} {
		if err := addID(obj[key]); err != nil {
			return sohuhaoCategoryCandidate{}, err
		}
	}
	for _, key := range []string{"text", "name", "label", "yixiaoerName"} {
		if err := addName(obj[key]); err != nil {
			return sohuhaoCategoryCandidate{}, err
		}
	}
	if raw, ok := obj["raw"].(map[string]interface{}); ok {
		if rawID, exists := raw["id"]; exists {
			if err := addID(rawID); err != nil {
				return sohuhaoCategoryCandidate{}, err
			}
		}
		if rawName, exists := raw["name"]; exists {
			if err := addName(rawName); err != nil {
				return sohuhaoCategoryCandidate{}, err
			}
		}
		if rawText, exists := raw["text"]; exists {
			if err := addName(rawText); err != nil {
				return sohuhaoCategoryCandidate{}, err
			}
		}
		if rawID, exists := raw["yixiaoerId"]; exists {
			if err := addID(rawID); err != nil {
				return sohuhaoCategoryCandidate{}, err
			}
		}
		if rawName, exists := raw["yixiaoerName"]; exists {
			if err := addName(rawName); err != nil {
				return sohuhaoCategoryCandidate{}, err
			}
		}
		if nestedRaw, ok := raw["raw"].(map[string]interface{}); ok {
			if nestedID, exists := nestedRaw["id"]; exists {
				if err := addID(nestedID); err != nil {
					return sohuhaoCategoryCandidate{}, err
				}
			}
		}
	}
	if candidate.id == "" && candidate.name == "" {
		return sohuhaoCategoryCandidate{}, fmt.Errorf("category candidate requires an id or name")
	}
	return candidate, nil
}

func firstCategoryString(obj map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value := categoryScalarString(obj[key]); value != "" {
			return value
		}
	}
	return ""
}

func categoryScalarString(value interface{}) string {
	text, ok := categoryScalar(value)
	if !ok {
		return ""
	}
	return text
}

func categoryScalar(value interface{}) (string, bool) {
	switch typed := value.(type) {
	case nil:
		return "", true
	case string:
		return strings.TrimSpace(typed), true
	case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return strings.TrimSpace(fmt.Sprint(typed)), true
	default:
		return "", false
	}
}

func newSohuhaoCategoryError(code, message, accountID string, requested, availablePaths []interface{}, cause string, retryable bool) *yxerrors.Error {
	details := map[string]interface{}{
		"accountId":      accountID,
		"requested":      clonePublishValue(requested),
		"availablePaths": clonePublishValue(availablePaths),
	}
	if cause != "" {
		details["cause"] = cause
	}
	hint := "请从 details.availablePaths 中选择有效的搜狐号视频分类路径。"
	switch code {
	case sohuhaoCategoryQueryFailed:
		hint = "请确认 apiKey 和搜狐号账号可用，然后重新执行 nextCommand 查询最新分类。"
	case sohuhaoCategoryDataInvalid:
		hint = "请重新执行 nextCommand 检查分类返回结构，确认分类包含有效的 raw.id。"
	case sohuhaoCategoryAmbiguous:
		hint = "分类名称不唯一，请改用分类 ID 或提供完整的父子分类路径。"
	case sohuhaoCategoryInvalid:
		hint = "请保持分类的 id、名称、raw.id 一致，并使用查询结果中的完整父子分类路径。"
	}
	nextCommand := fmt.Sprintf("yxer query categories %s --type video --json", accountID)
	return yxerrors.New(yxerrors.ValidationType, code, message, details).
		WithCategory(sohuhaoCategoryErrorCategory).
		WithHint(hint).
		WithNextCommand(nextCommand).
		WithRetryable(retryable)
}
