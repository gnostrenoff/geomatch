# Geomatch
## Run

Run the API server:
- requires Go 1.25 installed.
```bash
make run-api-local
```

Or with Docker:
```bash
make run-api-docker
```

### Send a request
Send a request to the server:
```bash
make send-local-request
```

### Run unit tests
```bash
make run-tests
```

## Code Architecture

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout) with `cmd/` for main applications, `internal/` for private code, and `configs/` for configuration.

### Core Abstractions

**GeoMatcher Interface**: The geomatching logic is abstracted through a `GeoMatcher` interface. Two implementations are currently available:
- **Haversine** (used by default): Accurate distance calculation using spherical geometry
- **Euclidean**: Fast approximation using simple coordinate differences
Any other geomatching algorithm can be plugged in.

**DataSource Interface**: Data retrieval is abstracted through an `EventDataSource` interface. Currently implemented with a CSV loader, but any data source can be plugged in.

### Trade-offs
- **Haversine vs Euclidean**: Haversine is more accurate for real-world distances but computationally expensive. Euclidean is faster but less accurate, especially over large distances or near the poles.
Current default is Haversine for accuracy, but Euclidean can be used for performance in low-latency scenarios.

### Future Enhancements

- Support more geomatching algorithms
- Support multi-core processing
- Add caching for repeated requests
- Observability: metrics, tracing, logging
- Health checks
- Graceful shutdown
- Profiling endpoints
- API Client generation/implementation
