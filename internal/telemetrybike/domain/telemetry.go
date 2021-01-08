package domain

// Telemetry definition
type Telemetry struct {
	StartTime string          `json:"startTime"`
	Data      []TelemetryBike `json:"data"`
}

// TelemetryBike definition
type TelemetryBike struct {
	Time  string `json:"time"`
	Gear  string `json:"gear"`
	RPM   int    `json:"rpm"`
	Speed int    `json:"speed"`
}
