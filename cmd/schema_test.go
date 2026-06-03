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
	// businessFields holds the platform-specific content fields directly.
	businessFields := data["businessFields"].(map[string]interface{})
	title := businessFields["title"].(map[string]interface{})
	if title["type"] != "string" || title["required"] != true {
		t.Fatalf("expected required string title in businessFields, got %#v", title)
	}
	// minimalTemplate provides a ready-to-edit skeleton using the standard envelope.
	template := data["minimalTemplate"].(map[string]interface{})
	if template["action"] != "publish" {
		t.Fatalf("expected minimalTemplate action=publish, got %#v", template)
	}
	templateArgs := template["publishArgs"].(map[string]interface{})
	templateForms := templateArgs["accountForms"].([]interface{})
	if len(templateForms) != 1 {
		t.Fatalf("expected single template account form, got %#v", templateForms)
	}
	// default (non-verbose) output must omit the debug-only schema views.
	for _, key := range []string{"fullDocument", "accountFormSchema", "contentPublishFormSchema", "businessSchema"} {
		if _, ok := data[key]; ok {
			t.Fatalf("expected default schema.get output to omit %q", key)
		}
	}
	if data["recommendedCommand"] != "yxer schema fields 抖音 video" {
		t.Fatalf("expected recommended schema.fields command, got %#v", data["recommendedCommand"])
	}
	guidance := data["guidance"].([]interface{})
	if len(guidance) < 3 {
		t.Fatalf("expected schema.get guidance, got %#v", guidance)
	}
}

func TestSchemaGetCommandVerboseOutputsDebugViews(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)
	cmd.Flags().Bool("verbose", false, "")
	if err := cmd.Flags().Set("verbose", "true"); err != nil {
		t.Fatal(err)
	}
	schemaGetVerbose = true
	t.Cleanup(func() {
		schemaGetVerbose = false
	})

	if err := schemaGetCmd.RunE(cmd, []string{"抖音", "video"}); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	accountFormSchema := data["accountFormSchema"].(map[string]interface{})
	accountFormProps := accountFormSchema["properties"].(map[string]interface{})
	if accountFormProps["platformAccountId"].(map[string]interface{})["required"] != true {
		t.Fatalf("expected accountFormSchema to require platformAccountId, got %#v", accountFormProps["platformAccountId"])
	}
	contentSchema := data["contentPublishFormSchema"].(map[string]interface{})
	contentProps := contentSchema["properties"].(map[string]interface{})
	if contentProps["title"].(map[string]interface{})["required"] != true {
		t.Fatalf("expected contentPublishFormSchema title to be required, got %#v", contentProps["title"])
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
	if data["recommendedResponse"] != "required + optional（按需查看 complex）" {
		t.Fatalf("expected grouped recommended response, got %#v", data["recommendedResponse"])
	}
	flatFields := data["flatFields"].([]interface{})
	if len(flatFields) == 0 {
		t.Fatal("expected compact flatFields view")
	}
	first := flatFields[0].(map[string]interface{})
	if first["path"] != "action" || first["required"] != true {
		t.Fatalf("expected required root field first in flatFields, got %#v", first)
	}
	foundTitle := false
	for _, entry := range flatFields {
		item := entry.(map[string]interface{})
		if item["path"] == "publishArgs.accountForms[].contentPublishForm.title" {
			foundTitle = true
			if item["type"] != "string" || item["required"] != true {
				t.Fatalf("expected title in flatFields to be required string, got %#v", item)
			}
			break
		}
	}
	if !foundTitle {
		t.Fatal("expected contentPublishForm.title in flatFields")
	}
	fields := data["fields"].(map[string]interface{})
	publishArgs := fields["publishArgs"].(map[string]interface{})
	accountForms := publishArgs["properties"].(map[string]interface{})["accountForms"].(map[string]interface{})
	title := accountForms["items"].(map[string]interface{})["properties"].(map[string]interface{})["contentPublishForm"].(map[string]interface{})["properties"].(map[string]interface{})["title"].(map[string]interface{})
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
