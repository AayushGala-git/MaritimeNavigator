package protocol

import (
	"encoding/json"
	"net"
	"project3/pkg/common"
)

// HandleVesselMessage handling messages from vessel
func HandleVesselMessage(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	var msg PositionMessage
	if err := decoder.Decode(&msg); err != nil {
		common.Logger.Println("Failed to decode message:", err)
		return
	}
	// Forward to ground station
	ForwardToGroundStation(msg)
}

// ForwardToGroundStation Forward the message to the ground station
func ForwardToGroundStation(msg PositionMessage) {
	msg.Type = ForwardedPosition
	data, err := json.Marshal(msg)
	if err != nil {
		common.Logger.Println("Failed to marshal message:", err)
		return
	}
	conn, err := net.Dial("tcp", common.AppConfig.GroundStationAddress)
	if err != nil {
		common.Logger.Println("Failed to connect to ground station:", err)
		return
	}
	defer conn.Close()
	conn.Write(data)
}
