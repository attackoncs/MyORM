package session

import (
	"myorm/clause"
	"reflect"
)

// 多次调Set构造每个子句，调Build按传入顺序构造出最终SQL语句，并调Raw和Exec执行
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// 查询，比较复杂
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem() //获取单个元素类型destType
	//reflect.New创建一个destType实例，作为Model的入参，映射出表结构RefTable
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	//构造出Select语句，查询所有记录rows
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		//利用反射创建destType实例dest，将其所有字段平铺，构造切片values
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		//将该行记录每一列的值依次赋值给 values 中的每一个字段
		if err := rows.Scan(values...); err != nil {
			return err
		}
		//dest 添加到切片 destSlice 中。循环直到所有的记录都添加到切片 destSlice 中
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
