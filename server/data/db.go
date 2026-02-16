package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var Db *sql.DB

func DbConnection() {
	var err error
	Db, err = sql.Open("sqlite", "./society.db")
	if err != nil {
		fmt.Println("Unable to connect db")
		panic(err)
	}

	Db.SetMaxOpenConns(1)
	if err = Db.Ping(); err != nil {
		log.Fatal("db is not available")
	}

	fmt.Println("connected succesfully")
}
