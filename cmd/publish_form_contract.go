package cmd

import (
	"fmt"
	"sort"

	platformutil "github.com/yixiaoer/yixiaoer-skill/internal/platform"
	"github.com/yixiaoer/yixiaoer-skill/internal/schema"
)

// buildPublishFormContract describes the same state boundaries an agent needs
// to reproduce the page flow without guessing payload paths.
func buildPublishFormContract(doc schema.Document) map[string]interface{} {
	fields := flattenFieldViews(buildStandardPublishFieldView(doc, doc.Properties))
	resourceFields := []string{}
	for _, key := range accountLevelResourceKeys(doc.Type) {
		resourceFields = append(resourceFields, "publishArgs.accountForms[]."+key)
	}
	sort.Strings(resourceFields)

	dynamic := buildDynamicFieldExamples(doc)
	queries := make([]string, 0, len(dynamic))
	for _, example := range dynamic {
		if example.QueryCommand != "" {
			queries = append(queries, example.QueryCommand)
		}
	}
	sort.Strings(queries)

	steps := []interface{}{
		map[string]interface{}{"id": "account", "kind": "account-selection", "required": true, "command": publishFormAccountSelectionCommand(doc)},
		map[string]interface{}{"id": "resources", "kind": "resource-upload", "fields": resourceFields, "command": "yxer upload --file <path> --json"},
		map[string]interface{}{"id": "platform-form", "kind": "field-entry", "fields": fields},
	}
	if len(dynamic) > 0 {
		steps = append(steps, map[string]interface{}{
			"id":      "dynamic-selection",
			"kind":    "query-selection",
			"fields":  dynamic,
			"queries": queries,
			"command": "yxer publish form choose <session.json> <field> --value-file <query.json> --id <candidate_id> --source-command \"yxer query ... --json\"",
		})
	}
	steps = append(steps, map[string]interface{}{"id": "review", "kind": "validation", "commands": []string{
		"yxer publish form verify <session.json>",
		"yxer publish form export <session.json> --output payload.json",
		fmt.Sprintf("yxer validate %s %s <payload.json>", doc.Platform, doc.Type),
		fmt.Sprintf("yxer publish %s %s <payload.json> --dry-run", doc.Type, doc.Platform),
		"yxer publish form review <session.json>",
	}})

	return map[string]interface{}{
		"version":              1,
		"platform":             doc.Platform,
		"platformName":         platformutil.ChineseName(doc.Platform),
		"type":                 doc.Type,
		"steps":                steps,
		"fields":               fields,
		"businessFields":       doc.Properties,
		"fieldPlacements":      buildFieldPlacements(doc),
		"dynamicFieldExamples": dynamic,
		"template":             buildPayloadTemplate(doc),
		"sourceOfTruth":        []string{"prepare", "schema fields", "schema get", "query results", "user-provided business values", "upload results", "session.sources"},
	}
}

func publishFormAccountSelectionCommand(doc schema.Document) string {
	if platformutil.CanonicalKey(doc.Platform) == "shipinhao" && doc.Type == "video" {
		return "yxer publish form account <session.json> --id <online_account_id>"
	}
	return fmt.Sprintf("yxer accounts list %s --status 1 --json", platformutil.ChineseName(doc.Platform))
}
