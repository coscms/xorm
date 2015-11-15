package xorm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/coscms/xorm/core"
)

type SelectRows []*ResultSet

func (this SelectRows) GetRow() (result *ResultSet) {
	if len(this) > 0 {
		result = this[0]
	}
	return
}

func (this SelectRows) GetOne() (result string) {
	results := this.GetRow()
	if results != nil {
		result = results.Get(0)
	}
	return
}

// =====================================
// 定义ResultSet
// =====================================
func NewResultSet() *ResultSet {
	return &ResultSet{
		Fields:    make([]string, 0),
		Values:    make([]string, 0),
		NameIndex: make(map[string]int),
		Length:    0,
	}
}

type ResultSet struct {
	Fields    []string
	Values    []string
	NameIndex map[string]int
	Length    int
}

func (r *ResultSet) Get(index int) string {
	if index >= r.Length {
		return ""
	}
	return r.Values[index]
}

func (r *ResultSet) GetByName(name string) string {
	if index, ok := r.NameIndex[name]; ok {
		return r.Get(index)
	}
	return ""
}

func (r *ResultSet) Set(index int, value string) bool {
	if index >= r.Length {
		return false
	}
	r.Values[index] = value
	return true
}

func (r *ResultSet) SetByName(name string, value string) bool {
	if index, ok := r.NameIndex[name]; ok {
		return r.Set(index, value)
	} else {
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
		stmt.Close()
	}
	if err != nil {
		return nil, err
	}
	return
}

/**
 * Exec a raw sql and return records as []*ResultSet
 * @param  string					SQL
 * @param  ...interface{}			params
 * @return []*ResultSet,error
 * @author AdamShen (swh@admpub.com)
 */
func (session *Session) Q(sqlStr string, paramStr ...interface{}) (resultsSlice []*ResultSet, err error) {

	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	resultsSlice = make([]*ResultSet, 0)
	rows, err := session.queryRows(sqlStr, paramStr...)
	if rows != nil {
		if err == nil {
			resultsSlice, err = rows2ResultSetSlice(rows)
		}
		rows.Close()
	}
	return
}

/**
 * 逐行执行回调函数
 * @param  func(*core.Rows) callback		callback func
 * @param  string sqlStr 					SQL
 * @param  ...interface{} paramStr			params
 * @return error
 * @author AdamShen (swh@admpub.com)
 * @example
 * QCallback(func(rows *core.Rows){
 * 	if err := rows.Scan(bean); err != nil {
 *		return
 *	}
 *	//.....
 * },"SELECT * FROM shop WHERE type=?","vip")
 */
func (session *Session) QCallback(callback func(*core.Rows, []string), sqlStr string, paramStr ...interface{}) (err error) {

	defer session.resetStatement()
	if session.IsAutoClose {
		defer session.Close()
	}

	rows, err := session.queryRows(sqlStr, paramStr...)
	if rows != nil {
		if err == nil {
			var fields []string
			fields, err = rows.Columns()
			if err != nil {
				return err
			}
			for rows.Next() {
				callback(rows, fields)
			}
		}
		rows.Close()
	}
	return
}

// =====================================
// 增加Engine结构体中的方法
// =====================================
func (this *Engine) RawQuery(sql string, paramStr ...interface{}) (resultsSlice []*ResultSet, err error) {
	session := this.NewSession()
	defer session.Close()
	resultsSlice, err = session.Q(sql, paramStr...)
	return
}

func (this *Engine) RawQueryCallback(callback func(*core.Rows, []string), sql string, paramStr ...interface{}) (err error) {
	session := this.NewSession()
	defer session.Close()
	err = session.QCallback(callback, sql, paramStr...)
	return
}

/**
 * 查询键值对
 */
func (this *Engine) RawQueryKv(key string, val string, sql string, paramStr ...interface{}) map[string]string {
	var results map[string]string = make(map[string]string, 0)
	err := this.RawQueryCallback(func(rows *core.Rows, fields []string) {
		var result map[string]string = make(map[string]string)
		LineAllFieldsProcessing(rows, fields, func(data string, index int, fieldName string) {
			result[fieldName] = data
		})
		if k, ok := result[key]; ok {
			if v, ok := result[val]; ok {
				results[k] = v
			}
		}
	}, sql, paramStr...)
	if err != nil {
		this.TagLogError("base", err)
	}
	return results
}

func (this *Engine) RawQueryKeySlice(key string, sql string, paramStr ...interface{}) map[string][]map[string]string {
	var results map[string][]map[string]string = make(map[string][]map[string]string, 0)
	err := this.RawQueryCallback(func(rows *core.Rows, fields []string) {
		var result map[string]string = make(map[string]string)
		LineAllFieldsProcessing(rows, fields, func(data string, index int, fieldName string) {
			result[fieldName] = data
		})
		if k, ok := result[key]; ok {
			if _, ok := results[k]; !ok {
				results[k] = make([]map[string]string, 0)
			}
			results[k] = append(results[k], result)
		}
	}, sql, paramStr...)
	if err != nil {
		this.TagLogError("base", err)
	}
	return results
}

// =======================
// 原生SQL查询
// =======================
func (this *Engine) GetRows(sql string, params ...interface{}) []*ResultSet {
	sql = this.ReplaceTablePrefix(sql)
	result, err := this.RawQuery(sql, params...)
	if err != nil {
		this.TagLogError("base", err)
	}
	return result
}

func (this *Engine) GetRow(sql string, params ...interface{}) (result *ResultSet) {
	sql = this.ReplaceTablePrefix(sql)
	results, err := this.RawQuery(sql+" LIMIT 1", params...)
	if err != nil {
		this.TagLogError("base", err)
		return
	}
	if len(results) > 0 {
		result = results[0]
	}
	return
}

func (this *Engine) GetOne(sql string, params ...interface{}) (result string) {
	results := this.GetRow(sql, params...)
	if results != nil {
		result = results.Get(0)
	}
	return
}

func (this *Engine) RawFetchAll(fields string, table string, where string, params ...interface{}) []map[string]string {
	if fields == "" {
		fields = "*"
	} else {
		fields = this.ReplaceTablePrefix(fields)
	}
	sql := `SELECT ` + fields + ` FROM ` + this.fullTableName(table) + ` WHERE ` + this.ReplaceTablePrefix(where)
	if len(params) == 1 {
		switch params[0].(type) {
		case []interface{}:
			return this.RawQuerySlice(sql, params[0].([]interface{})...)
		}
	}
	return this.RawQuerySlice(sql, params...)
}

func (this *Engine) RawFetch(fields string, table string, where string, params ...interface{}) (result map[string]string) {
	if fields == "" {
		fields = "*"
	} else {
		fields = this.ReplaceTablePrefix(fields)
	}
	sql := `SELECT ` + fields + ` FROM ` + this.fullTableName(table) + ` WHERE ` + this.ReplaceTablePrefix(where) + ` LIMIT 1`
	if len(params) == 1 {
		switch params[0].(type) {
		case []interface{}:
			results := this.RawQuerySlice(sql, params[0].([]interface{})...)
			if len(results) > 0 {
				result = results[0]
			}
			return
		}
	}
	results := this.RawQuerySlice(sql, params...)
	if len(results) > 0 {
		result = results[0]
	}
	return
}

// RawSelect("*","member","id=?",1)
// RawSelect("*","member","status=? AND sex=?",1,1)
// RawSelect("*","`~member` a,`~order` b","a.status=? AND b.status=?",1,1)
func (this *Engine) RawSelect(fields string, table string, where string, params ...interface{}) SelectRows {
	if fields == "" {
		fields = "*"
	} else {
		fields = this.ReplaceTablePrefix(fields)
	}
	sql := `SELECT ` + fields + ` FROM ` + this.fullTableName(table) + ` WHERE ` + this.ReplaceTablePrefix(where)
	if len(params) == 1 {
		switch params[0].(type) {
		case []interface{}:
			return this.GetRows(sql, params[0].([]interface{})...)
		}
	}
	return SelectRows(this.GetRows(sql, params...))
}

func (this *Engine) RawInsert(table string, sets map[string]interface{}) (lastId int64) {
	fields := ""
	values := ""
	params := make([]interface{}, 0)
	for k, v := range sets {
		if fields != "" {
			fields += ","
			values += ","
		}
		fields += this.Quote(k)
		values += "?"
		params = append(params, v)
	}
	sql := `INSERT INTO ` + this.fullTableName(table) + ` (` + fields + `) VALUES (` + values + `)`
	return this.RawExec(sql, true, params...)
}

func (this *Engine) RawReplace(table string, sets map[string]interface{}) int64 {
	fields := ""
	values := ""
	params := make([]interface{}, 0)
	for k, v := range sets {
		if fields != "" {
			fields += ","
			values += ","
		}
		fields += this.Quote(k)
		values += "?"
		params = append(params, v)
	}
	sql := `REPLACE INTO ` + this.fullTableName(table) + ` (` + fields + `) VALUES (` + values + `)`
	return this.RawExec(sql, false, params...)
}

func (this *Engine) RawUpdate(table string, sets map[string]interface{}, where string, args ...interface{}) int64 {
	set := ""
	params := make([]interface{}, 0)
	for k, v := range sets {
		if set != "" {
			set += ","
		}
		set += this.Quote(k) + "=?"
		params = append(params, v)
	}
	if len(args) > 0 {
		isAloneSlice := false
		if len(args) == 1 {
			switch args[0].(type) {
			case []interface{}:
				params = append(params, args[0].([]interface{})...)
				isAloneSlice = true
			}
		}
		if !isAloneSlice {
			for _, v := range args {
				params = append(params, v)
			}
		}
	}
	sql := `UPDATE ` + this.fullTableName(table) + ` SET ` + set + ` WHERE ` + where

	return this.RawExec(sql, false, params...)
}

func (this *Engine) RawDelete(table string, where string, params ...interface{}) int64 {
	sql := `DELETE FROM ` + this.fullTableName(table) + ` WHERE ` + where
	if len(params) == 1 {
		switch params[0].(type) {
		case []interface{}:
			return this.RawExec(sql, false, params[0].([]interface{})...)
		}
	}
	return this.RawExec(sql, false, params...)
}

func (this *Engine) ReplaceTablePrefix(sql string) (r string) {
	r = strings.Replace(sql, "~", this.TablePrefix, -1)
	return
}
func (this *Engine) fullTableName(table string) string {
	if table[0] != '`' && table[0] != '~' {
		table = this.Quote(this.TablePrefix + table + this.TableSuffix)
	}
	table = this.ReplaceTablePrefix(table)
	return table
}

func (this *Engine) QuoteValue(s string) string {
	return "'" + AddSlashes(s) + "'"
}

/**
 * 查询基于指定字段值为键名的map
 */
func (this *Engine) RawQueryKvs(key string, sql string, paramStr ...interface{}) map[string]map[string]string {
	if key == "" {
		key = "id"
	}
	var results map[string]map[string]string = make(map[string]map[string]string, 0)
	err := this.RawQueryCallback(func(rows *core.Rows, fields []string) {
		var result map[string]string = make(map[string]string)
		LineAllFieldsProcessing(rows, fields, func(data string, index int, fieldName string) {
			result[fieldName] = data
		})
		if k, ok := result[key]; ok {
			results[k] = result
		}
	}, sql, paramStr...)
	if err != nil {
		this.TagLogError("base", err)
	}
	return results
}

/**
 * 查询[]map[string]string
 */
func (this *Engine) RawQuerySlice(sql string, paramStr ...interface{}) []map[string]string {
	var results []map[string]string = make([]map[string]string, 0)
	err := this.RawQueryCallback(func(rows *core.Rows, fields []string) {
		var result map[string]string = make(map[string]string)
		LineAllFieldsProcessing(rows, fields, func(data string, index int, fieldName string) {
			result[fieldName] = data
		})
		results = append(results, result)
	}, sql, paramStr...)
	if err != nil {
		this.TagLogError("base", err)
	}
	return results
}

func (this *Engine) RawExec(sql string, retId bool, params ...interface{}) (affected int64) {
	if result, err := this.Exec(sql, params...); err != nil {
		if retId {
			affected, err = result.LastInsertId()
		} else {
			affected, err = result.RowsAffected()
		}
		if err != nil {
			this.TagLogError("base", err)
		}
	} else {
		this.TagLogError("base", err)
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
	LineAllFieldsProcessing(rows, fields, func(data string, index int, fieldName string) {
		//result[fieldName] = data
		result.NameIndex[fieldName] = len(result.Fields)
		result.Fields = append(result.Fields, fieldName)
		result.Values = append(result.Values, data)
	})
	result.Length = len(result.Values)
	return result, nil
}

//获取一行中每一列数据
func LineAllFieldsProcessing(rows *core.Rows, fields []string, fn func(data string, index int, fieldName string)) (err error) {
	length := len(fields)
	scanResultContainers := make([]interface{}, length)
	for i := 0; i < length; i++ {
		var resultContainer interface{}
		scanResultContainers[i] = &resultContainer
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
			fn(data, ii, key)
		} else {
			return err
		}
	}
	return nil
}

//根据core.Rows来查询结果
func getResultSliceByRows(rows *core.Rows, erre error) (resultsSlice []map[string][]byte, err error) {
	resultsSlice = make([]map[string][]byte, 0)
	if rows != nil {
		if erre == nil {
			resultsSlice, err = rows2maps(rows)
		}
		rows.Close()
	}
	return
}

//替换sql中的占位符
func BuildSqlResult(sqlStr string, args interface{}) string {
	if args, ok := args.([]interface{}); ok {
		for _, v := range args {
			val := ""
			switch v.(type) {
			case []interface{}:
				vals := v.([]interface{})
				for _, v := range vals {
					if val != "" {
						val += ","
					}
					rv := fmt.Sprintf("%v", v)
					rv = AddSlashes(rv)
					val += "'" + rv + "'"
				}
				val = strings.Replace(val, "'", `\'`, -1)
			default:
				val = fmt.Sprintf("%v", v)
				val = AddSlashes(val)
				val = "'" + val + "'"
			}
			sqlStr = strings.Replace(sqlStr, "?", val, 1)
		}
	}
	//fmt.Printf("%v\n", sqlStr)
	return sqlStr
}

func AddSlashes(s string) string {
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, "'", `\'`, -1)
	return s
}
