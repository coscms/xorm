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

type ReverseConfig struct {
	Engine    string `json:engine`
	Dsn       string `json:"dsn"`
	MultiFile bool   `json:multiFile`
	Lang      string `json:"lang"`
	PkgName   string `json:"pkgName"`
	GenJson   bool   `json:"genJson"`
	Prefix    string `json:"prefix"`
	Template  string `json:"template"`
}
