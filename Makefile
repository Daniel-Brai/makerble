DATABASE_URL := $(shell echo $${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/makerble_dev?sslmode=disable})

.PHONY: migrate-create
migrate-create:
	@echo "Creating migration files..."
	@migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-up
migrate-up:
	@echo "Running migrations up..."
	@migrate -path migrations -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	@echo "Running migrations down..."
	@migrate -path migrations -database "$(DATABASE_URL)" down

.PHONY: migrate-force
migrate-force:
	@echo "Forcing migration version..."
	@migrate -path migrations -database "$(DATABASE_URL)" force $(version)

.PHONY: migrate-version
migrate-version:
	@echo "Current migration version:"
	@migrate -path migrations -database "$(DATABASE_URL)" version

.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	@swag init --parseInternal -g cmd/api/main.go

.PHONY: run
run:
	@go run ./cmd/api/main.go

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out
