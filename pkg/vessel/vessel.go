package vessel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"project3/pkg/protocol"
	"project3/pkg/satellite"
	"time"
)

// VesselSimulator defines a single vessel
type VesselSimulator struct {
	VesselID  string
	Latitude  float64
	Longitude float64
}

// SimulateVessel handles individual vessel simulation
func SimulateVessel(vesselID, satelliteAddress string) {
	log.Printf("Simulating vessel %s sending updates to satellite at %s\n", vesselID, satelliteAddress)

	vessel := VesselSimulator{
		VesselID:  vesselID,
		Latitude:  rand.Float64()*180 - 90,
		Longitude: rand.Float64()*360 - 180,
	}

	msgID := 1

	for {
		vessel.Latitude += (rand.Float64() - 0.5) * 0.1
		vessel.Longitude += (rand.Float64() - 0.5) * 0.1

		vessel.Latitude = clamp(vessel.Latitude, -90, 90)
		vessel.Longitude = clamp(vessel.Longitude, -180, 180)

		msg := satellite.Message{
			ID:          msgID,
			Source:      vessel.VesselID,
			Destination: "GroundStation",
			Content: protocol.PositionMessage{
				Type:      protocol.PositionUpdate,
				VesselID:  vessel.VesselID,
				Latitude:  vessel.Latitude,
				Longitude: vessel.Longitude,
				Timestamp: time.Now(),
			},
			Priority: rand.Intn(10),
			TTL:      5,
		}

		msgID++

		err := sendToSatellite(msg, satelliteAddress)
		if err != nil {
			log.Printf("Failed to send update from vessel %s: %v", vessel.VesselID, err)
		}

		time.Sleep(5 * time.Second)
	}
}

// sendToSatellite sends a simulated update from a vessel to a satellite
func sendToSatellite(msg satellite.Message, satelliteAddress string) error {
	// Serialize the message to JSON
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	// Send the JSON payload via HTTP POST
	resp, err := http.Post(fmt.Sprintf("http://%s", satelliteAddress), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send message to satellite: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	log.Printf("Sent message to satellite at %s: %+v", satelliteAddress, msg)
	return nil
}

// clamp ensures values are within specified bounds
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
