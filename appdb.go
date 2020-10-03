package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DbPath is a path to SQLite DB

// Definition of DB queries
const (
	// URLs table
	sqlURLCreate = `CREATE TABLE IF NOT EXISTS urls_tab (
		id INTEGER PRIMARY KEY,
		url TEXT NOT NULL,
		interval INTEGER NOT NULL);`

	sqlURLHistCreate = `CREATE TABLE IF NOT EXISTS urls_history_tab (
		id INTEGER NOT NULL,
		response TEXT NOT NULL,
		duration REAL NOT NULL,
		created_at TEXT NOT NULL);`
)

// InitDb initilize SQLite DB
func InitDb() (*sql.DB, error) {
	var db *sql.DB
	const dbPath = "/tmp/app_db.db"
	// Check whether db exists
	if _, err := os.Stat(dbPath); err == nil { // when err == nil db exists
		if db, err = sql.Open("sqlite3", dbPath); err != nil {
			return db, fmt.Errorf("could not open database connection")
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

/*
	} else if os.IsNotExist(err) {
		//fmt.Printf("Database not exists: {dbPath=%s}\n", dbPath)
		//dbHand, err := sql.Open("sqlite3", ":memory:")
		if db, err = sql.Open("sqlite3", dbPath); err != nil {
			fmt.Println(err)
			fmt.Println("Could not open DB")
			return nil
		}

		// Create database structure
		if _, err = db.Exec(sqlURLCrt); err != nil {
			fmt.Println(err)
			fmt.Println("Could not create DB structure (urls_tab)")
			db.Close()
			return nil
		}

		if _, err = db.Exec(sqlHisCrt); err != nil {
			fmt.Println(err)
			fmt.Println("Could not create DB structure (history_tab)")
			db.Close()
			return nil
		}

	} else {
		// Impossible case...?
		fmt.Println("Unsupported case...")
		return nil
	}
	return db
}

// CloseDb - Function close connection to database
// @param db - database handler
//
func CloseDb(db *sql.DB) {
	if db != nil {
		db.Close()
	}

}

*/
