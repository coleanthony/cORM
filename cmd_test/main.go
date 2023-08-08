package main

import (
	"cORM"
	"cORM/log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine, err := cORM.NewEngine("sqlite3", "test.db")
	if err != nil {
		log.Error("connect database error")
	}
	defer func() { engine.Close() }()
	s := engine.NewSession()

	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, err := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()

	if err == nil {
		affected, _ := result.RowsAffected()
		log.Info(affected)
	}

	row := s.Raw("SELECT Name FROM User LIMIT 1").QueryRow()
	var name string
	if err := row.Scan(&name); err == nil {
		log.Info(name)
	}
}
