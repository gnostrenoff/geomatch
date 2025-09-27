package geomatch

type GeoMatcher interface {
	Match(pois []PointOfInterest) ([]MatchResult, error)
}
