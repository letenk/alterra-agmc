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

# Documentation
[Postman Documentation](https://documenter.getpostman.com/view/12132212/2s7YYu7icq)

# Todo
## Books
- [x] Get all books
- [x] Get book by id
- [x] Create new book
- [x] Update book by id
- [x] Delete book by id

## Users
- [x] Get all users
- [x] Get user by id
- [x] Create new user
- [x] Update user by id
- [x] Delete user by id
