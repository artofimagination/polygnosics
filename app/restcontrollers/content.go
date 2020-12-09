package restcontrollers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"polygnosics/app"
	"polygnosics/app/restcontrollers/auth"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/pkg/errors"

	"github.com/google/uuid"
)

const (
	DefaultUserAvatarPath    = "/assets/images/avatar.jpg"
	DefaultProductAvatarPath = "/assets/images/avatar.jpg"
)

const (
	UserAvatar     = "user_avatar"
	UserBackground = "user_background"

	ProductAvatar      = "product_avatar"
	ProductMainApp     = "main_app"
	ProductClientApp   = "client-app"
	ProductDescription = "product_description"
	ProductName        = "product_name"
	ProductRequires3D  = "requires_3d"
	ProductURL         = "product_url"
	ProductPublic      = "is_public"
)

const (
	Product = iota
	User
)

var ProductData *models.ProductData

func isRequired(formName string) bool {
	switch formName {
	case ProductMainApp,
		ProductName:
		return true
	default:
		return false
	}
}

func getUserContent() map[string]interface{} {
	p := make(map[string]interface{})
	p["assets"] = make(map[string]interface{})
	path := auth.UserData.Assets.GetImagePath(UserAvatar, DefaultUserAvatarPath)
	p["assets"].(map[string]interface{})[UserAvatar] = path
	p["texts"] = make(map[string]interface{})
	p["texts"].(map[string]interface{})["avatar-upload"] = "Upload your avatar"
	p["texts"].(map[string]interface{})["username"] = auth.UserData.Name

	return p
}

func getUserProductContent(userID *uuid.UUID) (map[string]interface{}, error) {

	products, err := app.ContextData.UserDBController.GetProductsByUserID(userID)
	if err != nil {
		return nil, err
	}

	p := make(map[string]interface{})

	productContent := make([]map[string]interface{}, len(products))
	for i, product := range products {
		content := make(map[string]interface{})
		content[ProductAvatar] = product.ProductData.Assets.GetImagePath(ProductAvatar, DefaultProductAvatarPath)
		content[ProductName] = product.ProductData.Details.GetURL(ProductName, "")
		productContent[i] = content
	}
	p["product"] = productContent

	return p, nil
}

func uploadUserFile(fileType string, defaultPath string, formName string, r *http.Request) (string, error) {
	if auth.UserData == nil {
		return "", errors.New("User is not configured")
	}

	if err := auth.UserData.Assets.SetImagePath(fileType); err != nil {
		return "", err
	}
	path := auth.UserData.Assets.GetImagePath(fileType, defaultPath)

	if err := createFile(path, formName, r); err != nil {
		if err2 := auth.UserData.Assets.ClearAsset(fileType); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return "", err
	}
	return path, nil
}

func uploadProductFile(fileType string, defaultPath string, formName string, r *http.Request) error {
	file, _, _ := r.FormFile(formName)
	if isRequired(formName) && file == nil {
		return fmt.Errorf("Missing form value for %s", formName)
	} else if !isRequired(formName) && file == nil {
		return nil
	}

	if ProductData == nil {
		return errors.New("Product is not configured")
	}

	if err := ProductData.Assets.SetImagePath(fileType); err != nil {
		return err
	}
	path := ProductData.Assets.GetImagePath(fileType, defaultPath)

	if err := createFile(path, formName, r); err != nil {
		if err2 := ProductData.Assets.ClearAsset(fileType); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return err
	}
	return nil
}

func createFile(destination string, formName string, r *http.Request) error {
	file, handler, err := r.FormFile(formName)
	if err != nil {
		return err
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create(destination)
	if err != nil {
		if err2 := dst.Close(); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return err
	}

	// Copy the uploaded file to the created file on the file system.
	if _, err := io.Copy(dst, file); err != nil {
		if err2 := dst.Close(); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return err
	}

	return nil
}
