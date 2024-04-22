// 特定SQL语句的转换，给定对象，转换为关系型数据库中的表结构
package schema

import (
	"go/ast"
	"myorm/dialect"
	"reflect"
)

// 字段
type Field struct {
	Name string //字段名
	Type string //类型Type
	Tag  string //约束条件
}

// 模式
type Schema struct {
	Model      interface{}       //对象
	Name       string            //表名
	Fields     []*Field          //字段
	FieldNames []string          //字段名
	fieldMap   map[string]*Field //字段名和字段的映射关系，无需遍历字段列表
}

// 返回schema的字段名对应的字段
func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

// 返回dest成员变量的值，从对象中找到对应值，按顺序平铺
func (s *Schema) RecordValues(dest interface{}) []interface{} {
	//间接获取，若用reflect.Elem对结构体会panic,而indirect可兼容结构体和指针
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

// 表名
type ITableName interface {
	TableName() string
}

// 解析表名dest和dialect，将任意对象解析为Schema实例
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	var tableName string
	t, ok := dest.(ITableName)
	if !ok {
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}
	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("myorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
