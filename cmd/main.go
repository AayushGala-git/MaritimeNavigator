package main

import (
	"project3/pkg/common"
	"project3/pkg/groundstation"
	"project3/pkg/satellite"
	"project3/pkg/vessel"
)

func main() {

	// Load configuration
	err := common.LoadConfig("config.json") // Specify the configuration file path
	if err != nil {
		common.Logger.Fatal("Failed to load config:", err) // If loading fails, log the error and exit
	}

	// Activate the ground station server
	go groundstation.StartServer()

	// Activate the satellite services
	go satellite.RunSimulation("config.json")

	// Activate the vessel simulator
	go vessel.RunSimulation("config.json")

	// Blocking the main thread
	select {}
}
