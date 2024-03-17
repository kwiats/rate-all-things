build:
	@go build -o bin/tit cmd/tit/main.go

run: build
	@./bin/tit

test: 
	@go test -v ./...