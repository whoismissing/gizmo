package structs

import (
	"encoding/json"
	"testing"
)

func TestInitializeStructs(t *testing.T) {
	t.Logf("TestInitializeStructs SUCCESS")

	var teamID uint = 1
	team1 := NewTeam(teamID)

	team1.Services = NewServices(6, teamID)
	team1.Services[0].ServiceCheck = WebService{URL: "team1.local"}
	team1.Services[1].ServiceCheck = DomainNameService{DomainName: "team1.local"}
	team1.Services[2].ServiceCheck = FileTransferService{Username: "anonymous"}
	team1.Services[3].ServiceCheck = SecureShellService{Username: "ccdc", Password: "changeme", Command: "ls"}
	team1.Services[4].ServiceCheck = PingService{}
	team1.Services[5].ServiceCheck = ExternalService{Filepath: "ls"}

	t.Logf("team 1 services = %+v\n", team1.Services)

	teams := []Team{team1}
	game := NewGame(teams)
	t.Logf("game = %+v\n", game)

	// This is a hack to add typing to JSON because the interface info is lost
	wwwJson, _ := json.Marshal(team1.Services[0].ServiceCheck)
	wwwServiceType := ServiceType{"www", wwwJson}

	dnsJson, _ := json.Marshal(team1.Services[1].ServiceCheck)
	dnsServiceType := ServiceType{"dns", dnsJson}

	ftpJson, _ := json.Marshal(team1.Services[2].ServiceCheck)
	ftpServiceType := ServiceType{"ftp", ftpJson}

	sshJson, _ := json.Marshal(team1.Services[3].ServiceCheck)
	sshServiceType := ServiceType{"ssh", sshJson}

	pingJson, _ := json.Marshal(team1.Services[4].ServiceCheck)
	pingServiceType := ServiceType{"ping", pingJson}

	extJson, _ := json.Marshal(team1.Services[5].ServiceCheck)
	extServiceType := ServiceType{"ext", extJson}

	team1.Services[0].ServiceType = wwwServiceType
	team1.Services[1].ServiceType = dnsServiceType
	team1.Services[2].ServiceType = ftpServiceType
	team1.Services[3].ServiceType = sshServiceType
	team1.Services[4].ServiceType = pingServiceType
	team1.Services[5].ServiceType = extServiceType

	finalJson, _ := json.Marshal(teams)
	t.Logf(string(finalJson))
}
