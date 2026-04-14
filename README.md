```
go mod init github.com/YOUR_USER_NAME/go-prject
```

## About test 
https://www.youtube.com/watch?v=8hQG7QlcLBk

migration:
```
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
```