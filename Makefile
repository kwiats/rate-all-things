postgres:
	docker-compose up --build -d

linters:
	golangci-lint run -v;gofmt -w .
	
build:
	@go build -o bin/tit cmd/tit/main.go

clear:
	clear
	
run: build linters clear
	@./bin/tit

test: 
	@go test -v ./...