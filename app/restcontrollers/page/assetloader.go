package page

import (
	"fmt"
	"polygnosics/app/models"
	"polygnosics/app/services/db/mysqldb"

	"github.com/google/uuid"
)

func generatePath(assetID *uuid.UUID) string {
	path := "./"
	increment := 4
	for i := 0; i < 15; i = i + increment {
		path = fmt.Sprintf("%s/%s", path, string(assetID.String()[i:i+increment]))
	}
	return path
}

func LoadAssetReferences(assetID *uuid.UUID) (*models.Asset, error) {
	asset, err := mysqldb.GetAsset(assetID)
	if err != nil {
		return nil, err
	}

	asset.Path = generatePath(assetID)
	return asset, nil
}

func SaveAssetReferences(asset *models.Asset) error {
	if err := mysqldb.UpdateAsset(asset); err != nil {
		return err
	}
	return nil
}
