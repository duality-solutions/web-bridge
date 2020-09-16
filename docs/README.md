# WebBridge Swagger Documentation

## Updating Swagger Documentation

```go
go get -u -v github.com/swaggo/swag/cmd/swag
swag init -g api/rest/router.go
```

New Swagger documentation files are saved in the docs directory

## Running Swagger Locally

```http
http://localhost:8080/swagger/index.html
```
