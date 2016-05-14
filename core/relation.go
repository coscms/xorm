//[SWH|+] 关联关系
package core

import (
	"fmt"
	"strings"
)

func NewRelation(table *Table) *Relation {
	return &Relation{
		Table:         table,
		Extends:       make([]*Table, 0),
		ExAlias:       make(map[string]string),
		RelTables:     []*RelTable{},
		nameRelTables: map[string]int{},
	}
}

func NewRelTable(joinType string, where string, tableName string) *RelTable {
	return &RelTable{
		JoinType:  joinType,
		Where:     where,
		TableName: tableName,
	}
}

type RelTable struct {
	JoinType  string
	Where     string
	TableName string
}

func (r *RelTable) Condition(args ...interface{}) string {
	return fmt.Sprintf(r.Where, args...)
}

func (r *RelTable) IsValid() bool {
	if r.JoinType == `` || r.TableName == `` {
		return false
	}
	return true
}

func (r *RelTable) String() string {
	if !r.IsValid() {
		return ``
	}
	return r.JoinType + ` JOIN ` + r.TableName + ` ON ` + r.Where
}

type Relation struct {
	*Table
	IsTable       bool
	Extends       []*Table          //join关联表。
	ExAlias       map[string]string //关联表真实表名和别名对照
	RelTables     []*RelTable
	nameRelTables map[string]int //名称对应的RelTables索引
}

func (r *Relation) Alias(rawName string) string {
	if alias, ok := r.ExAlias[rawName]; ok {
		return alias
	}
	return rawName
}

func (r *Relation) AliasGetByIndex(index int) string {
	if len(r.Extends) <= index {
		return ``
	}
	return r.Alias(r.Extends[index].Name)
}

func (r *Relation) GetRelTableByName(name string) *RelTable {
	if i, ok := r.nameRelTables[name]; ok {
		return r.RelTables[i]
	}
	return nil
}

func (r *Relation) AddExtend(t *Table, relTagStr string) {
	r.Extends = append(r.Extends, t)
	var rt *RelTable
	if relTagStr != `` {
		//example: `rel:"LEFT:A.id=B.uid"`
		//<JOIN类型>:<JOIN条件>
		relv := strings.Split(relTagStr, `:`)
		switch len(relv) {
		case 2:
			rt = NewRelTable(relv[0], relv[1], t.Name)
		case 1:
			rt = NewRelTable(`INNER`, relv[0], t.Name)
		}
	}
	r.nameRelTables[t.Name] = len(r.RelTables)
	r.RelTables = append(r.RelTables, rt)
}
