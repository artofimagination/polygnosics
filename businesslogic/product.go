package businesslogic

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/artofimagination/polygnosics/rest"

	"github.com/artofimagination/golang-docker/docker"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	ProductAvatarKey           = "avatar"
	ProductMainAppKey          = "main_app"
	ProductClientApp           = "client-app"
	ProductDescriptionKey      = "description"
	ProductShortDescriptionKey = "short_description"
	ProductNameKey             = "name"
	ProductRequires3DKey       = "requires_3d"
	ProductURLKey              = "url"
	ProductPublicKey           = "is_public"
	ProductPricingKey          = "pricing"
	ProductPriceKey            = "amount"
	ProductTagsKey             = "tags"
	ProductCategoriesKey       = "categories"
)

const (
	CreditCardNumberKey = "card_number"
	CreditCardExpiryKey = "expiry"
	CreditCardNameKey   = "name_on_card"
	CreditCardCVCKey    = "cvc"
)

const (
	PaymentTypeSingle = "Single Price"
	PaymentTypeSub    = "Subscription"
	PaymentTypeFree   = "Free"
)

const (
	CheckBoxUnChecked = "unchecked"
	CheckBoxChecked   = "checked"
)

const (
	CategoryMLKey           = "machine_learning"
	CategoryMLText          = "Machine Learning"
	CategoryCivilEngNameKey = "civil_eng"
	CategoryCivilEngText    = "Civil Engineering"
	CategoryMedicineKey     = "medicine"
	CategoryMedicineText    = "Medicine"
	CategoryChemistryKey    = "chemistry"
	CategoryChemistryText   = "Chemistry"
)

func CreateCategoriesMap() map[string]string {
	categoriesMap := make(map[string]string)
	categoriesMap[CategoryMLKey] = CategoryMLText
	categoriesMap[CategoryCivilEngNameKey] = CategoryCivilEngText
	categoriesMap[CategoryMedicineKey] = CategoryMedicineText
	return categoriesMap
}

func (c *Context) DeleteProduct(product *models.ProductData) error {
	projects, err := c.UserDBController.GetProjectsByProductID(&product.ID)
	if err != nil && err != dbcontrollers.ErrNoProjectForProduct {
		return err
	}

	for _, project := range projects {
		project := project
		if err := c.DeleteProject(&project); err != nil {
			return err
		}
	}

	if err := c.UserDBController.DeleteProduct(&product.ID); err != nil {
		return err
	}

	folder := c.DBModelFunctions.GetFilePath(product.Assets, models.BaseAssetPath, "")
	if err := removeFolder(folder); err != nil {
		return fmt.Errorf("Failed to delete product main app folder. %s", errors.WithStack(err))
	}
	return nil
}

func (c *Context) AddProduct(userID *uuid.UUID, r *rest.Request) (*models.ProductData, error) {
	product, err := c.UserDBController.CreateProduct(userID)
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

func (c *Context) uploadAsset(key string, defaultPath string, asset *models.Asset, r *rest.Request) error {
	fileExt := filepath.Ext(r.FormValue(key))
	if err := c.DBModelFunctions.SetFilePath(asset, key, fileExt); err != nil {
		return err
	}
	path := c.DBModelFunctions.GetFilePath(asset, key, defaultPath)
	if err := c.FileProcessor.UploadFile(key, path, r); err != nil {
		if err == http.ErrMissingFile {
			path := c.DBModelFunctions.GetFilePath(asset, key, defaultPath)
			if err := c.DBModelFunctions.SetFilePath(asset, key, path); err != nil {
				return err
			}
			return nil
		}
		if err2 := c.DBModelFunctions.ClearAsset(asset, key); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return err
	}
	return nil
}

func (c *Context) UploadFiles(asset *models.Asset, r *rest.Request) error {
	if err := c.uploadAsset(ProductAvatarKey, DefaultProductAvatarPath, asset, r); err != nil {
		return err
	}

	if err := c.uploadAsset(ProductMainAppKey, "", asset, r); err != nil {
		return err
	}

	if err := c.uploadAsset(ProductClientApp, "", asset, r); err != nil {
		return err
	}

	return nil
}

func (c *Context) CreateDockerImage(product *models.ProductData, userID *uuid.UUID) error {
	pathString := c.DBModelFunctions.GetFilePath(product.Assets, ProductMainAppKey, "")
	if err := untar(pathString); err != nil {
		return err
	}

	imageName := fmt.Sprintf("%s/%s", userID.String(), product.ID.String())
	sourcePath := path.Join(filepath.Dir(pathString), "build")
	if err := docker.CreateImage(sourcePath, imageName); err != nil {
		return err
	}
	pathString = path.Join(c.DBModelFunctions.GetFilePath(product.Assets, models.BaseAssetPath, ""), "build")
	if err := removeFolder(pathString); err != nil {
		return err
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

func (c *Context) storeProductCategories(details *models.Asset, r *http.Request) {
	categories := CreateCategoriesMap()
	categoryList := make([]string, 0)
	for k := range categories {
		if r.FormValue(k) == "checked" {
			categoryList = append(categoryList, k)
		}
	}
	c.DBModelFunctions.SetField(details, ProductCategoriesKey, categoryList)
}

// SetProductDetails sets the key-value content of product details based on form values.
func (c *Context) SetProductDetails(details *models.Asset, r *http.Request) {
	c.DBModelFunctions.SetField(details, ProductNameKey, r.FormValue(ProductNameKey))
	c.DBModelFunctions.SetField(details, ProductDescriptionKey, r.FormValue(ProductDescriptionKey))
	c.DBModelFunctions.SetField(details, ProductShortDescriptionKey, r.FormValue(ProductShortDescriptionKey))
	c.DBModelFunctions.SetField(details, ProductRequires3DKey, getBooleanString(r.FormValue(ProductRequires3DKey)))
	c.DBModelFunctions.SetField(details, ProductPublicKey, getBooleanString(r.FormValue(ProductPublicKey)))
	c.DBModelFunctions.SetField(details, ProductURLKey, r.FormValue(ProductURLKey))
	c.DBModelFunctions.SetField(details, ProductPricingKey, r.FormValue(ProductPricingKey))
	c.DBModelFunctions.SetField(details, ProductPriceKey, r.FormValue(ProductPriceKey))
	c.DBModelFunctions.SetField(details, ProductTagsKey, r.FormValue(ProductTagsKey))

	c.DBModelFunctions.SetField(details, CreditCardNameKey, r.FormValue(CreditCardNameKey))
	c.DBModelFunctions.SetField(details, CreditCardNumberKey, r.FormValue(CreditCardNumberKey))
	c.DBModelFunctions.SetField(details, CreditCardExpiryKey, r.FormValue(CreditCardExpiryKey))
	c.DBModelFunctions.SetField(details, CreditCardCVCKey, r.FormValue(CreditCardCVCKey))
	c.storeProductCategories(details, r)
}

func (c *Context) UpdateProductData(product *models.ProductData, r *http.Request) error {
	c.SetProductDetails(product.Details, r)

	if err := c.UserDBController.UpdateProductDetails(product); err != nil && err != dbcontrollers.ErrNoProductDetailUpdate {
		return err
	}

	if err := c.UserDBController.UpdateProductAssets(product); err != nil && err != dbcontrollers.ErrNoProductAssetUpdate {
		return err
	}
	return nil
}
