package gizmodbutils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func createGameTable(db *sql.DB) {
	_, _ = db.Exec(`
	CREATE TABLE IF NOT EXISTS "Game" (
		"GameStartTime"		INTEGER,
		"CurrentGameTime"	INTEGER,
		"GameID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE
	);
	`)

	return
}

func createTeamTable(db *sql.DB) {
	_, _ = db.Exec(`
	CREATE TABLE IF NOT EXISTS "Team" (
		"GameID"	INTEGER,
		"TeamID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		"TotalMissedChecks"	INTEGER,
		"TotalChecks"	INTEGER
	);
	`)

	return
}

func createServiceTable(db *sql.DB) {
	_, _ = db.Exec(`
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

	return
}

func createStatusTable(db *sql.DB) {
	_, _ = db.Exec(`
	CREATE TABLE IF NOT EXISTS "Status" (
		"ServiceID"	INTEGER,
		"StatusID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		"Time"	INTEGER,
		"State"	TEXT
	);
	`)

	return
}

func InitializeDatabase(db *sql.DB) {
	createGameTable(db)
	createTeamTable(db)
	createServiceTable(db)
	createStatusTable(db)
}
