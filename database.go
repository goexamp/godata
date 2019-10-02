package database

import (
	"database/sql"
	"log"
)

var db *sql.DB
var err error

func connect_db() {
	db, err = sql.Open("mysql", "root:10184902125410@/golang_db")

	if err != nil {
		log.Fatalln(err)
		log.Printf("Server:database wrong login and mysql password")
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}
