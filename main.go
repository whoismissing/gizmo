package main

import (
	"github.com/akamensky/argparse"
	_ "github.com/mattn/go-sqlite3"
	dbutils "github.com/whoismissing/gizmo/dbutils"
	config "github.com/whoismissing/gizmo/config"
	structs "github.com/whoismissing/gizmo/structs"

	"database/sql"
	"os"
	"fmt"
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

func main() {
	confName, dbName := parseArgs(argparse.NewParser("gizmo", "Service uptime scoreboard"))

	db, _ := sql.Open("sqlite3", dbName)
	dbutils.InitializeDatabase(db)

	teams := config.LoadConfig(confName)
	game := structs.NewGame(teams)

	// TODO: write code to insert values of game into SQL database

	for i := 0; i < len(teams); i++ {
		team := teams[i]

		// Concurrent service checks
		var wg sync.WaitGroup
		for j:= 0; j < len(team.Services); j++ {
			service := team.Services[j]
			wg.Add(1)
			go service.ServiceCheck.CheckHealth(j, &wg)
		}
		wg.Wait()

	}

	// Loop on a five minute timer until next service check
	time.Sleep(300 * time.Second)

	fmt.Printf("%+v\n", game)

}
