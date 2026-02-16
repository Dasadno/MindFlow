package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var Db *sql.DB

func DbConnection() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "server/data/society.db" // значение по умолчанию
	}
	Db, err = sql.Open("sqlite", dbPath)
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
