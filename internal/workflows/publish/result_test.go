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
