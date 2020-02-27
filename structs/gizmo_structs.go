package structs

import (
	check "github.com/whoismissing/gizmo/check"

	"encoding/json"
	"math/rand"
	"time"
	"sync"
	"fmt"
)

type Game struct {
	GameID int
	Teams []Team
	StartTime time.Time
}

type Team struct {
	TeamID uint
	TotalChecksMissed uint
	TotalChecksHit uint
	TotalChecksAttempted uint
	Services []Service
}

type Service struct {
	ServiceID int
	Name string
	Status bool
	ServiceType ServiceType
	ServiceCheck ServiceCheck
	HostIP string
	TeamID uint
	ChecksMissed uint
	ChecksHit uint
	ChecksAttempted uint
	PrevStatuses []Status
}

type Status struct {
	Time time.Time
	Status bool
}

/*
This struct is used as a hack to wrap a ServiceCheck object
so as not to lose type information when marshalled to JSON
*/
type ServiceType struct {
	Type string
	ServiceCheck json.RawMessage
}

type ServiceCheck interface {
	CheckHealth(service *Service, wg *sync.WaitGroup)
}

type WebService struct {
	URL string
}

type DomainNameService struct {
	DomainName string
}

type FileTransferService struct {
	Username string
	Password string
}

type SecureShellService struct {
	Username string
	Password string
	Command string
}

type PingService struct {

}

type ExternalService struct {
    ScriptPath string
}

func updateCheckCount(service *Service, status bool) {
	if status == true {
		(*service).ChecksHit += 1
	} else {
		(*service).ChecksMissed += 1
	}

	(*service).Status = status
	(*service).ChecksAttempted += 1

	// Save the result of the service check in []PrevStatuses
	currStatus := NewStatus(status)
	(*service).PrevStatuses = append((*service).PrevStatuses, currStatus)
	trunc := len((*service).PrevStatuses) - 1
	// Limit the history length of PrevStatuses
	if trunc > 9 {
		(*service).PrevStatuses = (*service).PrevStatuses[:trunc]
	}
}

func UpdateTeamCheckCount(team *Team) {
	services := (*team).Services
	for i := 0; i < len(services); i++ {
		service := services[i]
        if (service.Status) {
            (*team).TotalChecksHit += 1
        } else {
		    (*team).TotalChecksMissed += 1
        }
		(*team).TotalChecksAttempted += 1
	}
}

func (www WebService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := (*service).HostIP
	fmt.Printf("[ WWW ] targetip=%s\n", ip)

	status := check.Web(ip)
	updateCheckCount(service, status)
}

func (dns DomainNameService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := (*service).HostIP
	record := dns.DomainName
	fmt.Printf("[ DNS ] targetip=%s record=%s\n", ip, record)

	status := check.Dns(ip, record)
	updateCheckCount(service, status)
}

func (ftp FileTransferService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := (*service).HostIP
	user := ftp.Username
	pass := ftp.Password
	// TODO: obtain FTP filename from JSON config
	filename := "hello.txt"
	fmt.Printf("[ FTP ] targetip=%s filename=%s username=%s password=%s\n", ip, filename, user, pass)

	status := check.Ftp(ip, user, pass, filename)
	updateCheckCount(service, status)
}

func (ssh SecureShellService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := (*service).HostIP
	user := ssh.Username
	pass := ssh.Password
	fmt.Printf("[ SSH ] targetip=%s username=%s password=%s\n", ip, user, pass)

	status := check.Ssh(ip, user, pass)
	updateCheckCount(service, status)
}

func (ping PingService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := (*service).HostIP
	fmt.Printf("[ ping ] targetip=%s\n", ip)

	status := check.Ping(ip)
	updateCheckCount(service, status)
}


func (ext ExternalService) CheckHealth(service *Service, wg *sync.WaitGroup) {
    defer wg.Done()

    ip := (*service).HostIP
    filepath := ext.ScriptPath
    fmt.Printf("[ EXT ] targetip=%s filepath=%s\n", ip, filepath)

    status := check.External(ip, filepath)
    updateCheckCount(service, status)
}

/*
This function is used to unmarshal the 
specified ServiceCheck object type
wrapped in a ServiceType from the JSON config
*/
func LoadFromServiceType(serviceType ServiceType) ServiceCheck {
	data := serviceType.ServiceCheck
	switch stype := serviceType.Type; stype {
	case "www":
		var www WebService
		_ = json.Unmarshal(data, &www)
		return www
	case "dns":
		var dns DomainNameService
		_ = json.Unmarshal(data, &dns)
		return dns
	case "ftp":
		var ftp FileTransferService
		_ = json.Unmarshal(data, &ftp)
		return ftp
	case "ssh":
		var ssh SecureShellService
		_ = json.Unmarshal(data, &ssh)
		return ssh
	case "ping":
		var ping PingService
		_ = json.Unmarshal(data, &ping)
		return ping
    case "ext":
        var ext ExternalService
        _ = json.Unmarshal(data, &ext)
        return ext
	case "default":
        fmt.Println("LoadFromServiceType: unrecognized ServiceType")
		return nil
	}

	return nil
}

func NewGame(teams []Team) Game {
	game := Game{GameID:rand.Int(), Teams: teams, StartTime: time.Now()}
	return game
}

func NewTeam(newTeamID uint) Team {
	team := Team{TeamID: newTeamID, TotalChecksMissed: 0, TotalChecksHit: 0, TotalChecksAttempted: 0, Services: nil}
	return team
}

func NewTeams(numTeams uint) []Team {
	teams := make([]Team, numTeams)
	return teams
}

func NewService(serviceName string, newTeamID uint) Service {
	service := Service{Name: serviceName, Status: false, HostIP: "", TeamID: newTeamID, ChecksMissed: 0, ChecksHit: 0, ChecksAttempted: 0, PrevStatuses: nil}
	return service
}

func NewServices(numServices uint, teamID uint) []Service {
	services := make([]Service, numServices)
	for i:= uint(0); i < numServices; i++ {
		services[i].TeamID = teamID
	}
	return services
}

func NewStatus(initStatus bool) Status {
	status := Status{Time: time.Now(), Status: initStatus}
	return status
}

