package businesslogic

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"

	"github.com/artofimagination/golang-docker/docker"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	ProductAvatar      = "product_avatar"
	ProductMainApp     = "main_app"
	ProductClientApp   = "client-app"
	ProductDescription = "product_description"
	ProductName        = "product_name"
	ProductRequires3D  = "requires_3d"
	ProductURL         = "product_url"
	ProductPublic      = "is_public"
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

	folder := c.UserDBController.ModelFunctions.GetFilePath(product.Assets, ProductMainApp, "")
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

	err = c.UploadFile(product.Assets, ProductAvatar, DefaultProductAvatarPath, r)
	if err != nil {
		if errDelete := c.DeleteProduct(product); errDelete != nil {
			err = errors.Wrap(errors.WithStack(err), errDelete.Error())
			return nil, fmt.Errorf("Failed to delete product. %s", errors.WithStack(err))
		}
		return nil, fmt.Errorf("Failed to upload avatar. %s", errors.WithStack(err))
	}

	err = c.UploadFile(product.Assets, ProductMainApp, "", r)
	if err != nil {
		if errDelete := c.DeleteProduct(product); errDelete != nil {
			err = errors.Wrap(errors.WithStack(err), errDelete.Error())
			return nil, fmt.Errorf("Failed to delete product. %s", errors.WithStack(err))
		}
		return nil, fmt.Errorf("Failed to upload main app. %s", errors.WithStack(err))
	}

	if err := c.UploadFile(product.Assets, ProductClientApp, "", r); err != nil {
		if errDelete := c.DeleteProduct(product); errDelete != nil {
			err = errors.Wrap(errors.WithStack(err), errDelete.Error())
			return nil, fmt.Errorf("Failed to delete product. %s", errors.WithStack(err))
		}
		return nil, fmt.Errorf("Failed to upload client app. %s", errors.WithStack(err))
	}

	return product, nil
}

func (c *Context) CreateDockerImage(product *models.ProductData, userID *uuid.UUID) error {
	pathString := c.UserDBController.ModelFunctions.GetFilePath(product.Assets, ProductMainApp, "")
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
	if input == "" {
		return "No"
	}
	return input
}

// SetProductDetails sets the key-value content of product details based on form values.
func (c *Context) SetProductDetails(details *models.Asset, r *http.Request) {
	c.UserDBController.ModelFunctions.SetField(details, ProductName, r.FormValue("productName"))
	c.UserDBController.ModelFunctions.SetField(details, ProductDescription, r.FormValue("productDescription"))
	c.UserDBController.ModelFunctions.SetField(details, ProductRequires3D, r.FormValue("requires3D"))
	c.UserDBController.ModelFunctions.SetField(details, ProductPublic, getBooleanString(r.FormValue("publicProduct")))
	c.UserDBController.ModelFunctions.SetField(details, ProductURL, r.FormValue("productUrl"))
}

func (c *Context) UpdateProductData(product *models.ProductData, r *http.Request) error {
	c.SetProductDetails(product.Details, r)

	if err := c.UserDBController.UpdateProductDetails(product); err != nil {
		if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
			err = errors.Wrap(errors.WithStack(err), errDelete.Error())
			return fmt.Errorf("Failed to delete product. %s", errors.WithStack(err))
		}
		return fmt.Errorf("Failed to update product details. %s", errors.WithStack(err))
	}

	if err := c.UserDBController.UpdateProductAssets(product); err != nil {
		if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
			err = errors.Wrap(errors.WithStack(err), errDelete.Error())
			return fmt.Errorf("Failed to delete product. %s", errors.WithStack(err))
		}
		return fmt.Errorf("Failed to update product assets. %s", errors.WithStack(err))
	}
	return nil
}
