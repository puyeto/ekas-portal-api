package app

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// DBCon ...
var DBCon *dbx.DB

// InitializeDB initialize DB conn
func InitializeDB(dns string) *dbx.DB {
	db, err := dbx.MustOpen("mysql", dns)
	if err != nil {
		panic(err)
	}

	DBCon = db
	return db
}
