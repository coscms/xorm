package main

import (
	"text/template"

	"github.com/coscms/xorm/core"
)

func init() {
	GoLangTmpl.Funcs = template.FuncMap{"Mapper": mapper.Table2Obj,
		"Type":    typestring,
		"Tag":     tagForCoscms,
		"UnTitle": unTitle,
		"gt":      gt,
		"getCol":  getCol,
	}
	langTmpls["go"] = GoLangTmpl
}

func tagForCoscms(table *core.Table, col *core.Column) string {
	switch mapper.Table2Obj(col.Name) {
	case "Created":
		col.IsCreated = true
	case "Updated":
		col.IsUpdated = true
	}
	return tag(table, col)
}
