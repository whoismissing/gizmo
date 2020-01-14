package main

import (
	"database/sql"
	"github.com/akamensky/argparse"
	_ "github.com/mattn/go-sqlite3"
	dbutils "github.com/whoismissing/gizmo/dbutils"
	config "github.com/whoismissing/gizmo/config"
	structs "github.com/whoismissing/gizmo/structs"

	"os"
	"fmt"
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

	fmt.Printf("%+v\n", game)

}
