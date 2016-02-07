package models

// Schema represents a simplified database schema and possesses a list of Tables
type Schema struct {
	Tables  []*Table `db:"Tables"`
	Adapter string
}

// NewSchema returns a new schema pointer with the Tables field initialized
func NewSchema(adapter string) *Schema {
	return &Schema{Tables: make([]*Table, 0), Adapter: adapter}
}

// Table represents a sql Table with a Name and Rows, which give at least a Field Name and Type
type Table struct {
	Name string      `db:"Name"`
	Rows []*TableRow `db:"Rows"`
}

// TableRow represents a single row from a Table
type TableRow struct {
	Field   string      `db:"Field"`
	Type    string      `db:"Type"`
	Null    string      `db:"Null"`
	Key     string      `db:"Key"`
	Default interface{} `db:"Default"`
	Extra   string      `db:"Extra"`
}

type SQLiteTableRow struct {
	Field string `db:"name"`
	Type  string `db:"type"`
}
