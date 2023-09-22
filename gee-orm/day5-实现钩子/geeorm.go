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
