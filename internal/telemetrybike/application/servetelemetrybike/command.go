package servetelemetrybike

import "app/internal/telemetrybike/domain"

// ClientCommand message
type ClientCommand struct {
	Command string `json:"command"`
}

// WSCommandResponse server data message
type WSCommandResponse struct {
	Kind string               `json:"kind"`
	Data domain.TelemetryBike `json:"data"`
}

// WSEchoCommandResponse server echo response message
type WSEchoCommandResponse struct {
	Kind string `json:"kind"`
	Data Data   `json:"data"`
}

// Data command
type Data struct {
	Status string `json:"status"`
}
