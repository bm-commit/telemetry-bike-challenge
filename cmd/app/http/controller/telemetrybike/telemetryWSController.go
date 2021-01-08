package telemetrybike

import (
	"net/http"

	usecase "app/internal/telemetrybike/application/servetelemetrybike"
	"app/internal/telemetrybike/infraestructure/communication"

	"github.com/gorilla/websocket"
)

// Allows to convert HTTP connection to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocketController handler
func WebSocketController(uc usecase.Service) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		conn, _ := upgrader.Upgrade(res, req, nil) // Note: Connections support one concurrent reader and one concurrent writer.

		defer conn.Close()

		comm := communication.NewWSCommunication(conn)

		uc.SendInitialContent(comm)

		play := make(chan bool)
		stop := make(chan bool)
		reset := make(chan bool)

		go uc.ServeTelemetryBike(comm, play, stop, reset)

		for {
			// Read message from frontend client
			var msg usecase.ClientCommand
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}

			echoResponse := usecase.WSEchoCommandResponse{
				Kind: "echo",
				Data: usecase.Data{
					Status: msg.Command,
				},
			}

			// Write echo message back to frontend client
			if err := comm.Send(echoResponse); err != nil {
				return
			}

			switch msg.Command {
			case "play":
				play <- true
			case "stop":
				stop <- true
			case "reset":
				reset <- true
			}
		}
	}
}
