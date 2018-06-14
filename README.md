# Sqlite3CreateTableParser
:scroll: Advanced ~~PRAGMA~~ table_info through DDL parsing


## GOLANG port of [C library](https://github.com/marcobambini/sqlite-createtable-parser) done by
[@marcobambini](https://github.com/marcobambini).


## SQLite CREATE TABLE parser
A parser for sqlite create table sql statements.

SQLite is a very powerful software but it lacks an easy way to extract complete information about table and columns constraints. The built-in sql pragma:
```c
PRAGMA schema.table_info(table-name);
PRAGMA foreign_key_list(table-name);
```
provide incomplete information and a manual parsing is required in order to extract more useful information.

CREATE TABLE syntax diagrams can be found on the official [sqlite website](http://www.sqlite.org/lang_createtable.html).


## Usage

```go
package main

import "github.com/lempiy/Sqlite3CreateTableParser/parser"

//some fancy DDL
const ddl = `
CREATE TABLE contact_groups (
 contact_id integer,
 group_id integer,
 PRIMARY KEY (contact_id, group_id),
 FOREIGN KEY (contact_id) REFERENCES contacts (contact_id)
 ON DELETE CASCADE ON UPDATE NO ACTION,
 FOREIGN KEY (group_id) REFERENCES groups (group_id)
 ON DELETE CASCADE ON UPDATE NO ACTION
);
`

func main() {
    table, errCode := parser.ParseTable(sql, 0)
    if errCode != parser.ERROR_NONE {
        panic("Error during parsing sql")
    }
    // do stuff with received data
    fmt.Printf("%+v\n", table)
}
```


## Table info structs
```go
type Table struct {
	Name           string
	Schema         string
	IsTemporary    bool
	IsIfNotExists  bool
	IsWithoutRowid bool
	NumColumns     int
	Columns        []Column
	NumConstraint  int
	Constraints    []TableConstraint
}

type Column struct {
	Name                  string
	Type                  string
	Length                string
	ConstraintName        string
	IsPrimaryKey          bool
	IsAutoincrement       bool
	IsNotnull             bool
	IsUnique              bool
	PkOrder               OrderClause
	PkConflictClause      ConflictClause
	NotNullConflictClause ConflictClause
	UniqueConflictClause  ConflictClause
	CheckExpr             string
	DefaultExpr           string
	CollateName           string
	ForeignKeyClause      *ForeignKey
}

type TableConstraint struct {
	Type             ConstraintType
	Name             string
	NumIndexed       int
	IndexedColumns   []IdxColumn
	ConflictClause   ConflictClause
	CheckExpr        string
	ForeignKeyNum    int
	ForeignKeyName   []string
	ForeignKeyClause *ForeignKey
}

type ForeignKey struct {
	Table      string
	NumColumns int
	ColumnName []string
	OnDelete   FkAction
	OnUpdate   FkAction
	Match      string
	Deferrable FkDefType
}

type IdxColumn struct {
	Name        string
	CollateName string
	Order       OrderClause
}
```


## Limitations
- CREATE TABLE AS select-stmt syntax is not supported (SQL3ERROR_UNSUPPORTEDSQL is returned).
- EXPRESSIONS in column constraints (CHECK and DEFAULT constraint) and table constraint (CHECK constraint) are not supported (SQL3ERROR_UNSUPPORTEDSQL is returned).
