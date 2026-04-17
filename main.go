package main

import (
	"database/sql"
	"log"

	"ize-302/url-shortener/routes"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		url TEXT,
		code TEXT
	)
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Table 'urls' created successfully")

	routes.Requests()
}
