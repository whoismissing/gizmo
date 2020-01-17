package main

import (
	dbutils "github.com/whoismissing/gizmo/dbutils"
	config "github.com/whoismissing/gizmo/config"
	structs "github.com/whoismissing/gizmo/structs"
	web "github.com/whoismissing/gizmo/web"
	"github.com/akamensky/argparse"
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"net/http"
	"os"
	"fmt"
	"log"
	"time"
	"math/rand"
	"sync"
)

var scoreboardHTML string

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
	fmt.Fprintf(w, "%s", scoreboardHTML)
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

	scoreboardHTML = web.GenerateScoreboardHTML(teams)

	// Spin off separate thread for the web server so as not to block main
	go func() {
		fmt.Println("web server listening on :8080")
		http.HandleFunc("/", GetScoreboard)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Loop every three to five minutes until next service check
	min := 180 // 180 seconds = 3 minutes
	max := 300 // 300 seconds = 5 minutes
	for {

		for i := 0; i < len(teams); i++ {
			team := teams[i]

			ConcurrentServiceCheck(&team.Services)
			scoreboardHTML = web.GenerateScoreboardHTML(teams)

			structs.UpdateTeamCheckCount(&teams[i])
			dbutils.UpdateDatabase(db, game)
		}

		fmt.Println("=======================")
		sleepTime := rand.Intn(max - min) + min
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}
