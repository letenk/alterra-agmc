## up_build: stops docker-compose (if running), builds all projects and start docker compose file docker-compose.yml 
up_build:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images"
	docker-compose up -d
	@echo "Docker iamges build and started!"

## down: stop docker compose file docker-compose.yml
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## up_build_dev: stops docker-compose (if running), builds all projects and start docker compose file docker-compose.dev.yml for development 
up_build_dev:
	@echo "Stopping docker images (if running...)"
	docker-compose -f docker-compose.dev.yml down
	@echo "Building (when required) and starting docker images for development"
	docker-compose -f docker-compose.dev.yml up -d
	@echo "Docker iamges build and started!"

## down_dev: stop docker compose file docker-compose.dev.yml for development  
down_dev:
	@echo "Stopping docker compose..."
	docker-compose -f docker-compose.dev.yml down
	@echo "Done!"

## test: Run all test in this app
test:
	@echo "All tests are running..."
	go test -v ./...
	@echo "Test finished"

## test_cover: Run all test with coverage
test_cover:
	@echo "All test are running with coverage..."
	go test ./... -v -coverpkg=./...
	@echo "Test finished"


## test_cover_print: Run all test with coverage and print on CLI
test_cover_print:
	@echo "All test are running with coverage..."
	go test ./... -v -coverpkg=./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	@echo "Test finished"

## test_cover_print_html: Run all test with coverage and open on browser html
test_cover_print_html:
	@echo "All test are running with coverage..."
	go test ./... -v -coverpkg=./...
	go tool cover -html=coverage.out
	@echo "Test finished"