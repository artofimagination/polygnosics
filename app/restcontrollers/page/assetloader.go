package page

import (
	"fmt"
	"os"
	"path"
	"polygnosics/app/models"
	"polygnosics/app/services/db/mysqldb"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var splitRegexp = regexp.MustCompile(`(\S{4})`)

func generatePath(assetID *uuid.UUID) string {
	assetIDString := strings.Replace(assetID.String(), "-", "", -1)
	assetStringSplit := splitRegexp.FindAllString(assetIDString, -1)
	path := path.Join(assetStringSplit...)
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
