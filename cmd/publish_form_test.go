package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPublishFormSessionCanSetNestedAndArrayValues(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"抖音", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}
	initial, err := readPublishFormSession(sessionPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := setJSONPath(initial.Payload, "publishArgs.accountForms[0].contentPublishForm.title", "页面标题"); err != nil {
		t.Fatalf("direct path update: %v", err)
	}

	setTitle := newPublishFormSetCmd()
	valuePath := filepath.Join(t.TempDir(), "title.json")
	if err := os.WriteFile(valuePath, []byte(`"页面标题"`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := setTitle.Flags().Set("value-file", valuePath); err != nil {
		t.Fatal(err)
	}
	setTitle.SetArgs([]string{sessionPath, "publishArgs.accountForms[0].contentPublishForm.title"})
	parsed, parseErr := readFormValue("", valuePath)
	if parseErr != nil {
		t.Fatal(parseErr)
	}
	checkSession, readErr := readPublishFormSession(sessionPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if pathErr := setJSONPath(cloneJSONMap(checkSession.Payload), "publishArgs.accountForms[0].contentPublishForm.title", parsed); pathErr != nil {
		t.Fatalf("same flow path update: %v", pathErr)
	}
	setTitle.SetOut(&bytes.Buffer{})
	if err := setTitle.Execute(); err != nil {
		t.Logf("set error: %#v", err)
		t.Fatal(err)
	}

	setTag := newPublishFormSetCmd()
	tagPath := filepath.Join(t.TempDir(), "tag.json")
	if err := os.WriteFile(tagPath, []byte(`"标签"`), 0o644); err != nil {
		t.Fatal(err)
	}
	setTag.SetArgs([]string{sessionPath, "publishArgs.accountForms[0].contentPublishForm.tags[0]", "--value-file", tagPath})
	setTag.SetOut(&bytes.Buffer{})
	if err := setTag.Execute(); err != nil {
		t.Fatal(err)
	}

	session, err := readPublishFormSession(sessionPath)
	if err != nil {
		t.Fatal(err)
	}
	form := session.Payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if form["title"] != "页面标题" {
		t.Fatalf("unexpected title: %#v", form["title"])
	}
	if form["tags"].([]interface{})[0] != "标签" {
		t.Fatalf("unexpected tag: %#v", form["tags"])
	}
}

func TestPublishFormStartDryRunDoesNotWrite(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")
	var out bytes.Buffer
	cmd := newPublishFormStartCmd()
	cmd.SetArgs([]string{"抖音", "video", "--output", sessionPath, "--dry-run"})
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(sessionPath); !os.IsNotExist(err) {
		t.Fatalf("expected dry-run not to write session, stat error=%v", err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["action"] != "publish.form.start.dry-run" {
		t.Fatalf("unexpected action: %#v", response["action"])
	}
}

func TestPublishFormExportProducesStandardPayload(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")
	payloadPath := filepath.Join(t.TempDir(), "payload.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"知乎", "article", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}
	export := newPublishFormExportCmd()
	export.SetArgs([]string{sessionPath, "--output", payloadPath})
	export.SetOut(&bytes.Buffer{})
	if err := export.Execute(); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(payloadPath)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["action"] != "publish" || payload["publishType"] != "article" {
		t.Fatalf("unexpected exported payload: %#v", payload)
	}
}

func TestPublishFormExportNextCommandsPreserveLocalPublishMode(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")
	payloadPath := filepath.Join(t.TempDir(), "payload.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"抖音", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}
	session, err := readPublishFormSession(sessionPath)
	if err != nil {
		t.Fatal(err)
	}
	session.Payload["publishChannel"] = "local"
	session.Payload["clientId"] = "local-client-1"
	if err := writePublishFormSession(sessionPath, session); err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	export := newPublishFormExportCmd()
	export.SetArgs([]string{sessionPath, "--output", payloadPath, "--dry-run"})
	export.SetOut(&out)
	if err := export.Execute(); err != nil {
		t.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]interface{})
	next := data["next"].([]interface{})
	for _, raw := range next {
		command := raw.(string)
		if strings.Contains(command, "validate") || (strings.Contains(command, "publish") && !strings.Contains(command, "publish form")) {
			if !strings.Contains(command, "--publish-channel local") || !strings.Contains(command, "--client-id local-client-1") {
				t.Fatalf("expected local publish mode in next command, got %q", command)
			}
		}
	}
}

func TestPublishFormChooseSelectsCandidateAndRecordsSource(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")
	candidatesPath := filepath.Join(t.TempDir(), "categories.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"抖音", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(candidatesPath, []byte(`{"ok":true,"data":{"items":[{"yixiaoerId":"loc_sh","yixiaoerName":"上海","raw":{"id":"loc_sh"}},{"yixiaoerId":"loc_hz","yixiaoerName":"杭州","raw":{"id":"loc_hz"}}]}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	choose := newPublishFormChooseCmd()
	choose.SetArgs([]string{sessionPath, "location", "--value-file", candidatesPath, "--id", "loc_sh", "--source-command", "yxer query locations acc_1 --query 上海 --json"})
	choose.SetOut(&bytes.Buffer{})
	if err := choose.Execute(); err != nil {
		t.Fatal(err)
	}

	session, err := readPublishFormSession(sessionPath)
	if err != nil {
		t.Fatal(err)
	}
	form := session.Payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	location := form["location"].(map[string]interface{})
	if location["yixiaoerId"] != "loc_sh" {
		t.Fatalf("unexpected location: %#v", location)
	}
	if len(session.Sources) != 1 {
		t.Fatalf("expected one source record, got %#v", session.Sources)
	}
	if session.Sources[0].Kind != "query" || session.Sources[0].Path != "publishArgs.accountForms[0].contentPublishForm.location" {
		t.Fatalf("unexpected source record: %#v", session.Sources[0])
	}
	if session.Sources[0].ValueHash == "" || session.Sources[0].RawHash == "" || session.Sources[0].FetchedAt == "" {
		t.Fatalf("expected query source hashes and fetchedAt, got %#v", session.Sources[0])
	}
}

func TestPublishFormChooseDramaUsesExactShapeAndSkipsRawHash(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")
	candidatesPath := filepath.Join(t.TempDir(), "drama.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"视频号", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(candidatesPath, []byte(`{"ok":true,"data":{"items":[{"yixiaoerId":"event/1","yixiaoerImageUrl":"http://wxapp.tc.qq.com/cover","yixiaoerName":"风浪过后护妻安康"},{"yixiaoerId":"event/2","yixiaoerImageUrl":"http://wxapp.tc.qq.com/cover-2","yixiaoerName":"另一部剧"}]}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	choose := newPublishFormChooseCmd()
	choose.SetArgs([]string{sessionPath, "drama", "--value-file", candidatesPath, "--id", "event/1", "--source-command", "yxer query drama-tasks acc_1 --query 护妻 --json"})
	choose.SetOut(&bytes.Buffer{})
	if err := choose.Execute(); err != nil {
		t.Fatal(err)
	}

	session, err := readPublishFormSession(sessionPath)
	if err != nil {
		t.Fatal(err)
	}
	form := session.Payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	drama := form["drama"].(map[string]interface{})
	if drama["yixiaoerId"] != "event/1" || drama["yixiaoerName"] != "风浪过后护妻安康" {
		t.Fatalf("unexpected drama: %#v", drama)
	}
	if _, exists := drama["raw"]; exists {
		t.Fatalf("drama must not gain raw during choose: %#v", drama)
	}
	if len(session.Sources) != 1 {
		t.Fatalf("expected one drama source, got %#v", session.Sources)
	}
	source := session.Sources[0]
	if source.Path != "publishArgs.accountForms[0].contentPublishForm.drama" || source.ValueHash == "" || source.RawHash != "" || source.FetchedAt == "" {
		t.Fatalf("unexpected drama source provenance: %#v", source)
	}

	if _, err := validatePublishFormProvenance(session); err != nil {
		t.Fatalf("expected raw-free drama provenance to validate, got %v", err)
	}
}

func TestPublishFormChooseDramaRequiresRecognizedExactCandidates(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{name: "empty list", value: map[string]interface{}{"items": []interface{}{}}},
		{name: "wrong list type", value: map[string]interface{}{"items": map[string]interface{}{}}},
		{name: "multiple list keys", value: map[string]interface{}{"items": []interface{}{}, "results": []interface{}{}}},
		{name: "list and nested data", value: map[string]interface{}{"items": []interface{}{}, "data": map[string]interface{}{"items": []interface{}{}}}},
		{name: "scalar", value: "not a candidate"},
		{name: "extra field", value: map[string]interface{}{
			"yixiaoerId":       "event/1",
			"yixiaoerImageUrl": "http://wxapp.tc.qq.com/cover",
			"yixiaoerName":     "风浪过后护妻安康",
			"raw":              map[string]interface{}{"id": "event/1"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := selectPublishFormCandidateForField(tt.value, "drama", -1, ""); err == nil {
				t.Fatal("expected invalid drama candidate result")
			}
		})
	}
}

func TestPublishFormChooseDramaAcceptsSupportedCandidateEnvelopes(t *testing.T) {
	candidate := map[string]interface{}{
		"yixiaoerId":       "event/1",
		"yixiaoerImageUrl": "http://wxapp.tc.qq.com/cover",
		"yixiaoerName":     "风浪过后护妻安康",
	}
	candidates := []interface{}{candidate}
	tests := []struct {
		name  string
		value interface{}
	}{
		{name: "array", value: candidates},
		{name: "single object", value: candidate},
		{name: "items", value: map[string]interface{}{"items": candidates}},
		{name: "list", value: map[string]interface{}{"list": candidates}},
		{name: "dataList", value: map[string]interface{}{"dataList": candidates}},
		{name: "records", value: map[string]interface{}{"records": candidates}},
		{name: "results", value: map[string]interface{}{"results": candidates}},
		{name: "nested data", value: map[string]interface{}{"data": map[string]interface{}{"items": candidates}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selected, got, err := selectPublishFormCandidateForField(tt.value, "drama", -1, "")
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 || selected.(map[string]interface{})["yixiaoerId"] != "event/1" {
				t.Fatalf("unexpected selected drama candidate: %#v, candidates=%#v", selected, got)
			}
		})
	}
}

func TestPublishFormSetCannotWriteDramaField(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"视频号", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}

	setDrama := newPublishFormSetCmd()
	setDrama.SetArgs([]string{sessionPath, "publishArgs.accountForms[0].contentPublishForm.drama", "--value", `{"yixiaoerId":"event/1","yixiaoerImageUrl":"http://wxapp.tc.qq.com/cover","yixiaoerName":"风浪过后护妻安康"}`})
	setDrama.SetOut(&bytes.Buffer{})
	err := setDrama.Execute()
	if err == nil || !strings.Contains(err.Error(), "form set cannot write drama field") {
		t.Fatalf("expected direct drama set rejection, got %v", err)
	}
	session, readErr := readPublishFormSession(sessionPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	form := session.Payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})["contentPublishForm"].(map[string]interface{})
	if _, exists := form["drama"]; exists || len(session.Sources) != 0 {
		t.Fatalf("rejected set must not update payload or sources: %#v, %#v", form, session.Sources)
	}
}

func TestPublishFormChooseDramaRejectsPathOverrideAndWrongQuery(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")
	value := `{"yixiaoerId":"event/1","yixiaoerImageUrl":"http://wxapp.tc.qq.com/cover","yixiaoerName":"风浪过后护妻安康"}`

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"视频号", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}

	wrongQuery := newPublishFormChooseCmd()
	wrongQuery.SetArgs([]string{sessionPath, "drama", "--value", value, "--source-command", "yxer query collections acc_1 --type video --json"})
	wrongQuery.SetOut(&bytes.Buffer{})
	if err := wrongQuery.Execute(); err == nil || !strings.Contains(err.Error(), "drama source must come from yxer query drama-tasks") {
		t.Fatalf("expected drama query command rejection, got %v", err)
	}

	wrongPath := newPublishFormChooseCmd()
	wrongPath.SetArgs([]string{sessionPath, "drama", "--value", value, "--path", "publishArgs.accountForms[0].contentPublishForm.title", "--source-command", "yxer query drama-tasks acc_1 --json"})
	wrongPath.SetOut(&bytes.Buffer{})
	if err := wrongPath.Execute(); err == nil || !strings.Contains(err.Error(), "drama field path must match contract") {
		t.Fatalf("expected drama path rejection, got %v", err)
	}

	session, err := readPublishFormSession(sessionPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(session.Sources) != 0 {
		t.Fatalf("rejected drama choices must not record sources: %#v", session.Sources)
	}
}

func TestPublishFormProvenanceRequiresDramaQuerySource(t *testing.T) {
	payload := map[string]interface{}{
		"publishArgs": map[string]interface{}{
			"accountForms": []interface{}{
				map[string]interface{}{
					"platformAccountId": "acc_1",
					"contentPublishForm": map[string]interface{}{
						"drama": map[string]interface{}{
							"yixiaoerId":       "event/1",
							"yixiaoerImageUrl": "http://wxapp.tc.qq.com/cover",
							"yixiaoerName":     "风浪过后护妻安康",
						},
					},
				},
			},
		},
	}

	report, err := validatePublishFormProvenance(publishFormSession{Payload: payload})
	if err == nil {
		t.Fatal("expected missing drama query source error")
	}
	if report["valid"] != false {
		t.Fatalf("expected invalid provenance report, got %#v", report)
	}
	if !strings.Contains(fmt.Sprint(report["errors"]), "drama") {
		t.Fatalf("expected drama source error, got %#v", report["errors"])
	}
}

func TestPublishFormChooseRequiresTargetForMultipleAccountForms(t *testing.T) {
	session := publishFormSession{
		Kind:     "yxer.publish-form",
		Version:  1,
		Platform: "抖音",
		Type:     "video",
		Contract: map[string]interface{}{},
		Payload: map[string]interface{}{
			"publishArgs": map[string]interface{}{
				"accountForms": []interface{}{
					map[string]interface{}{"platformAccountId": "acc_1", "contentPublishForm": map[string]interface{}{}},
					map[string]interface{}{"platformAccountId": "acc_2", "contentPublishForm": map[string]interface{}{}},
				},
			},
		},
	}
	sessionPath := filepath.Join(t.TempDir(), "form.json")
	if err := writePublishFormSession(sessionPath, session); err != nil {
		t.Fatal(err)
	}

	choose := newPublishFormChooseCmd()
	choose.SetArgs([]string{sessionPath, "category", "--value", `{"yixiaoerId":"cat_food"}`, "--source-command", "yxer query categories acc_1 --type video --json"})
	choose.SetOut(&bytes.Buffer{})
	err := choose.Execute()
	if err == nil {
		t.Fatal("expected ambiguous target error")
	}
}

func TestPublishFormSetRejectsUndeclaredPath(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"抖音", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}

	setUnknown := newPublishFormSetCmd()
	setUnknown.SetArgs([]string{sessionPath, "publishArgs.accountForms[0].contentPublishForm.titel", "--value", `"错字字段"`})
	setUnknown.SetOut(&bytes.Buffer{})
	err := setUnknown.Execute()
	if err == nil {
		t.Fatal("expected undeclared path error")
	}
	if !strings.Contains(err.Error(), "form path is not declared in contract") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPublishFormChooseRejectsUndeclaredDynamicField(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"抖音", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}

	choose := newPublishFormChooseCmd()
	choose.SetArgs([]string{sessionPath, "not_a_dynamic_field", "--value", `{"yixiaoerId":"id_1","yixiaoerName":"名称","raw":{"id":"id_1"}}`, "--source-command", "yxer query locations acc_1 --json"})
	choose.SetOut(&bytes.Buffer{})
	err := choose.Execute()
	if err == nil {
		t.Fatal("expected undeclared dynamic field error")
	}
	if !strings.Contains(err.Error(), "form choose field is not query-backed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPublishFormSessionRecordsContractHashAndRejectsStaleHash(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"抖音", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}

	session, err := readPublishFormSession(sessionPath)
	if err != nil {
		t.Fatal(err)
	}
	if session.ContractHash == "" {
		t.Fatal("expected new session to include contractHash")
	}
	session.ContractHash = "sha256:stale"
	if err := writePublishFormSession(sessionPath, session); err != nil {
		t.Fatal(err)
	}
	if _, err := readPublishFormSession(sessionPath); err == nil || !strings.Contains(err.Error(), "publish form contract is stale") {
		t.Fatalf("expected stale contract error, got %v", err)
	}
}

func TestPublishFormChooseRequiresSourceCommand(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"抖音", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}

	choose := newPublishFormChooseCmd()
	choose.SetArgs([]string{sessionPath, "location", "--value", `{"yixiaoerId":"loc_sh","yixiaoerName":"上海","raw":{"id":"loc_sh"}}`})
	choose.SetOut(&bytes.Buffer{})
	err := choose.Execute()
	if err == nil || !strings.Contains(err.Error(), "form choose requires --source-command") {
		t.Fatalf("expected source command error, got %v", err)
	}
}

func TestPublishFormChooseRejectsCrossAccountQuerySource(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"抖音", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}

	setAccount := newPublishFormSetCmd()
	setAccount.SetArgs([]string{sessionPath, "publishArgs.accountForms[0].platformAccountId", "--value", `"acc_2"`})
	setAccount.SetOut(&bytes.Buffer{})
	if err := setAccount.Execute(); err != nil {
		t.Fatal(err)
	}

	choose := newPublishFormChooseCmd()
	choose.SetArgs([]string{sessionPath, "location", "--account-id", "acc_2", "--value", `{"yixiaoerId":"loc_sh","yixiaoerName":"上海","raw":{"id":"loc_sh"}}`, "--source-command", "yxer query locations acc_1 --query 上海 --json"})
	choose.SetOut(&bytes.Buffer{})
	err := choose.Execute()
	if err == nil || !strings.Contains(err.Error(), "form choose source account does not match target account") {
		t.Fatalf("expected cross-account source error, got %v", err)
	}
}

func TestPublishFormReviewRejectsSourcePayloadDrift(t *testing.T) {
	withRepoRoot(t)
	withGoBuildCache(t)
	sessionPath := filepath.Join(t.TempDir(), "form.json")

	start := newPublishFormStartCmd()
	start.SetArgs([]string{"抖音", "video", "--output", sessionPath})
	start.SetOut(&bytes.Buffer{})
	if err := start.Execute(); err != nil {
		t.Fatal(err)
	}

	choose := newPublishFormChooseCmd()
	choose.SetArgs([]string{sessionPath, "location", "--value", `{"yixiaoerId":"loc_sh","yixiaoerName":"上海","raw":{"id":"loc_sh"}}`, "--source-command", "yxer query locations acc_1 --query 上海 --json"})
	choose.SetOut(&bytes.Buffer{})
	if err := choose.Execute(); err != nil {
		t.Fatal(err)
	}

	session, err := readPublishFormSession(sessionPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := setJSONPath(session.Payload, "publishArgs.accountForms[0].contentPublishForm.location", map[string]interface{}{
		"yixiaoerId":   "loc_other",
		"yixiaoerName": "其他位置",
		"raw":          map[string]interface{}{"id": "loc_other"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := writePublishFormSession(sessionPath, session); err != nil {
		t.Fatal(err)
	}

	review := newPublishFormReviewCmd()
	review.SetArgs([]string{sessionPath, "--dry-run"})
	review.SetOut(&bytes.Buffer{})
	err = review.Execute()
	if err == nil || !strings.Contains(err.Error(), "publish form provenance validation failed") {
		t.Fatalf("expected provenance validation error, got %v", err)
	}
}

func TestReadFormValueAcceptsPlainTextAndCLIEnvelope(t *testing.T) {
	value, err := readFormValue("页面标题", "")
	if err != nil || value != "页面标题" {
		t.Fatalf("expected plain text value, got %#v, err=%v", value, err)
	}
	value, err = readFormValue(`{"ok":true,"data":{"yixiaoerId":"id_1"}}`, "")
	if err != nil {
		t.Fatal(err)
	}
	object := value.(map[string]interface{})
	if object["yixiaoerId"] != "id_1" {
		t.Fatalf("expected CLI data payload, got %#v", object)
	}
}
