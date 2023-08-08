package cORM

import (
	"cORM/dialect"
	"cORM/log"
	"cORM/session"
	"database/sql"
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
