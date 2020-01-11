package config

import (
	"io/ioutil"
	"encoding/json"
	structs "github.com/whoismissing/gizmo/structs"
	"fmt"
)

func UnmarshalTypedJson(structs.TypedJson) structs.ServiceType {

}

func LoadConfig(filename string) {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Failed to load config")
		panic(err)
	}

	var teams []structs.Team
	_ = json.Unmarshal(config, &teams)

	// TODO: handle unmarshalling of ServiceType interface

	/*
	if err != nil {
		fmt.Println("Failed to unmarshal JSON")
		panic(err)
	}
	*/

	fmt.Printf("%+v\n", teams)
}
