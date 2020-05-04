package app

import (
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// DBCon ...
var DBCon *dbx.DB

// SecondDBCon ...
var SecondDBCon *dbx.DB

// InitializeDB initialize DB conn
func InitializeDB(dns string) *dbx.DB {
	db, err := dbx.MustOpen("mysql", dns)
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(0)
	db.DB().SetConnMaxLifetime(time.Second * 10)

	DBCon = db
	return db
}
