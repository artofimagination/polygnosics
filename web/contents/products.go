package contents

import (
	"fmt"
	"net/http"

	"polygnosics/app/businesslogic"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

// Details and assets field keys
const (
	ProductAvatar           = "product_avatar"
	ProductMainApp          = "main_app"
	ProductClientApp        = "client-app"
	ProductDescription      = "product_description"
	ProductName             = "product_name"
	ProductPath             = "product_path"
	ProductFolder           = "product_folder"
	ProductRequires3D       = "requires_3d"
	ProductURL              = "product_url"
	ProductPublic           = "is_public"
	ProductOwnerNameKey     = "owner_name"
	ProductOwnerPageNameKey = "owner_page"
	ProductDetailPageKey    = "product_detail"
)

/// GenerateProductContent fills a string nested map with all product details and assets info
func (c *ContentController) generateProductContent(productData *models.ProductData) map[string]interface{} {
	content := make(map[string]interface{})
	content[ProductAvatar] = c.UserDBController.ModelFunctions.GetFilePath(productData.Assets, ProductAvatar, businesslogic.DefaultProductAvatarPath)
	content[ProductName] = c.UserDBController.ModelFunctions.GetField(productData.Details, ProductName, "")
	content[ProductPublic] = c.UserDBController.ModelFunctions.GetField(productData.Details, ProductPublic, "")
	content[ProductDescription] = c.UserDBController.ModelFunctions.GetField(productData.Details, ProductDescription, "")
	content[ProductPath] = fmt.Sprintf("/user-main/my-products/details?product=%s", productData.ID.String())
	content[NewProject] = fmt.Sprintf("/user-main/my-products/new-project-wizard?product=%s", productData.ID.String())
	return content
}

// SetProductDetails sets the key-value content of product details based on form values.
func (c *ContentController) SetProductDetails(details *models.Asset, r *http.Request) {
	c.UserDBController.ModelFunctions.SetField(details, ProductName, r.FormValue("productName"))
	c.UserDBController.ModelFunctions.SetField(details, ProductDescription, r.FormValue("productDescription"))
	c.UserDBController.ModelFunctions.SetField(details, ProductRequires3D, r.FormValue("requires3D"))
	c.UserDBController.ModelFunctions.SetField(details, ProductPublic, getBooleanString(r.FormValue("publicProduct")))
	c.UserDBController.ModelFunctions.SetField(details, ProductURL, r.FormValue("productUrl"))
}

// GetProductContent returns the selected product details and assets info.
func (c *ContentController) GetProductContent(productID *uuid.UUID) (map[string]interface{}, error) {
	product, err := c.UserDBController.GetProduct(productID)
	if err != nil {
		return nil, err
	}
	return c.generateProductContent(product), nil
}

// GetUserProductContent gathers the content of each product belonging to the specified user.
func (c *ContentController) GetUserProductContent(userID *uuid.UUID) ([]map[string]interface{}, error) {
	products, err := c.UserDBController.GetProductsByUserID(userID)
	if err != nil {
		return nil, err
	}

	productContent := make([]map[string]interface{}, len(products))
	for i, product := range products {
		productContent[i] = c.generateProductContent(&product.ProductData)
		productContent[i][ProductOwnerNameKey] = c.UserData.Name
	}

	return productContent, nil
}

// GetRecentProductsContent gathers the content of the latest 4 products
func (c *ContentController) GetRecentProductsContent(userID *uuid.UUID) ([]map[string]interface{}, error) {
	products, err := c.UserDBController.GetProductsByUserID(userID)
	if err != nil {
		return nil, err
	}

	limit := 4
	productContent := make([]map[string]interface{}, limit)
	for i, product := range products {
		if limit == 0 {
			break
		}
		limit--
		productContent[i] = c.generateProductContent(&product.ProductData)
		productContent[i][ProductOwnerNameKey] = c.UserData.Name
		productContent[i][ProductOwnerPageNameKey] = fmt.Sprintf("/user-main/profile?user=%s", c.UserData.ID)
		productContent[i][ProductDetailPageKey] = fmt.Sprintf("/user-main/product?product=%s", product.ProductData.ID)
	}

	return productContent, nil
}
