package main

import (
	"log"
	"net/http"
	"os"

	"app/cmd/app/http/controller/static"
	"app/cmd/app/http/controller/telemetrybike"
	service "app/internal/telemetrybike/application/servetelemetrybike"
	"app/internal/telemetrybike/infraestructure/persistence"

	"github.com/joho/godotenv"
)

var serverPort string
var publicDir string
var simfileDir string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}

	// Get environment variables
	serverPort = os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "3000"
	}

	// Get environment variables
	publicDir = os.Getenv("PUBLIC_DIR")
	if publicDir == "" {
		publicDir = "public/static/"
	}

	simfileDir = os.Getenv("SIMFILE_DIR")
	if simfileDir == "" {
		simfileDir = "data/"
	}
}

func main() {

	// Init in memory file repository
	repository := persistence.NewInFileTelemetryRepository("data/simfile.json")
	fileContent, err := repository.GetData()
	if err != nil {
		log.Fatalf("failed to get data from file: %s", err)
	}

	// Init telemetry bike service
	telemetryBikeService := service.NewService(fileContent)

	// Configure http routes
	static.RegisterRoute(publicDir)
	telemetrybike.RegisterRoutes(telemetryBikeService)

	// Init http server
	log.Println("Listening on port :" + serverPort)
	err = http.ListenAndServe("0.0.0.0:"+serverPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
