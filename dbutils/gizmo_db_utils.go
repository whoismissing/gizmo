package dbutils

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func createGameTable(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS "Game" (
		"GameStartTime"		INTEGER,
		"CurrentGameTime"	INTEGER,
		"GameID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE
	);
	`)

	if err != nil {
		log.Fatal(err)
	}

	return
}

func createTeamTable(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS "Team" (
		"GameID"	INTEGER,
		"TeamID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		"TotalMissedChecks"	INTEGER,
		"TotalChecks"	INTEGER
	);
	`)

	if err != nil {
		log.Fatal(err)
	}

	return
}

func createServiceTable(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS "Service" (
		"Name"	TEXT NOT NULL UNIQUE,
		"TeamID"	INTEGER,
		"ServiceID"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
		"HostIP"	TEXT,
		"NumberOfMissedChecks"	INTEGER,
		"NumberOfChecks"	INTEGER,
		"User"	TEXT,
		"Password"	TEXT,
		"Domain"	TEXT
	);
	`)

	if err != nil {
		log.Fatal(err)
	}

	return
}

func createStatusTable(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS "Status" (
		"ServiceID"	INTEGER,
		"StatusID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		"Time"	INTEGER,
		"State"	TEXT
	);
	`)

	if err != nil {
		log.Fatal(err)
	}

	return
}

func InitializeDatabase(db *sql.DB) {
	createGameTable(db)
	createTeamTable(db)
	createServiceTable(db)
	createStatusTable(db)
}
