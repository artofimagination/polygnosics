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
	ProductOwnerNameKey          = "owner_name"
	ProductOwnerPageNameKey      = "owner_page"
	ProductDetailPageKey         = "detail_path"
	ProductEditPageKey           = "edit_path"
	ProductPriceStringKey        = "price"
	ProductPublicTextKey         = "is_public_text"
	ProductRequires3DTextKey     = "requires_3d_text"
	ProductSupportsClientTextKey = "supports_client_text"
	ProductDeleteIDKey           = "product_to_delete"
	ProductDeleteTextKey         = "delete_text"
	ProductDeleteSuccessTextKey  = "delete_success_text"
	ProductDeleteURLKey          = "delete_url"
)

func generatePriceString(paymentType string, amount string) string {
	if paymentType == businesslogic.PaymentTypeFree {
		return paymentType
	} else if paymentType == businesslogic.PaymentTypeSub {
		return fmt.Sprintf("%s NZD/month", amount)
	} else {
		return fmt.Sprintf("%s NZD", amount)
	}
}

/// GenerateProductContent fills a string nested map with all product details and assets info
func (c *ContentController) generateProductContent(productData *models.ProductData) map[string]interface{} {
	content := make(map[string]interface{})
	content[UserMapKey] = make(map[string]interface{})
	content[businesslogic.ProductAvatarKey] = c.UserDBController.ModelFunctions.GetFilePath(productData.Assets, businesslogic.ProductAvatarKey, businesslogic.DefaultProductAvatarPath)
	content[businesslogic.ProductMainAppKey] = c.UserDBController.ModelFunctions.GetFilePath(productData.Assets, businesslogic.ProductMainAppKey, "")
	content[businesslogic.ProductClientApp] = c.UserDBController.ModelFunctions.GetFilePath(productData.Assets, businesslogic.ProductClientApp, "")
	content[businesslogic.ProductNameKey] = c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductNameKey, "")
	content[businesslogic.ProductURLKey] = c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductURLKey, "")
	content[businesslogic.ProductPublicKey] = convertToCheckboxValue(c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductPublicKey, "").(string))
	content[businesslogic.ProductRequires3DKey] = convertToCheckboxValue(c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductRequires3DKey, "").(string))
	content[businesslogic.ProductDescriptionKey] = c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductDescriptionKey, "")
	content[businesslogic.ProductShortDescriptionKey] = c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductShortDescriptionKey, "")
	content[businesslogic.ProductTagsKey] = c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductTagsKey, "")
	content[businesslogic.ProductCategoriesKey] = c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductCategoriesKey, "")
	pricingType := c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductPricingKey, "").(string)
	content[businesslogic.ProductPricingKey] = pricingType
	price := c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductPriceKey, "").(string)

	content[ProductSupportsClientTextKey] = "No"
	if c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductClientApp, "") == "" {
		content[ProductSupportsClientTextKey] = "Yes"
	}
	content[ProductPriceStringKey] = generatePriceString(pricingType, price)
	content[ProductPublicTextKey] = convertCheckedToYesNo(c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductPublicKey, "").(string))
	content[ProductRequires3DTextKey] = convertCheckedToYesNo(c.UserDBController.ModelFunctions.GetField(productData.Details, businesslogic.ProductRequires3DKey, "").(string))
	content[ProductDetailPageKey] = fmt.Sprintf("/user-main/my-products/details?item-id=%s", productData.ID.String())
	content[ProductEditPageKey] = fmt.Sprintf("/user-main/my-products/edit?item-id=%s", productData.ID.String())
	content[NewProject] = fmt.Sprintf("/user-main/my-products/new-project-wizard?item-id=%s", productData.ID.String())
	content[ProductDeleteIDKey] = productData.ID.String()
	content[ProductDeleteTextKey] = "It will delete all projects started from this product as well"
	content[ProductDeleteSuccessTextKey] = "Your product has been deleted"
	content[ProductDeleteURLKey] = "/user-main/my-products/delete"
	return content
}

// GetProductContent returns the selected product details and assets info.
func (c *ContentController) GetProductContent(productID *uuid.UUID) (map[string]interface{}, error) {
	product, err := c.UserDBController.GetProduct(productID)
	if err != nil {
		return nil, err
	}
	content := c.generateProductContent(product)
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

// GetProductsByCategory organizes product contents by categories. This is just a placeholder solution until proper search is introduced.
func (c *ContentController) GetProductsByCategory(userID *uuid.UUID) (map[string]interface{}, error) {
	products, err := c.UserDBController.GetProductsByUserID(userID)
	if err != nil {
		return nil, err
	}

	categorizedProducts := make(map[string]interface{})
	categories := businesslogic.CreateCategoriesMap()
	for k, v := range categories {
		if categorizedProducts[k] == nil {
			categorizedProducts[k] = make(map[string]interface{})
		}
		categorizedProducts[k].(map[string]interface{})["name"] = v
		for _, product := range products {
			for _, categoryKey := range c.UserDBController.ModelFunctions.GetField(product.ProductData.Details, businesslogic.ProductCategoriesKey, "").([]interface{}) {
				if categoryKey.(string) == k {
					_, ok := categorizedProducts[k].(map[string]interface{})["products"]
					if !ok {
						categorizedProducts[k].(map[string]interface{})["products"] = make([]map[string]interface{}, 0)
					}
					categorizedProducts[k].(map[string]interface{})["products"] = append(categorizedProducts[k].(map[string]interface{})["products"].([]map[string]interface{}), c.generateProductContent(&product.ProductData))
					break
				}
			}
		}
	}
	return categorizedProducts, nil
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
		productContent[i][ProductOwnerPageNameKey] = fmt.Sprintf("/user-main/profile?item-id=%s", c.UserData.ID)
		productContent[i][ProductDetailPageKey] = fmt.Sprintf("/user-main/product?item-id=%s", product.ProductData.ID)
	}

	return productContent, nil
}

func (c *ContentController) GetRecommendedProductsContent(userID *uuid.UUID) ([]map[string]interface{}, error) {
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
		productContent[i][ProductDetailPageKey] = fmt.Sprintf("/user-main/product?item-id=%s", product.ProductData.ID)
	}

	return productContent, nil
}
