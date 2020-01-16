package dbutils

import (
	structs "github.com/whoismissing/gizmo/structs"
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"log"
	"time"
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
		services[i].ServiceID = i
		initializeService(db, services[i], i)
	}

	return
}

func initializeTeam(db *sql.DB, team structs.Team, gameID int) {
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
	_, err = stmt.Exec(gameID, teamID, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	_ = tx.Commit()
	stmt.Close()

	services := team.Services
	initializeServices(db, services)

	return
}

func initializeTeams(db *sql.DB, teams []structs.Team, gameID int) {
	for i := 0; i < len(teams); i++ {
		initializeTeam(db, teams[i], gameID)
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
	gameID := game.GameID
	_, err = stmt.Exec(gameStartTime, gameStartTime, gameID)
	if err != nil {
		log.Fatal(err)
	}

	_ = tx.Commit()
	stmt.Close()

	teams := game.Teams
	initializeTeams(db, teams, gameID)

	return
}

func updateGameTable(db *sql.DB, game structs.Game) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
	UPDATE "Game" SET "CurrentGameTime"=?
	WHERE "GameID"=?;
	`)
	if err != nil {
		log.Fatal(err)
	}

	currentGameTime := time.Now().Unix()
	gameID := game.GameID
	_, err = stmt.Exec(currentGameTime, gameID)
	if err != nil {
		log.Fatal(err)
	}

	_ = tx.Commit()
	stmt.Close()

	return
}

func updateServiceTable(db *sql.DB, service structs.Service) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
	UPDATE "Service" SET 
		"NumberOfMissedChecks"=?,
		"NumberOfChecks"=?
	WHERE "ServiceID"=?;
	`)
	if err != nil {
		log.Fatal(err)
	}

	serviceID := service.ServiceID
	_, err = stmt.Exec(service.ChecksMissed, service.ChecksAttempted, serviceID)
	if err != nil {
		log.Fatal(err)
	}

	_ = tx.Commit()
	stmt.Close()

	return

}

func updateStatusTable(db *sql.DB, service structs.Service) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO "Status" (
		"ServiceID",
		"Time",
		"State"
	) VALUES ( ?, ?, ?);
	`)
	if err != nil {
		log.Fatal(err)
	}

	serviceID := service.ServiceID
	top := len(service.PrevStatuses) - 1

	/* if service.PrevStatuses is empty */
	if top < 0 {
		return
	}

	status := service.PrevStatuses[top]
	_, err = stmt.Exec(serviceID, status.Time, status.Status)
	if err != nil {
		log.Fatal(err)
	}

	_ = tx.Commit()
	stmt.Close()

	return

}


func updateServices(db *sql.DB, services []structs.Service) {
	for i := 0; i < len(services); i++ {
		updateServiceTable(db, services[i])
		updateStatusTable(db, services[i])
	}
}

func updateTeamTable(db *sql.DB, team structs.Team) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
	UPDATE "Team" SET 
		"TotalMissedChecks"=?,
		"TotalChecks"=?
	WHERE "TeamID"=?;
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(team.TotalChecksMissed, team.TotalChecksAttempted, team.TeamID)
	if err != nil {
		log.Fatal(err)
	}

	_ = tx.Commit()
	stmt.Close()

	updateServices(db, team.Services)

	return

}

func updateTeams(db *sql.DB, teams []structs.Team) {
	for i := 0; i < len(teams); i++ {
		updateTeamTable(db, teams[i])
	}
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

func UpdateDatabase(db *sql.DB, game structs.Game) {
	updateGameTable(db, game)
	updateTeams(db, game.Teams)
}
