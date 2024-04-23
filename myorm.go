package myorm

import (
	"database/sql"
	"fmt"
	"myorm/dialect"
	"myorm/mylog"
	"myorm/session"
	"strings"
)

// 核心数据结构，负责与用户的交互的入口
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// 创建engine实例，连接数据库并检查连接是否正常。
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		mylog.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		mylog.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver) //获取driver对应的dialect
	if !ok {
		mylog.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	mylog.Info("Connecct database success")
	return
}

// 关闭数据库
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		mylog.Error("Failed to close database")
	}
	mylog.Info("Close database success")
}

// 通过engine实例创建会话，与数据库进行交互
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

// 函数接口类型
type TxFunc func(*session.Session) (interface{}, error)

// https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
// 事务执行一段sql语句，若没错误，就自动提交
func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) //在Rollback后重新panic
		} else if err != nil {
			_ = s.Rollback() //err非空，不改变数据库数据
		} else {
			err = s.Commit() //err空，若commit返回错误则更新错误
		}
	}()
	return f(s)
}

// 返回两个字符串切片a中没有b中的元素的切片，新表-旧表=新增字段，旧表-新表=删除字段
func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// 迁移表，先求出两个表的字段切片的差集，用Alter语句新增字段，用创建新表并重命名方式删除字段
func (e *Engine) Migrate(value interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			mylog.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		//求出两个表的字段切片的差集
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		mylog.Infof("added cols %v,deleted cols %v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			//alter语句新增字段
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}
		//创建新表并重命名的方式删除字段
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}
