package connectors

import (
	"errors"
	"fmt"

	"github.com/alistanis/stg/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// Postgres is the postgres string
	Postgres = "postgres"
	// MySql is the mysql string
	MySql = "mysql"
	// SQLite is the sqlite3 string
	SQLite = "sqlite3"
)

// DB wraps a *sqlx.DB and an Adapter string
type DB struct {
	DB      *sqlx.DB
	Adapter string
}

// NewDB returns a new DB wrapper struct
func NewDB(db *sqlx.DB, adapter string) *DB {
	return &DB{DB: db, Adapter: adapter}
}

// GetSchema returns a schema based on the adapter being used
func (db *DB) GetSchema() (*models.Schema, error) {
	switch db.Adapter {
	case MySql:
		return db.getMysqlSchema()
	case Postgres:
		return db.getPostgresSchema()
	case SQLite:
		return db.getSqliteSchema()
	default:
		return nil, errors.New(fmt.Sprintf("Bad adapter string. Should be 'mysql', 'postgres', or 'sqlite3', was: %s", db.Adapter))
	}
}

// TableNames is a list of tables in a schema
type TableNames []string

// Open takes an Opener interface and attempts to open a connection to a database.
func Open(o Opener) (*sqlx.DB, error) {
	return o.Open()
}

// GetOpener returns an Opener interface based on the config given
func GetOpener(config Config) (Opener, error) {
	switch config.Adapter {
	case MySql:
		mysql := MySqlConfigFromConfig(config)
		return mysql, nil
	case Postgres:
		pg := PGConfigFromConfig(config)
		return pg, nil
	case SQLite:
		sqlite := SQLiteConfigFromConfig(config)
		return sqlite, nil
	default:
		return nil, errors.New(fmt.Sprintf("Bad adapter string. Should be 'mysql', 'postgres', or 'sqlite3', was: %s", config.Adapter))
	}
}

// Opener is an interface that can open a connection to a database
type Opener interface {
	// Sets all defaults for the particular type of database configuration
	SetDefaults()
	// Opens a connection to a database
	Open() (*sqlx.DB, error)
	// Returns the Name of the current database configuration
	GetName() string
	// Returns the Adapter of the interface
	GetAdapter() string
}

// Config is the base configuration struct for all database configurations. It is embedded in the more specific config
// types, PGConfig, MySqlConfig, and SQLiteConfig
type Config struct {
	Encoding string `yaml:"encoding" json:"encoding"`
	Adapter  string `yaml:"adapter" json:"adapter"`
	Database string `yaml:"database" json:"database"`
	Host     string `yaml:"host" json:"host"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Port     int    `yaml:"port" json:"port"`
	Socket   string `yaml:"socket" json:"socket"`
	Name     string `yaml:"name" json:"name"`
}

// GetName returns the name of the database configuration.
func (c *Config) GetName() string {
	return c.Name
}

// GetAdapter returns the adapter being used for this configuration.
func (c *Config) GetAdapter() string {
	return c.Adapter
}
