package xorm

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/coscms/xorm/core"
)

func indexNoCase(s, sep string) int {
	return strings.Index(strings.ToLower(s), strings.ToLower(sep))
}

func splitNoCase(s, sep string) []string {
	idx := indexNoCase(s, sep)
	if idx < 0 {
		return []string{s}
	}
	return strings.Split(s, s[idx:idx+len(sep)])
}

func splitNNoCase(s, sep string, n int) []string {
	idx := indexNoCase(s, sep)
	if idx < 0 {
		return []string{s}
	}
	return strings.SplitN(s, s[idx:idx+len(sep)], n)
}

func makeArray(elem string, count int) []string {
	res := make([]string, count)
	for i := 0; i < count; i++ {
		res[i] = elem
	}
	return res
}

func rValue(bean interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(bean))
}

func rType(bean interface{}) reflect.Type {
	sliceValue := reflect.Indirect(reflect.ValueOf(bean))
	//return reflect.TypeOf(sliceValue.Interface())
	return sliceValue.Type()
}

func structName(v reflect.Type) string {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Name()
}

func sliceEq(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	sort.Sort(sort.StringSlice(left))
	sort.Sort(sort.StringSlice(right))
	for i := 0; i < len(left); i++ {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func reflect2value(rawValue *reflect.Value) (str string, err error) {
	aa := reflect.TypeOf((*rawValue).Interface())
	vv := reflect.ValueOf((*rawValue).Interface())
	switch aa.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str = strconv.FormatInt(vv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str = strconv.FormatUint(vv.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		str = strconv.FormatFloat(vv.Float(), 'f', -1, 64)
	case reflect.String:
		str = vv.String()
	case reflect.Array, reflect.Slice:
		switch aa.Elem().Kind() {
		case reflect.Uint8:
			data := rawValue.Interface().([]byte)
			str = string(data)
		default:
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
	//时间类型
	case reflect.Struct:
		if aa == core.TimeType {
			str = rawValue.Interface().(time.Time).Format(time.RFC3339Nano)
		} else {
			err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
		}
	case reflect.Bool:
		str = strconv.FormatBool(vv.Bool())
	case reflect.Complex128, reflect.Complex64:
		str = fmt.Sprintf("%v", vv.Complex())
	/* TODO: unsupported types below
	   case reflect.Map:
	   case reflect.Ptr:
	   case reflect.Uintptr:
	   case reflect.UnsafePointer:
	   case reflect.Chan, reflect.Func, reflect.Interface:
	*/
	default:
		err = fmt.Errorf("Unsupported struct type %v", vv.Type().Name())
	}
	return
}

func value2Bytes(rawValue *reflect.Value) (data []byte, err error) {
	var str string
	str, err = reflect2value(rawValue)
	if err != nil {
		return
	}
	data = []byte(str)
	return
}

func value2String(rawValue *reflect.Value) (data string, err error) {
	data, err = reflect2value(rawValue)
	if err != nil {
		return
	}
	return
}

func rows2Strings(rows *core.Rows) (resultsSlice []*ResultSet, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := row2mapStr(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func rows2maps(rows *core.Rows) (resultsSlice []map[string][]byte, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := row2map(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func row2mapStr(rows *core.Rows, fields []string) (result *ResultSet, err error) {
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

func row2map(rows *core.Rows, fields []string) (resultsMap map[string][]byte, err error) {
	result := make(map[string][]byte)
	getRowByRows(rows, fields,func(data string,index int,fieldName string){
		result[fieldName] = []byte(data)
	})
	return result, nil
}

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
	}
	return false
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
