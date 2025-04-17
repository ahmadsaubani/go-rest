# _How To Install_

### Installation
```
Requirements: 
- go > 1.20.x
- postgre
```

#### Development
```sh
How to run :
- run gowatch or go main.go
```

```sh
To up migration : 
migrate -database "postgres://user:password@localhost:5432/go-rest?sslmode=disable" -path migrations/ up

To down migration : 
migrate -database "postgres://user:password@localhost:5432/go-rest?sslmode=disable" -path migrations/ down 
```

