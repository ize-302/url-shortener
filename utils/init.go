// Package utils containint ReadDB  GenerateRandomCode
package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
)

func ReadDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func GenerateRandomCode() string {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
