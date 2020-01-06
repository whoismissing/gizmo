package structs

// TODO: uppercase all struct value names

type Game struct {
	teams []Team
	time uint
}

type Team struct {
	teamID uint
	totalChecksMissed uint
	totalChecksHit uint
	totalChecksAttempted uint
	Services []Service
}

type Service struct {
	name string
	status bool
	hostIP string
	teamID uint
	checksMissed uint
	checksHit uint
	checksAttempted uint
	prevStatuses []Status
}

type Status struct {
	time uint
	status bool
}

/*
type ServiceType interface {
	ChooseService() ServiceType
}

type WebService struct {

}

func (s *WebService) ChooseService() WebService {
	return s
}
*/

func NewGame(initTime uint) Game {
	game := Game{teams: nil, time: initTime}
	return game
}

/*
func NewTeam(newTeamID uint) Team {
	team := Team{teamID: newTeamID, totalChecksMissed: 0, totalChecksHit: 0, totalChecksAttempted: 0, services: nil}
	return team
}
*/

func NewService(serviceName string, newTeamID uint) Service {
	service := Service{name: serviceName, status: false, hostIP: "", teamID: newTeamID, checksMissed: 0, checksHit: 0, checksAttempted: 0, prevStatuses: nil}
	return service
}

func NewStatus(initTime uint, initStatus bool) Status {
	status := Status{time: initTime, status: initStatus}
	return status
}

