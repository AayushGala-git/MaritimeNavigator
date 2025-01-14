package satellite

import (
	"fmt"
	"sync"
)

// TopologyManager manages the satellite network
type TopologyManager struct {
	Satellites map[string]*Satellite
	mu         sync.Mutex
}

// AddSatellite adds a new satellite to the topology
func (t *TopologyManager) AddSatellite(newSatellite *Satellite) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Satellites[newSatellite.ID] = newSatellite
	fmt.Printf("Satellite %s added to the topology.\n", newSatellite.ID)
}

// RemoveSatellite removes an existing satellite from the topology
func (t *TopologyManager) RemoveSatellite(satelliteID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.Satellites[satelliteID]; exists {
		delete(t.Satellites, satelliteID)
		fmt.Printf("Satellite %s removed from the topology.\n", satelliteID)
	} else {
		fmt.Printf("Satellite %s does not exist.\n", satelliteID)
	}
}

// UpdateLink updates the latency or packet loss between satellites
func (t *TopologyManager) UpdateLink(sourceID, targetID string, latency int, packetLoss float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	source, sourceExists := t.Satellites[sourceID]
	target, targetExists := t.Satellites[targetID]

	if sourceExists && targetExists {
		source.LatencyMap[targetID] = latency
		source.PacketLossMap[targetID] = packetLoss
		source.Neighbors = append(source.Neighbors, target)

		target.LatencyMap[sourceID] = latency
		target.PacketLossMap[sourceID] = packetLoss
		target.Neighbors = append(target.Neighbors, source)

		fmt.Printf("Link updated: %s <-> %s, Latency: %dms, PacketLoss: %.2f\n", sourceID, targetID, latency, packetLoss)
	} else {
		if !sourceExists {
			fmt.Printf("Source satellite %s does not exist.\n", sourceID)
		}
		if !targetExists {
			fmt.Printf("Target satellite %s does not exist.\n", targetID)
		}
	}
}
