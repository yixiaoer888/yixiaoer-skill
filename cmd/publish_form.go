package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type publishFormSession struct {
	Kind      string                 `json:"kind"`
	Version   int                    `json:"version"`
	PlanID    string                 `json:"planId,omitempty"`
	Platform  string                 `json:"platform"`
	Type      string                 `json:"type"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
	Contract  map[string]interface{} `json:"contract"`
	Payload   map[string]interface{} `json:"payload"`
	Sources   []publishFormSource    `json:"sources,omitempty"`
	Reviews   []publishFormReview    `json:"reviews,omitempty"`
}

type publishFormSource struct {
	Path          string      `json:"path"`
	Field         string      `json:"field,omitempty"`
	Kind          string      `json:"kind"`
	SourceCommand string      `json:"sourceCommand,omitempty"`
	Target        string      `json:"target,omitempty"`
	AccountIndex  int         `json:"accountIndex,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	UpdatedAt     string      `json:"updatedAt"`
}

type publishFormReview struct {
	PayloadHash string `json:"payloadHash"`
	ReviewedAt  string `json:"reviewedAt"`
}

func newPublishFormCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "form", Short: "按页面步骤维护可恢复的发布表单会话", Args: cobra.NoArgs}
	cmd.AddCommand(newPublishFormStartCmd())
	cmd.AddCommand(newPublishFormInspectCmd())
	cmd.AddCommand(newPublishFormSetCmd())
	cmd.AddCommand(newPublishFormChooseCmd())
	cmd.AddCommand(newPublishFormReviewCmd())
	cmd.AddCommand(newPublishFormExportCmd())
	return cmd
}

func newPublishFormStartCmd() *cobra.Command {
	var outputPath string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "start <platform> <type>",
		Short: "创建可恢复的发布表单会话",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			rt, err := app.Load()
			if err != nil {
				return err
			}
			doc, err := schema.NewValidator(rt.Config.SchemaDir).Schema(args[0], args[1])
			if err != nil {
				return yxerrors.Usage("schema not found", map[string]interface{}{"platform": args[0], "type": args[1]}).
					WithHint("请先用 yxer schema list 确认平台和类型。").WithNextCommand("yxer schema list")
			}
			target, err := absoluteFormPath(outputPath, "publish-form.json")
			if err != nil {
				return err
			}
			session := newPublishFormSession(doc)
			if dryRun {
				return output.Success(cmd.OutOrStdout(), "publish.form.start.dry-run", map[string]interface{}{"file": filepath.ToSlash(target), "session": session})
			}
			if err := writePublishFormSession(target, session); err != nil {
				return yxerrors.Internal("failed to write publish form session", err.Error())
			}
			return output.Success(cmd.OutOrStdout(), "publish.form.start", map[string]interface{}{"file": filepath.ToSlash(target), "session": session, "next": "yxer publish form inspect <session.json>"})
		},
	}
	cmd.Flags().StringVarP(&outputPath, "output", "o", "publish-form.json", "session file")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview the local session without writing it")
	return cmd
}

func newPublishFormInspectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "inspect <session.json>",
		Short: "查看表单步骤、字段和当前会话数据",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := readPublishFormSession(args[0])
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "publish.form.inspect", map[string]interface{}{
				"file": filepath.ToSlash(args[0]), "platform": session.Platform, "type": session.Type,
				"planId": session.PlanID, "contract": session.Contract, "payload": session.Payload,
				"sources": session.Sources, "reviews": session.Reviews,
			})
		},
	}
}

func newPublishFormSetCmd() *cobra.Command {
	var value, valueFile string
	var index int
	var sourceCommand string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "set <session.json> <path>",
		Short: "按页面字段路径设置一个 JSON 值",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(value) == "" && strings.TrimSpace(valueFile) == "" {
				return yxerrors.Usage("form set requires --value or --value-file", nil).
					WithHint("字符串也必须使用 JSON 格式，例如 --value '\"标题\"'；查询结果可先保存后使用 --value-file。").
					WithNextCommand("yxer publish form set --help")
			}
			parsed, err := readFormValue(value, valueFile)
			if err != nil {
				return yxerrors.Usage("form value is not valid JSON", err.Error())
			}
			if index >= 0 {
				items, ok := parsed.([]interface{})
				if !ok || index >= len(items) {
					return yxerrors.Usage("form selection index is out of range", map[string]interface{}{"index": index}).
						WithHint("--index 只适用于 query 返回的数组；请先 inspect 结果，再选择有效索引。")
				}
				parsed = items[index]
			}
			session, err := readPublishFormSession(args[0])
			if err != nil {
				return err
			}
			updated := cloneJSONMap(session.Payload)
			if err := setJSONPath(updated, args[1], parsed); err != nil {
				return yxerrors.Usage("form path cannot be updated", map[string]interface{}{"path": args[1], "cause": err.Error()}).
					WithHint("请使用 contract.fields 和 fieldPlacements 返回的路径；数组项使用 [0] 表示。").
					WithNextCommand("yxer publish form inspect " + args[0])
			}
			session.Payload = updated
			session.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			session.Sources = appendPublishFormSource(session.Sources, publishFormSource{
				Path:          args[1],
				Kind:          "manual",
				SourceCommand: sourceCommand,
				Value:         parsed,
				UpdatedAt:     session.UpdatedAt,
			})
			if dryRun {
				return output.Success(cmd.OutOrStdout(), "publish.form.set.dry-run", map[string]interface{}{"file": filepath.ToSlash(args[0]), "path": args[1], "value": parsed, "source": lastPublishFormSource(session.Sources), "payload": updated})
			}
			if err := writePublishFormSession(args[0], session); err != nil {
				return yxerrors.Internal("failed to write publish form session", err.Error())
			}
			return output.Success(cmd.OutOrStdout(), "publish.form.set", map[string]interface{}{"file": filepath.ToSlash(args[0]), "path": args[1], "value": parsed, "source": lastPublishFormSource(session.Sources), "payload": updated})
		},
	}
	cmd.Flags().StringVar(&value, "value", "", "JSON value to set")
	cmd.Flags().StringVar(&valueFile, "value-file", "", "file containing one JSON value")
	cmd.Flags().IntVar(&index, "index", -1, "select an item from an array value (for query results)")
	cmd.Flags().StringVar(&sourceCommand, "source-command", "", "command or note that produced this value")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview the updated session without writing it")
	return cmd
}

func newPublishFormChooseCmd() *cobra.Command {
	var value, valueFile string
	var index int
	var id, path, accountID, target, sourceCommand string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "choose <session.json> <field>",
		Short: "从 query 输出中选择候选并写入目标账号表单字段",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(value) == "" && strings.TrimSpace(valueFile) == "" {
				return yxerrors.Usage("form choose requires --value or --value-file", nil).
					WithHint("请把 yxer query ... --json 的输出通过 --value-file 传入，或将候选对象作为 --value 传入。")
			}
			rawValue, err := readFormValue(value, valueFile)
			if err != nil {
				return yxerrors.Usage("form choose value is not valid JSON", err.Error())
			}
			selected, candidates, err := selectPublishFormCandidate(rawValue, index, id)
			if err != nil {
				return err
			}
			session, err := readPublishFormSession(args[0])
			if err != nil {
				return err
			}
			accountIndex, resolvedTarget, err := resolvePublishFormAccountTarget(session.Payload, accountID, target)
			if err != nil {
				return err
			}
			resolvedPath := strings.TrimSpace(path)
			if resolvedPath == "" {
				resolvedPath, err = resolvePublishFormFieldPath(session, args[1], accountIndex)
				if err != nil {
					return err
				}
			}
			updated := cloneJSONMap(session.Payload)
			if err := setJSONPath(updated, resolvedPath, selected); err != nil {
				return yxerrors.Usage("form choose target path cannot be updated", map[string]interface{}{"path": resolvedPath, "cause": err.Error()}).
					WithHint("请确认字段存在于 form contract；多账号场景请使用 --account-id 或 --target 精确定位。").
					WithNextCommand("yxer publish form inspect " + args[0])
			}
			session.Payload = updated
			session.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			source := publishFormSource{
				Path:          resolvedPath,
				Field:         args[1],
				Kind:          "query",
				SourceCommand: sourceCommand,
				Target:        resolvedTarget,
				AccountIndex:  accountIndex,
				Value:         selected,
				UpdatedAt:     session.UpdatedAt,
			}
			session.Sources = appendPublishFormSource(session.Sources, source)
			if dryRun {
				return output.Success(cmd.OutOrStdout(), "publish.form.choose.dry-run", map[string]interface{}{
					"file": filepath.ToSlash(args[0]), "field": args[1], "path": resolvedPath, "selected": selected,
					"candidateCount": len(candidates), "source": source, "payload": updated,
				})
			}
			if err := writePublishFormSession(args[0], session); err != nil {
				return yxerrors.Internal("failed to write publish form session", err.Error())
			}
			return output.Success(cmd.OutOrStdout(), "publish.form.choose", map[string]interface{}{
				"file": filepath.ToSlash(args[0]), "field": args[1], "path": resolvedPath, "selected": selected,
				"candidateCount": len(candidates), "source": source, "payload": updated,
			})
		},
	}
	cmd.Flags().StringVar(&value, "value", "", "query result JSON, selected object JSON, or candidate array JSON")
	cmd.Flags().StringVar(&valueFile, "value-file", "", "file containing query result JSON")
	cmd.Flags().IntVar(&index, "index", -1, "candidate index to select when query result has multiple items")
	cmd.Flags().StringVar(&id, "id", "", "candidate id to select (matches yixiaoerId/id/value/key)")
	cmd.Flags().StringVar(&path, "path", "", "override target JSON path; defaults to contract dynamic field path")
	cmd.Flags().StringVar(&accountID, "account-id", "", "target platformAccountId/account_id when session has multiple accountForms")
	cmd.Flags().StringVar(&target, "target", "", "target selector, usually <platform>:<accountId> or <accountId>")
	cmd.Flags().StringVar(&sourceCommand, "source-command", "", "query command that produced this value")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview the selected write without updating the session")
	return cmd
}

func newPublishFormReviewCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "review <session.json>",
		Short: "生成会话摘要、payload hash 和后续 validate/dry-run 命令",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := readPublishFormSession(args[0])
			if err != nil {
				return err
			}
			hash, err := hashJSONValue(session.Payload)
			if err != nil {
				return yxerrors.Internal("failed to hash publish form payload", err.Error())
			}
			review := publishFormReview{PayloadHash: hash, ReviewedAt: time.Now().UTC().Format(time.RFC3339)}
			if !dryRun {
				session.Reviews = append(session.Reviews, review)
				session.UpdatedAt = review.ReviewedAt
				if err := writePublishFormSession(args[0], session); err != nil {
					return yxerrors.Internal("failed to write publish form session", err.Error())
				}
			}
			return output.Success(cmd.OutOrStdout(), actionForReview(dryRun), map[string]interface{}{
				"file":        filepath.ToSlash(args[0]),
				"planId":      session.PlanID,
				"platform":    session.Platform,
				"type":        session.Type,
				"payloadHash": hash,
				"sourceCount": len(session.Sources),
				"sources":     session.Sources,
				"next": []string{
					fmt.Sprintf("yxer publish form export %s --output payload.json", filepath.ToSlash(args[0])),
					fmt.Sprintf("yxer validate %s %s payload.json", session.Platform, session.Type),
					fmt.Sprintf("yxer publish %s %s payload.json --dry-run", session.Type, session.Platform),
					fmt.Sprintf("yxer publish %s %s payload.json", session.Type, session.Platform),
				},
			})
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview review without appending it to the session")
	return cmd
}

func newPublishFormExportCmd() *cobra.Command {
	var outputPath string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "export <session.json>",
		Short: "导出当前会话为标准 payload.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := readPublishFormSession(args[0])
			if err != nil {
				return err
			}
			target, err := absoluteFormPath(outputPath, "payload.json")
			if err != nil {
				return err
			}
			if dryRun {
				hash, _ := hashJSONValue(session.Payload)
				return output.Success(cmd.OutOrStdout(), "publish.form.export.dry-run", map[string]interface{}{"file": filepath.ToSlash(target), "payloadHash": hash, "sourceCount": len(session.Sources), "payload": session.Payload, "next": []string{
					fmt.Sprintf("yxer validate %s %s %s", session.Platform, session.Type, filepath.ToSlash(target)),
					fmt.Sprintf("yxer publish %s %s %s --dry-run", session.Type, session.Platform, filepath.ToSlash(target)),
					fmt.Sprintf("yxer publish %s %s %s", session.Type, session.Platform, filepath.ToSlash(target)),
				}})
			}
			raw, err := json.MarshalIndent(session.Payload, "", "  ")
			if err != nil {
				return yxerrors.Internal("failed to encode payload", err.Error())
			}
			if err := os.WriteFile(target, append(raw, '\n'), 0o644); err != nil {
				return yxerrors.Internal("failed to write payload", err.Error())
			}
			hash, _ := hashJSONValue(session.Payload)
			return output.Success(cmd.OutOrStdout(), "publish.form.export", map[string]interface{}{"file": filepath.ToSlash(target), "payloadHash": hash, "sourceCount": len(session.Sources), "payload": session.Payload, "next": []string{
				fmt.Sprintf("yxer validate %s %s %s", session.Platform, session.Type, filepath.ToSlash(target)),
				fmt.Sprintf("yxer publish %s %s %s --dry-run", session.Type, session.Platform, filepath.ToSlash(target)),
				fmt.Sprintf("yxer publish %s %s %s", session.Type, session.Platform, filepath.ToSlash(target)),
			}})
		},
	}
	cmd.Flags().StringVarP(&outputPath, "output", "o", "payload.json", "payload file")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview export without writing payload")
	return cmd
}

func newPublishFormSession(doc schema.Document) publishFormSession {
	now := time.Now().UTC().Format(time.RFC3339)
	return publishFormSession{Kind: "yxer.publish-form", Version: 1, PlanID: newPublishFormPlanID(doc, now), Platform: doc.Platform, Type: doc.Type, CreatedAt: now, UpdatedAt: now, Contract: buildPublishFormContract(doc), Payload: buildPayloadTemplate(doc)}
}

func newPublishFormPlanID(doc schema.Document, now string) string {
	sum := sha256.Sum256([]byte(doc.Platform + "\x00" + doc.Type + "\x00" + now))
	return "plan_" + hex.EncodeToString(sum[:])[:12]
}

func absoluteFormPath(raw, fallback string) (string, error) {
	path := strings.TrimSpace(raw)
	if path == "" {
		path = fallback
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", yxerrors.Usage("invalid form file path", err.Error())
	}
	return abs, nil
}

func readPublishFormSession(path string) (publishFormSession, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return publishFormSession{}, yxerrors.Usage("failed to read publish form session", err.Error()).WithHint("请先执行 yxer publish form start <platform> <type>。")
	}
	var session publishFormSession
	if err := json.Unmarshal(raw, &session); err != nil {
		return publishFormSession{}, yxerrors.Usage("publish form session is not valid JSON", err.Error())
	}
	if session.Kind != "yxer.publish-form" || session.Version != 1 || session.Payload == nil || session.Platform == "" || session.Type == "" {
		return publishFormSession{}, yxerrors.Usage("invalid publish form session", map[string]interface{}{"file": path}).WithHint("请使用 yxer publish form start 创建会话文件。")
	}
	return session, nil
}

func writePublishFormSession(path string, session publishFormSession) error {
	raw, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0o644)
}

func readFormValue(raw, file string) (interface{}, error) {
	if strings.TrimSpace(file) != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		raw = string(data)
	}
	var value interface{}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		// PowerShell commonly strips the quote characters around a simple
		// string argument. Preserve the page-like text input ergonomics while
		// still rejecting malformed structured JSON.
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "\"") {
			return nil, err
		}
		return raw, nil
	}
	if envelope, ok := value.(map[string]interface{}); ok && envelope["ok"] == true {
		if data, exists := envelope["data"]; exists {
			return data, nil
		}
	}
	return value, nil
}

func selectPublishFormCandidate(value interface{}, index int, id string) (interface{}, []interface{}, error) {
	candidates := publishFormCandidates(value)
	id = strings.TrimSpace(id)
	if len(candidates) == 0 {
		return value, nil, nil
	}
	if id != "" {
		for _, candidate := range candidates {
			if candidateMatchesID(candidate, id) {
				return candidate, candidates, nil
			}
		}
		return nil, candidates, yxerrors.Usage("form choose candidate id not found", map[string]interface{}{"id": id, "candidateCount": len(candidates)}).
			WithHint("请使用 query 输出中的 yixiaoerId/id/value/key，或改用 --index 指定候选。")
	}
	if index >= 0 {
		if index >= len(candidates) {
			return nil, candidates, yxerrors.Usage("form choose candidate index is out of range", map[string]interface{}{"index": index, "candidateCount": len(candidates)})
		}
		return candidates[index], candidates, nil
	}
	if len(candidates) == 1 {
		return candidates[0], candidates, nil
	}
	return nil, candidates, yxerrors.Usage("form choose has multiple candidates", map[string]interface{}{"candidateCount": len(candidates), "candidates": previewCandidates(candidates)}).
		WithHint("请使用 --id 或 --index 明确选择一个候选，避免把错误动态对象写入发布会话。")
}

func publishFormCandidates(value interface{}) []interface{} {
	switch typed := value.(type) {
	case []interface{}:
		return typed
	case map[string]interface{}:
		for _, key := range []string{"items", "list", "records", "results", "data"} {
			if items, ok := typed[key].([]interface{}); ok {
				return items
			}
		}
		if nested, ok := typed["data"].(map[string]interface{}); ok {
			return publishFormCandidates(nested)
		}
	}
	return nil
}

func candidateMatchesID(candidate interface{}, id string) bool {
	obj, ok := candidate.(map[string]interface{})
	if !ok {
		return strings.TrimSpace(fmt.Sprint(candidate)) == id
	}
	for _, key := range []string{"yixiaoerId", "id", "value", "key"} {
		if strings.TrimSpace(fmt.Sprint(obj[key])) == id {
			return true
		}
	}
	if raw, ok := obj["raw"].(map[string]interface{}); ok {
		for _, key := range []string{"id", "value", "key"} {
			if strings.TrimSpace(fmt.Sprint(raw[key])) == id {
				return true
			}
		}
	}
	return false
}

func previewCandidates(candidates []interface{}) []interface{} {
	limit := len(candidates)
	if limit > 10 {
		limit = 10
	}
	out := make([]interface{}, 0, limit)
	for i := 0; i < limit; i++ {
		obj, ok := candidates[i].(map[string]interface{})
		if !ok {
			out = append(out, map[string]interface{}{"index": i, "value": candidates[i]})
			continue
		}
		out = append(out, map[string]interface{}{
			"index":        i,
			"yixiaoerId":   firstStringField(obj, "yixiaoerId", "id", "value", "key"),
			"yixiaoerName": firstStringField(obj, "yixiaoerName", "name", "label", "text", "title"),
		})
	}
	return out
}

func firstStringField(obj map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		value := strings.TrimSpace(fmt.Sprint(obj[key]))
		if value != "" && value != "<nil>" {
			return value
		}
	}
	return ""
}

func resolvePublishFormAccountTarget(payload map[string]interface{}, accountID, target string) (int, string, error) {
	accountForms, err := publishFormAccountForms(payload)
	if err != nil {
		return 0, "", err
	}
	selector := strings.TrimSpace(accountID)
	rawTarget := strings.TrimSpace(target)
	if rawTarget != "" {
		parts := strings.Split(rawTarget, ":")
		selector = strings.TrimSpace(parts[len(parts)-1])
	}
	if selector != "" {
		for i, form := range accountForms {
			if publishFormAccountMatches(form, selector) {
				return i, firstNonEmptyAccountFormID(form), nil
			}
		}
		return 0, "", yxerrors.Usage("target accountForm not found", map[string]interface{}{"target": rawTargetOrAccount(rawTarget, selector)}).
			WithHint("请先设置 publishArgs.accountForms[].platformAccountId，或使用 inspect 查看当前会话里的账号表单。")
	}
	if len(accountForms) == 1 {
		return 0, firstNonEmptyAccountFormID(accountForms[0]), nil
	}
	return 0, "", yxerrors.Usage("form choose target is ambiguous", map[string]interface{}{"accountForms": len(accountForms)}).
		WithHint("当前会话包含多个 accountForms，请使用 --account-id 或 --target 精确定位写入目标。")
}

func publishFormAccountForms(payload map[string]interface{}) ([]map[string]interface{}, error) {
	publishArgs, _ := payload["publishArgs"].(map[string]interface{})
	if publishArgs == nil {
		return nil, yxerrors.Usage("publish form payload missing publishArgs", nil)
	}
	rawForms, _ := publishArgs["accountForms"].([]interface{})
	if len(rawForms) == 0 {
		return nil, yxerrors.Usage("publish form payload has no accountForms", nil)
	}
	forms := make([]map[string]interface{}, 0, len(rawForms))
	for _, raw := range rawForms {
		form, ok := raw.(map[string]interface{})
		if !ok {
			return nil, yxerrors.Usage("publish form accountForms item is not an object", nil)
		}
		forms = append(forms, form)
	}
	return forms, nil
}

func publishFormAccountMatches(form map[string]interface{}, selector string) bool {
	for _, key := range []string{"platformAccountId", "account_id", "accountId"} {
		if strings.TrimSpace(fmt.Sprint(form[key])) == selector {
			return true
		}
	}
	return false
}

func firstNonEmptyAccountFormID(form map[string]interface{}) string {
	for _, key := range []string{"platformAccountId", "account_id", "accountId"} {
		value := strings.TrimSpace(fmt.Sprint(form[key]))
		if value != "" && value != "<nil>" {
			return value
		}
	}
	return ""
}

func rawTargetOrAccount(target, accountID string) string {
	if strings.TrimSpace(target) != "" {
		return target
	}
	return accountID
}

func resolvePublishFormFieldPath(session publishFormSession, field string, accountIndex int) (string, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return "", yxerrors.Usage("form choose field is required", nil)
	}
	if path := dynamicFieldPath(session.Contract, field); path != "" {
		return expandAccountFormPath(path, accountIndex), nil
	}
	return fmt.Sprintf("publishArgs.accountForms[%d].contentPublishForm.%s", accountIndex, field), nil
}

func dynamicFieldPath(contract map[string]interface{}, field string) string {
	rawExamples, ok := contract["dynamicFieldExamples"]
	if !ok {
		return ""
	}
	switch examples := rawExamples.(type) {
	case map[string]dynamicFieldExample:
		if example, ok := examples[field]; ok {
			return example.Path
		}
	case map[string]interface{}:
		if raw, ok := examples[field].(map[string]interface{}); ok {
			return strings.TrimSpace(fmt.Sprint(raw["path"]))
		}
	}
	return ""
}

func expandAccountFormPath(path string, index int) string {
	replacement := fmt.Sprintf("accountForms[%d]", index)
	path = strings.Replace(path, "accountForms[]", replacement, 1)
	path = strings.Replace(path, "accountForms[].", replacement+".", 1)
	return path
}

func appendPublishFormSource(sources []publishFormSource, source publishFormSource) []publishFormSource {
	if source.UpdatedAt == "" {
		source.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	return append(sources, source)
}

func lastPublishFormSource(sources []publishFormSource) publishFormSource {
	if len(sources) == 0 {
		return publishFormSource{}
	}
	return sources[len(sources)-1]
}

func hashJSONValue(value interface{}) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func actionForReview(dryRun bool) string {
	if dryRun {
		return "publish.form.review.dry-run"
	}
	return "publish.form.review"
}

func cloneJSONMap(input map[string]interface{}) map[string]interface{} {
	raw, _ := json.Marshal(input)
	var output map[string]interface{}
	_ = json.Unmarshal(raw, &output)
	return output
}

func setJSONPath(root map[string]interface{}, path string, value interface{}) error {
	tokens, err := parseJSONPath(path)
	if err != nil || len(tokens) == 0 {
		if err == nil {
			err = fmt.Errorf("path is empty")
		}
		return err
	}
	var current interface{} = root
	for i, token := range tokens {
		last := i == len(tokens)-1
		switch typed := current.(type) {
		case map[string]interface{}:
			if last {
				typed[token] = value
				return nil
			}
			next, ok := typed[token]
			if !ok {
				if isArrayPathToken(tokens[i+1]) {
					next = []interface{}{}
				} else {
					next = map[string]interface{}{}
				}
				typed[token] = next
			}
			current = next
		case []interface{}:
			index, convErr := strconv.Atoi(token)
			if convErr != nil || index < 0 || index > len(typed) {
				return fmt.Errorf("array index %q is out of range", token)
			}
			if index == len(typed) {
				typed = append(typed, nil)
				if err := replaceJSONArray(root, tokens[:i], typed); err != nil {
					return err
				}
			}
			if last {
				typed[index] = value
				return nil
			}
			if typed[index] == nil {
				if isArrayPathToken(tokens[i+1]) {
					typed[index] = []interface{}{}
				} else {
					typed[index] = map[string]interface{}{}
				}
			}
			current = typed[index]
		default:
			return fmt.Errorf("path segment %q traverses a scalar", token)
		}
	}
	return nil
}

func isArrayPathToken(token string) bool {
	_, err := strconv.Atoi(token)
	return err == nil
}

// replaceJSONArray writes an expanded slice back through the already-traversed
// path. It keeps setJSONPath generic while avoiding pointers to interface values.
func replaceJSONArray(root map[string]interface{}, traversed []string, value []interface{}) error {
	if len(traversed) == 0 {
		return fmt.Errorf("cannot replace root array")
	}
	var current interface{} = root
	for i, token := range traversed {
		last := i == len(traversed)-1
		switch typed := current.(type) {
		case map[string]interface{}:
			next, ok := typed[token]
			if !ok {
				return fmt.Errorf("object key %q does not exist", token)
			}
			if last {
				typed[token] = value
				return nil
			}
			current = next
		case []interface{}:
			index, err := strconv.Atoi(token)
			if err != nil || index < 0 || index >= len(typed) {
				return fmt.Errorf("array index %q is out of range", token)
			}
			if last {
				typed[index] = value
				return nil
			}
			current = typed[index]
		default:
			return fmt.Errorf("path segment %q traverses a scalar", token)
		}
	}
	return nil
}

func parseJSONPath(path string) ([]string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}
	var tokens []string
	for i := 0; i < len(path); {
		if path[i] == '.' || path[i] == '[' || path[i] == ']' {
			return nil, fmt.Errorf("invalid path near %q", path[i:])
		}
		start := i
		for i < len(path) && path[i] != '.' && path[i] != '[' {
			i++
		}
		if start == i {
			return nil, fmt.Errorf("empty path segment")
		}
		tokens = append(tokens, path[start:i])
		for i < len(path) && path[i] == '[' {
			i++
			start = i
			for i < len(path) && path[i] != ']' {
				i++
			}
			if i >= len(path) || start == i {
				return nil, fmt.Errorf("invalid array segment")
			}
			tokens = append(tokens, path[start:i])
			i++
		}
		if i < len(path) {
			i++
			if i == len(path) {
				return nil, fmt.Errorf("path ends with a separator")
			}
		}
	}
	return tokens, nil
}
