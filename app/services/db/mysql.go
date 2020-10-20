package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // nolint:golint
	migrate "github.com/rubenv/sql-migrate"
)

var dbConnSystem = "root:password@tcp(172.18.0.1:3306)/core"

func BootstrapSystem() error {

	fmt.Printf("Executing MYSQL migration\n")
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations/mysql",
	}
	fmt.Printf("Getting migration files\n")

	db, err := sql.Open("mysql", dbConnSystem+"?parseTime=true")
	if err != nil {
		return err
	}
	fmt.Printf("DB connection open\n")

	n, err := migrate.Exec(db, "mysql", migrations, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}

func ConnectSystem() (*sql.DB, error) {
	fmt.Println("Connecting to MYSQL")

	db, err := sql.Open("mysql", dbConnSystem)

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	return db, nil
}
