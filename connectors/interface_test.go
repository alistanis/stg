package connectors

import (
	"testing"

	"github.com/alistanis/stg/generators"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTableNames(t *testing.T) {
	Convey("We can get table names and table descriptions from a mysql server", t, func() {
		config := Config{}
		config.Adapter = MySql
		config.Database = "world"
		config.Username = "ccooper"
		config.Port = 3306
		config.Host = "localhost"
		config.Encoding = "utf8"
		config.Password = ""

		opener, err := GetOpener(config)
		So(err, ShouldBeNil)
		sqlxDb, err := Open(opener)
		So(err, ShouldBeNil)
		db := NewDB(sqlxDb, opener.GetAdapter())
		schema, err := db.GetSchema()
		So(err, ShouldBeNil)
		So(schema, ShouldNotBeNil)
		generators.GenerateModels(schema)

	})

	Convey("We can get table names and table descriptions from a postgres server", t, func() {

		pgConfig := NewPGConfig()
		pgConfig.Adapter = Postgres
		pgConfig.Database = "world"
		pgConfig.Username = "ccooper"
		pgConfig.Port = 5432
		pgConfig.Host = "localhost"
		pgConfig.Encoding = "utf8"
		pgConfig.Password = ""
		pgConfig.SSLMode = "disable"
		sqlxDb, err := pgConfig.Open()
		So(err, ShouldBeNil)
		db := NewDB(sqlxDb, pgConfig.Adapter)
		schema, err := db.GetSchema()
		So(err, ShouldBeNil)
		So(schema, ShouldNotBeNil)
		generators.GenerateModels(schema)

	})

	Convey("We can get table names and table descriptions from sqlite3", t, func() {
		lConf := NewSQLiteConfig()
		lConf.Adapter = SQLite
		lConf.Database = "/Users/ccooper/work/thunderbirds/src/github.com/alistanis/stg/connectors/world.sql"
		sqlxDb, err := lConf.Open()
		So(err, ShouldBeNil)
		db := NewDB(sqlxDb, lConf.Adapter)
		schema, err := db.GetSchema()
		So(err, ShouldBeNil)
		So(schema, ShouldNotBeNil)
		generators.GenerateModels(schema)

	})
}
