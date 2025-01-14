package groundstation

import (
	"encoding/json"
	"io"
	"net/http"
	"project3/pkg/common"
	"project3/pkg/satellite"
)

// StartServer starts the HTTP server for the ground station
func StartServer() {
	http.HandleFunc("/", handleSatelliteMessage) // Matches the root endpoint used by the satellite

	address := common.AppConfig.GroundStationAddress
	common.Logger.Printf("Ground station HTTP server started at %s\n", address)

	if err := http.ListenAndServe(address, nil); err != nil {
		common.Logger.Fatalf("Failed to start Ground Station server: %v\n", err)
	}
}

// handleSatelliteMessage handles incoming HTTP POST requests from satellites
func handleSatelliteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		common.Logger.Printf("Failed to read request body: %v\n", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var msg satellite.Message // Use the correct struct to unmarshal the satellite's payload
	if err := json.Unmarshal(body, &msg); err != nil {
		common.Logger.Printf("Failed to decode JSON message: %v\n", err)
		http.Error(w, "Failed to decode JSON message", http.StatusBadRequest)
		return
	}

	common.Logger.Printf("Received message: %+v\n", msg)

	// Example: Storing the message or processing it
	SaveToDatabase(msg)

	// Respond to the satellite to confirm receipt
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message received and processed"))
}
