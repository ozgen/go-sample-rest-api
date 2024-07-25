build:
	@go build -o bin/app cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/app

debug:
	@dlv debug --headless --listen=:2345 --log --api-version=2 ./cmd/main.go


migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

swagger:
	swag init -d ./,./service/user,./service/camerametadata --generalInfo service/user/routes.go --output docs/
