// 封装go标准库的原生事务的支持，也方便统一打印日志
package session

import "myorm/mylog"

// 开启事务的函数封装，赋值s.tx
func (s *Session) Begin() (err error) {
	mylog.Info("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {
		mylog.Error(err)
		return
	}
	return
}

// 提交事务的函数封装
func (s *Session) Commit() (err error) {
	mylog.Info("transaction ccommit")
	if err = s.tx.Commit(); err != nil {
		mylog.Error(err)
	}
	return
}

// 回滚事务的封装
func (s *Session) Rollback() (err error) {
	mylog.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		mylog.Error(err)
	}
	return
}
