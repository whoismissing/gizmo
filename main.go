package main

import (
	//"database/sql"
	//db_utils "github.com/whoismissing/gizmo/dbutils"

	//_ "github.com/mattn/go-sqlite3"
	structs "github.com/whoismissing/gizmo/structs"
	"fmt"
)

func main() {
	//db, _ := sql.Open("sqlite3", "./gizmo.db")
	//db_utils.InitializeDatabase(db)

	game := structs.NewGame(22) // {teams: nil, time: 1}
	fmt.Printf("game = %+v\n", game)

	teams := make([]structs.Team, 5)
	teams[0].Services = make([]structs.Service, 5)
	fmt.Printf("teams = %+v\n", teams[0])

	/*
	team1 := structs.NewTeam(69)
	fmt.Printf("team = %+v\n", team1)

	service := structs.NewService("www", 69)
	fmt.Printf("service = %+v\n", service)

	status := structs.NewStatus(23, false)
	fmt.Printf("status = %+v\n", status)
	*/
}
