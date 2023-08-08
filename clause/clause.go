package clause

import (
	"strings"
)

type Type int
type Clause struct {
	sql     map[Type]string
	sqlVals map[Type][]interface{}
}

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

// 实现结构体 Clause 拼接各个独立的子句
// Set 方法根据 Type 调用对应的 generator，生成该子句对应的 SQL 语句。
// Build 方法根据传入的 Type 的顺序，构造出最终的 SQL 语句。
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVals = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVals[name] = vars
}

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vals []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vals = append(vals, c.sqlVals[order]...)
		}
	}
	return strings.Join(sqls, " "), vals
}
