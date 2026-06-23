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
				"page":      1,
				"size":      20,
				"totalPage": 1,
				"totalSize": 1,
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

func TestAccountsMapsPlatformKeyToChineseQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("platform"); got != "抖音" {
			t.Fatalf("unexpected platform query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"platformAccountId": "acc_1", "name": "抖音账号", "status": 1},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	if _, err := client.Accounts("douyin"); err != nil {
		t.Fatal(err)
	}
}

func TestAccountsAcceptsTopLevelListResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("platform"); got != "视频号" {
			t.Fatalf("unexpected platform query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"list": []map[string]interface{}{
				{"platformAccountId": "acc_1", "platformAccountName": "视频号账号", "loginStatus": 1},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	accounts, err := client.Accounts("视频号")
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) != 1 {
		t.Fatalf("expected one account, got %d", len(accounts))
	}
	if id := AccountID(accounts[0]); id != "acc_1" {
		t.Fatalf("unexpected account id: %s", id)
	}
}

func TestAccountsMapsLegacyShipinghaoAliasToCanonicalChineseQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("platform"); got != "视频号" {
			t.Fatalf("unexpected platform query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"platformAccountId": "acc_1", "platformAccountName": "视频号账号", "status": 1},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	if _, err := client.Accounts("shipinghao"); err != nil {
		t.Fatal(err)
	}
}

func TestAccountsPageUsesExplicitPageAndSize(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if got := r.URL.Query().Get("platform"); got != "抖音" {
			t.Fatalf("unexpected platform query: %s", got)
		}
		if got := r.URL.Query().Get("page"); got != "3" {
			t.Fatalf("unexpected page query: %s", got)
		}
		if got := r.URL.Query().Get("size"); got != "50" {
			t.Fatalf("unexpected size query: %s", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"list": []map[string]interface{}{
					{"platformAccountId": "acc_101", "platformAccountName": "账号101", "status": 1},
				},
				"page":      3,
				"size":      50,
				"totalPage": 3,
				"totalSize": 120,
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	accounts, meta, err := client.AccountsPage("抖音", 3, 50)
	if err != nil {
		t.Fatal(err)
	}
	if requests != 1 {
		t.Fatalf("expected one request, got %d", requests)
	}
	if len(accounts) != 1 {
		t.Fatalf("expected one account, got %d", len(accounts))
	}
	if meta.page != 3 || meta.size != 50 || meta.total != 120 {
		t.Fatalf("unexpected meta: %+v", meta)
	}
}

func TestAccountsDoesNotFetchNextPageWithoutPaginationMeta(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if got := r.URL.Query().Get("platform"); got != "抖音" {
			t.Fatalf("unexpected platform query: %s", got)
		}
		if got := r.URL.Query().Get("size"); got != "20" {
			t.Fatalf("unexpected size query: %s", got)
		}
		if page := r.URL.Query().Get("page"); page != "1" {
			t.Fatalf("unexpected page query: %s", page)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]interface{}{
				{"platformAccountId": "acc_1", "platformAccountName": "账号1", "status": 1},
				{"platformAccountId": "acc_2", "platformAccountName": "账号2", "status": 1},
			},
		})
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	accounts, err := client.Accounts("抖音")
	if err != nil {
		t.Fatal(err)
	}
	if requests != 1 {
		t.Fatalf("expected one request, got %d", requests)
	}
	if len(accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts))
	}
}

func TestAccountsStopsWhenPaginatedMetaReportsAllRowsFetched(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		page := r.URL.Query().Get("page")
		switch page {
		case "1":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"data": []map[string]interface{}{
						{"platformAccountId": "acc_1", "platformAccountName": "账号1", "status": 1},
						{"platformAccountId": "acc_2", "platformAccountName": "账号2", "status": 1},
					},
					"page":      1,
					"size":      2,
					"totalPage": 2,
					"totalSize": 3,
				},
			})
		case "2":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"data": []map[string]interface{}{
						{"platformAccountId": "acc_3", "platformAccountName": "账号3", "status": 1},
					},
					"page":      2,
					"size":      2,
					"totalPage": 2,
					"totalSize": 3,
				},
			})
		default:
			t.Fatalf("unexpected page query: %s", page)
		}
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	accounts, err := client.Accounts("")
	if err != nil {
		t.Fatal(err)
	}
	if requests != 2 {
		t.Fatalf("expected two requests, got %d", requests)
	}
	if len(accounts) != 3 {
		t.Fatalf("expected three accounts, got %d", len(accounts))
	}
}

func TestAccountsFetchesNextPageWhenTotalPageIndicatesMorePages(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		page := r.URL.Query().Get("page")
		switch page {
		case "1":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"list": []map[string]interface{}{
						{"platformAccountId": "acc_1", "platformAccountName": "账号1", "status": 1},
					},
					"page":      1,
					"size":      20,
					"totalPage": 2,
					"totalSize": 2,
				},
			})
		case "2":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"list": []map[string]interface{}{
						{"platformAccountId": "acc_2", "platformAccountName": "账号2", "status": 1},
					},
					"page":      2,
					"size":      20,
					"totalPage": 2,
					"totalSize": 2,
				},
			})
		default:
			t.Fatalf("unexpected page query: %s", page)
		}
	}))
	defer server.Close()

	client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
	accounts, err := client.Accounts("")
	if err != nil {
		t.Fatal(err)
	}
	if requests != 2 {
		t.Fatalf("expected two requests, got %d", requests)
	}
	if len(accounts) != 2 {
		t.Fatalf("expected two accounts, got %d", len(accounts))
	}
}
