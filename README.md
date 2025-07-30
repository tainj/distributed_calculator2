# Миграции

```bash
migrate -path ./migrations -database "postgres://calculator:verydifficlutpassword@localhost:5432/calculator_db?sslmode=disable" up
```

# Кодогенерация

```bash
protoc -I ./proto --go_out ./pkg/api/ --go_opt paths=source_relative --go-grpc_out ./pkg/api/ --go-grpc_opt paths=source_relative --grpc-gateway_out ./pkg/api/ --grpc-gateway_opt paths=source_relative ./proto/calculator.proto
```

# Подключение к бд 

```bash
psql -h localhost -p 5432 -U calculator -d calculator_db
```