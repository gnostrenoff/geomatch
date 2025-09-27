package haversine

import (
	"log/slog"
	"math"

	"github.com/gnostenoff/geomatch/internal/geomatch"
	"github.com/gnostenoff/geomatch/internal/pkg/datasource"
)

type Matcher struct {
	datasource datasource.EventDataSource
}

func NewMatcher(datasource datasource.EventDataSource) *Matcher {
	return &Matcher{
		datasource: datasource,
	}
}

func (m *Matcher) Match(pois []geomatch.PointOfInterest) ([]geomatch.MatchResult, error) {
	events, err := m.datasource.Get()
	if err != nil {
		slog.Error("failed to get events from datasource", slog.String("error", err.Error()))
		return nil, err
	}

	if len(pois) == 0 || len(events) == 0 {
		return []geomatch.MatchResult{}, nil
	}

	poiStats := make(map[int]*geomatch.MatchResult, len(pois))
	for i, poi := range pois {
		poiStats[i] = &geomatch.MatchResult{
			PointOfInterest: poi,
			Impressions:     0,
			Clicks:          0,
		}
	}

	for _, event := range events {
		closestPOIIndex := findClosestPOI(event.Lat, event.Lon, pois)

		if closestPOIIndex >= 0 {
			result := poiStats[closestPOIIndex]
			switch event.Type {
			case "imp":
				result.Impressions++
			case "click":
				result.Clicks++
			}
		}
	}

	results := make([]geomatch.MatchResult, 0, len(pois))
	for i := 0; i < len(pois); i++ {
		results = append(results, *poiStats[i])
	}

	return results, nil
}

func findClosestPOI(eventLat, eventLon float64, pois []geomatch.PointOfInterest) int {
	if len(pois) == 0 {
		return -1
	}

	minDistance := math.MaxFloat64
	closestIndex := -1

	for i, poi := range pois {
		distance := haversineDistance(eventLat, eventLon, poi.Lat, poi.Lon)
		if distance < minDistance {
			minDistance = distance
			closestIndex = i
		}
	}

	return closestIndex
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}
