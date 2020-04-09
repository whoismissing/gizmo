package dbutils

import (
	structs "github.com/whoismissing/gizmo/structs"
    debug "github.com/whoismissing/gizmo/debug"
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"log"
	"time"
    "fmt"
)

func createGameTable(db *sql.DB) {
    debug.LogBegin()
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

    debug.LogEnd()
	return
}

func createTeamTable(db *sql.DB) {
    debug.LogBegin()
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

    debug.LogEnd()
	return
}

func createServiceTable(db *sql.DB) {
    debug.LogBegin()
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS "Service" (
		"Name"	TEXT NOT NULL,
		"TeamID"	INTEGER,
		"ServiceID"	INTEGER AUTO_INCREMENT,
		"HostIP"	TEXT,
		"NumberOfMissedChecks"	INTEGER,
		"NumberOfChecks"	INTEGER
	);
	`)

	if err != nil {
		log.Fatal(err)
	}

    debug.LogEnd()
	return
}

func createStatusTable(db *sql.DB) {
    debug.LogBegin()
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS "Status" (
		"ServiceID"	INTEGER,
        "TeamID" INTEGER,
		"StatusID"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Time"	INTEGER,
		"State"	TEXT
	);
	`)

	if err != nil {
		log.Fatal(err)
	}

    debug.LogEnd()
	return
}

func insertNewTeam(db *sql.DB, team structs.Team) {
    debug.LogBegin()
    sqlStatement := `SELECT GameID FROM GAME`

    row := db.QueryRow(sqlStatement)
    var gameID int

    switch err := row.Scan(&gameID); err {
    case sql.ErrNoRows:
        fmt.Println("insertNewTeam: No rows returned")
    case nil: // success!
        initializeTeam(db, team, gameID)
    default:
        panic(err)
    }

    debug.LogEnd()
}

func insertNewService(db *sql.DB, service structs.Service) {
    debug.LogBegin()
    sqlStatement := `SELECT MAX(ServiceID) FROM Service`

    row := db.QueryRow(sqlStatement)
    var lastServiceID int

    switch err := row.Scan(&lastServiceID); err {
    case sql.ErrNoRows:
        fmt.Println("insertNewService: No rows returned")
    case nil: // success!
        initializeService(db, service, lastServiceID + 1)
    default:
        panic(err)
    }

    debug.LogEnd()
}

func initializeService(db *sql.DB, service structs.Service, serviceID int) {
    debug.LogBegin()
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

    debug.LogEnd()
	return
}

func initializeServices(db *sql.DB, services []structs.Service) {
    debug.LogBegin()

	for i := 0; i < len(services); i++ {
		services[i].ServiceID = i
		initializeService(db, services[i], i)
	}

    debug.LogEnd()
	return
}

func initializeTeam(db *sql.DB, team structs.Team, gameID int) {
    debug.LogBegin()

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

    debug.LogEnd()
	return
}

func initializeTeams(db *sql.DB, teams []structs.Team, gameID int) {
    debug.LogBegin()

	for i := 0; i < len(teams); i++ {
		initializeTeam(db, teams[i], gameID)
	}

    debug.LogEnd()
}

func initializeGame(db *sql.DB, game structs.Game) {
    debug.LogBegin()

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

    debug.LogEnd()
	return
}

func updateGameTable(db *sql.DB, game structs.Game) {
    debug.LogBegin()

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

    debug.LogEnd()
	return
}

func updateServiceTable(db *sql.DB, service structs.Service) {
    debug.LogBegin()

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

    debug.LogEnd()
	return

}

func updateStatusTable(db *sql.DB, service structs.Service) {
    debug.LogBegin()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`
	INSERT INTO "Status" (
		"ServiceID",
        "TeamID",
		"Time",
		"State"
	) VALUES ( ?, ?, ?, ?);
	`)
	if err != nil {
		log.Fatal(err)
	}

	serviceID := service.ServiceID
    teamID := service.TeamID
	top := len(service.PrevStatuses) - 1

	/* if service.PrevStatuses is empty */
	if top < 0 {
		return
	}

	status := service.PrevStatuses[top]
	_, err = stmt.Exec(serviceID, teamID, status.Time, status.Status)
	if err != nil {
		log.Fatal(err)
	}

	_ = tx.Commit()
	stmt.Close()

    debug.LogEnd()
	return

}


func updateServices(db *sql.DB, services []structs.Service) {
    debug.LogBegin()

	for i := 0; i < len(services); i++ {
		updateServiceTable(db, services[i])
		updateStatusTable(db, services[i])
	}

    debug.LogEnd()
}

func updateTeamTable(db *sql.DB, team structs.Team) {
    debug.LogBegin()

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

    debug.LogEnd()
	return
}

func updateTeams(db *sql.DB, teams []structs.Team) {
    debug.LogBegin()

	for i := 0; i < len(teams); i++ {
		updateTeamTable(db, teams[i])
	}

    debug.LogEnd()
}

func countTablesFromSQLMaster(db *sql.DB) int64 {
    debug.LogBegin()

    sqlStatement := `SELECT COUNT(*) FROM sqlite_master`
    row := db.QueryRow(sqlStatement)
    var count int64

    switch err := row.Scan(&count); err {
    case sql.ErrNoRows:
        fmt.Println("countTablesFromSQLMaster: No rows returned")
    case nil: // success!
        debug.LogEnd()
        return count
    default:
        panic(err)
    }

    debug.LogEnd()
    return -1
}

func CreateDatabase(db *sql.DB) {
    debug.LogBegin()

	createGameTable(db)
	createTeamTable(db)
	createServiceTable(db)
	createStatusTable(db)

    debug.LogEnd()
}

func InitializeDatabase(db *sql.DB, game structs.Game) {
    debug.LogBegin()
	initializeGame(db, game)
    debug.LogEnd()
}

func UpdateDatabase(db *sql.DB, game structs.Game) {
    debug.LogBegin()
	updateGameTable(db, game)
	updateTeams(db, game.Teams)
    debug.LogEnd()
}

func DoesDatabaseExist(db *sql.DB) bool {
    debug.LogBegin()
    if countTablesFromSQLMaster(db) < 4 {
        debug.LogEnd()
        return false
    }
    debug.LogEnd()
    return true
}

func UpdateGameFromDatabase(db *sql.DB, game *structs.Game) {
    debug.LogBegin()

    sqlStatement := `SELECT GameID, GameStartTime FROM Game`

    row := db.QueryRow(sqlStatement)
    var unix_time int64

    switch err := row.Scan(&game.GameID, &unix_time); err {
    case sql.ErrNoRows:
        fmt.Println("UpdateGameFromDatabase: No rows returned")
    case nil: // success!
        game.StartTime = time.Unix(unix_time, 0)
    default:
        panic(err)
    }

    debug.LogEnd()
}

func LoadTeamChecksFromDatabase(db *sql.DB, team *structs.Team) {
    debug.LogBegin()

    sqlStatement := "SELECT TotalMissedChecks, TotalChecks FROM Team WHERE TeamID=?"
    row := db.QueryRow(sqlStatement, team.TeamID)

    var totalChecksMissed int
    var totalChecksAttempted int
    switch err := row.Scan(&totalChecksMissed, &totalChecksAttempted); err {
    case sql.ErrNoRows:
        fmt.Println("LoadTeamChecksFromDatabase: No rows returned")
        insertNewTeam(db, *team)
    case nil: // success!
        team.TotalChecksMissed = uint(totalChecksMissed)
        team.TotalChecksHit = uint(totalChecksAttempted) - uint(totalChecksMissed)
        team.TotalChecksAttempted = uint(totalChecksAttempted)
    default:
    }

    debug.LogEnd()
}

func LoadServiceFromDatabase(db *sql.DB, service *structs.Service) {
    debug.LogBegin()

    sqlStatement := "SELECT ServiceID, NumberOfMissedChecks, NumberOfChecks FROM Service WHERE Name=?"
    row := db.QueryRow(sqlStatement, service.Name)

    var serviceID int
    var checksMissed int
    var checksAttempted int
    switch err := row.Scan(&serviceID, &checksMissed, &checksAttempted); err {
    case sql.ErrNoRows:
        fmt.Println("LoadServiceFromDatabase: No rows returned")
        insertNewService(db, *service)
    case nil: // success!
        service.ServiceID = serviceID
        service.ChecksAttempted = uint(checksAttempted)
        service.ChecksMissed = uint(checksMissed)
        service.ChecksHit = uint(checksAttempted) - uint(checksMissed)
    default:
    }

    debug.LogEnd()
}

func LoadGameFromDatabase(db *sql.DB, game *structs.Game) {
    debug.LogBegin()

    for i := 0; i < len(game.Teams); i++ {
        LoadTeamChecksFromDatabase(db, &game.Teams[i])
        for j := 0; j < len(game.Teams[i].Services); j++ {
            LoadServiceFromDatabase(db, &game.Teams[i].Services[j])
        }
    }

    debug.LogEnd()
}
