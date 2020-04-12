// Package structs provides the objects that represent the database data that
// is recorded. There is a single Game, a Game contains Teams, and each Team
// has Services.
package structs

import (
	check "github.com/whoismissing/gizmo/check"

	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Game describes the instance of a session and can contain multiple teams.
// Game is a singleton so there should only be one game serviced at a time.
type Game struct {
	GameID    int
	Teams     []Team
	StartTime time.Time
}

// Team describes the unique instance of a single team which has multiple services.
// A Team contains the ratio of service uptime across all services.
// TeamID is the unique key for representing a Team.
type Team struct {
	TeamID               uint
	TotalChecksMissed    uint
	TotalChecksHit       uint
	TotalChecksAttempted uint
	Services             []Service
}

// Service describes the unique instance of a service that a team is supporting.
// A Service contains the ratio of service uptime for a single service and the
// previous statuses (currently limited to 10 previous statuses).
// ServiceType is a wrapper for ServiceCheck that is used as a hack to unmarshal
// and obtain the ServiceCheck object from the JSON config.
// ServiceCheck is the actual object that represents the type of service check
// [ www / dns / ftp / ssh / ping ] and will provide the corresponding method of
// checking.
// The Service.Name and combination of ServiceID and TeamID are the unique keys for
// representing a Service. ServiceID's on their own are NOT unique.
type Service struct {
	ServiceID       int
	Name            string // primary key in the SQL database
	Status          bool
	ServiceType     ServiceType
	ServiceCheck    ServiceCheck
	HostIP          string
	TeamID          uint
	ChecksMissed    uint
	ChecksHit       uint
	ChecksAttempted uint
	PrevStatuses    []Status
}

// Status represents a single instance of a binary check [ up / down ] and records
// the time of the check.
type Status struct {
	Time   time.Time
	Status bool
}

// ServiceType is a wrapper for ServiceCheck that is used as a hack so
// as not to lose type information when a ServiceCheck (interface) object is
// marshalled to JSON.
type ServiceType struct {
	Type         string
	ServiceCheck json.RawMessage
}

// ServiceCheck is an interface with the method CheckHealth() that is implemented
// by types WebService, DomainNameService, etc.
type ServiceCheck interface {
	CheckHealth(service *Service, wg *sync.WaitGroup)
}

// WebService represents a web ServiceCheck and contains the URL string to check.
type WebService struct {
	URL string
}

// DomainNameService represents a dns ServiceCheck and contains the domain name
// string to use during a service check.
type DomainNameService struct {
	DomainName string
}

// FileTransferService represents a ftp ServiceCheck and contains the username
// and password strings to use during a service check.
type FileTransferService struct {
	Username string
	Password string
	Filename string
}

// SecureShellService represents a ssh ServiceCheck and contains the username,
// password, and command strings to use during a ssh service check.
type SecureShellService struct {
	Username string
	Password string
	Command  string
}

// PingService represents a ping ServiceCheck and is empty because the
// Service.HostIP field is used for its service check.
type PingService struct {
}

// ExternalService represents an externally written service check and contains
// the script path to execute. The Service.HostIP is passed as the first argument
// to the script.
type ExternalService struct {
	ScriptPath string
}

// updateCheckCount() increments the ChecksHit or ChecksMissed counts of a
// Service depending on the current status.
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

// UpdateTeamCheckCount() increments a Team's TotalChecksHit or TotalChecksMissed
// depending on each team.Service's current status.
func UpdateTeamCheckCount(team *Team) {
	services := (*team).Services
	for i := 0; i < len(services); i++ {
		service := services[i]
		if service.Status {
			(*team).TotalChecksHit += 1
		} else {
			(*team).TotalChecksMissed += 1
		}
		(*team).TotalChecksAttempted += 1
	}
}

// www.CheckHealth() is the CheckHealth() method implemented by the WebService
// object that calls the corresponding web checking method from package check.
func (www WebService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := (*service).HostIP
	fmt.Printf("[ WWW ] teamid=%d name=%s targetip=%s url=%s\n", service.TeamID, service.Name, ip, www.URL)

	// ip or url must be in format: http://192.168.1.1 for web check
	// so we prepend http://
	httpIP := "http://" + ip

	var status bool
	if www.URL != "" {
		status = check.Web(www.URL)
	} else {
		status = check.Web(httpIP)
	}
	updateCheckCount(service, status)
}

// dns.CheckHealth() is the CheckHealth() method implemented by the DomainNameService
// object that calls the corresponding dns checking method from package check.
func (dns DomainNameService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := (*service).HostIP
	record := dns.DomainName
	fmt.Printf("[ DNS ] teamid=%d name=%s targetip=%s record=%s\n", service.TeamID, service.Name, ip, record)

	status := check.Dns(ip, record)
	updateCheckCount(service, status)
}

// ftp.CheckHealth() is the CheckHealth() method implemented by the FileTransferService
// object that calls the corresponding ftp checking method from package check.
func (ftp FileTransferService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := (*service).HostIP
	user := ftp.Username
	pass := ftp.Password
	filename := ftp.Filename
	fmt.Printf("[ FTP ] teamid=%d name=%s targetip=%s filename=%s username=%s password=%s filename=%s\n", service.TeamID, service.Name, ip, filename, user, pass, filename)

	status := check.Ftp(ip, user, pass, filename)
	updateCheckCount(service, status)
}

// ssh.CheckHealth() is the CheckHealth() method implemented by the SecureShellService
// object that calls the corresponding ssh checking method from package check.
func (ssh SecureShellService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := (*service).HostIP
	user := ssh.Username
	pass := ssh.Password
	command := ssh.Command
	fmt.Printf("[ SSH ] teamID=%d name=%s targetip=%s username=%s password=%s command=%s\n", service.TeamID, service.Name, ip, user, pass, command)

	status := check.Ssh(ip, user, pass, command)
	updateCheckCount(service, status)
}

// ping.CheckHealth() is the CheckHealth() method implemented by the PingService
// object that calls the corresponding ping checking method from package check.
func (ping PingService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := (*service).HostIP
	fmt.Printf("[ ping ] teamID=%d name=%s targetip=%s\n", service.TeamID, service.Name, ip)

	status := check.Ping(ip)
	updateCheckCount(service, status)
}

// ext.CheckHealth() is the CheckHealth() method implemented by the ExternalService
// object that calls the corresponding external checking method from package check.
func (ext ExternalService) CheckHealth(service *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	ip := (*service).HostIP
	filepath := ext.ScriptPath
	fmt.Printf("[ EXT ] teamID=%d name=%s targetip=%s filepath=%s\n", service.TeamID, service.Name, ip, filepath)

	status := check.External(ip, filepath)
	updateCheckCount(service, status)
}

// LoadFromServiceType() unmarshals the specified ServiceCheck object given its
// corresponding ServiceType wrapper.
// This is used to obtain the type information lost from the JSON config.
func LoadFromServiceType(serviceType ServiceType) ServiceCheck {
	data := serviceType.ServiceCheck
	switch stype := serviceType.Type; stype {
	case "www":
		var www WebService
		_ = json.Unmarshal(data, &www)
		// TODO: do some verification of www.URL to ensure it is either "" or
		// is prepended by http:// or https://
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

// NewGame() returns a new Game object given an array of Teams.
func NewGame(teams []Team) Game {
	game := Game{GameID: rand.Int(), Teams: teams, StartTime: time.Now()}
	return game
}

// NewTeam() returns a new Team object given a teamID.
func NewTeam(newTeamID uint) Team {
	team := Team{TeamID: newTeamID, TotalChecksMissed: 0, TotalChecksHit: 0, TotalChecksAttempted: 0, Services: nil}
	return team
}

// NewTeams() returns a user-specified number of empty Team objects in an array.
func NewTeams(numTeams uint) []Team {
	teams := make([]Team, numTeams)
	return teams
}

// NewService() returns a new Service object given a service name and corresponding teamID.
func NewService(serviceName string, newTeamID uint) Service {
	service := Service{Name: serviceName, Status: false, HostIP: "", TeamID: newTeamID, ChecksMissed: 0, ChecksHit: 0, ChecksAttempted: 0, PrevStatuses: nil}
	return service
}

// NewServices() returns a user-specified number of empty Service objects in an array given
// a corresponding teamID.
func NewServices(numServices uint, teamID uint) []Service {
	services := make([]Service, numServices)
	for i := uint(0); i < numServices; i++ {
		services[i].TeamID = teamID
	}
	return services
}

// NewStatus() returns a Status object given a boolean.
func NewStatus(initStatus bool) Status {
	status := Status{Time: time.Now(), Status: initStatus}
	return status
}
