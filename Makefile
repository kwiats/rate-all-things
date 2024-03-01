build:
	@go build -o bin/rate-all-things cmd/rate-all-things/main.go

run: build
	@./bin/rate-all-things

test: 
	@go test -v ./...