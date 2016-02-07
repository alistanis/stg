package connectors

import (
	"fmt"

	"github.com/alistanis/stg/models"
	"github.com/jmoiron/sqlx"
)

// MySQLConfig represents a MySql database configuration struct
type MySQLConfig struct {
	Config   `yaml:",inline"`
	Protocol string `yaml:"protocol" json:"protocol"`
}

// NewMySqlConfig returns a pointer to a MySql configuration struct with UseDefaults set to true
func NewMySqlConfig() *MySQLConfig {
	config := &MySQLConfig{}
	config.SetDefaults()
	return config
}

// MySqlConfigFromConfig returns a pointer to a MySQLConfig with default values
func MySqlConfigFromConfig(config Config) *MySQLConfig {
	mysql := &MySQLConfig{}
	mysql.Config = config
	mysql.SetDefaults()
	return mysql
}

// open assembles a datasource string and attempts to open a connection to a MySQL server
func (m *MySQLConfig) Open() (*sqlx.DB, error) {
	var (
		pass, host, proto = "", "", ""
	)

	if m.Host != "" && m.Host != "localhost" && m.Host != "127.0.0.1" {
		host = fmt.Sprintf("(%s:%d)", m.Host, m.Port)
	}

	if m.Protocol != "" {
		proto = m.Protocol
	} else {
		if host != "" { // if localhost, we default to unix domain sockets
			proto = "tcp"
		}
	}
	if m.Password != "" {
		pass = ":" + m.Password
	}
	dsn := m.Username + pass + "@" + proto + host + "/" + m.Database
	return sqlx.Open(m.Adapter, dsn)
}

// SetDefaults sets default MySQL settings
func (m *MySQLConfig) SetDefaults() {
	if m.Encoding == "" {
		m.Encoding = "utf8"
	}
	if m.Adapter == "" {
		m.Adapter = MySql
	}
	if m.Port == 0 {
		m.Port = 3306
	}
}

// getMysqlSchema returns a mysql schema struct pointer
func (db *DB) getMysqlSchema() (*models.Schema, error) {
	schema := models.NewSchema(MySql)
	tableNames := TableNames{}
	err := db.DB.Select(&tableNames, "show tables;")
	if err != nil {
		return nil, err
	}
	for _, name := range tableNames {
		table := &models.Table{}
		rows, err := db.DB.Queryx(fmt.Sprintf("describe %s;", name))
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
