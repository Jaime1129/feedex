.PHONY: setup test

setup:
	@echo "Setting up the project..."
	go get -u ./...

test: setup
	@echo "Running tests..."
	go test -v ./...

run: setup
	@echo "Starting server..."
	go run cmd/main.go

swagger: setup
	@echo "Generating swagger"
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/main.go