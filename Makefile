.PHONY: help run build tidy fmt vet migrate migrate-status clean

APP := production-tracer
BIN := bin

help:
	@echo "Targets:"
	@echo "  run             run the API (go run ./cmd/api)"
	@echo "  build           build api + migrate into ./bin"
	@echo "  tidy            go mod tidy"
	@echo "  fmt / vet       format / vet the code"
	@echo "  migrate         apply the schema (GORM AutoMigrate)"
	@echo "  migrate-status  list models registered for migration"

run:
	go run ./cmd/api

build:
	@mkdir -p $(BIN)
	go build -o $(BIN)/api ./cmd/api
	go build -o $(BIN)/migrate ./cmd/migrate

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

migrate:
	go run ./cmd/migrate up

migrate-status:
	go run ./cmd/migrate status

clean:
	rm -rf $(BIN)
