package dbutils

import (
	structs "github.com/whoismissing/gizmo/structs"

	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
)

var test_db = "./test_db.db"

func TestInitializeDatabaseSuccess(t *testing.T) {
	t.Logf("Test InitializeDatabase SUCCESS")

	db, err := sql.Open("sqlite3", test_db)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database exists?: ", DoesDatabaseExist(db))

	CreateDatabase(db)

	fmt.Println("Database exists?: ", DoesDatabaseExist(db))

	game := structs.NewGame(nil)
	InitializeDatabase(db, game)

	fmt.Println("before: ", game)
	UpdateGameFromDatabase(db, &game)
	fmt.Println("after: ", game)

	err = os.Remove(test_db)

	if err != nil {
		log.Fatal(err)
	}
}

func TestLoadTeamChecksFromDatabase(t *testing.T) {
	db, _ := sql.Open("sqlite3", "gizmo.db")
	team := structs.NewTeam(1)
	fmt.Println("before: ", team)

	/*
	   LoadTeamChecksFromDatabase(db, &team)
	   fmt.Println("after: ", team)

	   team.Services = structs.NewServices(1, team.TeamID)
	   team.Services[0].ServiceCheck = structs.WebService{URL:"team1.local"}
	   team.Services[0].Name = "www"

	   service := team.Services[0]
	   fmt.Println("before: ", service)
	   LoadServiceFromDatabase(db, &service)
	   fmt.Println("after: ", service)
	*/

	teams := []structs.Team{team}
	game := structs.NewGame(teams)
	fmt.Println("before: ", game)
	LoadGameFromDatabase(db, &game)
	fmt.Println("after: ", game)
}
