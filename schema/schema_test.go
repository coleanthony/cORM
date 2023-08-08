package schema

import (
	"cORM/dialect"
	"testing"
)

type Test struct {
	Name string `corm:"PRIMARY KEY"`
	Age  int
}

var Testdial, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&Test{}, Testdial)
	if schema.Name != "Test" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse tag")
	}
}
