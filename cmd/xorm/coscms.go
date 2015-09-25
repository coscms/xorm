package main

import (
	"github.com/coscms/xorm/core"
)

func init() {
	GoLangTmpl.Funcs["Tag"] = tagForCoscms
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
