package mysqldb

import (
	"polygnosics/app/models"

	"github.com/google/uuid"
)

func AddAsset() (*uuid.UUID, error) {
	queryString := "INSERT INTO assets (id, refs) VALUES (UUID_TO_BIN(?), ?)"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	newID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	assets := models.Asset{}
	query, err := db.Query(queryString, newID, assets.References)
	if err != nil {
		return nil, err
	}

	defer query.Close()
	return &newID, nil
}

func GetAsset(assetID *uuid.UUID) (*models.Asset, error) {
	asset := models.Asset{}
	queryString := "SELECT refs FROM assets WHERE id = UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := db.QueryRow(queryString, *assetID)

	if err != nil {
		return nil, err
	}

	if err := query.Scan(&asset.References); err != nil {
		return nil, err
	}

	return &asset, nil
}

func DeleteAsset(assetID *uuid.UUID) error {
	query := "DELETE FROM assets WHERE id=UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query, *assetID)
	if err != nil {
		return err
	}
	return nil
}
