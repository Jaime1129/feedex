.PHONY: setup test run swagger build mock

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

build: setup
	@echo "Building"
	go build -o feedex cmd/main.go

mock:
	@echo "Generating mock fiels"
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen@latest
	mockgen -source=internal/components/bn_price_cli.go -destination=mock/components/bn_price_cli.go
	mockgen -source=internal/components/eth_scan_cli.go -destination=mock/components/eth_scan_cli.go
	mockgen -source=internal/repository/trx_fee_repo.go -destination=mock/repository/trx_fee_repo.go