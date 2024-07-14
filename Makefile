build:
	@go build -o bin/app cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/app
