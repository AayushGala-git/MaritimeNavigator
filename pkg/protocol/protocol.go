package protocol

import (
	"time"
)

// MessageType Defining message types
type MessageType string

const (
	PositionUpdate    MessageType = "position"
	ForwardedPosition MessageType = "forwarded_position"
)

// PositionMessage Position Message Structure
type PositionMessage struct {
	Type      MessageType `json:"type"`
	VesselID  string      `json:"vesselID"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
	Timestamp time.Time   `json:"timestamp"`
}
