package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestTaobaoGuangheSchemasAreDiscoverable(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	for _, publishType := range []string{"video", "imageText"} {
		t.Run(publishType, func(t *testing.T) {
			var out bytes.Buffer
			cmd := newSchemaGetCmd()
			cmd.SetOut(&out)
			cmd.SetArgs([]string{"淘宝光合", publishType})
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
			var response map[string]interface{}
			if err := json.Unmarshal(out.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			data := response["data"].(map[string]interface{})
			if data["key"] != "taobaoguanghe/"+publishType {
				t.Fatalf("unexpected schema key: %#v", data["key"])
			}
			examples := data["dynamicFieldExamples"].(map[string]interface{})
			cart := examples["shopping_cart"].(map[string]interface{})
			want := "yxer query taobao-guanghe-goods <account_id> --type " + publishType + " --source <source> --json"
			if cart["queryCommand"] != want {
				t.Fatalf("unexpected query command: %#v", cart)
			}
		})
	}
}

func TestTaobaoGuangheImageTextTemplateUsesFirstImageCover(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	var out bytes.Buffer
	cmd := newSchemaGetCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"taobao-guanghe", "imageText"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var response map[string]interface{}
	_ = json.Unmarshal(out.Bytes(), &response)
	form := response["data"].(map[string]interface{})["minimalTemplate"].(map[string]interface{})["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
	if _, ok := form["cover"]; ok {
		t.Fatalf("did not expect separate cover: %#v", form)
	}
	if len(form["images"].([]interface{})) != 1 {
		t.Fatalf("expected image placeholder: %#v", form)
	}
}
