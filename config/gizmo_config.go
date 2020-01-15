package config

import (
	structs "github.com/whoismissing/gizmo/structs"

	"io/ioutil"
	"encoding/json"
	"fmt"
)

func LoadServiceChecks(services []structs.Service) {
	for i := 0; i < len(services); i++ {
		serviceType := services[i].ServiceType
		serviceCheck := structs.LoadFromServiceType(serviceType)
		if services[i].Name == "" {
			services[i].Name = serviceType.Type
		}
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

func LoadConfig(filename string) []structs.Team {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to load config")
		panic(err)
	}

	teams := LoadTeams(config)

	return teams
}
