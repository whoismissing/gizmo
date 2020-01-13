package config

import (
	"io/ioutil"
	"encoding/json"
	structs "github.com/whoismissing/gizmo/structs"
	"fmt"
)

func LoadServiceChecks(services []structs.Service) {
	for i := 0; i < len(services); i++ {
		serviceType := services[i].ServiceType
		serviceCheck := structs.LoadFromServiceType(serviceType)
		services[i].ServiceCheck = serviceCheck
	}
}

func LoadTeams(config []byte) []structs.Team {
	var teams []structs.Team
	_ = json.Unmarshal(config, &teams)

	for i := 0; i < len(teams); i++ {
		services := teams[i].Services
		LoadServiceChecks(services)
	}

	return teams
}

func LoadConfig(filename string) {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to load config")
		panic(err)
	}

	teams := LoadTeams(config)

	team1 := teams[0]
	fmt.Printf("services = %+v\n", team1.Services)
}
