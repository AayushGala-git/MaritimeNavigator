package vessel

import (
	"fmt"
	"log"
	"project3/pkg/common"
	"project3/pkg/satellite"
	"sync"
)

// RunSimulation starts the vessel simulation process
func RunSimulation(configPath string) {
	// Load configuration
	err := common.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	err = common.ValidateConfig()
	if err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Create a topology manager for the satellites
	manager := &satellite.TopologyManager{Satellites: make(map[string]*satellite.Satellite)}

	// Dynamically create satellites based on the configuration
	for _, satConfig := range common.AppConfig.Satellites {
		sat := &satellite.Satellite{
			ID:            satConfig.ID,
			Port:          satConfig.Port,
			LatencyMap:    make(map[string]int),
			PacketLossMap: make(map[string]float64),
			Status:        "Active",
		}
		manager.AddSatellite(sat)
	}

	// Simulate vessels
	var wg sync.WaitGroup
	for _, vesselConfig := range common.AppConfig.Vessels {
		wg.Add(1)
		go func(vConfig common.VesselConfig) {
			defer wg.Done()
			satelliteAddress := fmt.Sprintf("127.0.0.1:%d", manager.Satellites[vConfig.Satellite].Port)
			SimulateVessel(vConfig.ID, satelliteAddress)
		}(vesselConfig)
	}

	// Wait for all vessel simulations to complete
	wg.Wait()
	log.Println("Vessel simulation completed successfully.")
}
