// sqlite3的支持
package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct {
}

var _ Dialect = (*sqlite3)(nil)

// 包在第一次加载时，将sqlite3数据库名和结构体指针添加到全局哈希表
func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

// 利用反射将Go语言类型映射为SQLite数据类型
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// 返回在SQLite中判断表tableName是否存在的SQL语句
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}
