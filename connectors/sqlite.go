package connectors

import (
	"fmt"

	"github.com/alistanis/stg/models"
	"github.com/jmoiron/sqlx"
)

// SQLiteConfig represents a SQLite configuration
type SQLiteConfig struct {
	Config `yaml:",inline"`
}

// NewSQLiteConfig returns a SQliteConfig struct with UseDefaults set to true
func NewSQLiteConfig() *SQLiteConfig {
	config := &SQLiteConfig{}
	config.SetDefaults()
	return config
}

// SQLiteConfigFromConfig returns a *SQLiteConfig with defaults
func SQLiteConfigFromConfig(config Config) *SQLiteConfig {
	sqlite := &SQLiteConfig{}
	sqlite.Config = config
	sqlite.SetDefaults()
	return sqlite
}

// open pens a connection to the SQLite database
func (lite *SQLiteConfig) Open() (*sqlx.DB, error) {
	return sqlx.Open(lite.Adapter, lite.Database)
}

// SetDefaults sets SQLite default settings
func (lite *SQLiteConfig) SetDefaults() {
	if lite.Encoding == "" {
		lite.Encoding = "utf8"
	}
	if lite.Adapter == "" {
		lite.Adapter = SQLite
	}
}

// getSqliteSchema returns a schema from sqlite
func (db *DB) getSqliteSchema() (*models.Schema, error) {
	schema := models.NewSchema(SQLite)
	tableNames := TableNames{}
	err := db.DB.Select(&tableNames, "SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		return nil, err
	}
	for _, name := range tableNames {
		table := &models.Table{}
		rows, err := db.DB.Unsafe().Queryx(fmt.Sprintf("PRAGMA table_info(%s);", name))

		if err != nil {
			return nil, err
		}
		for rows.Next() {
			tableRow := &models.SQLiteTableRow{}
			err = rows.StructScan(tableRow)
			if err != nil {
				return nil, err
			}
			table.Rows = append(table.Rows, &models.TableRow{Field: tableRow.Field, Type: tableRow.Type})
		}
		table.Name = name
		schema.Tables = append(schema.Tables, table)
	}
	return schema, nil
}
