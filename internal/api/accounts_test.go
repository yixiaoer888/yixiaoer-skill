package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yixiaoer/yixiaoer-skill/internal/config"
)

func TestAccountsUsesExpectedEndpointAndPlatformQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/platform/accounts" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("platform"); got != "抖音" {
			t.Fatalf("unexpected platform query: %s", got)
		}
		if got := r.Header.Get("Authorization"); got != "test-key" {
			t.Fatalf("unexpected authorization header: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"platformAccountId": "acc_1", "name": "抖音账号", "status": 1},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	accounts, err := client.Accounts("抖音")
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) != 1 {
		t.Fatalf("expected one account, got %d", len(accounts))
	}
	if id := AccountID(accounts[0]); id != "acc_1" {
		t.Fatalf("unexpected account id: %s", id)
	}
	if status := AccountStatus(accounts[0]); status != 1 {
		t.Fatalf("unexpected status: %d", status)
	}
}

func TestAccountsAcceptsNestedPaginatedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"data": []map[string]interface{}{
					{"platformAccountId": "acc_1", "platformAccountName": "抖音账号", "loginStatus": 1},
				},
				"total": 1,
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	accounts, err := client.Accounts("抖音")
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) != 1 {
		t.Fatalf("expected one account, got %d", len(accounts))
	}
	if id := AccountID(accounts[0]); id != "acc_1" {
		t.Fatalf("unexpected account id: %s", id)
	}
	if status := AccountStatus(accounts[0]); status != 1 {
		t.Fatalf("unexpected status: %d", status)
	}
}
