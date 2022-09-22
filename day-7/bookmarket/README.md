# How To Run
1. Create a new database, with any name

2. Setup `.env` file for connection to database
```go
DB_SOURCE=db_user:db_password@tcp(db_host:db_port)/db_name?charset=utf8mb4&parseTime=True&loc=Local
```
Example:
```go
DB_SOURCE=root:root@tcp(127.0.0.1:3306)/bookmarket?charset=utf8mb4&parseTime=True&loc=Local
```

3. Run app
```go
go run main.go
```

# How To Run Test
1. Create a new database for test, with any name

2. Setup `.env` file for connection to database
```go
DB_SOURCE_TEST=db_user:db_password@tcp(db_host:db_port)/db_name?charset=utf8mb4&parseTime=True&loc=Local
```
Example:
```go
DB_SOURCE_TEST=root:root@tcp(127.0.0.1:3306)/bookmarket_test?charset=utf8mb4&parseTime=True&loc=Local
```

3. Run Test

- Run all tests
```go
go test -v ./...
```

- Run all tests with coverage
```go
go test ./... -v -coverpkg=./...
```

- Run all test with coverage and print on CLI
```go
go test ./... -v -coverpkg=./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

- Run all test with coverage and print on browser
```go
go test ./... -v -coverpkg=./...
go tool cover -html=coverage.out
```

# Documentation
[Postman Documentation](https://documenter.getpostman.com/view/12132212/2s7ZDyu7Ki)

# Todo
- [x] Hexagonal Architecture