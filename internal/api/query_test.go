package api

import "testing"

func TestPlatformDocFileNameUsesShipinhaoAliasForImageText(t *testing.T) {
	if got := platformDocFileName("shipinhao", "imageText"); got != "shipinhao.md" {
		t.Fatalf("expected shipinhao imageText doc file, got %q", got)
	}
	if got := platformDocFileName("douyin", "imageText"); got != "douyin.md" {
		t.Fatalf("expected default platform doc file, got %q", got)
	}
}
