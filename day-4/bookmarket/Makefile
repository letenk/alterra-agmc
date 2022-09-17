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