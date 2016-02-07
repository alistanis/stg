package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alistanis/stg/connectors"
	"github.com/alistanis/stg/generators"
	"github.com/jmoiron/sqlx"
)

var (
	path     string
	write    bool
	postgres bool
	mysql    bool
	sqlite   bool
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.StringVar(&path, "path", "", "the path to output source to. if not provided, the current directory is used.")
	flag.BoolVar(&write, "w", false, "write to the current directory, otherwise show output")
	flag.BoolVar(&postgres, "pg", false, "use postgres")
	flag.BoolVar(&mysql, "mysql", false, "use mysql")
	flag.Parse()
	var sqlxDb *sqlx.DB
	adapter := ""
	if postgres {
		pgConfig := connectors.NewPGConfig()
		pgConfig.Adapter = connectors.Postgres
		pgConfig.Database = "world"
		pgConfig.Username = "ccooper"
		pgConfig.Port = 5432
		pgConfig.Host = "localhost"
		pgConfig.Encoding = "utf8"
		pgConfig.Password = ""
		pgConfig.SSLMode = "disable"
		var err error
		sqlxDb, err = pgConfig.Open()
		if err != nil {
			fmt.Println(err)
			return -1
		}
		adapter = pgConfig.Adapter
	} else if mysql {
		config := connectors.Config{}
		config.Adapter = connectors.MySql
		config.Database = "world"
		config.Username = "ccooper"
		config.Port = 3306
		config.Host = "localhost"
		config.Encoding = "utf8"
		config.Password = ""
		opener, err := connectors.GetOpener(config)
		if err != nil {
			fmt.Println(err)
			return -1
		}
		sqlxDb, err = connectors.Open(opener)
		if err != nil {
			fmt.Println(err)
			return -1
		}
		adapter = connectors.MySql
	} else {
		lConf := connectors.NewSQLiteConfig()
		lConf.Adapter = connectors.SQLite
		lConf.Database = "/Users/ccooper/work/thunderbirds/src/github.com/alistanis/stg/connectors/world.sql"
		var err error
		sqlxDb, err = lConf.Open()
		if err != nil {
			fmt.Println(err)
			return -1
		}
		adapter = lConf.Adapter
	}

	db := connectors.NewDB(sqlxDb, adapter)
	schema, err := db.GetSchema()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	if write {
		if path != "" {
			generators.Write(schema, path)
		} else {
			binDir, binErr := filepath.Abs(filepath.Dir(os.Args[0]))
			if binErr != nil {
				fmt.Println(err)
				return -1
			}
			generators.Write(schema, binDir)
		}
	} else {
		modelsSrc, err := generators.GenerateModels(schema)
		if err != nil {
			fmt.Println(err)
			return -1
		}
		for _, src := range modelsSrc {
			fmt.Println(string(src))
		}
	}

	return 0
}
