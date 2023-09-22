// 拼接各个独立的子句

// 这个包的主要目的是将不同类型的 SQL 子句和相关的变量进行管理和拼接，使构建复杂的 SQL 查询语句变得更加方便。
// 在使用时，可以调用 Set 方法来设置不同类型的 SQL 子句，然后调用 Build 方法来构建完整的查询语句

package clause

import "strings"

type Clause struct {
	sql     map[Type]string        // 用于存储不同类型的 SQL 子句
	sqlVars map[Type][]interface{} // 用于存储与每个子句相关的变量
}

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// Set 根据 Type 调用对应的 generator，生成该子句对应的 SQL 语句
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// Build 构建最终的 SQL 查询语句
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	// 可变参数 orders，存储的是类型，表示需要构建的 SQL 子句的顺序
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
