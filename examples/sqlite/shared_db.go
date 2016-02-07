package sqlite

import "github.com/jmoiron/sqlx"

var (
	db *sqlx.DB
)

func SetDb(d *sqlx.DB) {
	db = d
}
