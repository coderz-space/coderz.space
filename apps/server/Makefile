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