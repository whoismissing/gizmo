package main

import (
    structs "github.com/whoismissing/gizmo/structs"

    "encoding/json"
    "strings"
    "os"
    "bufio"
    "fmt"
)

var (
    teams = make([]structs.Team, 0) // alternative syntax is teams []structs.Team
    game = structs.NewGame(teams) // game is a singleton
)

func printUsage() {
   fmt.Println("Usage: ", os.Args[0], "config_filename")
}

func readIntFromUser() int {
    var userInput int
    _, _ = fmt.Scanln(&userInput) // numItemsRead, err = fmt.Scanln()

    return userInput
}

func readStringFromUser() string {
    reader := bufio.NewReader(os.Stdin)
    var userInput string
    userInput, _ = reader.ReadString('\n')

    return userInput
}

func getWebServiceType() structs.ServiceType {
    var www structs.WebService
    fmt.Printf("\t\t\tEnter a URL: ")
    url := readStringFromUser()
    url = strings.TrimSuffix(url, "\n")
    www.URL = url

    wwwJson, _ := json.Marshal(www)
    wwwServiceType := structs.ServiceType{"www", wwwJson}

    return wwwServiceType
}

func getDnsServiceType() structs.ServiceType {
    var dns structs.DomainNameService
    fmt.Printf("\t\t\tEnter a domain name: ")
    domain := readStringFromUser()
    domain = strings.TrimSuffix(domain, "\n")
    dns.DomainName = domain

    dnsJson, _ := json.Marshal(dns)
    dnsServiceType := structs.ServiceType{"dns", dnsJson}

    return dnsServiceType
}

func getFtpServiceType() structs.ServiceType {
    var ftp structs.FileTransferService
    fmt.Printf("\t\t\tEnter a username: ")
    username := readStringFromUser()
    username = strings.TrimSuffix(username, "\n")
    ftp.Username = username

    ftpJson, _ := json.Marshal(ftp)
    ftpServiceType := structs.ServiceType{"ftp", ftpJson}

    return ftpServiceType
}

func getSshServiceType() structs.ServiceType {
    var ssh structs.SecureShellService

    fmt.Printf("\t\t\tEnter a username: ")
    username := readStringFromUser()
    username = strings.TrimSuffix(username, "\n")

    fmt.Printf("\t\t\tEnter a password: ")
    password := readStringFromUser()
    password = strings.TrimSuffix(password, "\n")
    ssh.Command = "ls"
    ssh.Username = username
    ssh.Password = password

    sshJson, _ := json.Marshal(ssh)
    sshServiceType := structs.ServiceType{"ssh", sshJson}

    return sshServiceType
}

func getPingServiceType() structs.ServiceType {
    var ping structs.PingService

    pingJson, _ := json.Marshal(ping)
    pingServiceType := structs.ServiceType{"ping", pingJson}

    return pingServiceType
}

func getExternalServiceType() structs.ServiceType {
    var ext structs.ExternalService
    fmt.Printf("\t\t\tEnter a program filepath: ")
    path := readStringFromUser()
    path = strings.TrimSuffix(path, "\n")
    ext.ScriptPath = path

    extJson, _ := json.Marshal(ext)
    extServiceType := structs.ServiceType{"ext", extJson}

    return extServiceType
}

func getServiceTypeFromRaw(rawType string) (structs.ServiceType, bool) {
    // uninitialized, doesn't matter what type
    var def structs.ServiceType

    rawType = strings.TrimSuffix(rawType, "\n")
    switch stype := rawType; stype {
    case "www":
        wwwServiceType := getWebServiceType()
        return wwwServiceType, true
    case "dns":
        dnsServiceType := getDnsServiceType()
        return dnsServiceType, true
    case "ftp":
        ftpServiceType := getFtpServiceType()
        return ftpServiceType, true
    case "ssh":
        sshServiceType := getSshServiceType()
        return sshServiceType, true
    case "ping":
        pingServiceType := getPingServiceType()
        return pingServiceType, true
    case "ext":
        extServiceType := getExternalServiceType()
        return extServiceType, true
    default:
        fmt.Println("Unrecognized service type")
        return def, false
    }

    return def, false
}

func verifyRawServiceType(rawType string) bool {
    retval := false
    rawType = strings.TrimSuffix(rawType, "\n")
    switch stype := rawType; stype {
    case "www":
        fallthrough
    case "dns":
        fallthrough
    case "ftp":
        fallthrough
    case "ssh":
        fallthrough
    case "ping":
        fallthrough
    case "ext":
        retval = true
        break
    default:
        retval = false
        break
    }

    return retval
}

func readServiceTypeFromUser() string {
    for {
        fmt.Println("\tEnter the new value, no default option")
        fmt.Printf("\tServiceType [www/dns/ext/ftp/ssh/ping]: ")
        rawType := readStringFromUser()
        if verifyRawServiceType(rawType) == true {
            return rawType
        } else {
            fmt.Println("\t\t[!] Unrecognized ServiceType")
            fmt.Printf("\t\tServiceType is required - try again? [y/n]: ")
            another := readStringFromUser()

            if strings.Contains(another, "n") {
                break
            }
            if strings.Contains(another, "y") {
                continue
            } else { // not "y"
                break
            }

        } // end else
    } // end for

    return ""
}

/*
Prompt user to add a service by providing:
1. ServiceName
2. HostIP
3. ServiceType
Depending on the type, additional information may be needed
*/
func addService(serviceID int, teamID int) {
    fmt.Printf("\tAdding service '%d' ...\n", serviceID)

    fmt.Println("\tEnter the new value, press ENTER for default")
    fmt.Printf("\tServiceName [team%d-service%d]: ", teamID, serviceID)
    serviceName := readStringFromUser()
    if strings.Compare(serviceName, "\n") == 0 {
        defaultName := fmt.Sprintf("%s%d-%s%d", "team", teamID, "service", serviceID)
        fmt.Println("\t\tService name is default", defaultName)
        serviceName = defaultName
    }

    fmt.Printf("\tHostIP [127.0.0.1]: ")
    hostIP := readStringFromUser()
    if strings.Compare(hostIP, "\n") == 0 {
        defaultIP := "127.0.0.1"
        fmt.Println("\t\tHost ip is default", defaultIP)
        hostIP = defaultIP
    }

    rawType := readServiceTypeFromUser()
    if rawType == "" {
        fmt.Println("\t[!] No rawType given")
        return
    }

    /* function call to handle raw service type */
    serviceType, status := getServiceTypeFromRaw(rawType)
    if status != true {
        os.Exit(1)
    }

    newService := structs.NewService(serviceName, uint(teamID))
    newService.ServiceID = serviceID
    newService.HostIP = hostIP
    newService.ServiceType = serviceType
    teams[teamID].Services = append(teams[teamID].Services, newService)
}

func addServices(team structs.Team) {
    serviceID := 0
    for {
        addService(serviceID, int(team.TeamID))
        fmt.Printf("Add another service? [y/n] ")
        another := readStringFromUser()

        if strings.Contains(another, "n") {
            break
        }
        if strings.Contains(another, "y") {
            serviceID += 1
            continue
        } else { // not "y"
            break
        }
    } // end for
}

func addTeam(teamID int) {
    fmt.Printf("Adding team '%d' ...\n", teamID)
    newTeam := structs.NewTeam(uint(teamID))
    newTeam.Services = make([]structs.Service, 0)

    teams = append(teams, newTeam)
    addServices(newTeam)
}

func addTeams() {
    teamID := 0
    for {
        addTeam(teamID)
        fmt.Printf("Add another team? [y/n] ")
        another := readStringFromUser()
        if strings.Contains(another, "n") {
            break
        }
        if strings.Contains(another, "y") {
            teamID += 1
            continue
        } else { // not "y"
            break
        }
    } // end for
}

func promptUser() {
    fmt.Println("Creating a new Game config ...")
    addTeams()
}

func writeToFile(filename string, data string) {
    fd, err := os.Create(filename)
    if err != nil {
        fmt.Println(err)
        return
    }
    _, err = fd.WriteString(data)
    if err != nil {
        fmt.Println(err)
        fd.Close()
        return
    }
    err = fd.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
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

    // pretty print json
    fmt.Println("JSON config: ")
    jsonConfig, _ := json.MarshalIndent(teams, "", "    ")
    fmt.Println(string(jsonConfig))
    fmt.Println()

    fmt.Printf("Wrote config to file '%s'\n", filename)
    writeToFile(filename, string(jsonConfig))
}
