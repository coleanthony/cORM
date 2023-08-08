package cORM

import (
	"cORM/dialect"
	"cORM/log"
	"cORM/session"
	"database/sql"
	"fmt"
	"strings"
)

// 实现Engine，负责连接/测试数据库，关闭连接等
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// 连接数据库，返回 *sql.DB
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Error("get dialect error")
		return
	}

	e = &Engine{
		db:      db,
		dialect: dial,
	}
	log.Info("Connect to database successfully")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Close database error")
	}
	log.Info("Close database successfully")
}

// 通过 Engine 实例创建会话，进而与数据库进行交互
func (engine *Engine) NewSession() *session.Session {
	s := session.New(engine.db, engine.dialect)
	return s
}

type TxFunc func(*session.Session) (interface{}, error)

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	return f(s)
}

func difference(a []string, b []string) (diff []string) {
	mymap := make(map[string]bool)
	for _, v := range b {
		mymap[v] = true
	}
	for _, v := range a {
		if _, ok := mymap[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

/*
1.1 新增字段
ALTER TABLE table_name ADD COLUMN col_name, col_type;
大部分数据支持使用 ALTER 关键字新增字段，或者重命名字段。

1.2 删除字段
对于 SQLite 来说，删除字段并不像新增字段那么容易，一个比较可行的方法需要执行下列几个步骤：
CREATE TABLE new_table AS SELECT col1, col2, ... from old_table
DROP TABLE old_table
ALTER TABLE new_table RENAME TO old_table;
第一步：从 old_table 中挑选需要保留的字段到 new_table 中。
第二步：删除 old_table。
第三步：重命名 new_table 为 old_table。
*/

func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s does not exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("add cols %v,del cols %v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			sqlstr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlstr).Exec(); err != nil {
				return
			}
		}
		if len(delCols) == 0 {
			return
		}
		tmp := "tmp_" + table.Name
		fieldstr := strings.Join(table.FieldNames, ", ")

		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", tmp, fieldstr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}
