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

	DBCon = db
	return db
}

// InitializeSecondDB connect to second db
func InitializeSecondDB(dns string) *dbx.DB {
	db, err := dbx.MustOpen("mysql", dns)
	if err != nil {
		panic(err)
	}

	SecondDBCon = db
	return db
}
