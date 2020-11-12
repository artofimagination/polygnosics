package mysqldb

import (
	"encoding/json"
	"fmt"
	"polygnosics/app/models"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func AddAsset() (*uuid.UUID, error) {
	queryString := "INSERT INTO user_assets (id, refs) VALUES (UUID_TO_BIN(?), ?)"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	newID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	binary, err := json.Marshal(models.References{})
	if err != nil {
		return nil, err
	}
	query, err := db.Query(queryString, newID, binary)
	if err != nil {
		return nil, err
	}

	defer query.Close()
	return &newID, nil
}

func UpdateAsset(asset *models.Asset) error {
	queryString := "UPDATE user_assets set refs = ? where id = UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}

	defer db.Close()

	query, err := db.Query(queryString, asset.References, asset.ID)
	if err != nil {
		return err
	}

	defer query.Close()
	return nil
}

func GetAsset(assetID *uuid.UUID) (*models.Asset, error) {
	asset := models.Asset{}
	queryString := "SELECT refs FROM user_assets WHERE id = UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := db.QueryRow(queryString, *assetID)
	if err != nil {
		return nil, err
	}

	refs := json.RawMessage{}
	if err := query.Scan(&refs); err != nil {
		return nil, errors.Wrap(errors.WithStack(err), fmt.Sprintf("Asset %s not found", assetID.String()))
	}

	if err := json.Unmarshal(refs, &asset.References); err != nil {
		return nil, err
	}

	return &asset, nil
}

func DeleteAsset(assetID *uuid.UUID) error {
	query := "DELETE FROM user_assets WHERE id=UUID_TO_BIN(?)"
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
