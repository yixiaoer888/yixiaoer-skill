package publish

import "testing"

func TestExecuteEnvelopeWrapsPublishResult(t *testing.T) {
	result := map[string]interface{}{"taskSetId": "task_set_1"}
	svc := Service{}

	got, err := svc.wrapExecuteEnvelope(result, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Action != "publish" {
		t.Fatalf("expected publish action, got %q", got.Action)
	}
	if got.Data["taskSetId"] != "task_set_1" {
		t.Fatalf("expected publish result in envelope, got %+v", got.Data)
	}
}

func TestDryRunEnvelopeWrapsDryRunResult(t *testing.T) {
	svc := Service{}
	result := DryRunResult{
		Platform:      "抖音",
		PublishType:   "video",
		PublishBody:   map[string]interface{}{"publishType": "video"},
		PublishMode:   "cloud",
		ClientID:      "",
		AccountIDs:    []string{"acc_001"},
		PlatformDraft: false,
		YixiaoerDraft: false,
		SchemaChecked: true,
	}

	got, err := svc.wrapDryRunEnvelope(result, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Action != "publish.dry-run" {
		t.Fatalf("expected publish.dry-run action, got %q", got.Action)
	}
	if got.Data["dryRun"] != true {
		t.Fatalf("expected dryRun marker, got %+v", got.Data)
	}
	meta, ok := got.Data["meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected meta object, got %+v", got.Data["meta"])
	}
	if meta["platform"] != "抖音" || meta["publishType"] != "video" {
		t.Fatalf("unexpected dry-run meta: %+v", meta)
	}
}
