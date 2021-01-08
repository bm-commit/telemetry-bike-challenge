package telemetrybike

import (
	"net/http"

	service "app/internal/telemetrybike/application/servetelemetrybike"
)

// RegisterRoutes of websocket telemetry bike
func RegisterRoutes(serv service.Service) {
	http.HandleFunc("/replay", WebSocketController(serv))
}
