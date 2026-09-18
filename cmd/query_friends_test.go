package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFriendsCommandSearchesByRedIDAndReturnsStranger(t *testing.T) {
	configureAPIKey(t, "test-key")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform-accounts/acc_xhs/friends" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("keyWord"); got != "26439197554" {
			t.Fatalf("unexpected gateway keyWord: %q", got)
		}
		if got := r.URL.Query().Get("keyword"); got != "" {
			t.Fatalf("client-only keyword must not be sent to the gateway: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"list": []map[string]interface{}{
					{
						"yixiaoerId":       "698747370000000021023230",
						"yixiaoerName":     "默默",
						"yixiaoerImageUrl": "https://example.com/avatar.jpg",
						"raw": map[string]interface{}{
							"red_id":        "26439197554",
							"user_id":       "698747370000000021023230",
							"user_nickname": "默默",
						},
					},
				},
			},
		})
	}))
	defer server.Close()
	useTestAPIBaseURL(t, server.URL)
	withRepoRoot(t)

	var stdout bytes.Buffer
	cmd := newFriendsCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"acc_xhs", "--red-id", "26439197554"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatalf("stdout must contain JSON only: %v; output=%s", err, stdout.String())
	}
	data := response["data"].(map[string]interface{})
	items := data["list"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected one stranger search result, got %#v", data["list"])
	}
	item := items[0].(map[string]interface{})
	if item["yixiaoerName"] != "默默" {
		t.Fatalf("unexpected search result: %#v", item)
	}
	raw := item["raw"].(map[string]interface{})
	if raw["red_id"] != "26439197554" {
		t.Fatalf("expected complete raw object to survive filtering, got %#v", raw)
	}
}
