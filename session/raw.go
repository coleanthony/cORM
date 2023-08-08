package session

import (
	"cORM/clause"
	"cORM/dialect"
	"cORM/log"
	"cORM/schema"
	"database/sql"
	"strings"
)

// 用于实现与数据库的交互
type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder
	sqlVals  []interface{}
	tx       *sql.Tx
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	s := &Session{
		db:      db,
		dialect: dialect,
	}
	return s
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVals = nil
	s.clause = clause.Clause{}
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVals = append(s.sqlVals, values...)
	return s
}

// 封装 Exec()、Query() 和 QueryRow() 三个原生方法。
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVals)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVals...); err != nil {
		log.Error(err)
	}
	return
}

// 从db中取得一条记录
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVals)
	return s.DB().QueryRow(s.sql.String(), s.sqlVals...)
}

// 从db中取得一列记录
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVals)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVals...); err != nil {
		log.Error(err)
	}
	return
}
