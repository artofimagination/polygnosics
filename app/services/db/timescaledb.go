package db

import (
	"database/sql"
	"fmt"

	// Linter does not like this import. Forcing to ignore it.
	"github.com/rubenv/sql-migrate" // nolint: goimports
)

var host = "172.18.0.1"
var port = "5432"
var user = "root"
var password = "password"
var dbName = "data"

func BootstrapData() error {
	fmt.Printf("Executing TimeScaleDB migration\n")

	dbConnData := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations/timescaledb",
	}
	fmt.Printf("Getting migration files\n")

	db, err := sql.Open("postgres", dbConnData)
	if err != nil {
		return err
	}
	fmt.Printf("DB connection open\n")

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}

func ConnectData() (*sql.DB, error) {
	fmt.Println("Connecting to TimescaleDB")

	dbConnData := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err := sql.Open("postgres", dbConnData)

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	return db, nil
}
