// 抽象出各个数据库差异的部分

package dialect

import "reflect"

// 全局变量，用于存储不同数据库方言的映射
var dialectsMap = map[string]Dialect{}

// 接口，它定义了处理不同数据库方言所需的方法
type Dialect interface {
	// 用于将 Go 语言中的数据类型映射到数据库特定的数据类型
	DataTypeOf(typ reflect.Value) string
	// 用于生成检查数据库表是否存在的 SQL 查询语句
	TableExistSQL(tableName string) (string, []interface{})
}

// RegisterDialect 用于注册数据库方言到 dialectsMap 中
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect 获取已注册的数据库方言
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
