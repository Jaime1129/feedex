.PHONY: setup test

setup:
	@echo "Setting up the project..."
	go get -u ./...

test: setup
	@echo "Running tests..."
	go test -v ./...
