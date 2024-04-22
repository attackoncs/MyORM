// 数据库表相关操作
package session

import (
	"fmt"
	"myorm/mylog"
	"myorm/schema"
	"reflect"
	"strings"
)

// 解析耗时将解析结果保存在成员变量refTable中，只要结构体名不变不更新refTable值
func (s *Session) Model(value interface{}) *Session {
	// nil or different model, update refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// 返回refTable的值
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		mylog.Error("Model is not set")
	}
	return s.refTable
}

// 创建表，利用reftable返回的数据库表和字段的信息，拼接sql语句，调原生sql语句执行
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// 删除表
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

// 若表存在则返回true
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
