package telemetrybike

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	controller "app/cmd/app/http/controller/telemetrybike"
	service "app/internal/telemetrybike/application/servetelemetrybike"
	"app/internal/telemetrybike/infraestructure/persistence"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTelemetryBikeCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Telemetry Bike Suite")
}

func setupServer(srv service.Service) *http.ServeMux {
	handler := http.NewServeMux()
	handler.HandleFunc("/replay", controller.WebSocketController(srv))
	return handler
}

var _ = Describe("Server", func() {
	var server *httptest.Server
	var wsocket *websocket.Conn

	BeforeEach(func() {
		repository := persistence.NewInFileTelemetryRepository("simfile_mock.json")
		fileContent, err := repository.GetData()
		Expect(err).ShouldNot(HaveOccurred())

		telemetryBikeService := service.NewService(fileContent)

		server = httptest.NewServer(setupServer(telemetryBikeService))
		u := "ws" + strings.TrimPrefix(server.URL, "http") + "/replay"
		wsocket, _, err = websocket.DefaultDialer.Dial(u, nil)
		Expect(err).ShouldNot(HaveOccurred())

	})

	AfterEach(func() {
		wsocket.Close()
		server.Close()
	})

	Context("When client connect to websocket in ws://HOST/reply", func() {
		It("Returns echo message when commands received", func() {
			// Omit initial data and stop command from server
			wsocket.ReadMessage()
			wsocket.ReadMessage()

			err := wsocket.WriteMessage(websocket.TextMessage, []byte(`{"command":"reset"}`))
			Expect(err).ShouldNot(HaveOccurred())

			_, stopCommand, err := wsocket.ReadMessage()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(stopCommand)).To(MatchJSON(`{
				"kind": "echo",
				"data": {
					"status": "reset"
				}
			}`))
		})

		It("Returns initial data and stop command json", func() {
			_, initialData, err := wsocket.ReadMessage()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(initialData)).To(MatchJSON(`{
				"kind": "data",
				"data": {
					"time": "09:01:00.011",
					"gear": "1",
					"rpm": 10895,
					"speed": 0
				}
			}`))
			_, stopCommand, err := wsocket.ReadMessage()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(stopCommand)).To(MatchJSON(`{
				"kind": "status",
				"data": {
					"status": "stop"
				}
			}`))
		})

		It("begins sending data if client sends PLAY command", func() {
			// Omit initial data and stop command from server
			wsocket.ReadMessage()
			wsocket.ReadMessage()

			err := wsocket.WriteMessage(websocket.TextMessage, []byte(`{"command":"play"}`))
			Expect(err).ShouldNot(HaveOccurred())
			// Omit echo PLAY message
			wsocket.ReadMessage()

			_, firstData, err := wsocket.ReadMessage()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(firstData)).To(MatchJSON(`{
				"kind": "data",
				"data": {
					"time": "09:01:00.053",
					"gear": "1",
					"rpm": 11012,
					"speed": 0
				}
			}`))
			_, secondData, err := wsocket.ReadMessage()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(secondData)).To(MatchJSON(`{
				"kind": "data",
				"data": {
					"time": "09:01:00.586",
					"gear": "1",
					"rpm": 10930,
					"speed": 0
				}
			}`))
		})

		It("reset internal state if client sends reset command", func() {
			// Omit initial data and stop command from server
			wsocket.ReadMessage()
			wsocket.ReadMessage()

			err := wsocket.WriteMessage(websocket.TextMessage, []byte(`{"command":"play"}`))
			Expect(err).ShouldNot(HaveOccurred())
			// Omit echo PLAY message
			wsocket.ReadMessage()

			_, firstData, err := wsocket.ReadMessage()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(firstData)).To(MatchJSON(`{
				"kind": "data",
				"data": {
					"time": "09:01:00.053",
					"gear": "1",
					"rpm": 11012,
					"speed": 0
				}
			}`))

			err = wsocket.WriteMessage(websocket.TextMessage, []byte(`{"command":"reset"}`))
			Expect(err).ShouldNot(HaveOccurred())
			// Omit echo RESET message
			wsocket.ReadMessage()

			_, initialData, err := wsocket.ReadMessage()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(initialData)).To(MatchJSON(`{
				"kind": "data",
				"data": {
					"time": "09:01:00.011",
					"gear": "1",
					"rpm": 10895,
					"speed": 0
				}
			}`))
		})
	})
})
