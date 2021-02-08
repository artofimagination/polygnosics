package businesslogic

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"

	"github.com/artofimagination/golang-docker/docker"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	ProductAvatarKey   = "avatar"
	ProductMainAppKey  = "main_app"
	ProductClientApp   = "client-app"
	ProductDescription = "description"
	ProductName        = "name"
	ProductRequires3D  = "requires_3d"
	ProductURL         = "url"
	ProductPublic      = "is_public"
)

const (
	CheckBoxUnChecked = "unchecked"
	CheckBoxChecked   = "checked"
)

func (c *Context) DeleteProduct(product *models.ProductData) error {
	projects, err := c.UserDBController.GetProjectsByProductID(&product.ID)
	if err != nil {
		return err
	}

	for _, project := range projects {
		if err := c.UserDBController.DeleteProject(&project.ID); err != nil {
			return err
		}
	}

	if err := c.UserDBController.DeleteProduct(&product.ID); err != nil {
		return err
	}

	folder := c.UserDBController.ModelFunctions.GetFilePath(product.Assets, ProductMainAppKey, "")
	dir, _ := filepath.Split(folder)
	if err := removeContents(dir); err != nil {
		return fmt.Errorf("Failed to delete product. %s", errors.WithStack(err))
	}
	return nil
}

func (c *Context) AddProduct(userID *uuid.UUID, productName string, r *http.Request) (*models.ProductData, error) {
	product, err := c.UserDBController.CreateProduct(
		productName,
		userID,
		GeneratePath)
	if err != nil {
		return nil, err
	}

	if err := c.UploadFiles(product.Assets, r); err != nil {
		if errDelete := c.DeleteProduct(product); errDelete != nil {
			err = errors.Wrap(errors.WithStack(err), errDelete.Error())
			return nil, fmt.Errorf("Failed to delete product. %s", errors.WithStack(err))
		}
		return nil, fmt.Errorf("Failed to upload asset. %s", errors.WithStack(err))
	}

	return product, nil
}

func (c *Context) UploadFiles(assets *models.Asset, r *http.Request) error {
	if err := c.UploadFile(assets, ProductAvatarKey, DefaultProductAvatarPath, r); err != nil {
		return err
	}

	if err := c.UploadFile(assets, ProductMainAppKey, "", r); err != nil {
		return err
	}

	if err := c.UploadFile(assets, ProductClientApp, "", r); err != nil {
		return err
	}

	return nil
}

func (c *Context) CreateDockerImage(product *models.ProductData, userID *uuid.UUID) error {
	pathString := c.UserDBController.ModelFunctions.GetFilePath(product.Assets, ProductMainAppKey, "")
	if err := untar(pathString); err != nil {
		if errDelete := c.DeleteProduct(product); errDelete != nil {
			err = errors.Wrap(errors.WithStack(err), errDelete.Error())
			return fmt.Errorf("Failed to delete product. %s", errors.WithStack(err))
		}
		return fmt.Errorf("Failed to decompress main app. %s", errors.WithStack(err))
	}

	imageName := fmt.Sprintf("%s/%s", userID.String(), product.ID.String())
	sourcePath := path.Join(filepath.Dir(pathString), "build")
	if err := docker.CreateImage(sourcePath, imageName); err != nil {
		if errDelete := c.DeleteProduct(product); errDelete != nil {
			err = errors.Wrap(errors.WithStack(err), errDelete.Error())
			return fmt.Errorf("Failed to delete product. %s", errors.WithStack(err))
		}
		return fmt.Errorf("Failed to create product image. %s", errors.WithStack(err))
	}
	return nil
}

// getBooleanString returns a check box stat Yes/No string
func getBooleanString(input string) string {
	if input == "" || input == CheckBoxUnChecked {
		return CheckBoxUnChecked
	}
	return CheckBoxChecked
}

// SetProductDetails sets the key-value content of product details based on form values.
func (c *Context) SetProductDetails(details *models.Asset, r *http.Request) {
	c.UserDBController.ModelFunctions.SetField(details, ProductName, r.FormValue(ProductName))
	c.UserDBController.ModelFunctions.SetField(details, ProductDescription, r.FormValue(ProductDescription))
	c.UserDBController.ModelFunctions.SetField(details, ProductRequires3D, getBooleanString(r.FormValue(ProductRequires3D)))
	c.UserDBController.ModelFunctions.SetField(details, ProductPublic, getBooleanString(r.FormValue(ProductPublic)))
	c.UserDBController.ModelFunctions.SetField(details, ProductURL, r.FormValue(ProductURL))
}

func (c *Context) UpdateProductData(product *models.ProductData, r *http.Request) error {
	c.SetProductDetails(product.Details, r)

	if err := c.UserDBController.UpdateProductDetails(product); err != nil && err != dbcontrollers.ErrMissingProductDetail {
		return err
	}

	if err := c.UserDBController.UpdateProductAssets(product); err != nil && err != dbcontrollers.ErrMissingProductAsset {
		return err
	}
	return nil
}
