// 用于放置操作数据库表相关的代码

package session

import (
	"fmt"
	"geeorm/log"
	"reflect"
	"strings"

	"geeorm/schema"
)

// Model 用于设置或更新会话对象的 refTable 字段
func (s *Session) Model(value interface{}) *Session {
	// 检查 s.refTable 是否为 nil，或者传入的 value 与已有的模型类型是否不同。如果是其中之一，说明需要更新 refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// RefTable 用于获取会话对象的 refTable 字段，即包含了与数据库表格相关的元数据信息的 schema.Schema 实例
func (s *Session) RefTable() *schema.Schema {
	// 这是为了确保在使用 RefTable 之前必须先调用 Model 方法来设置 refTable
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// CreateTable 创建数据库表格
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// DropTable 删除数据库表格
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

// HasTable 检查数据库中是否存在指定的表格
func (s *Session) HasTable() bool {
	// 获取检查表格是否存在的SQL语句，以及相应的参数值
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	// 执行该SQL语句，得到一个包含表格名称的结果
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	// 结果扫描（Scan）到 tmp 变量中，如果 tmp 的值与表格的名称相同，就表示表格存在，返回 true，否则返回 false
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}