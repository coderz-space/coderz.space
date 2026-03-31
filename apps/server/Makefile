# Load .env file :
-include .env
export


.PHONY: build run clear-logs swagger sqlc

MAIN="cmd/main.go"

# build the server
build:
	@echo "Building server..."
	@go build -o bin/main ${MAIN}

# run the development server
run:
	@echo "Running server..."
	@go run ${MAIN}

# clear logs files
clear-logs:
	@echo "Clearing logs..."
	@rm -rf logs


# generates swagger docs :
swagger:
	@echo "Generating swagger docs..."
	swag init -o ./swagger --parseDependency --parseInternal -g cmd/main.go

# clear swagger docs  : 
clear-swagger:
	@echo "Clearing swagger docs..."
	@rm -rf swagger


# Database : 
sqlc:
	@echo "Generating sqlc queries..."
	sqlc generate -f db/sqlc.yaml


# Migrations : 
migrate-up:
	@echo "Applying migrations (up)"
	docker compose run --rm migrate -path /migrations -database "${DB_DSN}" up

migrate-down:
	@echo "Rolling back latest migration"
	docker compose run --rm migrate -path /migrations -database "${DB_DSN}" down 1



# Docker : 
docker-up:
	@echo "Starting docker containers..."
	@docker compose up -d


docker-down:
	@echo "Stopping docker containers..."
	@docker compose down