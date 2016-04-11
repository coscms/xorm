//[SWH|+] 关联关系
package core

func NewRelation(table *Table) *Relation {
	return &Relation{
		Table:   table,
		Extends: make([]*Table, 0),
		ExAlias: make(map[string]string),
	}
}

type Relation struct {
	*Table
	HasColumnField bool
	Extends        []*Table          //join关联表。
	ExAlias        map[string]string //关联表真实表名和别名对照
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

func (r *Relation) AddExtend(t *Table) {
	r.Extends = append(r.Extends, t)
}
