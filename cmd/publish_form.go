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
	"github.com/yixiaoer/yixiaoer-skill/internal/api"
	"github.com/yixiaoer/yixiaoer-skill/internal/app"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
	accountsflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/accounts"
	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type publishFormSession struct {
	Kind         string                 `json:"kind"`
	Version      int                    `json:"version"`
	PlanID       string                 `json:"planId,omitempty"`
	Platform     string                 `json:"platform"`
	Type         string                 `json:"type"`
	CreatedAt    string                 `json:"createdAt"`
	UpdatedAt    string                 `json:"updatedAt"`
	ContractHash string                 `json:"contractHash,omitempty"`
	Contract     map[string]interface{} `json:"contract"`
	Payload      map[string]interface{} `json:"payload"`
	Sources      []publishFormSource    `json:"sources,omitempty"`
	Reviews      []publishFormReview    `json:"reviews,omitempty"`
}

type publishFormSource struct {
	Path          string      `json:"path"`
	Field         string      `json:"field,omitempty"`
	Kind          string      `json:"kind"`
	SourceCommand string      `json:"sourceCommand,omitempty"`
	Target        string      `json:"target,omitempty"`
	AccountIndex  int         `json:"accountIndex,omitempty"`
	Value         interface{} `json:"value,omitempty"`
	ValueHash     string      `json:"valueHash,omitempty"`
	RawHash       string      `json:"rawHash,omitempty"`
	FetchedAt     string      `json:"fetchedAt,omitempty"`
	UpdatedAt     string      `json:"updatedAt"`
}

type publishFormReview struct {
	PayloadHash string `json:"payloadHash"`
	ReviewedAt  string `json:"reviewedAt"`
}

func newPublishFormCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "form", Short: "按页面步骤维护可恢复的发布表单会话", Args: cobra.NoArgs}
	cmd.AddCommand(newPublishFormStartCmd())
	cmd.AddCommand(newPublishFormAccountCmd())
	cmd.AddCommand(newPublishFormInspectCmd())
	cmd.AddCommand(newPublishFormSetCmd())
	cmd.AddCommand(newPublishFormChooseCmd())
	cmd.AddCommand(newPublishFormVerifyCmd())
	cmd.AddCommand(newPublishFormReviewCmd())
	cmd.AddCommand(newPublishFormExportCmd())
	return cmd
}

func newPublishFormAccountCmd() *cobra.Command {
	var id string
	var index int
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "account <session.json>",
		Short: "查询在线账号并选择发布目标账号",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := readPublishFormSession(args[0])
			if err != nil {
				return err
			}
			rt, err := app.Load()
			if err != nil {
				return err
			}
			accounts, err := accountsflow.NewService(rt).ListWithOptions(session.Platform, "", 1, accountsflow.ListOptions{Size: 1000, All: true})
			if err != nil {
				return err
			}
			if len(accounts) == 0 {
				return yxerrors.Usage("no online account is available", map[string]interface{}{"platform": session.Platform}).
					WithHint("请先登录或恢复一个 status=1 的账号，再继续填写视频号发布资料。").
					WithNextCommand(fmt.Sprintf("yxer accounts list %s --status 1 --json", session.Platform))
			}
			selected, err := selectPublishFormAccount(accounts, id, index)
			if err != nil {
				return err
			}
			forms, err := publishFormAccountForms(session.Payload)
			if err != nil {
				return err
			}
			if len(forms) != 1 {
				return yxerrors.Usage("account selection requires a single account form", map[string]interface{}{"accountForms": len(forms)}).
					WithHint("当前会话包含多个 accountForms，请先拆分会话后再选择账号。")
			}
			updated := cloneJSONMap(session.Payload)
			updatedForms, _ := publishFormAccountForms(updated)
			updatedForms[0]["platformAccountId"] = api.AccountID(selected)
			now := time.Now().UTC().Format(time.RFC3339)
			source := publishFormSource{
				Path:          "publishArgs.accountForms[0].platformAccountId",
				Field:         "platformAccountId",
				Kind:          "account",
				SourceCommand: fmt.Sprintf("yxer accounts list %s --status 1 --json", session.Platform),
				Target:        api.AccountID(selected),
				Value:         selected,
				UpdatedAt:     now,
			}
			session.Payload = updated
			session.UpdatedAt = now
			session.Sources = replacePublishFormAccountSource(session.Sources, source)
			if dryRun {
				return output.Success(cmd.OutOrStdout(), "publish.form.account.dry-run", map[string]interface{}{"file": filepath.ToSlash(args[0]), "selected": selected, "source": source, "payload": updated})
			}
			if err := writePublishFormSession(args[0], session); err != nil {
				return yxerrors.Internal("failed to write publish form session", err.Error())
			}
			return output.Success(cmd.OutOrStdout(), "publish.form.account", map[string]interface{}{"file": filepath.ToSlash(args[0]), "selected": selected, "source": source, "payload": updated})
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "online platform account ID")
	cmd.Flags().IntVar(&index, "index", -1, "online account index")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview account selection without updating the session")
	return cmd
}

func selectPublishFormAccount(accounts []map[string]interface{}, id string, index int) (map[string]interface{}, error) {
	id = strings.TrimSpace(id)
	if id != "" {
		for _, account := range accounts {
			if api.AccountID(account) == id {
				return account, nil
			}
		}
		return nil, yxerrors.Usage("online account id was not found", map[string]interface{}{"id": id, "candidateCount": len(accounts)}).
			WithHint("请使用 accounts list 返回的 status=1 账号 ID。")
	}
	if index >= 0 {
		if index >= len(accounts) {
			return nil, yxerrors.Usage("online account index is out of range", map[string]interface{}{"index": index, "candidateCount": len(accounts)})
		}
		return accounts[index], nil
	}
	if len(accounts) == 1 {
		return accounts[0], nil
	}
	return nil, yxerrors.Usage("form account selection has multiple candidates", map[string]interface{}{"candidateCount": len(accounts), "candidates": previewPublishFormAccounts(accounts)}).
		WithHint("请让用户选择一个在线账号后重试，并传入 --id 或 --index。")
}

func previewPublishFormAccounts(accounts []map[string]interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(accounts))
	for i, account := range accounts {
		out = append(out, map[string]interface{}{"index": i, "platformAccountId": api.AccountID(account), "platformAccountName": accountsflow.AccountName(account), "status": api.AccountStatus(account)})
	}
	return out
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
			if err := validatePublishFormWritePath(session, args[1]); err != nil {
				return err
			}
			if isCanonicalDramaPublishFormPath(args[1]) {
				return yxerrors.Usage("form set cannot write drama field", map[string]interface{}{"path": args[1]}).
					WithHint("剧集字段必须使用 form choose 从 yxer query drama-tasks 结果选择，不能手工拼接或直接 set。")
			}
			if err := requireShipinhaoVideoAccountSelection(session, args[1]); err != nil {
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
			if strings.TrimSpace(sourceCommand) == "" {
				return yxerrors.Usage("form choose requires --source-command", nil).
					WithHint("动态字段候选必须记录产生该候选的 yxer query 命令，便于 review/export 前校验来源。").
					WithNextCommand("yxer publish form choose <session.json> <field> --value-file <query.json> --id <candidate_id> --source-command \"yxer query ... --json\"")
			}
			rawValue, err := readFormValue(value, valueFile)
			if err != nil {
				return yxerrors.Usage("form choose value is not valid JSON", err.Error())
			}
			selected, candidates, err := selectPublishFormCandidateForField(rawValue, args[1], index, id)
			if err != nil {
				return err
			}
			session, err := readPublishFormSession(args[0])
			if err != nil {
				return err
			}
			if err := requireShipinhaoVideoAccountSelection(session, args[1]); err != nil {
				return err
			}
			accountIndex, resolvedTarget, err := resolvePublishFormAccountTarget(session.Payload, accountID, target)
			if err != nil {
				return err
			}
			if err := validatePublishFormSourceAccount(sourceCommand, resolvedTarget); err != nil {
				return err
			}
			resolvedPath := strings.TrimSpace(path)
			if resolvedPath == "" {
				resolvedPath, err = resolvePublishFormFieldPath(session, args[1], accountIndex)
				if err != nil {
					return err
				}
			}
			if strings.EqualFold(strings.TrimSpace(args[1]), "drama") || isCanonicalDramaPublishFormPath(resolvedPath) {
				expectedPath, pathErr := resolvePublishFormFieldPath(session, args[1], accountIndex)
				if pathErr != nil {
					return pathErr
				}
				if strings.TrimSpace(resolvedPath) != strings.TrimSpace(expectedPath) || !isCanonicalDramaPublishFormPath(resolvedPath) {
					return yxerrors.Usage("drama field path must match contract", map[string]interface{}{"path": resolvedPath, "expected": expectedPath}).
						WithHint("剧集只能写入 contract.dynamicFieldExamples.drama 声明的 publishArgs.accountForms[].contentPublishForm.drama 路径。")
				}
				if err := validateDramaSourceCommand(sourceCommand); err != nil {
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

func newPublishFormVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify <session.json>",
		Short: "校验表单会话来源记录和当前 payload 是否一致",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session, err := readPublishFormSession(args[0])
			if err != nil {
				return err
			}
			if err := requireShipinhaoVideoAccountSelection(session, "verify"); err != nil {
				return err
			}
			report, err := validatePublishFormProvenance(session)
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), "publish.form.verify", map[string]interface{}{
				"file":       filepath.ToSlash(args[0]),
				"planId":     session.PlanID,
				"platform":   session.Platform,
				"type":       session.Type,
				"provenance": report,
			})
		},
	}
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
			if err := requireShipinhaoVideoAccountSelection(session, "review"); err != nil {
				return err
			}
			provenance, err := validatePublishFormProvenance(session)
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
			channel, clientID := publishModeFromPayload(session.Payload)
			return output.Success(cmd.OutOrStdout(), actionForReview(dryRun), map[string]interface{}{
				"file":        filepath.ToSlash(args[0]),
				"planId":      session.PlanID,
				"platform":    session.Platform,
				"type":        session.Type,
				"payloadHash": hash,
				"sourceCount": len(session.Sources),
				"sources":     session.Sources,
				"provenance":  provenance,
				"next": []string{
					fmt.Sprintf("yxer publish form export %s --output payload.json", filepath.ToSlash(args[0])),
					validateNextCommand(session.Platform, session.Type, "payload.json", channel, clientID),
					publishNextCommand(session.Type, session.Platform, "payload.json", channel, clientID, true),
					publishNextCommand(session.Type, session.Platform, "payload.json", channel, clientID, false),
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
			if err := requireShipinhaoVideoAccountSelection(session, "export"); err != nil {
				return err
			}
			provenance, err := validatePublishFormProvenance(session)
			if err != nil {
				return err
			}
			target, err := absoluteFormPath(outputPath, "payload.json")
			if err != nil {
				return err
			}
			channel, clientID := publishModeFromPayload(session.Payload)
			if dryRun {
				hash, _ := hashJSONValue(session.Payload)
				return output.Success(cmd.OutOrStdout(), "publish.form.export.dry-run", map[string]interface{}{"file": filepath.ToSlash(target), "payloadHash": hash, "sourceCount": len(session.Sources), "provenance": provenance, "payload": session.Payload, "next": []string{
					validateNextCommand(session.Platform, session.Type, filepath.ToSlash(target), channel, clientID),
					publishNextCommand(session.Type, session.Platform, filepath.ToSlash(target), channel, clientID, true),
					publishNextCommand(session.Type, session.Platform, filepath.ToSlash(target), channel, clientID, false),
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
			return output.Success(cmd.OutOrStdout(), "publish.form.export", map[string]interface{}{"file": filepath.ToSlash(target), "payloadHash": hash, "sourceCount": len(session.Sources), "provenance": provenance, "payload": session.Payload, "next": []string{
				validateNextCommand(session.Platform, session.Type, filepath.ToSlash(target), channel, clientID),
				publishNextCommand(session.Type, session.Platform, filepath.ToSlash(target), channel, clientID, true),
				publishNextCommand(session.Type, session.Platform, filepath.ToSlash(target), channel, clientID, false),
			}})
		},
	}
	cmd.Flags().StringVarP(&outputPath, "output", "o", "payload.json", "payload file")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview export without writing payload")
	return cmd
}

func newPublishFormSession(doc schema.Document) publishFormSession {
	now := time.Now().UTC().Format(time.RFC3339)
	contract := buildPublishFormContract(doc)
	hash, _ := hashJSONValue(contract)
	return publishFormSession{Kind: "yxer.publish-form", Version: 1, PlanID: newPublishFormPlanID(doc, now), Platform: doc.Platform, Type: doc.Type, CreatedAt: now, UpdatedAt: now, ContractHash: hash, Contract: contract, Payload: buildPayloadTemplate(doc)}
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
	if err := validatePublishFormContractFresh(session); err != nil {
		return publishFormSession{}, err
	}
	return session, nil
}

func validatePublishFormContractFresh(session publishFormSession) error {
	if strings.TrimSpace(session.ContractHash) == "" {
		return nil
	}
	rt, err := app.Load()
	if err != nil {
		return err
	}
	doc, err := schema.NewValidator(rt.Config.SchemaDir).Schema(session.Platform, session.Type)
	if err != nil {
		return yxerrors.Usage("publish form contract cannot be refreshed", map[string]interface{}{"platform": session.Platform, "type": session.Type}).
			WithHint("当前会话引用的平台或类型已无法解析，请重新执行 publish form start 创建新会话。").
			WithNextCommand("yxer publish form start " + session.Platform + " " + session.Type)
	}
	currentHash, _ := hashJSONValue(buildPublishFormContract(doc))
	if currentHash != session.ContractHash {
		return yxerrors.Usage("publish form contract is stale", map[string]interface{}{"expected": currentHash, "actual": session.ContractHash, "platform": session.Platform, "type": session.Type}).
			WithHint("schema 或表单契约已变化，请重新执行 publish form start，避免旧字段路径继续写入。").
			WithNextCommand("yxer publish form start " + session.Platform + " " + session.Type)
	}
	return nil
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
	return selectPublishFormCandidateFromCandidates(value, candidates, index, id)
}

func selectPublishFormCandidateForField(value interface{}, field string, index int, id string) (interface{}, []interface{}, error) {
	if strings.EqualFold(strings.TrimSpace(field), "drama") {
		candidates, err := dramaPublishFormCandidates(value)
		if err != nil {
			return nil, nil, err
		}
		return selectPublishFormCandidateFromCandidates(value, candidates, index, id)
	}
	return selectPublishFormCandidate(value, index, id)
}

func selectPublishFormCandidateFromCandidates(original interface{}, candidates []interface{}, index int, id string) (interface{}, []interface{}, error) {
	id = strings.TrimSpace(id)
	if len(candidates) == 0 {
		return original, nil, nil
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

func dramaPublishFormCandidates(value interface{}) ([]interface{}, error) {
	return dramaPublishFormCandidatesAtDepth(value, 0)
}

func dramaPublishFormCandidatesAtDepth(value interface{}, depth int) ([]interface{}, error) {
	switch typed := value.(type) {
	case []interface{}:
		if len(typed) == 0 {
			return nil, dramaCandidateError("form choose has no candidates", map[string]interface{}{"reason": "empty array"})
		}
		for _, candidate := range typed {
			if !isExactDramaCandidate(candidate) {
				return nil, dramaCandidateError("drama query result contains an incomplete or extra-field candidate", map[string]interface{}{"candidate": candidate})
			}
		}
		return typed, nil
	case map[string]interface{}:
		knownKeys := []string{"items", "list", "dataList", "records", "results"}
		present := make([]string, 0, 1)
		for _, key := range knownKeys {
			if _, exists := typed[key]; exists {
				present = append(present, key)
			}
		}
		if len(present) == 1 {
			if _, hasNestedData := typed["data"]; hasNestedData {
				return nil, dramaCandidateError("drama query result contains multiple candidate lists", map[string]interface{}{"keys": []string{present[0], "data"}})
			}
		}
		if len(present) > 1 {
			return nil, dramaCandidateError("drama query result contains multiple candidate lists", map[string]interface{}{"keys": present})
		}
		if len(present) == 1 {
			items, ok := typed[present[0]].([]interface{})
			if !ok {
				return nil, dramaCandidateError("drama query result candidate list must be an array", map[string]interface{}{"key": present[0]})
			}
			return dramaPublishFormCandidatesAtDepth(items, depth)
		}
		if isExactDramaCandidate(typed) {
			return []interface{}{typed}, nil
		}
		if nested, exists := typed["data"]; exists {
			if depth >= 1 {
				return nil, dramaCandidateError("drama query result has too many nested data envelopes", nil)
			}
			return dramaPublishFormCandidatesAtDepth(nested, depth+1)
		}
	}
	return nil, dramaCandidateError("drama query result is not a recognized candidate envelope", map[string]interface{}{"value": value})
}

func isExactDramaCandidate(value interface{}) bool {
	obj, ok := value.(map[string]interface{})
	if !ok || len(obj) != 3 {
		return false
	}
	for _, key := range []string{"yixiaoerId", "yixiaoerImageUrl", "yixiaoerName"} {
		text, ok := obj[key].(string)
		if !ok || strings.TrimSpace(text) == "" {
			return false
		}
	}
	return true
}

func dramaCandidateError(message string, details interface{}) error {
	return yxerrors.Usage(message, details).
		WithHint("请使用 yxer query drama-tasks <account_id> --json 的完整返回值，并选择只含 yixiaoerId、yixiaoerImageUrl、yixiaoerName 的候选对象。")
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
	return "", yxerrors.Usage("form choose field is not query-backed", map[string]interface{}{"field": field}).
		WithHint("choose 只能写 dynamicFieldExamples 中声明的动态字段；普通文本或枚举字段请使用 publish form set。").
		WithNextCommand("yxer publish form inspect <session.json>")
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

func validatePublishFormWritePath(session publishFormSession, path string) error {
	allowed := allowedPublishFormPaths(session.Contract)
	normalized := normalizePublishFormPath(path)
	if allowed[normalized] {
		return nil
	}
	return yxerrors.Usage("form path is not declared in contract", map[string]interface{}{"path": path}).
		WithHint("请使用 publish form inspect 返回的 contract.fields、fieldPlacements 或 dynamicFieldExamples 中声明的路径，避免拼错字段后写入无效 payload。").
		WithNextCommand("yxer publish form inspect <session.json>")
}

func allowedPublishFormPaths(contract map[string]interface{}) map[string]bool {
	allowed := map[string]bool{}
	addAllowedPath := func(path string) {
		path = normalizePublishFormPath(path)
		if path != "" {
			allowed[path] = true
		}
	}
	if fields, ok := contract["fields"].([]flatFieldView); ok {
		for _, field := range fields {
			addAllowedPath(field.Path)
		}
	}
	if fields, ok := contract["fields"].([]interface{}); ok {
		for _, item := range fields {
			if obj, ok := item.(map[string]interface{}); ok {
				addAllowedPath(strings.TrimSpace(fmt.Sprint(obj["path"])))
			}
		}
	}
	if placements, ok := contract["fieldPlacements"].(map[string]fieldPlacementView); ok {
		for _, placement := range placements {
			for _, path := range placement.InputPaths {
				addAllowedPath(path)
			}
		}
	}
	if placements, ok := contract["fieldPlacements"].(map[string]interface{}); ok {
		for _, raw := range placements {
			placement, _ := raw.(map[string]interface{})
			if placement == nil {
				continue
			}
			if inputPaths, ok := placement["inputPaths"].([]interface{}); ok {
				for _, path := range inputPaths {
					addAllowedPath(strings.TrimSpace(fmt.Sprint(path)))
				}
			}
		}
	}
	if examples, ok := contract["dynamicFieldExamples"].(map[string]dynamicFieldExample); ok {
		for _, example := range examples {
			addAllowedPath(example.Path)
		}
	}
	if examples, ok := contract["dynamicFieldExamples"].(map[string]interface{}); ok {
		for _, raw := range examples {
			if obj, ok := raw.(map[string]interface{}); ok {
				addAllowedPath(strings.TrimSpace(fmt.Sprint(obj["path"])))
			}
		}
	}
	return allowed
}

func normalizePublishFormPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	parseable := strings.ReplaceAll(path, "[]", "[0]")
	tokens, err := parseJSONPath(parseable)
	if err != nil || len(tokens) == 0 {
		return path
	}
	for i, token := range tokens {
		if _, err := strconv.Atoi(token); err == nil {
			tokens[i] = "[]"
		}
	}
	return strings.Join(tokens, ".")
}

func appendPublishFormSource(sources []publishFormSource, source publishFormSource) []publishFormSource {
	if source.UpdatedAt == "" {
		source.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if source.FetchedAt == "" {
		source.FetchedAt = source.UpdatedAt
	}
	if source.ValueHash == "" && source.Value != nil {
		source.ValueHash, _ = hashJSONValue(source.Value)
	}
	if source.RawHash == "" && !isCanonicalDramaPublishFormPath(source.Path) {
		source.RawHash, _ = hashDynamicRaw(source.Value)
	}
	return append(sources, source)
}

func validatePublishFormProvenance(session publishFormSession) (map[string]interface{}, error) {
	var errors []map[string]interface{}
	for i, source := range session.Sources {
		sourceErrors := validatePublishFormSource(session.Payload, source)
		for _, detail := range sourceErrors {
			detail["index"] = i
			errors = append(errors, detail)
		}
	}
	for _, missing := range missingDramaSourceErrors(session.Payload, session.Sources) {
		errors = append(errors, missing)
	}
	report := map[string]interface{}{
		"valid":       len(errors) == 0,
		"sourceCount": len(session.Sources),
		"errors":      errors,
	}
	if len(errors) > 0 {
		return report, yxerrors.Usage("publish form provenance validation failed", errors).
			WithHint("来源记录和当前 payload 不一致，或 query 来源信息不完整；请重新执行 publish form choose/set 修复对应字段。").
			WithNextCommand("yxer publish form inspect <session.json>")
	}
	return report, nil
}

func validatePublishFormSource(payload map[string]interface{}, source publishFormSource) []map[string]interface{} {
	var errors []map[string]interface{}
	add := func(code, message string) {
		errors = append(errors, map[string]interface{}{"code": code, "path": source.Path, "message": message})
	}
	if strings.TrimSpace(source.Path) == "" {
		add("missing_path", "source path is required")
		return errors
	}
	current, err := getJSONPathValue(payload, source.Path)
	if err != nil {
		add("path_not_found", err.Error())
		return errors
	}
	currentHash, _ := hashJSONValue(current)
	expectedValueHash := strings.TrimSpace(source.ValueHash)
	if expectedValueHash == "" && source.Value != nil {
		expectedValueHash, _ = hashJSONValue(source.Value)
	}
	if expectedValueHash != "" && currentHash != expectedValueHash {
		add("value_hash_mismatch", "current payload value no longer matches the recorded source value")
	}
	if source.Kind == "query" {
		if strings.TrimSpace(source.SourceCommand) == "" {
			add("missing_source_command", "query source must include the yxer query command that produced it")
		}
		if strings.TrimSpace(source.FetchedAt) == "" && strings.TrimSpace(source.UpdatedAt) == "" {
			add("missing_fetched_at", "query source must include fetchedAt or updatedAt")
		}
		if isCanonicalDramaPublishFormPath(source.Path) {
			if err := validateDramaSourceCommand(source.SourceCommand); err != nil {
				add("invalid_source_command", err.Error())
			}
		} else {
			expectedRawHash := strings.TrimSpace(source.RawHash)
			if expectedRawHash == "" {
				expectedRawHash, _ = hashDynamicRaw(source.Value)
			}
			currentRawHash, _ := hashDynamicRaw(current)
			if expectedRawHash == "" {
				add("missing_raw_hash", "query source must include hashable raw data")
			} else if currentRawHash != expectedRawHash {
				add("raw_hash_mismatch", "current payload raw data no longer matches the recorded query raw data")
			}
		}
	}
	return errors
}

func isCanonicalDramaPublishFormPath(path string) bool {
	return normalizePublishFormPath(path) == "publishArgs.accountForms.[].contentPublishForm.drama"
}

func validateDramaSourceCommand(command string) error {
	resource := querySourceCommandResource(command)
	if resource == "drama-tasks" {
		return nil
	}
	return yxerrors.Usage("drama source must come from yxer query drama-tasks", map[string]interface{}{"sourceCommand": command}).
		WithHint("请先执行 yxer query drama-tasks <account_id> [--query 关键词] --json，再使用 form choose 选择剧集。")
}

func querySourceCommandResource(command string) string {
	parts := strings.Fields(command)
	for i := 0; i+1 < len(parts); i++ {
		if parts[i] == "query" || strings.HasSuffix(parts[i], "query") {
			return strings.TrimSpace(parts[i+1])
		}
	}
	return ""
}

func missingDramaSourceErrors(payload map[string]interface{}, sources []publishFormSource) []map[string]interface{} {
	publishArgs, _ := payload["publishArgs"].(map[string]interface{})
	accountForms, _ := publishArgs["accountForms"].([]interface{})
	var errors []map[string]interface{}
	for i, rawForm := range accountForms {
		form, _ := rawForm.(map[string]interface{})
		cpf, _ := form["contentPublishForm"].(map[string]interface{})
		if _, exists := cpf["drama"]; !exists {
			continue
		}
		path := fmt.Sprintf("publishArgs.accountForms[%d].contentPublishForm.drama", i)
		found := false
		for _, source := range sources {
			if source.Kind == "query" && strings.TrimSpace(source.Path) == path {
				found = true
				break
			}
		}
		if !found {
			errors = append(errors, map[string]interface{}{
				"code":    "missing_drama_source",
				"path":    path,
				"message": "drama field must have a matching query source record",
			})
		}
	}
	return errors
}

func validatePublishFormSourceAccount(sourceCommand, target string) error {
	sourceAccountID := querySourceCommandAccountID(sourceCommand)
	target = strings.TrimSpace(target)
	if sourceAccountID == "" || target == "" || isTemplatePlaceholder(target) {
		return nil
	}
	if sourceAccountID != target {
		return yxerrors.Usage("form choose source account does not match target account", map[string]interface{}{"sourceAccountId": sourceAccountID, "targetAccountId": target}).
			WithHint("动态字段必须用目标账号执行 query 后再选择，不能跨账号复用前置对象。").
			WithNextCommand("yxer query <resource> " + target + " --json")
	}
	return nil
}

func requireShipinhaoVideoAccountSelection(session publishFormSession, path string) error {
	if platformutil.CanonicalKey(session.Platform) != "shipinhao" || strings.TrimSpace(session.Type) != "video" {
		return nil
	}
	// Drama selection carries its own target account in the query command and
	// remains compatible with existing query-first sessions.
	if strings.Contains(strings.ToLower(path), "drama") {
		return nil
	}
	forms, err := publishFormAccountForms(session.Payload)
	if err != nil {
		return err
	}
	if len(forms) == 1 && !isTemplatePlaceholder(firstNonEmptyAccountFormID(forms[0])) {
		return nil
	}
	return yxerrors.Usage("shipinhao video account selection is required first", map[string]interface{}{"path": path}).
		WithHint("请先执行 yxer publish form account <session.json> --id <online_account_id>，确认账号有效后再填写视频资料。")
}

func replacePublishFormAccountSource(sources []publishFormSource, source publishFormSource) []publishFormSource {
	filtered := sources[:0]
	for _, existing := range sources {
		if existing.Kind == "account" && existing.Path == source.Path {
			continue
		}
		filtered = append(filtered, existing)
	}
	return appendPublishFormSource(filtered, source)
}

func querySourceCommandAccountID(command string) string {
	parts := strings.Fields(command)
	for i := 0; i+2 < len(parts); i++ {
		if parts[i] == "query" || strings.HasSuffix(parts[i], "query") {
			candidate := strings.TrimSpace(parts[i+2])
			if candidate == "" || strings.HasPrefix(candidate, "-") || isTemplatePlaceholder(candidate) {
				return ""
			}
			return candidate
		}
	}
	return ""
}

func isTemplatePlaceholder(value string) bool {
	value = strings.TrimSpace(value)
	return strings.HasPrefix(value, "<") && strings.HasSuffix(value, ">")
}

func getJSONPathValue(root map[string]interface{}, path string) (interface{}, error) {
	tokens, err := parseJSONPath(path)
	if err != nil || len(tokens) == 0 {
		if err == nil {
			err = fmt.Errorf("path is empty")
		}
		return nil, err
	}
	var current interface{} = root
	for _, token := range tokens {
		switch typed := current.(type) {
		case map[string]interface{}:
			next, ok := typed[token]
			if !ok {
				return nil, fmt.Errorf("object key %q does not exist", token)
			}
			current = next
		case []interface{}:
			index, convErr := strconv.Atoi(token)
			if convErr != nil || index < 0 || index >= len(typed) {
				return nil, fmt.Errorf("array index %q is out of range", token)
			}
			current = typed[index]
		default:
			return nil, fmt.Errorf("path segment %q traverses a scalar", token)
		}
	}
	return current, nil
}

func hashDynamicRaw(value interface{}) (string, error) {
	rawValues := collectDynamicRawValues(value)
	if len(rawValues) == 0 {
		return "", nil
	}
	if len(rawValues) == 1 {
		return hashJSONValue(rawValues[0])
	}
	values := make([]interface{}, len(rawValues))
	copy(values, rawValues)
	return hashJSONValue(values)
}

func collectDynamicRawValues(value interface{}) []interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		if raw, ok := typed["raw"].(map[string]interface{}); ok && raw != nil {
			return []interface{}{raw}
		}
		var result []interface{}
		for _, child := range typed {
			result = append(result, collectDynamicRawValues(child)...)
		}
		return result
	case []interface{}:
		var result []interface{}
		for _, child := range typed {
			result = append(result, collectDynamicRawValues(child)...)
		}
		return result
	default:
		return nil
	}
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

func publishModeFromPayload(payload map[string]interface{}) (string, string) {
	channel := strings.TrimSpace(fmt.Sprint(payload["publishChannel"]))
	if channel == "" || channel == "<nil>" {
		channel = "cloud"
	}
	clientID := strings.TrimSpace(fmt.Sprint(payload["clientId"]))
	if clientID == "<nil>" {
		clientID = ""
	}
	return channel, clientID
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
