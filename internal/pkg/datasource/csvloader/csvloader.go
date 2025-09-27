package csvloader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gnostenoff/geomatch/internal/pkg/datasource"
)

/*
CSVLoader implements the EventDataSource interface to load events from a CSV file.
*/

type EventCSVLoader struct {
	filePath string
	events   []datasource.Event
}

func NewEventCSVLoader(filePath string) *EventCSVLoader {
	return &EventCSVLoader{filePath: filePath}
}

func (l *EventCSVLoader) Get() ([]datasource.Event, error) {
	return l.events, nil
}

func (l *EventCSVLoader) Load() error {
	file, err := os.Open(l.filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file %s: %w", l.filePath, err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	latIndex, lonIndex, typeIndex := -1, -1, -1
	for i, field := range header {
		switch field {
		case "lat":
			latIndex = i
		case "lon":
			lonIndex = i
		case "event_type":
			typeIndex = i
		}
	}

	if latIndex == -1 || lonIndex == -1 || typeIndex == -1 {
		return fmt.Errorf("CSV header missing required fields: need lat, lon, event_type")
	}

	var events []datasource.Event
	lineNum := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record at line %d: %w", lineNum+1, err)
		}

		lat, err := strconv.ParseFloat(record[latIndex], 64)
		if err != nil {
			return fmt.Errorf("invalid latitude at line %d: %w", lineNum+1, err)
		}

		lon, err := strconv.ParseFloat(record[lonIndex], 64)
		if err != nil {
			return fmt.Errorf("invalid longitude at line %d: %w", lineNum+1, err)
		}

		events = append(events, datasource.Event{
			Lat:  lat,
			Lon:  lon,
			Type: record[typeIndex],
		})

		lineNum++
	}

	l.events = events
	return nil
}
