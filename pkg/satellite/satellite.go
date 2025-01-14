package satellite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"project3/pkg/protocol"
	"sync"
	"time"
)

// Satellite represents a satellite node
type Satellite struct {
	ID                string
	Port              int
	Neighbors         []*Satellite
	LatencyMap        map[string]int
	PacketLossMap     map[string]float64
	Status            string // "Active" or "Failed"
	GroundStationAddr string
	mu                sync.Mutex
}

// Message represents a communication message with TTL
type Message struct {
	ID          int                      `json:"id"`
	Source      string                   `json:"source"`
	Destination string                   `json:"destination"`
	Content     protocol.PositionMessage `json:"content"` // Change here
	Priority    int                      `json:"priority"`
	TTL         int                      `json:"ttl"`
}

// Listen starts the satellite HTTP server
func (s *Satellite) Listen() error {
	// Create a new ServeMux for this satellite
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
			return
		}

		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			log.Printf("Satellite %s failed to decode message: %v", s.ID, err)
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Log the received message
		log.Printf("Satellite %s received message: %+v", s.ID, msg)

		// Forward message if not the destination and TTL > 0
		if msg.Destination != s.ID && msg.TTL > 0 {
			msg.TTL-- // Decrement TTL
			go s.ForwardMessage(&msg, make(map[string]bool))
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Message received successfully")
	})

	address := fmt.Sprintf(":%d", s.Port)
	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	log.Printf("Satellite %s listening on %s", s.ID, address)
	return server.ListenAndServe()
}

// ForwardMessage forwards a message to neighbors or the ground station
func (s *Satellite) ForwardMessage(msg *Message, visited map[string]bool) {
	s.mu.Lock()
	if visited[s.ID] {
		s.mu.Unlock()
		return
	}
	visited[s.ID] = true
	s.mu.Unlock()

	if msg.TTL <= 0 {
		fmt.Printf("Message expired at Satellite %s. Stopping forwarding.\n", s.ID)
		return
	}

	// Forward to ground station if the destination is "GroundStation"
	if msg.Destination == "GroundStation" {
		go func() {
			err := sendToGroundStation(s.GroundStationAddr, msg)
			if err != nil {
				fmt.Printf("Failed to send message to Ground Station: %v\n", err)
			} else {
				fmt.Printf("Message successfully sent to Ground Station from Satellite %s\n", s.ID)
			}
		}()
		return
	}

	// Forward to neighbors
	for _, neighbor := range s.Neighbors {
		if neighbor.Status == "Failed" {
			fmt.Printf("Neighbor Satellite %s is down. Skipping...\n", neighbor.ID)
			continue
		}

		latency := s.LatencyMap[neighbor.ID]
		packetLoss := s.PacketLossMap[neighbor.ID]

		go func(neighbor *Satellite, latency int, packetLoss float64) {
			time.Sleep(time.Duration(latency) * time.Millisecond)
			if rand.Float64() > packetLoss {
				url := fmt.Sprintf("http://localhost:%d", neighbor.Port)
				body, err := json.Marshal(msg)
				if err != nil {
					fmt.Printf("Failed to encode message for Satellite %s: %v\n", neighbor.ID, err)
					return
				}

				resp, err := http.Post(url, "application/json", bytes.NewReader(body))
				if err != nil {
					fmt.Printf("Failed to send message to Satellite %s: %v\n", neighbor.ID, err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusOK {
					fmt.Printf("Message successfully sent from %s to %s (TTL: %d)\n", s.ID, neighbor.ID, msg.TTL)
				} else {
					fmt.Printf("Satellite %s returned status %d\n", neighbor.ID, resp.StatusCode)
				}
			} else {
				fmt.Printf("Message lost between %s and %s\n", s.ID, neighbor.ID)
			}
		}(neighbor, latency, packetLoss)
	}
}

// sendToGroundStation sends a message to the ground station using HTTP
func sendToGroundStation(address string, msg *Message) error {
	url := fmt.Sprintf("http://%s", address)
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send message to Ground Station: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ground Station returned status %d", resp.StatusCode)
	}

	return nil
}
