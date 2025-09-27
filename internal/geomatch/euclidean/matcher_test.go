package euclidean

import (
	"errors"
	"math"
	"testing"

	"github.com/gnostenoff/geomatch/internal/geomatch"
	"github.com/gnostenoff/geomatch/internal/pkg/datasource"
)

type mockDataSource struct {
	events []datasource.Event
	err    error
}

func (m *mockDataSource) Get() ([]datasource.Event, error) {
	return m.events, m.err
}

func (m *mockDataSource) Load() error {
	return nil
}

func TestMatcher_Match(t *testing.T) {
	tests := []struct {
		name    string
		events  []datasource.Event
		pois    []geomatch.PointOfInterest
		dsErr   error
		want    []geomatch.MatchResult
		wantErr bool
	}{
		{
			name: "single POI with multiple events",
			events: []datasource.Event{
				{Lat: 48.8566, Lon: 2.3522, Type: "imp"},
				{Lat: 48.8567, Lon: 2.3523, Type: "click"},
				{Lat: 48.8568, Lon: 2.3524, Type: "imp"},
			},
			pois: []geomatch.PointOfInterest{
				{Name: "Eiffel Tower", Lat: 48.8566, Lon: 2.3522},
			},
			want: []geomatch.MatchResult{
				{
					PointOfInterest: geomatch.PointOfInterest{Name: "Eiffel Tower", Lat: 48.8566, Lon: 2.3522},
					Impressions:     2,
					Clicks:          1,
				},
			},
		},
		{
			name: "multiple POIs closest match",
			events: []datasource.Event{
				{Lat: 48.8566, Lon: 2.3522, Type: "imp"},
				{Lat: 40.7589, Lon: -73.9851, Type: "click"},
			},
			pois: []geomatch.PointOfInterest{
				{Name: "Eiffel Tower", Lat: 48.8566, Lon: 2.3522},
				{Name: "Times Square", Lat: 40.7589, Lon: -73.9851},
			},
			want: []geomatch.MatchResult{
				{
					PointOfInterest: geomatch.PointOfInterest{Name: "Eiffel Tower", Lat: 48.8566, Lon: 2.3522},
					Impressions:     1,
					Clicks:          0,
				},
				{
					PointOfInterest: geomatch.PointOfInterest{Name: "Times Square", Lat: 40.7589, Lon: -73.9851},
					Impressions:     0,
					Clicks:          1,
				},
			},
		},
		{
			name:   "empty events",
			events: []datasource.Event{},
			pois: []geomatch.PointOfInterest{
				{Name: "Eiffel Tower", Lat: 48.8566, Lon: 2.3522},
			},
			want: []geomatch.MatchResult{},
		},
		{
			name: "empty POIs",
			events: []datasource.Event{
				{Lat: 48.8566, Lon: 2.3522, Type: "imp"},
			},
			pois: []geomatch.PointOfInterest{},
			want: []geomatch.MatchResult{},
		},
		{
			name: "unknown event type",
			events: []datasource.Event{
				{Lat: 48.8566, Lon: 2.3522, Type: "unknown"},
			},
			pois: []geomatch.PointOfInterest{
				{Name: "Eiffel Tower", Lat: 48.8566, Lon: 2.3522},
			},
			want: []geomatch.MatchResult{
				{
					PointOfInterest: geomatch.PointOfInterest{Name: "Eiffel Tower", Lat: 48.8566, Lon: 2.3522},
					Impressions:     0,
					Clicks:          0,
				},
			},
		},
		{
			name:    "datasource error",
			events:  nil,
			pois:    []geomatch.PointOfInterest{{Name: "Test", Lat: 0, Lon: 0}},
			dsErr:   errors.New("datasource failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &mockDataSource{events: tt.events, err: tt.dsErr}
			m := NewMatcher(ds)

			got, err := m.Match(tt.pois)

			if (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("Match() got %d results, want %d", len(got), len(tt.want))
					return
				}

				for i, result := range got {
					expected := tt.want[i]
					if result.Name != expected.Name ||
						result.Lat != expected.Lat ||
						result.Lon != expected.Lon ||
						result.Impressions != expected.Impressions ||
						result.Clicks != expected.Clicks {
						t.Errorf("Match() result[%d] = %+v, want %+v", i, result, expected)
					}
				}
			}
		})
	}
}

func TestMatcher_euclideanDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
		delta    float64
	}{
		{
			name:     "same point",
			lat1:     48.8566,
			lon1:     2.3522,
			lat2:     48.8566,
			lon2:     2.3522,
			expected: 0.0,
			delta:    0.001,
		},
		{
			name:     "simple coordinates",
			lat1:     0.0,
			lon1:     0.0,
			lat2:     3.0,
			lon2:     4.0,
			expected: 5.0,
			delta:    0.001,
		},
		{
			name:     "negative coordinates",
			lat1:     -1.0,
			lon1:     -1.0,
			lat2:     2.0,
			lon2:     3.0,
			expected: 5.0,
			delta:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := euclideanDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			if math.Abs(got-tt.expected) > tt.delta {
				t.Errorf("euclideanDistance() = %v, want %v ± %v", got, tt.expected, tt.delta)
			}
		})
	}
}

func TestMatcher_findClosestPOI(t *testing.T) {
	tests := []struct {
		name     string
		eventLat float64
		eventLon float64
		pois     []geomatch.PointOfInterest
		want     int
	}{
		{
			name:     "single POI",
			eventLat: 48.8566,
			eventLon: 2.3522,
			pois: []geomatch.PointOfInterest{
				{Name: "Eiffel Tower", Lat: 48.8566, Lon: 2.3522},
			},
			want: 0,
		},
		{
			name:     "multiple POIs - closest first",
			eventLat: 48.8566,
			eventLon: 2.3522,
			pois: []geomatch.PointOfInterest{
				{Name: "Eiffel Tower", Lat: 48.8566, Lon: 2.3522},
				{Name: "Times Square", Lat: 40.7589, Lon: -73.9851},
			},
			want: 0,
		},
		{
			name:     "multiple POIs - closest second",
			eventLat: 40.7589,
			eventLon: -73.9851,
			pois: []geomatch.PointOfInterest{
				{Name: "Eiffel Tower", Lat: 48.8566, Lon: 2.3522},
				{Name: "Times Square", Lat: 40.7589, Lon: -73.9851},
			},
			want: 1,
		},
		{
			name:     "empty POIs",
			eventLat: 48.8566,
			eventLon: 2.3522,
			pois:     []geomatch.PointOfInterest{},
			want:     -1,
		},
		{
			name:     "very close POIs",
			eventLat: 48.8566,
			eventLon: 2.3522,
			pois: []geomatch.PointOfInterest{
				{Name: "Point A", Lat: 48.8567, Lon: 2.3523},
				{Name: "Point B", Lat: 48.8566, Lon: 2.3522},
				{Name: "Point C", Lat: 48.8565, Lon: 2.3521},
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findClosestPOI(tt.eventLat, tt.eventLon, tt.pois)
			if got != tt.want {
				t.Errorf("findClosestPOI() = %v, want %v", got, tt.want)
			}
		})
	}
}
