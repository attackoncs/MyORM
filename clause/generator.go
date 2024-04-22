// 子句（clause）的生成规则
package clause

import (
	"fmt"
	"strings"
)

// 函数类型
type generator func(values ...interface{}) (string, []interface{})

// 全局子句、函数键值对
var generators map[Type]generator

// 注册子句、函数键值对
func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
}

// 生成num个占位符字符串，并用", "分割
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

// insert子句，规则如：INSERT INTO $tableName ($fields)
func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

// values子句，规则如：VALUES ($v1), ($v2), ...
func _values(values ...interface{}) (string, []interface{}) {
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

// select子句，规则如：SELECT $fields FROM $tableName
func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

// limit子句，规则如：LIMIT $num
func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

// where子句，规则如：WHERE $desc
func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

// order子句
func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}
