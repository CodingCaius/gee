// Engine 是 GeeORM 与用户交互的入口
// 负责交互前的准备工作（比如连接/测试数据库），交互后的收尾工作（关闭连接）等
// 封装了数据库连接的创建、关闭以及会话管理的功能

package geeorm

import (
	"database/sql"

	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
)

// Engine 结构体是整个库的入口点，用于与数据库进行交互。它包含一个指向数据库连接的指针
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewEngine 构造函数，用于创建并初始化一个 Engine 实例
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

// Close 关闭数据库连接
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// NewSession 创建一个新的数据库会话，该会话将用于执行数据库操作
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}


// 自定义的函数类型，它接受一个 *session.Session 类型的参数，并返回一个 interface{} 类型的结果和一个 error 类型的错误。
// 这个函数类型用于表示一个数据库事务操作
type TxFunc func(*session.Session) (interface{}, error)

// 在 geeorm.go 中为用户提供傻瓜式/一键式使用的事务接口
// 用于执行数据库事务
func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		// 首先，使用 recover() 函数来捕获可能在事务操作中引发的panic
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			// commit失败需要再回滚一次
			defer func ()  {
				if err != nil {
				_ = s.Rollback()
				}
			} ()
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}
