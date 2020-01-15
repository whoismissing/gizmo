package dbutils

// TODO: Implement gameID usage in InitializeTeams

import (
	structs "github.com/whoismissing/gizmo/structs"
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"math/rand"
	"log"
	"fmt"
)

func createGameTable(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS "Game" (
		"GameStartTime"		INTEGER,
		"CurrentGameTime"	INTEGER,
		"GameID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT
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
		"TeamID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
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
		"Name"	TEXT NOT NULL,
		"TeamID"	INTEGER,
		"ServiceID"	INTEGER PRIMARY KEY AUTOINCREMENT,
		"HostIP"	TEXT,
		"NumberOfMissedChecks"	INTEGER,
		"NumberOfChecks"	INTEGER
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
		"StatusID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Time"	INTEGER,
		"State"	TEXT
	);
	`)

	if err != nil {
		log.Fatal(err)
	}

	return
}

func initializeService(db *sql.DB, service structs.Service, serviceID int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO "Service" (
		"Name",
		"TeamID",
		"ServiceID",
		"HostIP",
		"NumberOfMissedChecks",
		"NumberOfChecks"
	) VALUES ( ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("service = +%v\n", service)
	name := service.Name
	teamID := service.TeamID
	hostIP := service.HostIP
	missedChecks := service.ChecksMissed
	totalChecks := service.ChecksAttempted
	_, err = stmt.Exec(name, teamID, serviceID, hostIP, missedChecks, totalChecks)
	if err != nil {
		log.Fatal(err)
	}

	_ = tx.Commit()
	stmt.Close()

	return
}

func initializeServices(db *sql.DB, services []structs.Service) {

	for i := 0; i < len(services); i++ {
		initializeService(db, services[i], i)
	}

	return
}

func initializeTeam(db *sql.DB, team structs.Team) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO "Team" (
		"GameID",
		"TeamID",
		"TotalMissedChecks",
		"TotalChecks"
	) VALUES ( ?, ?, ?, ?);
	`)

	teamID := team.TeamID
	_, err = stmt.Exec(0, teamID, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	_ = tx.Commit()
	stmt.Close()

	services := team.Services
	initializeServices(db, services)

	return
}

func initializeTeams(db *sql.DB, teams []structs.Team) {
	for i := 0; i < len(teams); i++ {
		initializeTeam(db, teams[i])
	}
}

func initializeGame(db *sql.DB, game structs.Game) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO "Game" (
		"GameStartTime",
		"CurrentGameTime",
		"GameID"
	) VALUES ( ?, ?, ? );
	`)
	if err != nil {
		log.Fatal(err)
	}

	gameStartTime := game.StartTime.Unix()
	gameID := rand.Int()
	_, err = stmt.Exec(gameStartTime, gameStartTime, gameID)
	if err != nil {
		log.Fatal(err)
	}

	_ = tx.Commit()
	stmt.Close()

	teams := game.Teams
	initializeTeams(db, teams)

	return
}

func CreateDatabase(db *sql.DB) {
	createGameTable(db)
	createTeamTable(db)
	createServiceTable(db)
	createStatusTable(db)
}

func InitializeDatabase(db *sql.DB, game structs.Game) {
	initializeGame(db, game)
}
