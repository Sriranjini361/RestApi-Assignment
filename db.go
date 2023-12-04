package main

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	DBDriver   = "mysql"
	DBUser     = "root"
	DBPassword = "9994570668@sri"
	DBName     = "app"
)

var db *sql.DB

func init() {

	connStr := fmt.Sprintf("%s:%s@tcp(192.168.246.77:3306)/%s", DBUser, DBPassword, DBName)

	var err error
	db, err = sql.Open(DBDriver, connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}
}
