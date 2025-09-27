package datasource

type EventDataSource interface {
	Get() ([]Event, error)
}
