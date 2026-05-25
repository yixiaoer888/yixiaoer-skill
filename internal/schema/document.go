package schema

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type Entry struct {
	Key        string `json:"key"`
	Platform   string `json:"platform"`
	Type       string `json:"type"`
	File       string `json:"file"`
	RootSchema string `json:"rootSchema"`
}

type Document struct {
	Key                  string                  `json:"key"`
	Platform             string                  `json:"platform"`
	Type                 string                  `json:"type"`
	File                 string                  `json:"file"`
	RootSchema           string                  `json:"rootSchema"`
	Title                string                  `json:"title,omitempty"`
	Required             []string                `json:"required,omitempty"`
	AdditionalProperties bool                    `json:"additionalProperties"`
	Properties           map[string]PropertyView `json:"properties,omitempty"`
}

type PropertyView struct {
	Type       string                  `json:"type,omitempty"`
	Required   bool                    `json:"required,omitempty"`
	Format     string                  `json:"format,omitempty"`
	Const      interface{}             `json:"const,omitempty"`
	Default    interface{}             `json:"default,omitempty"`
	Enum       []interface{}           `json:"enum,omitempty"`
	MinLength  *int                    `json:"minLength,omitempty"`
	MaxLength  *int                    `json:"maxLength,omitempty"`
	MinItems   *int                    `json:"minItems,omitempty"`
	MaxItems   *int                    `json:"maxItems,omitempty"`
	Minimum    *float64                `json:"minimum,omitempty"`
	Maximum    *float64                `json:"maximum,omitempty"`
	Properties map[string]PropertyView `json:"properties,omitempty"`
	Items      *PropertyView           `json:"items,omitempty"`
}

type Catalog struct {
	SchemaDir   string  `json:"schemaDir"`
	RootSchemas []Entry `json:"rootSchemas"`
	Platforms   []Entry `json:"platforms"`
}

func rootSchemaPath(publishType string) string {
	switch TypeKey(publishType) {
	case "account":
		return "schemas/account.schema.json"
	default:
		return "schemas/publish.schema.json"
	}
}

func RootSchemaEntries() []Entry {
	return []Entry{
		{
			Key:        "root/account",
			Platform:   "root",
			Type:       "account",
			File:       "schemas/account.schema.json",
			RootSchema: "schemas/account.schema.json",
		},
		{
			Key:        "root/publish",
			Platform:   "root",
			Type:       "publish",
			File:       "schemas/publish.schema.json",
			RootSchema: "schemas/publish.schema.json",
		},
	}
}

func buildEntry(base string) Entry {
	name := strings.TrimSuffix(base, ".schema.json")
	dot := strings.LastIndex(name, ".")
	platform := name[:dot]
	publishType := DisplayType(name[dot+1:])
	return Entry{
		Key:        fmt.Sprintf("%s/%s", platform, publishType),
		Platform:   platform,
		Type:       publishType,
		File:       filepath.ToSlash(filepath.Join("schemas", "platforms", base)),
		RootSchema: rootSchemaPath(publishType),
	}
}

func buildDocument(entry Entry, schemaDoc map[string]interface{}) Document {
	required := requiredFieldNames(schemaDoc)
	requiredSet := map[string]bool{}
	for _, name := range required {
		requiredSet[name] = true
	}
	doc := Document{
		Key:                  entry.Key,
		Platform:             entry.Platform,
		Type:                 entry.Type,
		File:                 entry.File,
		RootSchema:           entry.RootSchema,
		Title:                stringValue(schemaDoc["title"]),
		Required:             required,
		AdditionalProperties: boolValue(schemaDoc["additionalProperties"]),
		Properties:           map[string]PropertyView{},
	}
	if rawProps, ok := schemaDoc["properties"].(map[string]interface{}); ok {
		keys := make([]string, 0, len(rawProps))
		for key := range rawProps {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			child, ok := rawProps[key].(map[string]interface{})
			if !ok {
				continue
			}
			doc.Properties[key] = buildPropertyView(child, requiredSet[key])
		}
	}
	if len(doc.Properties) == 0 {
		doc.Properties = nil
	}
	return doc
}

func buildPropertyView(schemaDoc map[string]interface{}, required bool) PropertyView {
	view := PropertyView{
		Type:     stringValue(schemaDoc["type"]),
		Required: required,
		Format:   stringValue(schemaDoc["format"]),
		Const:    schemaDoc["const"],
		Default:  schemaDoc["default"],
		Enum:     interfaceSlice(schemaDoc["enum"]),
		MinLength: intPointer(numberValue(schemaDoc["minLength"])),
		MaxLength: intPointer(numberValue(schemaDoc["maxLength"])),
		MinItems:  intPointer(numberValue(schemaDoc["minItems"])),
		MaxItems:  intPointer(numberValue(schemaDoc["maxItems"])),
		Minimum:   floatPointer(numberValue(schemaDoc["minimum"])),
		Maximum:   floatPointer(numberValue(schemaDoc["maximum"])),
	}
	if rawProps, ok := schemaDoc["properties"].(map[string]interface{}); ok {
		requiredSet := map[string]bool{}
		for _, name := range requiredFieldNames(schemaDoc) {
			requiredSet[name] = true
		}
		keys := make([]string, 0, len(rawProps))
		for key := range rawProps {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		view.Properties = map[string]PropertyView{}
		for _, key := range keys {
			child, ok := rawProps[key].(map[string]interface{})
			if !ok {
				continue
			}
			view.Properties[key] = buildPropertyView(child, requiredSet[key])
		}
		if len(view.Properties) == 0 {
			view.Properties = nil
		}
	}
	if rawItems, ok := schemaDoc["items"].(map[string]interface{}); ok {
		itemView := buildPropertyView(rawItems, false)
		view.Items = &itemView
	}
	return view
}

func requiredFieldNames(schemaDoc map[string]interface{}) []string {
	items, ok := schemaDoc["required"].([]interface{})
	if !ok {
		return nil
	}
	names := make([]string, 0, len(items))
	for _, item := range items {
		name := fmt.Sprint(item)
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names
}

func stringValue(value interface{}) string {
	if typed, ok := value.(string); ok {
		return typed
	}
	return ""
}

func boolValue(value interface{}) bool {
	typed, _ := value.(bool)
	return typed
}

func interfaceSlice(value interface{}) []interface{} {
	typed, _ := value.([]interface{})
	return typed
}

func numberValue(value interface{}) (float64, bool) {
	return number(value)
}

func intPointer(value float64, ok bool) *int {
	if !ok {
		return nil
	}
	n := int(value)
	return &n
}

func floatPointer(value float64, ok bool) *float64 {
	if !ok {
		return nil
	}
	n := value
	return &n
}
