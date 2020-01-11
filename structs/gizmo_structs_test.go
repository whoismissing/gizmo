package structs

import (
	"testing"
	"encoding/json"
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

	// This is a hack to add typing to JSON because the interface info is lost
	wwwJson, _ := json.Marshal(team1.Services[0].ServiceType)
	wwwTypedJson := TypedJson{"www", wwwJson}

	dnsJson, _ := json.Marshal(team1.Services[1].ServiceType)
	dnsTypedJson := TypedJson{"dns", dnsJson}

	ftpJson, _ := json.Marshal(team1.Services[2].ServiceType)
	ftpTypedJson := TypedJson{"ftp", ftpJson}

	sshJson, _ := json.Marshal(team1.Services[3].ServiceType)
	sshTypedJson := TypedJson{"ssh", sshJson}

	pingJson, _ := json.Marshal(team1.Services[4].ServiceType)
	pingTypedJson := TypedJson{"ping", pingJson}

	team1.Services[0].ObjectLoad = wwwTypedJson
	team1.Services[1].ObjectLoad = dnsTypedJson
	team1.Services[2].ObjectLoad = ftpTypedJson
	team1.Services[3].ObjectLoad = sshTypedJson
	team1.Services[4].ObjectLoad = pingTypedJson

	finalJson, _ := json.Marshal(teams)
	t.Logf(string(finalJson))
}
