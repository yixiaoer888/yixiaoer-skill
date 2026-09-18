# 微信公众号图文发布支持 Implementation Plan

> This plan records the approved in-chat requirements for the `yixiaoer-cli-weixingongzhonghao` branch.

**Goal:** Add `imageText` publishing for `微信公众号` while keeping the existing `article` platform-form contract unchanged.

**Architecture:** Reuse the standard `publishArgs.accountForms[].contentPublishForm` flow. Add a dedicated platform schema and a narrow normalization/preflight adapter for defaults, first-image cover derivation, and the two-hour scheduling rule.

**Contract:**

- `images` is required and accepts multiple uploaded resources; `images[0]` is the cover.
- `title` has a 20-character maximum; `description` is a separate image-text description field.
- `notifySubscribers`: `0` no group send (default), `1` group send.
- `sex`: `0` all (default), `1` male, `2` female.
- `scheduledTime` is optional; when present it must be at least two hours after the current time.
- `commentType`: `0` closed (default), `1` followers, `2` followers for at least seven days, `3` everyone.
- `declaration`: `0` no declaration (default), followed by the five requested creation-source choices.
- `recommend`: `true` allow platform recommendation (default), `false` disallow.
- `pubType`: `1` publish directly (default), `0` save to the platform draft box.

## Implementation Steps

1. Add red tests for alias resolution, schema fields/defaults, required uploaded images, first-image cover derivation, and schedule validation.
2. Add the `weixin.account.imageText` schema and platform alias/cover metadata.
3. Normalize the default settings and enforce the two-hour schedule boundary in image-text preflight.
4. Add platform docs and query/schema discovery coverage, including `WeiXinGongZhongHao` aliases.
5. Run targeted and full Go tests plus CLI dry-run/schema smoke checks.

## Global Constraints

- Keep stdout JSON-only and preserve structured CLI errors.
- Do not change the existing `weixin.account` `article` form or its `platformForms` placement.
- Keep all write flows behind existing `--dry-run` support.
