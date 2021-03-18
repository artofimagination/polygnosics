package userdb

import (
	"encoding/json"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/polygnosics/rest"
	"github.com/google/uuid"
)

const (
	productPathAdd           = "/add-product"
	productPathDelete        = "/delete-product"
	productPathGet           = "/get-product"
	productPathUpdateDetails = "/update-product-details"
	productPathUpdateAssets  = "/update-product-assets"
)

func (c *RESTController) CreateProduct(owner *uuid.UUID) (*models.ProductData, error) {
	params := make(map[string]interface{})
	params["owner_id"] = owner.String()
	data, err := rest.Post(rest.UserDBAddress, productPathAdd, params)
	if err != nil {
		return nil, err
	}

	product := &models.ProductData{}
	if err := json.Unmarshal(data.([]byte), &product); err != nil {
		return nil, err
	}

	return product, nil
}

func (c *RESTController) DeleteProduct(productID *uuid.UUID) error {
	params := make(map[string]interface{})
	params["id"] = productID.String()
	_, err := rest.Post(rest.UserDBAddress, productPathDelete, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *RESTController) GetProduct(productID *uuid.UUID) (*models.ProductData, error) {
	params := fmt.Sprintf("?id=%s", productID.String())
	data, err := rest.Get(rest.UserDBAddress, productPathGet, params)
	if err != nil {
		return nil, err
	}

	productData := &models.ProductData{}
	if err := json.Unmarshal(data.([]byte), &productData); err != nil {
		return nil, err
	}
	return productData, nil
}

func (c *RESTController) UpdateProductDetails(productData *models.ProductData) error {
	params := make(map[string]interface{})
	productDataBytes, err := json.Marshal(productData)
	if err != nil {
		return err
	}
	params["product-data"] = productDataBytes
	_, err = rest.Post(rest.UserDBAddress, productPathUpdateDetails, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *RESTController) UpdateProductAssets(productData *models.ProductData) error {
	params := make(map[string]interface{})
	productDataBytes, err := json.Marshal(productData)
	if err != nil {
		return err
	}
	params["product-data"] = productDataBytes
	_, err = rest.Post(rest.UserDBAddress, productPathUpdateAssets, params)
	if err != nil {
		return err
	}
	return nil
}
