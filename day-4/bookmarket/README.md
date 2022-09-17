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
[Postman Documentation](https://documenter.getpostman.com/view/12132212/2s7YmqEmWn)

# Todo
## Integration Testing Login
- [x]  Login success (200).
- [x] Login failed, validation error (400).
- [x] Login failed, user or password incorrect (400).

## Integration Testing Users
### Create User
- [x] Create user success (201)
- [x] Create user failed, email already exists (400)
- [x] Create user failed, validation error (400).

### Get all users
- [x] Get all user success (200).
- [x] Get all user with data empty slice, because data is not available (200).
- [x] Get all user failed, unauthorized (400 and 401)

### Get user find by id
- [x] Get user find by id success (200).
- [x] Get user find by id failed, user not found (400).
- [x] Get user find by id failed, unauthorized (400 and 401).

### Update user
- [x] Update user success (200).
- [x] Update user failed, validation error (400).
- [x] Update user failed, not access. Because only can updating the data self him (400).
- [x] Update user failed, unauthorized (400 and 401).

### Delete User
- [x] Delete user success (200).
- [x] Delete user failed, not access. Because only can deleting the data self him (400) .
- [x] Delete user failed, unauthorized (400 and 401).

## Integration Testing Books
### Create Book
- [x] Create book success (200).
- [x] Create book failed, validation error (400).
- [x] Create book failed, unauthorized (400 and 401).

### Get all Books
- [x] Get all books success (200).

### Get Book by id
- [x] Get book by id success (200).
- [x] Get book by id failed, book not found (400).

### Update Book
- [x] Update book success (200).
- [x] Update book failed, book not found (400).
- [x] Update book failed, validation error (400).
- [x] Update book failed, unauthorized (400 and 401). 

### Delete Book
- [x] Delete Book success (200).
- [x] Delete book failed, book not found (400). 
- [x] Delete book failed, unauthorized (400 and 401)

