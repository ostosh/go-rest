package example

import (
	"github.com/jmoiron/sqlx"
)

var (
	conn *sqlx.DB
)

//Return active sql connection
func GetConnection() *sqlx.DB {
	return conn
}

//Set active sql connection
func SetConnection(db *sqlx.DB) {
	conn = db
}
