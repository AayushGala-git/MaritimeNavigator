package groundstation

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"project3/pkg/common"
	"project3/pkg/protocol"
	"project3/pkg/satellite"
	"sync"
)

// Mutex to ensure thread-safe file writes
var fileMutex sync.Mutex

// SaveToDatabase saves the message to a local JSON file
func SaveToDatabase(msg satellite.Message) {
	// Serialize the message
	data, err := json.Marshal(msg)
	if err != nil {
		common.Logger.Println("Failed to marshal message to JSON:", err)
		return
	}

	// Thread-safe file access
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// Open file in append mode or create if it doesn't exist
	file, err := os.OpenFile("database.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		common.Logger.Println("Failed to open database file:", err)
		return
	}
	defer file.Close()

	// Write the serialized message to the file
	_, err = file.Write(append(data, '\n'))
	if err != nil {
		common.Logger.Println("Failed to write to database:", err)
	} else {
		common.Logger.Println("Message successfully saved to database:", msg)
	}
}

// LoadFromDatabase Load all data (available for API interface for scalable functionality)
func LoadFromDatabase() ([]protocol.PositionMessage, error) {
	data, err := ioutil.ReadFile("database.json")
	if err != nil {
		return nil, err
	}
	var messages []protocol.PositionMessage
	lines := splitLines(string(data))
	for _, line := range lines {
		var msg protocol.PositionMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			common.Logger.Println("Failed to unmarshal line:", err)
			continue
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func splitLines(s string) []string {
	var lines []string
	var line string
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(r)
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}
