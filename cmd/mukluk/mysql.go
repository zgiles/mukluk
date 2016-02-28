package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func mysqlStart(dbconfig mysqlconfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbconfig.Connectstring)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
	// defer db.Close()
}
