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
		if err := c.DeleteProject(&project); err != nil {
			return err
		}
	}

	if err := c.UserDBController.DeleteProduct(&product.ID); err != nil {
		return err
	}

	folder := c.UserDBController.ModelFunctions.GetFilePath(product.Assets, models.BaseAssetPath, "")
	if err := removeFolder(folder); err != nil {
		return fmt.Errorf("Failed to delete product main app folder. %s", errors.WithStack(err))
	}
	return nil
}

func (c *Context) AddProduct(userID *uuid.UUID, ProductNameKey string, r *http.Request) (*models.ProductData, error) {
	product, err := c.UserDBController.CreateProduct(
		ProductNameKey,
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
		return err
	}

	imageName := fmt.Sprintf("%s/%s", userID.String(), product.ID.String())
	sourcePath := path.Join(filepath.Dir(pathString), "build")
	if err := docker.CreateImage(sourcePath, imageName); err != nil {
		return err
	}
	pathString = path.Join(c.UserDBController.ModelFunctions.GetFilePath(product.Assets, models.BaseAssetPath, ""), "build")
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
	for k, _ := range categories {
		if r.FormValue(k) == "checked" {
			categoryList = append(categoryList, k)
		}
	}
	c.UserDBController.ModelFunctions.SetField(details, ProductCategoriesKey, categoryList)
}

// SetProductDetails sets the key-value content of product details based on form values.
func (c *Context) SetProductDetails(details *models.Asset, r *http.Request) {
	c.UserDBController.ModelFunctions.SetField(details, ProductNameKey, r.FormValue(ProductNameKey))
	c.UserDBController.ModelFunctions.SetField(details, ProductDescriptionKey, r.FormValue(ProductDescriptionKey))
	c.UserDBController.ModelFunctions.SetField(details, ProductShortDescriptionKey, r.FormValue(ProductShortDescriptionKey))
	c.UserDBController.ModelFunctions.SetField(details, ProductRequires3DKey, getBooleanString(r.FormValue(ProductRequires3DKey)))
	c.UserDBController.ModelFunctions.SetField(details, ProductPublicKey, getBooleanString(r.FormValue(ProductPublicKey)))
	c.UserDBController.ModelFunctions.SetField(details, ProductURLKey, r.FormValue(ProductURLKey))
	c.UserDBController.ModelFunctions.SetField(details, ProductPricingKey, r.FormValue(ProductPricingKey))
	c.UserDBController.ModelFunctions.SetField(details, ProductPriceKey, r.FormValue(ProductPriceKey))
	c.UserDBController.ModelFunctions.SetField(details, ProductTagsKey, r.FormValue(ProductTagsKey))

	c.UserDBController.ModelFunctions.SetField(details, CreditCardNameKey, r.FormValue(CreditCardNameKey))
	c.UserDBController.ModelFunctions.SetField(details, CreditCardNumberKey, r.FormValue(CreditCardNumberKey))
	c.UserDBController.ModelFunctions.SetField(details, CreditCardExpiryKey, r.FormValue(CreditCardExpiryKey))
	c.UserDBController.ModelFunctions.SetField(details, CreditCardCVCKey, r.FormValue(CreditCardCVCKey))
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
