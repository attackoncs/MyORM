package session

import (
	"database/sql"
	"myorm/dialect"
	"myorm/mylog"
	"myorm/schema"
	"strings"
)

// 核心结构 Session
type Session struct {
	db       *sql.DB //sql.Open连接数据库成功后返回的指针
	dialect  dialect.Dialect
	refTable *schema.Schema
	sql      strings.Builder //拼接sql语句
	sqlVars  []interface{}   //sql语句中占位符的对应值
}

// 创建session实例
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db,
		dialect: dialect}
}

// 清理session
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

// 获得数据库句柄指针
func (s *Session) DB() *sql.DB {
	return s.db
}

// 改变sql和sqlVars的值
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// 封装Exec方法
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	mylog.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		mylog.Error(err)
	}
	return
}

// 封装查询一行QueryRow方法
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	mylog.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// 封装查询多行Query方法
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	mylog.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		mylog.Error(err)
	}
	return
}
