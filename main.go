package main

import (
	//"database/sql"
	//db_utils "github.com/whoismissing/gizmo/dbutils"

	//_ "github.com/mattn/go-sqlite3"
	structs "github.com/whoismissing/gizmo/structs"
	"fmt"
	"reflect"
)

func main() {
	//db, _ := sql.Open("sqlite3", "./gizmo.db")
	//db_utils.InitializeDatabase(db)

	game := structs.NewGame(nil) // {teams: nil, time: 1}
	fmt.Printf("game = %+v\n", game)

	teams := structs.NewTeams(5)
	teams[0].Services = make([]structs.Service, 5)
	teams[1].Services = make([]structs.Service, 3)
	fmt.Printf("teams = %+v\n", teams)

	www := structs.WebService{URL: "google.com"}
	fmt.Printf("www = %+v\n", www)

	fmt.Println("type ", reflect.TypeOf(www))

	service := structs.NewService("www", 69)
	service.ServiceType = www
	fmt.Printf("service = %+v\n", service)
	fmt.Println("type ", reflect.TypeOf(www))

	/*
	team1 := structs.NewTeam(69)
	fmt.Printf("team = %+v\n", team1)

	service := structs.NewService("www", 69)
	fmt.Printf("service = %+v\n", service)

	status := structs.NewStatus(23, false)
	fmt.Printf("status = %+v\n", status)
	*/
}
