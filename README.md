# account-test

## Installation

Ensure Golang is installed: https://go.dev/doc/install

```bash
go mod tidy
go install github.com/golang/mock/mockgen@v1.6.0
```

## Setup

Fill in the values in the dev.env with the values for your postgres server

Example:
```cgo
DB_HOST: localhost
DB_PORT: 5432
DB_USERNAME: postgres
DB_PASSWORD: 
DB_NAME: postgres
DB_SCHEMA: public
```
## Usage

```cgo
go run server.go dev #start the app

mockgen -source=./internal/core/ports/ports.go -destination=./internal/mocks/ports/ports.go #generates mock implementation for unit test

go test ./internal/... -count=1 #run test cases
```