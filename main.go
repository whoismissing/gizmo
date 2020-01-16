package main

import (
	dbutils "github.com/whoismissing/gizmo/dbutils"
	config "github.com/whoismissing/gizmo/config"
	structs "github.com/whoismissing/gizmo/structs"
	"github.com/akamensky/argparse"
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"net/http"
	"os"
	"fmt"
	"time"
	"sync"
)

func parseArgs(parser *argparse.Parser) (string, string) {
	conf := parser.String("i", "input", &argparse.Options{Required: true, Help: "Input config filename"})

	defaultFilename := "gizmo.db"
	dbName := parser.String("o", "output", &argparse.Options{Required: false, Default: defaultFilename, Help: "Output database filename"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	return *conf, *dbName
}

func GetScoreboard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func ConcurrentServiceCheck(servicesPtr *[]structs.Service) {
	var wg sync.WaitGroup
	services := *servicesPtr
	for i:= 0; i < len(services); i++ {
		service := services[i]
		wg.Add(1)
		go service.ServiceCheck.CheckHealth(&services[i], &wg)
	}
	wg.Wait()
}

func main() {
	confName, dbName := parseArgs(argparse.NewParser("gizmo", "Service uptime scoreboard"))

	db, _ := sql.Open("sqlite3", dbName)
	dbutils.CreateDatabase(db)

	teams := config.LoadConfig(confName)
	game := structs.NewGame(teams)

	dbutils.InitializeDatabase(db, game)

	//http.HandleFunc("/", GetScoreboard)
	//http.ListenAndServe(":8080", nil)

	for i := 0; i < len(teams); i++ {
		team := teams[i]

		ConcurrentServiceCheck(&team.Services)

		structs.UpdateTeamCheckCount(&teams[i])
		dbutils.UpdateDatabase(db, game)
	}

	time.Sleep(1 * time.Second)
	// Loop on a five minute timer until next service check
	//time.Sleep(300 * time.Second)

	fmt.Printf("%+v\n", game)
}
