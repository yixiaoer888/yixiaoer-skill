package schema

import base "github.com/yixiaoer/yixiaoer-skill/internal/schema"

type Validator = base.Validator
type Document = base.Document
type Entry = base.Entry
type Catalog = base.Catalog
type PropertyView = base.PropertyView
type Result = base.Result

func NewValidator(schemaDir string) Validator {
	return base.NewValidator(schemaDir)
}
