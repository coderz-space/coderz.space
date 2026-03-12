.PHONY: build run clear-logs swagger

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