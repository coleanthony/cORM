package dialect

import "reflect"

//将 Go 语言的类型映射为数据库中的类型
//不同数据库支持的数据类型也是有差异的，即使功能相同，在 SQL 语句的表达上也可能有差异。
//ORM 框架往往需要兼容多种数据库，因此我们需要将差异的这一部分提取出来，每一种数据库分别实现，实现最大程度的复用和解耦。

var dialectsMap = map[string]Dialect{}

// DataTypeOf 用于将 Go 语言的类型转换为该数据库的数据类型。
// TableExistSQL 返回某个表是否存在的 SQL 语句，参数是表名(table)。
type Dialect interface {
	DataTypeof(typ reflect.Value) string
	TableExistSQL(tablename string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (Dialect, bool) {
	dialect, ok := dialectsMap[name]
	return dialect, ok
}
