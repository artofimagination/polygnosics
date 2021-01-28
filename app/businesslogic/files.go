package businesslogic

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Default file paths
const (
	DefaultUserAvatarPath    = "/assets/images/avatar.jpg"
	DefaultProductAvatarPath = "/assets/images/avatar.jpg"
	DefaultProjectAvatarPath = "/assets/images/avatar.jpg"
)

var splitRegexp = regexp.MustCompile(`(\S{4})`)

func GeneratePath(assetID *uuid.UUID) (string, error) {
	assetIDString := strings.Replace(assetID.String(), "-", "", -1)
	assetStringSplit := splitRegexp.FindAllString(assetIDString, -1)
	assetPath := path.Join(assetStringSplit...)
	rootPath := os.Getenv("USER_STORE_DOCKER")
	assetPath = path.Join(rootPath, assetPath)
	if err := os.MkdirAll(assetPath, os.ModePerm); err != nil {
		return "", err
	}
	return assetPath, nil
}

func (c *Context) UploadFile(asset *models.Asset, fileType string, defaultPath string, formName string, r *http.Request) error {
	file, handler, err := r.FormFile(formName)
	if err == http.ErrMissingFile {
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	defer file.Close()

	if err := c.UserDBController.ModelFunctions.SetFilePath(asset, fileType, filepath.Ext(handler.Filename)); err != nil {
		return err
	}
	path := c.UserDBController.ModelFunctions.GetFilePath(asset, fileType, defaultPath)

	// Create file
	dst, err := os.Create(path)
	if err != nil {
		if err2 := dst.Close(); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		if err2 := c.UserDBController.ModelFunctions.ClearAsset(asset, fileType); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return err
	}

	// Copy the uploaded file to the created file on the file system.
	if _, err := io.Copy(dst, file); err != nil {
		if err2 := dst.Close(); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		if err2 := c.UserDBController.ModelFunctions.ClearAsset(asset, fileType); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return err
	}

	return nil
}
