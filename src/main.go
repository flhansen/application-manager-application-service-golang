package main

import (
	"flhansen/application-manager/application-service/src/service"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	os.Exit(runApplication())
}

func runApplication() int {
	args := os.Args[1:]
	configPath := args[0]

	fileContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("An error occured while reading the configuration file %s: %v\n", configPath, err)
		return 1
	}

	var serviceConfig service.ApplicationServiceConfig
	if err := yaml.Unmarshal(fileContent, &serviceConfig); err != nil {
		fmt.Printf("An error occured while unmarshalling the configuration file content: %v\n", err)
		return 1
	}

	s, err := service.NewService(serviceConfig)
	if err != nil {
		fmt.Printf("An error occured while creating the service: %v\n", err)
		return 1
	}

	if err := s.Start(); err != nil {
		fmt.Printf("An error occured while starting the service: %v\n", err)
		return 1
	}

	return 0
}
