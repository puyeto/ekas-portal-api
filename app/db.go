package app

import (
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

	// db.DB().SetMaxIdleConns(0)
	// db.DB().SetConnMaxLifetime(time.Second * 10)
	db.DB().SetConnMaxLifetime(0)
	db.DB().SetMaxIdleConns(3)
	db.DB().SetMaxOpenConns(100)

	DBCon = db
	return db
}
