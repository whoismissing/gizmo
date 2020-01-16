package structs

import (
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
	CheckHealth(id int, wg *sync.WaitGroup)
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

func (www WebService) CheckHealth(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("www CheckHealth()", id)
}

func (dns DomainNameService) CheckHealth(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("dns CheckHealth()", id)
}

func (ftp FileTransferService) CheckHealth(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("ftp CheckHealth()", id)
}

func (ssh SecureShellService) CheckHealth(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("ssh CheckHealth()", id)
}

func (ping PingService) CheckHealth(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("ping CheckHealth()", id)
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

