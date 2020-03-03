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
    "net/http"
    "path/filepath"
    "fmt"
    "log"
    "net"
    "os"
    "time"
    "sync"
)

// Global variable representing HTML to be used by GetScoreboard()
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

/*
addScriptServicesToTeam() adds all of the scripts 
in the script directory as services to the database and team structure
*/
func addScriptServicesToTeam(db *sql.DB, directory string, team *structs.Team) {
    /* get a list of all files (recursively) in the script directory */
    err := filepath.Walk(directory,
        func(path string, info os.FileInfo, err error) error {

            if err != nil {
                return err
            }

            if fileExists(path) {
                ext := structs.ExternalService{ ScriptPath: path }
                extJson, _ := json.Marshal(ext)
                extServiceType := structs.ServiceType{"ext", extJson}

                service := structs.NewService(path, (*team).TeamID)
                serviceCheck := structs.LoadFromServiceType(extServiceType)
                service.ServiceType = extServiceType
                service.ServiceCheck = serviceCheck
                // Check for service in database and insert into the database
                // if the service does NOT exist
                dbutils.LoadServiceFromDatabase(db, &service)

                // if service is not in the list then, add service to team
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

/*
addScriptServicesToTeams() adds all of the scripts
in the script directory as services to all teams
*/
func addScriptServicesToTeams(db *sql.DB, directory string, teams *[]structs.Team) {
    for i := 0; i < len(*teams); i++ {
        addScriptServicesToTeam(db, directory, &(*teams)[i])
    }
}

/* 
GetLocalIP() returns the non loopback local IP of the host
*/
func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback then return it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

/*
GetScoreboard() will be return scoreboardHTML to the http server handler
*/
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

    if !fileExists(confName) {
        fmt.Println("[!] Configuration file: ", confName, "does not exist!")
        os.Exit(1)
    }

    db, _ := sql.Open("sqlite3", dbName)

    // Load objects from JSON config
    teams := config.LoadConfig(confName)

    if scriptDir != "" {
        // if the user specifies a script directory, then addScriptServicesToTeams()
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

    for {

        for i := 0; i < len(teams); i++ {
            team := teams[i]

            if scriptDir != "" {
                addScriptServicesToTeam(db, scriptDir, &teams[i])
            }

            ConcurrentServiceCheck(&team.Services)
            scoreboardHTML = web.GenerateScoreboardHTML(teams)

            structs.UpdateTeamCheckCount(&teams[i])
            dbutils.UpdateDatabase(db, game)
        }

        fmt.Println("=======================")
        sleepTime := 3
        time.Sleep(time.Duration(sleepTime) * time.Second)
    }
}
