package xorm

import (
	"reflect"

	"github.com/coscms/xorm/core"
)


// =====================================
// 定义ResultSet
// =====================================
func NewResultSet() *ResultSet {
	return &ResultSet{
		Fields:		make([]string,0),
		Values:		make([]string,0),
		NameIndex:	make(map[string]int),
		Length:		0,
	}
}

type ResultSet struct{
	Fields		[]string
	Values		[]string
	NameIndex	map[string]int
	Length		int
}

func (r *ResultSet) Get(index int) string {
	if index >= r.Length {
		return ""
	}
	return r.Values[index]
}

func (r *ResultSet) GetByName(name string) string {
	if index,ok:=r.NameIndex[name];ok {
		return r.Get(index)
	}
	return ""
}

func (r *ResultSet) Set(index int, value string) bool {
	if index >= r.Length {
		return false
	}
	r.Values[index]=value
	return true
}

func (r *ResultSet) SetByName(name string, value string) bool {
	if index,ok:=r.NameIndex[name];ok {
		return r.Set(index, value)
	}else{
		r.NameIndex[name] = len(r.Values)
		r.Fields = append(r.Fields, name)
		r.Values = append(r.Values, value)
		r.Length = len(r.Values)
	}
	return true
}



// =====================================
// 增加Session结构体中的方法
// =====================================
func (session *Session) queryRows(sqlStr string, paramStr ...interface{}) (rows *core.Rows, err error) {
	session.queryPreprocess(&sqlStr, paramStr...)

	if session.IsAutoCommit {
		return session.innerQueryRows(session.DB(), sqlStr, paramStr...)
	}
	return session.txQueryRows(session.Tx, sqlStr, paramStr...)
}

func (session *Session) txQueryRows(tx *core.Tx, sqlStr string, params ...interface{}) (rows *core.Rows, err error) {
	rows, err = tx.Query(sqlStr, params...)
	if err != nil {
		return nil, err
	}
	return
}

func (session *Session) innerQueryRows(db *core.DB, sqlStr string, params ...interface{}) (rows *core.Rows, err error) {
	stmt, rows, err := session.Engine.LogSQLQueryTime(sqlStr, params, func() (*core.Stmt, *core.Rows, error) {
		stmt, err := db.Prepare(sqlStr)
		if err != nil {
			return stmt, nil, err
		}
		rows, err := stmt.Query(params...)

		return stmt, rows, err
	})
	if stmt != nil {
		defer stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	return
}


// Exec a raw sql and return records as []*ResultSet
func (session *Session) Q(sqlStr string, paramStr ...interface{}) (resultsSlice []*ResultSet, err error) {

	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	resultsSlice = make([]*ResultSet,0)
	rows,err:=session.queryRows(sqlStr, paramStr...)
	if rows!=nil && err==nil {
		resultsSlice,err=rows2ResultSetSlice(rows)
	}
	if rows!=nil{
		defer rows.Close()
	}
	return
}



// =====================================
// 函数
// =====================================
func rows2ResultSetSlice(rows *core.Rows) (resultsSlice []*ResultSet, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := row2ResultSet(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func row2ResultSet(rows *core.Rows, fields []string) (result *ResultSet, err error) {
	//result := make(map[string]string)
	result = NewResultSet()
	getRowByRows(rows, fields, func(data string,index int,fieldName string){
		//result[fieldName] = data
		result.NameIndex[fieldName] = len(result.Fields)
		result.Fields = append(result.Fields, fieldName)
		result.Values = append(result.Values, data)
	})
	result.Length = len(result.Values)
	return result, nil
}

//获取一行中每一列数据
func getRowByRows(rows *core.Rows, fields []string, fn func(data string,index int,fieldName string)) (err error) {
	length := len(fields)
	scanResultContainers := make([]interface{}, length)
	for i := 0; i < length; i++ {
		var scanResultContainer interface{}
		scanResultContainers[i] = &scanResultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return err
	}
	for ii, key := range fields {
		rawValue := reflect.Indirect(reflect.ValueOf(scanResultContainers[ii]))
		//if row is null then ignore
		if rawValue.Interface() == nil {
			continue
		}
		if data, err := value2String(&rawValue); err == nil {
			fn(data,ii,key)
		} else {
			return err
		}
	}
	return nil
}

//根据core.Rows来查询结果
func getResultSliceByRows(rows *core.Rows,erre error)(resultsSlice []map[string][]byte,err error){
	resultsSlice=make([]map[string][]byte,0)
	if rows!=nil && erre==nil {
		resultsSlice,err = rows2maps(rows)
	}
	if rows!=nil {
		defer rows.Close()
	}
	return
}
