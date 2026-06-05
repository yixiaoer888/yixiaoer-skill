package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
)

type Validator struct {
	SchemaDir string
}

type Result struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
}

type validationTarget struct {
	Value  interface{}
	Prefix string
}

func NewValidator(schemaDir string) Validator {
	return Validator{SchemaDir: schemaDir}
}

func (v Validator) Validate(platform, publishType string, payload map[string]interface{}) Result {
	raw, _, err := v.readSchema(platform, publishType)
	if err != nil {
		return basicValidate(payload)
	}
	var schema map[string]interface{}
	if err := json.Unmarshal(stripBOM(raw), &schema); err != nil {
		return basicValidate(payload)
	}
	sanitizeSchemaDocument(schema)
	targets := validationTargets(publishType, payload)
	var errors []string
	for _, target := range targets {
		errors = append(errors, validateValue(schema, target.Value, "/", target.Prefix)...)
	}
	return Result{Valid: len(errors) == 0, Errors: errors}
}

func (v Validator) Schema(platform, publishType string) (Document, error) {
	raw, path, err := v.readSchema(platform, publishType)
	if err != nil {
		return Document{}, err
	}
	var schema map[string]interface{}
	if err := json.Unmarshal(stripBOM(raw), &schema); err != nil {
		return Document{}, err
	}
	sanitizeSchemaDocument(schema)
	entry := buildEntry(filepath.Base(path))
	return buildDocument(entry, schema), nil
}

func (v Validator) List() ([]Entry, error) {
	pattern := filepath.Join(v.SchemaDir, "platforms", "*.schema.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	entries := make([]Entry, 0, len(files))
	for _, file := range files {
		base := filepath.Base(file)
		name := strings.TrimSuffix(base, ".schema.json")
		dot := strings.LastIndex(name, ".")
		if dot < 0 {
			continue
		}
		entries = append(entries, buildEntry(base))
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Type == entries[j].Type {
			return entries[i].Platform < entries[j].Platform
		}
		return entries[i].Type < entries[j].Type
	})
	return entries, nil
}

func (v Validator) Catalog() (Catalog, error) {
	entries, err := v.List()
	if err != nil {
		return Catalog{}, err
	}
	return Catalog{
		SchemaDir:   filepath.ToSlash(v.SchemaDir),
		RootSchemas: RootSchemaEntries(),
		Platforms:   entries,
	}, nil
}

func (v Validator) Fields(platform, publishType string) (map[string]PropertyView, error) {
	doc, err := v.Schema(platform, publishType)
	if err != nil {
		return nil, err
	}
	return doc.Properties, nil
}

func (v Validator) readSchema(platform, publishType string) ([]byte, string, error) {
	var lastErr error
	for _, key := range schemaPlatformKeys(platform, publishType) {
		schemaPath := filepath.Join(v.SchemaDir, "platforms", fmt.Sprintf("%s.%s.schema.json", key, TypeKey(publishType)))
		raw, err := os.ReadFile(schemaPath)
		if err == nil {
			return raw, schemaPath, nil
		}
		lastErr = err
	}
	return nil, "", lastErr
}

func schemaPlatformKeys(platformName, publishType string) []string {
	trimmed := strings.TrimSpace(platformName)
	normalized := strings.ToLower(trimmed)
	canonicalKey := platformutil.CanonicalKey(platformName)
	keys := []string{canonicalKey, normalized}
	if canonicalKey == "xhs" {
		keys = append(keys, "xiaohongshu")
	}
	if canonicalChineseName := platformutil.ChineseName(platformName); canonicalChineseName != "" {
		keys = append(keys, canonicalChineseName, strings.ToLower(canonicalChineseName))
	}
	keys = append(keys, trimmed)
	return uniqueStrings(keys)
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	var result []string
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}

func TypeKey(publishType string) string {
	return publishType
}

func DisplayType(typeKey string) string {
	return typeKey
}

func stripBOM(raw []byte) []byte {
	return []byte(strings.TrimPrefix(string(raw), "\uFEFF"))
}

func sanitizeSchemaDocument(schema map[string]interface{}) {
	for _, key := range []string{"$schema", "$id"} {
		delete(schema, key)
	}
	for _, value := range schema {
		switch typed := value.(type) {
		case map[string]interface{}:
			sanitizeSchemaDocument(typed)
		case []interface{}:
			for _, item := range typed {
				if child, ok := item.(map[string]interface{}); ok {
					sanitizeSchemaDocument(child)
				}
			}
		}
	}
}

func validationTargets(publishType string, payload map[string]interface{}) []validationTarget {
	if publishArgs, ok := payload["publishArgs"].(map[string]interface{}); ok {
		return validationTargets(publishType, normalizeValidationPayload(publishType, publishArgs))
	}
	if accountForms, ok := payload["accountForms"].([]interface{}); ok {
		var targets []validationTarget
		for i, form := range accountForms {
			formMap, ok := form.(map[string]interface{})
			if !ok {
				continue
			}
			if cpf, ok := formMap["contentPublishForm"]; ok {
				targets = append(targets, validationTarget{
					Value:  cpf,
					Prefix: fmt.Sprintf("accountForms[%d].contentPublishForm: ", i),
				})
			}
		}
		if len(targets) > 0 {
			return targets
		}
	}
	if _, ok := payload["formType"]; ok {
		return []validationTarget{{Value: payload}}
	}
	if _, ok := payload["title"]; ok {
		return []validationTarget{{Value: payload}}
	}
	if _, ok := payload["description"]; ok {
		return []validationTarget{{Value: payload}}
	}
	return []validationTarget{{Value: payload}}
}

func normalizeValidationPayload(publishType string, payload map[string]interface{}) map[string]interface{} {
	if TypeKey(publishType) != "article" {
		return payload
	}
	content, _ := payload["content"].(string)
	if strings.TrimSpace(content) == "" {
		return payload
	}
	accountForms, ok := payload["accountForms"].([]interface{})
	if !ok || len(accountForms) == 0 {
		return payload
	}
	for _, item := range accountForms {
		form, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if cpf == nil {
			continue
		}
		if _, exists := cpf["content"]; !exists {
			cpf["content"] = content
		}
	}
	return payload
}

func validateValue(schema map[string]interface{}, value interface{}, pathLabel, prefix string) []string {
	var errors []string
	if expectedType, ok := schema["type"].(string); ok && !matchesType(value, expectedType) {
		return append(errors, fmt.Sprintf("%s%s: expected %s", prefix, pathLabel, expectedType))
	}

	if constValue, ok := schema["const"]; ok && fmt.Sprint(value) != fmt.Sprint(constValue) {
		errors = append(errors, fmt.Sprintf("%s%s: must equal %v", prefix, pathLabel, constValue))
	}

	if enumValues, ok := schema["enum"].([]interface{}); ok && !inEnum(value, enumValues) {
		errors = append(errors, fmt.Sprintf("%s%s: must be one of %v", prefix, pathLabel, enumValues))
	}

	switch typed := value.(type) {
	case map[string]interface{}:
		errors = append(errors, validateObject(schema, typed, pathLabel, prefix)...)
	case []interface{}:
		errors = append(errors, validateArray(schema, typed, pathLabel, prefix)...)
	case string:
		errors = append(errors, validateString(schema, typed, pathLabel, prefix)...)
	}
	return errors
}

func validateObject(schema map[string]interface{}, value map[string]interface{}, pathLabel, prefix string) []string {
	var errors []string
	if required, ok := schema["required"].([]interface{}); ok {
		for _, item := range required {
			key := fmt.Sprint(item)
			if isOptionalResourceMetadata(key) {
				continue
			}
			if _, exists := value[key]; !exists {
				errors = append(errors, fmt.Sprintf("%s%s: missing required field \"%s\"", prefix, pathLabel, key))
			}
		}
	}

	properties := map[string]interface{}{}
	if rawProps, ok := schema["properties"].(map[string]interface{}); ok {
		properties = rawProps
	}
	for key, child := range value {
		childSchema, exists := properties[key]
		if !exists {
			if isCLICommonOptionalField(key) {
				continue
			}
			if additional, ok := schema["additionalProperties"].(bool); ok && !additional {
				errors = append(errors, fmt.Sprintf("%s%s: unexpected field \"%s\" (not in schema)", prefix, pathLabel, key))
			}
			continue
		}
		childMap, ok := childSchema.(map[string]interface{})
		if !ok {
			continue
		}
		errors = append(errors, validateValue(childMap, child, joinPath(pathLabel, key), prefix)...)
	}
	return errors
}

func isOptionalResourceMetadata(key string) bool {
	return key == "size" || key == "width" || key == "height"
}

func isCLICommonOptionalField(key string) bool {
	switch key {
	case "scheduledTime",
		"video",
		"images",
		"cover",
		"coverKey",
		"content",
		"mediaId",
		"platformName",
		"publishContentId",
		"fps",
		"isAppContent",
		"publishChannel",
		"clientId":
		return true
	default:
		return false
	}
}

func validateArray(schema map[string]interface{}, value []interface{}, pathLabel, prefix string) []string {
	var errors []string
	if minItems, ok := number(schema["minItems"]); ok && len(value) < int(minItems) {
		errors = append(errors, fmt.Sprintf("%s%s: must have at least %d items", prefix, pathLabel, int(minItems)))
	}
	if maxItems, ok := number(schema["maxItems"]); ok && len(value) > int(maxItems) {
		errors = append(errors, fmt.Sprintf("%s%s: must have at most %d items", prefix, pathLabel, int(maxItems)))
	}
	itemSchema, ok := schema["items"].(map[string]interface{})
	if !ok {
		return errors
	}
	for i, child := range value {
		errors = append(errors, validateValue(itemSchema, child, fmt.Sprintf("%s/%d", strings.TrimRight(pathLabel, "/"), i), prefix)...)
	}
	return errors
}

func validateString(schema map[string]interface{}, value, pathLabel, prefix string) []string {
	var errors []string
	if minLength, ok := number(schema["minLength"]); ok && len([]rune(value)) < int(minLength) {
		errors = append(errors, fmt.Sprintf("%s%s: must NOT have fewer than %d characters", prefix, pathLabel, int(minLength)))
	}
	if maxLength, ok := number(schema["maxLength"]); ok && len([]rune(value)) > int(maxLength) {
		errors = append(errors, fmt.Sprintf("%s%s: must NOT have more than %d characters", prefix, pathLabel, int(maxLength)))
	}
	return errors
}

func matchesType(value interface{}, expected string) bool {
	switch expected {
	case "object":
		_, ok := value.(map[string]interface{})
		return ok
	case "array":
		_, ok := value.([]interface{})
		return ok
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		_, ok := value.(float64)
		return ok
	case "integer":
		n, ok := value.(float64)
		return ok && n == float64(int64(n))
	case "boolean":
		_, ok := value.(bool)
		return ok
	default:
		return true
	}
}

func inEnum(value interface{}, enumValues []interface{}) bool {
	for _, allowed := range enumValues {
		if fmt.Sprint(value) == fmt.Sprint(allowed) {
			return true
		}
	}
	return false
}

func number(value interface{}) (float64, bool) {
	n, ok := value.(float64)
	return n, ok
}

func joinPath(parent, key string) string {
	if parent == "/" {
		return "/" + key
	}
	return strings.TrimRight(parent, "/") + "/" + key
}

func basicValidate(payload map[string]interface{}) Result {
	var errors []string
	accountForms, hasAccountForms := payload["accountForms"].([]interface{})
	if !hasAccountForms {
		if _, ok := payload["formType"]; ok {
			if payload["content"] == nil && payload["description"] == nil {
				errors = append(errors, "Inner form payload must include content or description")
			}
			return Result{Valid: len(errors) == 0, Errors: errors}
		}
		errors = append(errors, "Missing required field: accountForms")
		return Result{Valid: false, Errors: errors}
	}
	for i, item := range accountForms {
		form, ok := item.(map[string]interface{})
		if !ok {
			errors = append(errors, fmt.Sprintf("accountForms[%d]: must be an object", i))
			continue
		}
		if form["platformAccountId"] == nil && form["account_id"] == nil {
			errors = append(errors, fmt.Sprintf("accountForms[%d]: missing platformAccountId", i))
		}
		if form["contentPublishForm"] == nil {
			errors = append(errors, fmt.Sprintf("accountForms[%d]: missing contentPublishForm", i))
		}
	}
	return Result{Valid: len(errors) == 0, Errors: errors}
}
