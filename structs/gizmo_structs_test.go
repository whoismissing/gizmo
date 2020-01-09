package structs

import (
	"testing"
)

func TestInitializeStructs(t *testing.T) {
	t.Logf("TestInitializeStructs SUCCESS")

	var teamID uint = 1
	team1 := NewTeam(teamID)

	team1.Services = NewServices(5, teamID)
	team1.Services[0].ServiceType = WebService{URL:"team1.local"}
	team1.Services[1].ServiceType = DomainNameService{DomainName:"team1.local"}
	team1.Services[2].ServiceType = FileTransferService{Username:"anonymous"}
	team1.Services[3].ServiceType = SecureShellService{Username:"ccdc", Password:"changeme", Command: "ls"}
	team1.Services[4].ServiceType = PingService{}

	t.Logf("team 1 services = %+v\n", team1.Services)

	teams := []Team{team1}
	game := NewGame(teams)
	t.Logf("game = %+v\n", game)
}
