// Package config provides the methods for obtaining the structs described in
// package structs by reading a user-specified config filename.
package config

import (
	structs "github.com/whoismissing/gizmo/structs"

	"encoding/json"
	"fmt"
	"io/ioutil"
)

// LoadServiceChecks() obtains the ServiceCheck objects for each Service in
// services assuming their ServiceType has been obtained.
func LoadServiceChecks(services []structs.Service) {
	for i := 0; i < len(services); i++ {
		// ServiceType is a wrapper for ServiceCheck as ServiceCheck is an
		// interface whose type information is lost when dumped to JSON
		// This is a hack!
		serviceType := services[i].ServiceType
		serviceCheck := structs.LoadFromServiceType(serviceType)
		if services[i].Name == "" {
			services[i].Name = serviceType.Type
		}
		services[i].ServiceCheck = serviceCheck
	}
}

// LoadTeams() unmarshals an array of Team objects given the raw bytes of the
// config data provided.
func LoadTeams(config []byte) []structs.Team {
	var teams []structs.Team
	_ = json.Unmarshal(config, &teams)

	for i := 0; i < len(teams); i++ {
		services := teams[i].Services
		LoadServiceChecks(services)
	}

	return teams
}

// LoadConfig() reads the data of a user-specified filename and unmarshals the
// data, returning complete list of Team objects and corresponding Service and
// ServiceCheck objects.
func LoadConfig(filename string) []structs.Team {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to load config")
		panic(err)
	}

	teams := LoadTeams(config)

	return teams
}
