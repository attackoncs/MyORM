// 隔离不同数据库之间的差异，便于扩展
package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

// 屏蔽不同数据库差异，目前包含2个方法
type Dialect interface {
	//将go语言数据类型转为数据库数据类型
	DataTypeOf(typ reflect.Value) string
	//返回某个表是否存在的SQL语句，参数是表名
	TableExistSQL(tableName string) (string, []interface{})
}

// 注册dialect实例，增加对某个数据库的支持
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// 获取dialect实例，也就是数据库driver名
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
