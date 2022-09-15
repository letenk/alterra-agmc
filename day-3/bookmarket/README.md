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

4. Create a user use endpoint `create user`.
5. Login use endpoint `login` with the user created above, for get a token and use token each `Authorization` endpoint is require as a `Bearer Token`.

# Documentation
[Postman Documentation](https://documenter.getpostman.com/view/12132212/2s7YfLfaXg)

# Todo
## Middleware Logs
- [x] Implement log middleware

## Middleware JWT
- [x] Implement JWT (Create a endpoint login for getting the token)

### Books
- [x] Endpoint Get all books (/books) and Get book by id (/books/:id) Not Authenticated

### Users
- [x] Endpoint Create new user (/users) Not Authenticated

# Middleware Validator
- [x] Implement validator (Ex: empty string(not required), not email format)

# Additional Task
- [x] Make as user only able to do PUT and Delete on him self.
- [x] Do not create a user with same email on the database.
