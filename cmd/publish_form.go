package cmd

import (
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
	Platform  string                 `json:"platform"`
	Type      string                 `json:"type"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
	Contract  map[string]interface{} `json:"contract"`
	Payload   map[string]interface{} `json:"payload"`
}

func newPublishFormCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "form", Short: "按页面步骤维护可恢复的发布表单会话", Args: cobra.NoArgs}
	cmd.AddCommand(newPublishFormStartCmd())
	cmd.AddCommand(newPublishFormInspectCmd())
	cmd.AddCommand(newPublishFormSetCmd())
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
				"contract": session.Contract, "payload": session.Payload,
			})
		},
	}
}

func newPublishFormSetCmd() *cobra.Command {
	var value, valueFile string
	var index int
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
			if dryRun {
				return output.Success(cmd.OutOrStdout(), "publish.form.set.dry-run", map[string]interface{}{"file": filepath.ToSlash(args[0]), "path": args[1], "value": parsed, "payload": updated})
			}
			if err := writePublishFormSession(args[0], session); err != nil {
				return yxerrors.Internal("failed to write publish form session", err.Error())
			}
			return output.Success(cmd.OutOrStdout(), "publish.form.set", map[string]interface{}{"file": filepath.ToSlash(args[0]), "path": args[1], "value": parsed, "payload": updated})
		},
	}
	cmd.Flags().StringVar(&value, "value", "", "JSON value to set")
	cmd.Flags().StringVar(&valueFile, "value-file", "", "file containing one JSON value")
	cmd.Flags().IntVar(&index, "index", -1, "select an item from an array value (for query results)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview the updated session without writing it")
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
				return output.Success(cmd.OutOrStdout(), "publish.form.export.dry-run", map[string]interface{}{"file": filepath.ToSlash(target), "payload": session.Payload})
			}
			raw, err := json.MarshalIndent(session.Payload, "", "  ")
			if err != nil {
				return yxerrors.Internal("failed to encode payload", err.Error())
			}
			if err := os.WriteFile(target, append(raw, '\n'), 0o644); err != nil {
				return yxerrors.Internal("failed to write payload", err.Error())
			}
			return output.Success(cmd.OutOrStdout(), "publish.form.export", map[string]interface{}{"file": filepath.ToSlash(target), "payload": session.Payload, "next": []string{
				fmt.Sprintf("yxer validate %s %s %s", session.Platform, session.Type, filepath.ToSlash(target)),
				fmt.Sprintf("yxer publish %s %s %s --dry-run", session.Type, session.Platform, filepath.ToSlash(target)),
			}})
		},
	}
	cmd.Flags().StringVarP(&outputPath, "output", "o", "payload.json", "payload file")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview export without writing payload")
	return cmd
}

func newPublishFormSession(doc schema.Document) publishFormSession {
	now := time.Now().UTC().Format(time.RFC3339)
	return publishFormSession{Kind: "yxer.publish-form", Version: 1, Platform: doc.Platform, Type: doc.Type, CreatedAt: now, UpdatedAt: now, Contract: buildPublishFormContract(doc), Payload: buildPayloadTemplate(doc)}
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
