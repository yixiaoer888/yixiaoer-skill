package cmd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestSchemaListCommandOutputsAgentDiscoverableItems(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := schemaListCmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	if data["count"].(float64) == 0 {
		t.Fatal("expected schema list count")
	}
	items := data["items"].([]interface{})
	found := false
	for _, item := range items {
		entry := item.(map[string]interface{})
		if entry["platform"] == "douyin" && entry["type"] == "video" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected douyin video schema in items")
	}
}

func TestSchemaGetCommandOutputsSchemaForChinesePlatformAlias(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := schemaGetCmd.RunE(cmd, []string{"抖音", "video"}); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	if data["platform"] != "douyin" || data["type"] != "video" {
		t.Fatalf("unexpected schema metadata: %#v", data)
	}
	if data["key"] != "douyin/video" {
		t.Fatalf("unexpected schema key: %#v", data["key"])
	}
	if data["rootSchema"] != "schemas/publish.schema.json" {
		t.Fatalf("unexpected root schema: %#v", data["rootSchema"])
	}
	schemaDoc := data["document"].(map[string]interface{})
	required := schemaDoc["required"].([]interface{})
	if len(required) == 0 {
		t.Fatal("expected required fields in schema")
	}
	properties := schemaDoc["properties"].(map[string]interface{})
	title := properties["title"].(map[string]interface{})
	if title["type"] != "string" {
		t.Fatalf("expected title field view, got %#v", title)
	}
	if _, ok := schemaDoc["$schema"]; ok {
		t.Fatal("expected document output to omit raw $schema metadata")
	}
}

func TestSchemaCatalogCommandOutputsRootSchemasAndPlatforms(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := schemaCatalogCmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	rootSchemas := data["rootSchemas"].([]interface{})
	if len(rootSchemas) < 2 {
		t.Fatalf("expected root schemas, got %#v", data)
	}
	platforms := data["platforms"].([]interface{})
	if len(platforms) == 0 {
		t.Fatal("expected platform schema entries")
	}
}

func TestSchemaFieldsCommandOutputsFieldView(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)

	if err := schemaFieldsCmd.RunE(cmd, []string{"抖音", "video"}); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	if data["key"] != "douyin/video" {
		t.Fatalf("unexpected schema key: %#v", data["key"])
	}
	fields := data["fields"].(map[string]interface{})
	title := fields["title"].(map[string]interface{})
	if title["required"] != true {
		t.Fatalf("expected title to be required, got %#v", title)
	}
}

func withGoBuildCache(t *testing.T) {
	t.Helper()
	repoRoot, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOCACHE", filepath.Join(repoRoot, ".gocache"))
	t.Setenv("GOMODCACHE", filepath.Join(repoRoot, ".gomodcache"))
}
