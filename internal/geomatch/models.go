package geomatch

type PointOfInterest struct {
	Name string
	Lat  float64
	Lon  float64
}

type MatchResult struct {
	PointOfInterest
	Impressions int
	Clicks      int
}
