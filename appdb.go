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

	sqlURLSel = "SELECT id FROM urls_tab WHERE url=?;"
	sqlURLIns = `INSERT INTO urls_tab (url, interval) VALUES (?, ?);`
	sqlURLDel = `DELETE FROM urls_tab WHERE ID=?;`
	sqlURLUpd = `UPDATE urls_tab SET interval=? WHERE url=?;`

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
func InsertRow(db *sql.DB, sqlSel, sqlIns, sqlUpd string, item ItemAdd) (int, error) {
	// Check if URL is already in DB and eventually return its ID
	rows, err := db.Query(sqlSel, item.URL)
	if err != nil {
		return -1, fmt.Errorf("could not exec sql select query {%+v}", err)
	}
	defer rows.Close()

	id := -1
	for rows.Next() {
		rows.Scan(&id)
	}

	if id == -1 { // Add new entry
		if _, err := db.Exec(sqlIns, item.URL, item.Interval); err != nil {
			return -1, fmt.Errorf("could not add new db entry {%+v}", err)
		}
		// Check what ID was set
		rows, err := db.Query(sqlSel, item.URL)
		if err != nil {
			return -1, fmt.Errorf("could not exec sql select query {%+v}", err)
		}
		defer rows.Close()

		id = -1
		for rows.Next() {
			rows.Scan(&id)
		}
	} else { // Update old entry
		fmt.Printf("Update {%s}", sqlUpd)
		if _, err := db.Exec(sqlUpd, item.Interval, item.URL); err != nil {
			return -1, fmt.Errorf("could not add new db entry {%+v}", err)
		}
	}
	return id, nil
}

// DeleteRow removes entry from db by ID
func DeleteRow(db *sql.DB, sqlDel string, id int) (int, error) {
	if _, err := db.Exec(sqlDel, id); err != nil {
		fmt.Printf("Colud not delete row - %s\n", err.Error())
		return -1, fmt.Errorf("could not remove db entry {%+v}", err)
	}
	return 0, nil
}

// UpdateRow update database entry
func UpdateRow(db *sql.DB, sqlUpd string, rd ReqData) (int, error) {
	if _, err := db.Exec(sqlUpd, rd.URL, rd.Interval); err != nil {
		return -1, fmt.Errorf("could not removupdate db entry {%+v}", err)
	}
	return 0, nil
}
