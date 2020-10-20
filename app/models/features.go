package models

import (
	"encoding/json"

	"polygnosics/app/services/db/mysqldb"
)

// Feature describes the available simulation features.
type Feature struct {
	ID     int             `json:"id"`
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

// GetAllFeatures returns all fetures stored in the database.
func GetAllFeatures() ([]Feature, error) {
	queryString := "select id, name, config from features"
	db, err := mysqldb.ConnectSystem()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query, err := db.Query(queryString)

	if err != nil {
		return nil, err
	}

	defer query.Close()

	features := []Feature{}
	for query.Next() {
		var feature Feature
		err := query.Scan(&feature.ID, &feature.Name, &feature.Config)
		if err != nil {
			return nil, err
		}
		features = append(features, feature)
	}

	return features, nil
}
