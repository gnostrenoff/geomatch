package benchmarks

import (
	"testing"

	"github.com/gnostenoff/geomatch/internal/geomatch"
	"github.com/gnostenoff/geomatch/internal/geomatch/euclidean"
	"github.com/gnostenoff/geomatch/internal/geomatch/haversine"
	"github.com/gnostenoff/geomatch/internal/pkg/datasource"
)

type mockDataSource struct {
	events []datasource.Event
}

func (m *mockDataSource) Get() ([]datasource.Event, error) {
	return m.events, nil
}

func BenchmarkMatcher_Euclidean(b *testing.B) {
	events := generateBenchmarkEvents(1000)
	pois := generateBenchmarkPOIs(50)

	ds := &mockDataSource{events: events}
	matcher := euclidean.NewMatcher(ds)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = matcher.Match(pois)
	}
}

func BenchmarkMatcher_Haversine(b *testing.B) {
	events := generateBenchmarkEvents(1000)
	pois := generateBenchmarkPOIs(50)

	ds := &mockDataSource{events: events}
	matcher := haversine.NewMatcher(ds)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = matcher.Match(pois)
	}
}

func generateBenchmarkEvents(count int) []datasource.Event {
	events := make([]datasource.Event, count)
	for i := 0; i < count; i++ {
		events[i] = datasource.Event{
			Lat:  48.8566 + float64(i%10)*0.001,
			Lon:  2.3522 + float64(i%10)*0.001,
			Type: []string{"imp", "click"}[i%2],
		}
	}
	return events
}

func generateBenchmarkPOIs(count int) []geomatch.PointOfInterest {
	pois := make([]geomatch.PointOfInterest, count)
	for i := 0; i < count; i++ {
		pois[i] = geomatch.PointOfInterest{
			Name: "POI " + string(rune(i)),
			Lat:  48.8566 + float64(i)*0.01,
			Lon:  2.3522 + float64(i)*0.01,
		}
	}
	return pois
}
