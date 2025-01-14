package satellite

import (
	"log"
	"time"

	"project3/pkg/common"
)

// RunSimulation sets up the satellite network simulation
func RunSimulation(configPath string) {
	// Load configuration
	err := common.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create a topology manager for the satellites
	manager := &TopologyManager{Satellites: make(map[string]*Satellite)}

	// Dynamically create satellites based on the configuration
	for _, satConfig := range common.AppConfig.Satellites {
		satellite := &Satellite{
			ID:                satConfig.ID,
			Port:              satConfig.Port,
			GroundStationAddr: common.AppConfig.GroundStationAddress,
			LatencyMap:        make(map[string]int),
			PacketLossMap:     make(map[string]float64),
			Status:            "Active",
		}
		manager.AddSatellite(satellite)
	}

	// Setup links between satellites
	for _, satConfig := range common.AppConfig.Satellites {
		sourceSatellite := manager.Satellites[satConfig.ID]
		for _, neighbor := range satConfig.Neighbors {
			manager.UpdateLink(sourceSatellite.ID, neighbor.ID, neighbor.Latency, neighbor.PacketLoss)
		}
	}

	// Start satellite listeners
	for _, satellite := range manager.Satellites {
		go satellite.Listen()
	}

	// Allow listeners to start
	time.Sleep(time.Second)
	log.Println("Satellite network simulation started successfully.")
}
