package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Definition of DB queries
const (
	// URLs table
	sqlURLCreate = `CREATE TABLE IF NOT EXISTS urls_tab (
		id INTEGER PRIMARY KEY,
		url TEXT NOT NULL,
		interval INTEGER NOT NULL);`

	sqlURLSelMaxID = "SELECT MAX(id) FROM urls_tab;"
	sqlURLIns      = `INSERT INTO urls_tab (url, interval) VALUES (?, ?);`

	sqlURLHistCreate = `CREATE TABLE IF NOT EXISTS urls_history_tab (
		id INTEGER NOT NULL,
		response TEXT NOT NULL,
		duration REAL NOT NULL,
		created_at TEXT NOT NULL);`
)

// ReqData is a placeholder for incoming request data
type ReqData struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}

// InitDb initilize SQLite DB
func InitDb() (*sql.DB, error) {
	var db *sql.DB
	const dbPath = "/tmp/app_db.db"
	// Check whether db exists
	if _, err := os.Stat(dbPath); err == nil { // when err == nil db exists
		if db, err = sql.Open("sqlite3", dbPath); err != nil {
			return db, fmt.Errorf("could not open database connection")
		}
	} else if os.IsNotExist(err) { // db is not available, create a new one
		if db, err = sql.Open("sqlite3", dbPath); err != nil {
			return db, fmt.Errorf("could not open database connection")
		}
		defer db.Close()
		// Create tables: 1) URLs, 2) URLs' history
		if _, err = db.Exec(sqlURLCreate); err != nil {
			return db, fmt.Errorf("could not create urls table")
		}
		if _, err = db.Exec(sqlURLHistCreate); err != nil {
			return db, fmt.Errorf("could not create history urls table")
		}
	}
	return db, nil
}

// InsertRow adds a new entry into URLs table
func InsertRow(db *sql.DB, s string, rd ReqData) (int, error) {
	if _, err := db.Exec(s, rd.URL, rd.Interval); err != nil {
		return -1, fmt.Errorf("could not add new db entry {%+v}", err)
	}

	id := -1
	err := db.QueryRow(sqlURLSelMaxID).Scan(&id)
	switch err {
	case nil:
		return id, nil
	case sql.ErrNoRows:
		return -1, fmt.Errorf("could not get index info")
	default:
		return -1, fmt.Errorf("could not finalize add operation")
	}
}
