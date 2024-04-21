package myorm

import (
	"database/sql"
	"myorm/mylog"
	"myorm/session"
)

// 核心数据结构，负责与用户的交互的入口
type Engine struct {
	db *sql.DB
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
	e = &Engine{db: db}
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
	return session.New(e.db)
}
