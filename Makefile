run-api-local:
	ENV=local go run cmd/api.go

run-api-docker:
	docker build -t geomatch-api .
	docker run --rm -e ENV=local --name geomatch-api-container -p 8080:8080 geomatch-api

send-local-request:
	curl --location 'localhost:8080/events' \
		--header 'Content-Type: application/json' \
		--data '{"points_of_interest":[{"lat":48.86,"lon":2.35,"name":"Chatelet"},{"lat":48.8759992,"lon":2.3481253,"name":"Arc de triomphe"}]}' \
		| jq

run-tests:
	go test ./...

format:
	gofmt -s -w .

run-benchmarks:
	go test -bench=. -benchmem -benchtime=5s ./internal/geomatch/benchmarks