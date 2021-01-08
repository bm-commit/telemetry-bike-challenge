package persistence

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"app/internal/telemetrybike/domain"
)

const timeLayout = "2006-01-02T15:04:05.999+01:00"

type repository struct {
	Filename string
}

// NewInFileTelemetryRepository return a new instance to load data from I/O
func NewInFileTelemetryRepository(filename string) domain.TelemetryBikeRepository {
	return &repository{
		Filename: filename,
	}
}

func (r *repository) GetData() (*domain.Telemetry, error) {
	data, err := ioutil.ReadFile(r.Filename)
	if err != nil {
		log.Fatalf("failed reading data from file: %s", err)
		return nil, err
	}
	var fileContent domain.Telemetry
	err = json.Unmarshal([]byte(data), &fileContent)
	if err != nil {
		log.Fatalf("failed to unmarshal data from file: %s", err)
		return nil, err
	}

	for i, entry := range fileContent.Data {
		t, err := time.Parse(timeLayout, entry.Time)
		if err != nil {
			log.Fatalf("failed to parse time entry from file: %s", err)
			return nil, err
		}
		fileContent.Data[i].Time = t.Format("15:04:05.000")
	}

	return &fileContent, nil
}
