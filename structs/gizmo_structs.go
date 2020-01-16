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

func updateCheckCount(service *Service, status bool) {
	if status == true {
		(*service).ChecksHit += 1
	} else {
		(*service).ChecksMissed += 1
	}

	(*service).Status = status
	(*service).ChecksAttempted += 1

	currStatus := NewStatus(status)
	(*service).PrevStatuses = append((*service).PrevStatuses, currStatus)
	trunc := len((*service).PrevStatuses) - 1
	if trunc > 9 {
		(*service).PrevStatuses = (*service).PrevStatuses[:trunc]
	}
}

func UpdateTeamCheckCount(team *Team) {
	services := (*team).Services
	for i := 0; i < len(services); i++ {
		service := services[i]

		(*team).TotalChecksHit += service.ChecksHit
		(*team).TotalChecksMissed += service.ChecksMissed
		(*team).TotalChecksAttempted += service.ChecksAttempted
	}
}

func (www WebService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("www CheckHealth()")
}

func (dns DomainNameService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("dns CheckHealth()")
}

func (ftp FileTransferService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("ftp CheckHealth()")
}

func (ssh SecureShellService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("ssh CheckHealth()")
}

func (ping PingService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("ping CheckHealth()")

	ip := (*service).HostIP
	status := check.Ping(ip)

	updateCheckCount(service, status)
}

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
	case "default":
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

