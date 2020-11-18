package configloader

import (
	"errors"
	"os"
)

type DBConfig struct {
	Address      string
	Port         string
	Username     string
	Password     string
	MigrationDir string
}

// LoadDBConfigFromEnv loads the database connection details from environment variables.
// Currently supported DB configs: Mysql, Postgres Timescale, MongoDB
func LoadDBConfigFromEnv(db string) (*DBConfig, error) {
	configuration := DBConfig{}

	switch {
	case db == "MYSQL":
		configuration.Address = os.Getenv("MYSQL_DB_ADDRESS")
		configuration.Port = os.Getenv("MYSQL_DB_PORT")
		configuration.Username = os.Getenv("MYSQL_DB_USER")
		configuration.Password = os.Getenv("MYSQL_DB_PASSWORD")
		configuration.MigrationDir = os.Getenv("MYSQL_DB_MIGRATION_DIR")
	case db == "Timescale":
		configuration.Address = os.Getenv("TIMESCALE_DB_ADDRESS")
		configuration.Port = os.Getenv("TIMESCALE_DB_PORT")
		configuration.Username = os.Getenv("TIMESCALE_DB_USER")
		configuration.Password = os.Getenv("TIMESCALE_DB_PASSWORD")
		configuration.MigrationDir = os.Getenv("TIMESCALE_DB_MIGRATION_DIR")
	case db == "Mongo":
		configuration.Address = os.Getenv("MONGO_DB_ADDRESS")
		configuration.Port = os.Getenv("MONGO_DB_PORT")
		configuration.Username = os.Getenv("MONGO_DB_USER")
		configuration.Password = os.Getenv("MONGO_DB_PASSWORD")
	default:
		return nil, errors.New("Unknown DB type")
	}

	if configuration.Address == "" {
		return nil, errors.New("Address is empty")
	}

	if configuration.Port == "" {
		return nil, errors.New("Port is empty")
	}

	if configuration.Username == "" {
		return nil, errors.New("Username is empty")
	}

	if configuration.Password == "" {
		return nil, errors.New("Password is empty")
	}

	return &configuration, nil
}
