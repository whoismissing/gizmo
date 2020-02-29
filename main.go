package main

import (
	dbutils "github.com/whoismissing/gizmo/dbutils"
	config "github.com/whoismissing/gizmo/config"
	structs "github.com/whoismissing/gizmo/structs"
	web "github.com/whoismissing/gizmo/web"
	"github.com/akamensky/argparse"
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
    "encoding/json"
    "net"
	"net/http"
	"os"
    "path/filepath"
	"fmt"
	"log"
	"time"
	//"math/rand"
	"sync"
)

var scoreboardHTML string

func parseArgs(parser *argparse.Parser) (string, string, string) {
	conf := parser.String("i", "input", &argparse.Options{Required: true, Help: "Input config filename"})

	defaultFilename := "gizmo.db"
	dbName := parser.String("o", "output", &argparse.Options{Required: false, Default: defaultFilename, Help: "Output database filename"})

	scriptDir := parser.String("s", "script_directory", &argparse.Options{Required: false, Default: "", Help: "Script directory"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	return *conf, *dbName, *scriptDir
}

// shamelessly ripped from https://golangcode.com/check-if-a-file-exists/
func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func isServiceAlreadyInList(team structs.Team, service structs.Service) bool {
    for i := 0; i < len(team.Services); i++ {
        if service.Name == team.Services[i].Name {
            return true
        }
    }
    return false
}

func addScriptServicesToTeam(db *sql.DB, directory string, team *structs.Team) {
    /* list all files (recursively) in the script directory */
    err := filepath.Walk(directory,
        func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if fileExists(path) {
            // TODO: Check for service in database
            // TODO: Insert into database if service does NOT exist
            ext := structs.ExternalService{ ScriptPath: path }
            extJson, _ := json.Marshal(ext)
            extServiceType := structs.ServiceType{"ext", extJson}

            service := structs.NewService(path, (*team).TeamID)
            serviceCheck := structs.LoadFromServiceType(extServiceType)
            service.ServiceType = extServiceType
            service.ServiceCheck = serviceCheck
            dbutils.LoadServiceFromDatabase(db, &service)
            // TODO: Add Service to all teams
            if !isServiceAlreadyInList(*team, service) {
                (*team).Services = append((*team).Services, service)
                fmt.Println("[+] team ", (*team).TeamID, "loaded new script: ", path)
            }
        }
        return nil
    })
    if err != nil {
        log.Println(err)
    }
}

func addScriptServicesToTeams(db *sql.DB, directory string, teams *[]structs.Team) {
    for i := 0; i < len(*teams); i++ {
        addScriptServicesToTeam(db, directory, &(*teams)[i])
    }
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback then display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

func GetScoreboard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", scoreboardHTML)
}

func ConcurrentServiceCheck(servicesPtr *[]structs.Service) {
	// Spin off a separate goroutine for each service check
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
	confName, dbName, scriptDir := parseArgs(argparse.NewParser("gizmo", "Service uptime scoreboard"))

    // TODO: Check if specified config / database files exist
    if !fileExists(confName) {
        fmt.Println("[!] Configuration file: ", confName, "does not exist!")
        os.Exit(1)
    }

	db, _ := sql.Open("sqlite3", dbName)

	teams := config.LoadConfig(confName)

    if scriptDir != "" {
        if _, err := os.Stat(scriptDir); !os.IsNotExist(err) {
            addScriptServicesToTeams(db, scriptDir, &teams)
        } else {
            fmt.Println("[!] Script directory: ", scriptDir, "does not exist!")
            os.Exit(1)
        }
    }

	game := structs.NewGame(teams)

    if dbutils.DoesDatabaseExist(db) {
        fmt.Println("[+] Database already exists - loading")
        dbutils.UpdateGameFromDatabase(db, &game)
        dbutils.LoadGameFromDatabase(db, &game)
    } else {
        fmt.Println("[+] Database is empty - creating and initializing")
	    dbutils.CreateDatabase(db)
	    dbutils.InitializeDatabase(db, game)
    }

	scoreboardHTML = web.GenerateScoreboardHTML(teams)

	// Spin off separate thread for the web server so as not to block main
	go func() {
        local_ip := GetLocalIP()
		fmt.Printf("web server listening on %s:8080\n", local_ip)
		http.HandleFunc("/", GetScoreboard)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

    time.Sleep(time.Duration(1) * time.Second)

	// Loop every three to five minutes until next service check
	//min := 180 // 180 seconds = 3 minutes
	//max := 300 // 300 seconds = 5 minutes
	for {

		for i := 0; i < len(teams); i++ {
			team := teams[i]

            addScriptServicesToTeam(db, scriptDir, &teams[i])
			ConcurrentServiceCheck(&team.Services)
			scoreboardHTML = web.GenerateScoreboardHTML(teams)

			structs.UpdateTeamCheckCount(&teams[i])
			dbutils.UpdateDatabase(db, game)
		}

		fmt.Println("=======================")
		sleepTime := 3 //rand.Intn(max - min) + min
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}
