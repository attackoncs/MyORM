package myorm

import (
	"database/sql"
	"myorm/dialect"
	"myorm/mylog"
	"myorm/session"
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
