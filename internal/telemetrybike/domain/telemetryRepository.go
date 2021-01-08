package domain

// TelemetryBikeRepository definition to get telemetry data
type TelemetryBikeRepository interface {
	GetData() (*Telemetry, error)
}
