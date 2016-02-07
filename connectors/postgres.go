package connectors

import (
	"bytes"
	"fmt"

	"github.com/alistanis/stg/models"
	"github.com/jmoiron/sqlx"
)

// PGConfig represents a Postgres configuration struct
type PGConfig struct {
	Config            `yaml:",inline"`
	SSLMode           string `yaml:"ssl_mode" json:"ssl_mode"`
	ConnectionTimeout int    `yaml:"connection_timeout" json:"connection_timeout"`
}

// NewPGConfig returns a pointer to a postgres configuration struct with UseDefaults set to true
func NewPGConfig() *PGConfig {
	config := &PGConfig{}
	config.SetDefaults()
	return config
}

// PGConfigFromConfig sets defaults specific to postgres.
func PGConfigFromConfig(config Config) *PGConfig {
	pg := &PGConfig{}
	pg.Config = config
	pg.SetDefaults()
	return pg
}

// open assembles a datasource string and attempts to open a connection to a postgres/redshift server
func (pg *PGConfig) Open() (*sqlx.DB, error) {
	buffer := make([]byte, 0)
	writer := bytes.NewBuffer(buffer)

	writer.WriteString(fmt.Sprintf("user=%s dbname=%s", pg.Username, pg.Database))

	if pg.Password != "" {
		writer.WriteString(fmt.Sprintf(" password=%s", pg.Password))
	}
	writer.WriteString(fmt.Sprintf(" host=%s", pg.Host))
	if pg.Port != 0 {
		writer.WriteString(fmt.Sprintf(" port=%d", pg.Port))
	}
	if pg.ConnectionTimeout != 0 {
		writer.WriteString(fmt.Sprintf(" connect_timeout=%d", pg.ConnectionTimeout))
	}
	writer.WriteString(fmt.Sprintf(" sslmode="))
	if pg.SSLMode != "" {
		writer.WriteString(pg.SSLMode)
	} else {
		writer.WriteString("disable")
	}
	return sqlx.Open(pg.Adapter, writer.String())
}

// SetDefaults sets defaults specific to postgres.
func (pg *PGConfig) SetDefaults() {
	if pg.Encoding == "" {
		pg.Encoding = "utf8"
	}
	if pg.Adapter == "" {
		pg.Adapter = Postgres
	}
	if pg.Port == 0 {
		pg.Port = 5439
	}
	if pg.SSLMode == "" {
		pg.SSLMode = "require"
	}
}

// getPostgresSchema returns a postgres schema struct
func (db *DB) getPostgresSchema() (*models.Schema, error) {
	schema := models.NewSchema(Postgres)
	tableNames := TableNames{}
	err := db.DB.Select(&tableNames, `SELECT
	c.relname as "Name"
	FROM pg_catalog.pg_class c
	LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
	WHERE c.relkind IN ('r','')
	AND n.nspname <> 'pg_catalog'
	AND n.nspname <> 'information_schema'
	AND n.nspname !~ '^pg_toast'
	AND pg_catalog.pg_table_is_visible(c.oid);`)
	if err != nil {
		return nil, err
	}
	for _, name := range tableNames {
		table := &models.Table{}
		rows, err := db.DB.Queryx(fmt.Sprintf(`SELECT a.attname as "Field", format_type(a.atttypid, a.atttypmod) AS "Type"
FROM pg_attribute a
JOIN pg_class b ON (a.attrelid = b.relfilenode)
WHERE b.relname = '%s' and a.attstattarget = -1;`, name))

		if err != nil {
			return nil, err
		}
		for rows.Next() {
			tableRow := &models.TableRow{}
			err = rows.StructScan(tableRow)
			if err != nil {
				return nil, err
			}
			table.Rows = append(table.Rows, tableRow)
		}
		table.Name = name
		schema.Tables = append(schema.Tables, table)
	}
	return schema, nil
}
