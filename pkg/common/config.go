package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// SatelliteConfig holds individual satellite configuration
type SatelliteConfig struct {
	ID        string           `json:"id"`
	Port      int              `json:"port"`
	Neighbors []NeighborConfig `json:"neighbors"`
}

// NeighborConfig defines a satellite's connection to its neighbors
type NeighborConfig struct {
	ID         string  `json:"id"`
	Latency    int     `json:"latency"`
	PacketLoss float64 `json:"packet_loss"`
}

// VesselConfig defines a vessel and its associated satellite
type VesselConfig struct {
	ID        string `json:"id"`
	Satellite string `json:"satellite"` // Associated satellite ID
}

// Config holds the overall configuration
type Config struct {
	GroundStationAddress string            `json:"ground_station_address"`
	Satellites           []SatelliteConfig `json:"satellites"`
	Vessels              []VesselConfig    `json:"vessels"`
}

var (
	// AppConfig global configuration
	AppConfig Config
)

// LoadConfig loads configuration from a JSON file into AppConfig
func LoadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read configuration file at %s: %v", path, err)
		return err
	}

	err = json.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Printf("Failed to parse configuration JSON: %v", err)
		return err
	}

	log.Printf("Configuration successfully loaded from %s", path)
	return nil
}

// ValidateConfig ensures the configuration is valid and complete
func ValidateConfig() error {
	if AppConfig.GroundStationAddress == "" {
		return fmt.Errorf("ground station address is missing")
	}
	if len(AppConfig.Satellites) == 0 {
		return fmt.Errorf("no satellites configured")
	}
	for _, satellite := range AppConfig.Satellites {
		if satellite.ID == "" {
			return fmt.Errorf("a satellite is missing an ID")
		}
		if satellite.Port == 0 {
			return fmt.Errorf("satellite %s is missing a port", satellite.ID)
		}
		for _, neighbor := range satellite.Neighbors {
			if neighbor.ID == "" {
				return fmt.Errorf("satellite %s has a neighbor with a missing ID", satellite.ID)
			}
			if neighbor.Latency <= 0 {
				return fmt.Errorf("invalid latency between satellite %s and neighbor %s", satellite.ID, neighbor.ID)
			}
			if neighbor.PacketLoss < 0 || neighbor.PacketLoss > 1 {
				return fmt.Errorf("invalid packet loss rate between satellite %s and neighbor %s", satellite.ID, neighbor.ID)
			}
		}
	}
	if len(AppConfig.Vessels) == 0 {
		return fmt.Errorf("no vessels configured")
	}
	for _, vessel := range AppConfig.Vessels {
		if vessel.ID == "" {
			return fmt.Errorf("a vessel is missing an ID")
		}
		if vessel.Satellite == "" {
			return fmt.Errorf("vessel %s is missing an associated satellite", vessel.ID)
		}
	}
	return nil
}
