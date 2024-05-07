package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// This is added here for ease of setup. In a normal dev setup, table creation should be done outside of code either through a pipeline or DB team
var schema = `
	CREATE TABLE IF NOT EXISTS %s.account(
		id VARCHAR PRIMARY KEY NOT NULL,
		balance float NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
`

const DriverName = "postgres"

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	Schema   string
}

// Init function takes in DB config object and returns a wrapper to sql/DB object.
func Init(dbConfig *DBConfig) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Name)
	client, err := sqlx.Open(DriverName, dataSource)
	if err != nil {
		log.Println("err - ", err.Error())
		return nil, err
	}

	// verifies connection is db is working
	if err := client.Ping(); err != nil {
		log.Println("err 2 - ", err.Error())
		return nil, err
	}

	schema = fmt.Sprintf(schema, dbConfig.Schema)
	client.MustExec(schema)

	return client, nil
}
