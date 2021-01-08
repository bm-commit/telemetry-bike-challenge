package domain

// TelemetryBikeCommunication definition to send data to the client
type TelemetryBikeCommunication interface {
	Send(interface{}) error
}
