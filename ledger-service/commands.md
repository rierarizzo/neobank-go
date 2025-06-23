# Commands

- Generate proto files:

```shell
protoc --proto_path=proto/ledger --go_out=proto/ledger --go_opt=paths=source_relative --go-grpc_out=proto/ledger --go-grpc_opt=paths=source_relative proto/ledger/ledger.proto
```

- Migrate database files:

```shell
migrate -source file://./migrations -database postgres://ledger_user:ledger_pw@localhost:5432/ledger?sslmode=disable up
```