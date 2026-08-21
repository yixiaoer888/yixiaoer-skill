package platform

import "testing"

func TestWeixinAccountAliasesResolveToCanonicalPlatform(t *testing.T) {
	for _, alias := range []string{"WeiXinGongZhongHao", "weixingongzhonghao", "微信公众号"} {
		if got := CanonicalKey(alias); got != "weixin.account" {
			t.Fatalf("CanonicalKey(%q) = %q, want weixin.account", alias, got)
		}
		if got := ChineseName(alias); got != "微信公众号" {
			t.Fatalf("ChineseName(%q) = %q, want 微信公众号", alias, got)
		}
	}
}

func TestWeixinAccountImageTextUsesFirstImageAsCover(t *testing.T) {
	if !ImageTextUsesFirstImageAsCover("WeiXinGongZhongHao") {
		t.Fatal("expected 微信公众号 imageText to derive cover from images[0]")
	}
}
