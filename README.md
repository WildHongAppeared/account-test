# account-test

## Installation

Ensure Golang is installed: https://go.dev/doc/install

```bash
go mod tidy
go install github.com/golang/mock/mockgen@v1.6.0 ## mockgen for mock test generation. No need to install if not generating mock
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

## Assumption
1. The precision of calculation for transaction is set to 5 floating point as seen in the question sheet to prevent precision error
2. Line 11-29 and 57-58 in postgres/db.go file is added for ease of setting up database tables. For a actual code repository in a professional setting, it is assumed that the database tables setup will be handled either through separate automation scripts or database teams
3. Account IDs are currently upper bound to 32 characters only and currently allows freetext. 
4. Balance and amount values are returned as string type as seen in the question sheet but calculation are performed in float64 after parsing
