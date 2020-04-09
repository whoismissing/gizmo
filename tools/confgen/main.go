package main

import (
    structs "github.com/whoismissing/gizmo/structs"

    //"encoding/json"
    "os"
    "bufio"
    "fmt"
)

var (
    //teams []structs.Team
    teams = make([]structs.Team, 0)
    game = structs.NewGame(teams) // game is a singleton
)

func printUsage() {
   fmt.Println("Usage: ", os.Args[0], "config_filename") 
}

func readIntFromUser() int {
    var userInput int
    _, _ = fmt.Scan(&userInput)

    return userInput
}

func readStringFromUser() string {
    reader := bufio.NewReader(os.Stdin)
    var userInput string
    userInput, _ = reader.ReadString('\n')

    return userInput
}

/*
Adding user `hello' ...
Adding new group `hello' (1001) ...
Adding new user `hello' (1001) with group `hello' ...
Creating home directory `/home/hello' ...
Copying files from `/etc/skel' ...
Enter new UNIX password: 
Retype new UNIX password: 
passwd: password updated successfully
Changing the user information for hello
Enter the new value, or press ENTER for the default
    Full Name []: 
    Room Number []: 
    Work Phone []: 
    Home Phone []: 
    Other []: 
Is the information correct? [Y/n]
*/

func getServiceTypeFromRaw(rawType string) structs.ServiceType {
    switch stype := rawType; stype {
    case "www":
        var www WebService
        return www
    case "dns":
        return dns
    case "ftp":
        return ftp
    case "ssh":
        return ssh
    case "ping":
        return ping
    case "ext":
        return ext
    case "default":
        fmt.Println("Unrecognized service type")
        return nil
    }

    return nil
}

/*
Prompt user to add a service by providing:
1. ServiceName
2. ServiceType
Depending on the type, additional information may be needed
3. HostIP
*/
func addService(defaultServiceID int, teamID int) {
    fmt.Printf("\tAdding service '%d' ...\n", defaultID)
    fmt.Println("\tEnter the new value, or press ENTER for the default")

    fmt.Printf("ServiceName [team%d-service%d] ...\n", defaultID, defaultID)
    serviceName := readStringFromUser()

    fmt.Printf("ServiceID [%d]: ", defaultID)
    serviceID := readIntFromUser()

    /* required - no default option */
    fmt.Println("\tEnter the new value, no default option")
    fmt.Println("ServiceType [www/dns/ext/ftp/ssh/ping] ...")
    rawType := readStringFromUser()

    /* function call to handle raw service type */
    serviceType := getServiceTypeFromRaw(rawType)

    newService := structs.NewService(serviceName, uint(serviceID))
    teams[teamID].Services = append(teams[teamID].Services, newService)
}

func addServices(team structs.Team) {
    addService(0, team.TeamID)
}

/*
Prompt user to add a team by providing:
1. TeamID   - default autoincrement
2. TeamName - default
3. Services - addServices()
*/
func addTeam(defaultID int) {
    fmt.Printf("Adding team '%d' ...\n", defaultID)
    fmt.Println("Enter the new value, or press ENTER for the default")
    fmt.Printf("TeamID [%d]: ", defaultID)
    teamID := readIntFromUser()
    newTeam := structs.NewTeam(uint(teamID))
    newTeam.Services = make([]structs.Service, 0)

    addServices(newTeam)
    teams = append(teams, newTeam)
}

func addTeams() {
    addTeam(0)
}

func promptUser() {
    fmt.Println("Creating a new Game config ...")
    addTeams()
}

func main() {

    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

    filename := os.Args[1]

    fmt.Println("Config Filename:", filename)

    promptUser()

    game.Teams = teams
    fmt.Printf("Game = %+v\n", game)
    //fmt.Printf("teams = %+v\n", teams)
}
