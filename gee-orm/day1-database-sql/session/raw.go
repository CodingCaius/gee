// 该包负责与数据库的交互

package session

import (
	"database/sql"
	"geeorm/log"
	"strings"
)

type Session struct {
	db      *sql.DB // 使用 sql.Open() 方法连接数据库成功之后返回的指针
	sql     strings.Builder // 用于构建SQL查询语句
	sqlVars []interface{} // 用于存储SQL查询语句中的参数
}

// 用于创建一个新的 Session 对象
func New(db *sql.DB) *Session {
	return &Session{db: db}
}

// 清除 Session 对象中的 SQL 查询语句和参数
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

// 返回 Session 对象中的数据库连接对象 db
func (s *Session) DB() *sql.DB {
	return s.db
}

// 用于构建原始的 SQL 查询语句
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// 封装三个原生方法

// Exec 用于执行原始的 SQL 查询或操作
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRow 仅返回查询结果的第一行
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows gets a list of records from db
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}