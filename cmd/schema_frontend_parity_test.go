package cmd

import (
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
)

func TestDouyinSchemasExposeCurrentDeclarationOptions(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	validator := schema.NewValidator("schemas")
	for _, publishType := range []string{"video", "imageText"} {
		doc, err := validator.Schema("抖音", publishType)
		if err != nil {
			t.Fatal(err)
		}
		declaration, ok := doc.Properties["declaration"]
		if !ok {
			t.Fatalf("%s schema has no declaration field", publishType)
		}
		if !containsSchemaEnum(declaration.Enum, float64(7)) || !containsSchemaEnum(declaration.Enum, float64(8)) {
			t.Fatalf("%s declaration enum does not include current frontend values: %#v", publishType, declaration.Enum)
		}
	}
}

func containsSchemaEnum(values []interface{}, expected float64) bool {
	for _, value := range values {
		if number, ok := value.(float64); ok && number == expected {
			return true
		}
	}
	return false
}
