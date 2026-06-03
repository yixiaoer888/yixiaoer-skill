# AGENTS.md

## Purpose

`yxer` is consumed by both humans and AI agents. Command design, output shape, and error messages must remain stable and machine-readable.

## Output Contract

- `stdout` is for JSON data only.
- `stderr` is for diagnostics, warnings, prompts, and human guidance.
- Do not mix explanatory text into JSON responses.

## Error Contract

- Return structured `yxerrors.Error` values instead of bare `fmt.Errorf` for agent-facing failures.
- Every actionable error should include:
  - `code`
  - `message`
  - `category`
  - `hint` when a repair path exists
  - `nextCommand` when a concrete next step exists
  - `retryable` when retry is safe

## Workflow Rules

- Write operations should provide a `--dry-run` mode before the real side effect.
- Skills and workflows must treat CLI output as the source of truth; do not invent payload fields.
- Dynamic platform objects must come from CLI query results, not handwritten JSON.

## Test Expectations

- Agent-facing output changes require tests for JSON shape stability.
- New or changed write flows should cover `--dry-run` output before real execution.
