package servetelemetrybike

import (
	"time"

	"app/internal/telemetrybike/domain"
)

const delay = 100

// Service usecase method definitions
type Service interface {
	SendInitialContent(domain.TelemetryBikeCommunication) error
	ServeTelemetryBike(domain.TelemetryBikeCommunication, <-chan bool, <-chan bool, <-chan bool)
}

type service struct {
	telemetryBike *domain.Telemetry
}

// NewService return a new instance of serve telemetry bike usecase
func NewService(data *domain.Telemetry) Service {
	return &service{
		telemetryBike: data,
	}
}

func (s *service) SendInitialContent(c domain.TelemetryBikeCommunication) error {
	var initialContent WSCommandResponse = WSCommandResponse{
		Kind: "data",
		Data: s.telemetryBike.Data[0],
	}

	if err := c.Send(initialContent); err != nil {
		return err
	}

	echoResponse := WSEchoCommandResponse{
		Kind: "status",
		Data: Data{
			Status: "stop",
		},
	}

	if err := c.Send(echoResponse); err != nil {
		return err
	}

	return nil
}

func (s *service) ServeTelemetryBike(c domain.TelemetryBikeCommunication, play, stop, reset <-chan bool) {

	var index int = 1
	var serveTelemetryData bool = false

	for {
		select {
		case <-play:
			serveTelemetryData = true
		case <-stop:
			serveTelemetryData = false
		case <-reset:
			serveTelemetryData = false
			index = 1
			if err := s.SendInitialContent(c); err != nil {
				return
			}
		default:
			if serveTelemetryData {
				if index < len(s.telemetryBike.Data) {
					telemetryData := WSCommandResponse{
						Kind: "data",
						Data: s.telemetryBike.Data[index],
					}
					if err := c.Send(telemetryData); err != nil {
						return
					}
					index++
				} else {
					serveTelemetryData = false
					stopNotify := WSEchoCommandResponse{
						Kind: "status",
						Data: Data{
							Status: "stop",
						},
					}
					if err := c.Send(stopNotify); err != nil {
						return
					}
				}

			}
		}
		time.Sleep(delay * time.Millisecond)
	}
}
