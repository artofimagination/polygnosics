package page

import (
	"fmt"
	"os"
	"polygnosics/app/models"
	"polygnosics/app/services/db/mysqldb"
	"strings"

	"github.com/google/uuid"
)

func generatePath(assetID *uuid.UUID) string {
	path := ""
	increment := 4
	assetIDString := assetID.String()
	assetIDString = strings.Replace(assetIDString, "-", "", -1)
	for i := 0; i < 31; i = i + increment {
		path = fmt.Sprintf("%s/%s", path, assetIDString[i:i+increment])
	}
	return path
}

func LoadAssetReferences(assetID *uuid.UUID) (*models.Asset, error) {
	asset, err := mysqldb.GetAsset(assetID)
	if err != nil {
		return nil, err
	}

	absolutePath := fmt.Sprintf("/user-assets%s", generatePath(assetID))
	asset.Path = absolutePath
	if err := os.MkdirAll(asset.Path, os.ModePerm); err != nil {
		return nil, err
	}
	return asset, nil
}

func SaveAssetReferences(asset *models.Asset) error {
	if err := mysqldb.UpdateAsset(asset); err != nil {
		return err
	}
	return nil
}
