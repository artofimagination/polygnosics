package contents

import (
	"fmt"

	"polygnosics/app/businesslogic"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

// Details and assets field keys
const (
	ProductMapKey                = "product"
	ProductDeleteIDKey           = "product_to_delete"
	ProductOwnerNameKey          = "owner_name"
	ProductOwnerPageNameKey      = "owner_page"
	ProductDetailPageKey         = "detail_path"
	ProductEditPageKey           = "edit_path"
	Product3rdPartyDetailPageKey = "3rdparty_detail_path"
)

/// GenerateProductContent fills a string nested map with all product details and assets info
func (c *ContentController) generateProductContent(productData *models.ProductData) map[string]interface{} {
	content := make(map[string]interface{})
	content[UserMapKey] = make(map[string]interface{})
	content[businesslogic.ProductAvatarKey] = c.UserDBController.ModelFunctions.GetFilePath(productData.Assets, businesslogic.ProductAvatarKey, businesslogic.DefaultProductAvatarPath)
	content[businesslogic.ProductMainAppKey] = c.UserDBController.ModelFunctions.GetFilePath(productData.Assets, businesslogic.ProductMainAppKey, "")
	content[businesslogic.ProductClientApp] = c.UserDBController.ModelFunctions.GetFilePath(productData.Assets, businesslogic.ProductClientApp, "")
	content[businesslogic.ProductName] = c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductName, "")
	content[businesslogic.ProductURL] = c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductURL, "")
	content[businesslogic.ProductPublic] = convertToCheckboxValue(c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductPublic, ""))
	content[businesslogic.ProductRequires3D] = convertToCheckboxValue(c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductRequires3D, ""))
	content[businesslogic.ProductDescription] = c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductDescription, "")
	content[ProductDetailPageKey] = fmt.Sprintf("/user-main/my-products/details?product=%s", productData.ID.String())
	content[ProductEditPageKey] = fmt.Sprintf("/user-main/my-products/edit?product=%s", productData.ID.String())
	content[NewProject] = fmt.Sprintf("/user-main/my-products/new-project-wizard?product=%s", productData.ID.String())
	content[ProductDeleteIDKey] = productData.ID.String()
	return content
}

// GetProductContent returns the selected product details and assets info.
func (c *ContentController) GetProductContent(productID *uuid.UUID) (map[string]interface{}, error) {
	content := make(map[string]interface{})
	product, err := c.UserDBController.GetProduct(productID)
	if err != nil {
		return nil, err
	}
	content[ProductMapKey] = c.generateProductContent(product)
	return content, nil
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
	if len(products) < limit {
		limit = len(products)
	}
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
