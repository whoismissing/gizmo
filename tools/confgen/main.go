// Tools/confgen provides the command-line tool to help users generate a JSON
// configuration file for gizmo in a user-prompting fashion similar to the
// adduser command.
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

// readIntFromUser() scans user input until a new-line char and returns
// the user-provided int
func readIntFromUser() int {
    var userInput int
    _, _ = fmt.Scanln(&userInput) // numItemsRead, err = fmt.Scanln()

    return userInput
}

// readStringFromUser() reads user input until a new-line char and returns
// the user-provided string (including the new-line char)
func readStringFromUser() string {
    reader := bufio.NewReader(os.Stdin)
    var userInput string
    userInput, _ = reader.ReadString('\n')

    return userInput
}

// getWebServiceType() prompts the user to enter a URL and returns the
// newly-created web ServiceType object.
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

// getDnsServiceType() prompts the user to enter a domain name and returns the 
// newly-created dns ServiceType object.
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

// getFtpServiceType() prompts the user for a username and returns the newly-created
// ftp ServiceType object.
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

// getSshServiceType() prompts the user for a username and password and returns the
// newly-created ssh ServiceType object. 'ls' is used as the default command.
func getSshServiceType() structs.ServiceType {
    var ssh structs.SecureShellService

    fmt.Printf("\t\t\tEnter a username: ")
    username := readStringFromUser()
    username = strings.TrimSuffix(username, "\n")

    fmt.Printf("\t\t\tEnter a password: ")
    password := readStringFromUser()
    password = strings.TrimSuffix(password, "\n")
    ssh.Command = "ls" // TODO: allow the user to provide a custom command
    ssh.Username = username
    ssh.Password = password

    sshJson, _ := json.Marshal(ssh)
    sshServiceType := structs.ServiceType{"ssh", sshJson}

    return sshServiceType
}

// getPingServiceType() returns a ping ServiceType object.
func getPingServiceType() structs.ServiceType {
    var ping structs.PingService

    pingJson, _ := json.Marshal(ping)
    pingServiceType := structs.ServiceType{"ping", pingJson}

    return pingServiceType
}

// getExternalServiceType() prompts the user for a program filepath and returns
// the newly-created External ServiceType object.
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

// getServiceTypeFromRaw() calls the corresponding ServiceType prompt to get further
// user-provided information to obtain the ServiceType in question, given an initial
// raw service type. The raw service types are: www, dns, ftp, ssh, ping, ext.
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

// verifyRawServiceType() trims the new-line from the user-provided input
// and returns true for the supported options: www, dns, ftp, ssh, ping, ext.
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

// readServiceTypeFromUser() prompts the user to provide a raw service type,
// reprompting when the input is unrecognized. Returns an empty-string on exit with 'n',
// otherwise returns the service type the user chose as a string.
// Options are: www, dns, ext, ftp, ssh, ping.
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

// addService() prompts the user to add a service by providing
// 1. ServiceName
// 2. HostIP
// 3. ServiceType
// Depending on the type, additional information is requested from
// the user for the corresponding type chosen
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
    serviceName = strings.TrimSuffix(serviceName, "\n")

    fmt.Printf("\tHostIP [127.0.0.1]: ")
    hostIP := readStringFromUser()
    if strings.Compare(hostIP, "\n") == 0 {
        defaultIP := "127.0.0.1"
        fmt.Println("\t\tHost ip is default", defaultIP)
        hostIP = defaultIP
    }
    hostIP = strings.TrimSuffix(hostIP, "\n")

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

// addServices() prompts the user to continuously add services
// until they select 'n' for 'no'.
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

// addTeam() prompts the user to add a team to the config. Upon doing so,
// the user will be prompted to add services to the corresponding team.
func addTeam(teamID int) {
    fmt.Printf("Adding team '%d' ...\n", teamID)
    newTeam := structs.NewTeam(uint(teamID))
    newTeam.Services = make([]structs.Service, 0)

    teams = append(teams, newTeam)
    addServices(newTeam)
}

// addTeams() prompts the user to continuously add teams until 'n'
// is selected for 'no'.
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

// promptUser() is the entry point for prompting the user to the
// addTeams() loop.
func promptUser() {
    fmt.Println("Creating a new Game config ...")
    addTeams()
}

// writeToFile() simply writes all specified data to the specified filename
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

// main() is the entrypoint of the program, the user is prompted
// for information and the corresponding objects will be created.
// Once the user prompts are finished, the final array of Team objects 
// is marshalled to JSON and written to the user-specified filename
// completing the creation of the JSON configuration file for usage
// in gizmo.
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
