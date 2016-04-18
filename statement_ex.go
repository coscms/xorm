package xorm

import (
	"fmt"
	"reflect"

	"github.com/coscms/xorm/core"
)

// == ORDER BY ==
type orderBy []*orderByParam

type orderByParam struct {
	Field string //Column name
	Sort  string // ASC/DESC
}

// == Fields ==
type fields []string //for Omit

// == JOIN ==
type joinTables []*joinParam

func (j *joinTables) New(stmt *Statement) *joinParam {
	join := NewJoinParam(stmt)
	*j = append(*j, join)
	return join
}

func (j *joinTables) Add(join *joinParam) {
	*j = append(*j, join)
}

func (j joinTables) String() string {
	joinStr := ``
	t := ``
	for _, join := range j {
		joinStr += t + join.String()
		t = ` `
	}
	return joinStr
}

func (j joinTables) Args() []interface{} {
	args := make([]interface{}, 0)
	for _, join := range j {
		args = append(args, join.Args...)
	}
	return args
}

func NewJoinParam(stmt *Statement) *joinParam {
	return &joinParam{
		Args:      make([]interface{}, 0),
		statement: stmt,
	}
}

type joinParam struct {
	Operator string //LEFT/RIGHT/INNER...
	Table    string
	Alias    string
	ONStr    string
	Args     []interface{}

	statement *Statement
}

func (j *joinParam) String() string {
	joinStr := j.Operator + ` JOIN ` + j.statement.Engine.Quote(j.Table)
	if len(j.Alias) == 0 && j.statement.relation != nil {
		j.Alias, _ = j.statement.relation.ExAlias[j.Table]
	}
	if len(j.Alias) > 0 {
		joinStr += ` AS ` + j.statement.Engine.Quote(j.Alias)
	}
	if len(j.ONStr) > 0 {
		joinStr += ` ON ` + j.ONStr
	}
	return joinStr
}

// == Extends Statement ==

// Join The joinOP should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
func (statement *Statement) join(joinOP string, tablename interface{}, condition string, args ...interface{}) *Statement {
	join := statement.joinTables.New(statement)
	join.Operator = joinOP
	join.Table = ``
	join.Alias = ``
	join.ONStr = condition
	join.Args = args

	switch tablename.(type) {
	case []string:
		t := tablename.([]string)
		if len(t) > 1 {
			join.Table = t[0]
			join.Alias = t[1]
		} else if len(t) == 1 {
			join.Table = t[0]
		}
	case []interface{}:
		t := tablename.([]interface{})
		l := len(t)
		var table string
		if l > 0 {
			f := t[0]
			v := rValue(f)
			t := v.Type()
			if t.Kind() == reflect.String {
				table = f.(string)
			} else if t.Kind() == reflect.Struct {
				r := statement.Engine.autoMapType(v)
				table = r.Name
			} else {
				table = fmt.Sprintf("%v", f)
			}
			join.Table = table
		}
		if l > 1 {
			join.Alias = fmt.Sprintf("%v", t[1])
		}
	case string:
		join.Table = tablename.(string)
	default:
		v := rValue(tablename)
		t := v.Type()
		if t.Kind() == reflect.Struct {
			r := statement.Engine.autoMapType(v)
			join.Table = r.Name
		} else {
			join.Table = fmt.Sprintf("%v", tablename)
		}
	}
	return statement
}

func (statement *Statement) JoinStr() string {
	if statement.joinGenerated {
		return statement.joinStr
	}
	statement.joinStr = statement.joinTables.String()
	statement.joinArgs = statement.joinTables.Args()
	statement.joinGenerated = true
	return statement.joinStr
}

func (statement *Statement) SetRelation(r *core.Relation) {
	statement.relation = r
	if r != nil && len(r.Extends) > 0 {
		if !r.IsTable {
			statement.RefTable = r.Extends[0]
			if len(statement.TableAlias) == 0 {
				name := statement.RefTable.Name
				statement.TableAlias, _ = r.ExAlias[name]
			}
		} else {
			if len(statement.TableAlias) == 0 {
				statement.TableAlias = r.Table.Type.Name()
			}
		}
	}
}
