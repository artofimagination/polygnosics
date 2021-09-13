package resourcesdb

import (
	"encoding/json"
	"fmt"

	"github.com/artofimagination/polygnosics/models"
	"github.com/google/uuid"
)

const (
	GetCategoriesURI = "/get-categories"

	AddResourceURI            = "/add-resource"
	GetResourceURI            = "/get-resource-by-id"
	GetResourcesByCategoryURI = "/get-resources-by-category"
	GetResourcesByIDsURI      = "/get-resources-by-ids"
	UpdateResourceURI         = "/update-resource"
	DeleteResourceURI         = "/delete-resource"
)

func (c *RESTController) DeleteResource(id string) error {
	data := make(map[string]string)
	data["id"] = id
	_, err := c.Post(DeleteResourceURI, data)
	if err != nil {

		return err
	}
	return nil
}

// AddResource sends a new add request to the database server.
// The added resource then will be updated in @resource (updated with valid ID).
func (c *RESTController) AddResource(categoryID int, data interface{}) (*uuid.UUID, error) {
	resource := &models.Resource{
		Category: categoryID,
	}
	if err := c.modelFunc.ToResourceContent(resource, data); err != nil {
		return nil, err
	}
	data, err := c.Post(AddResourceURI, resource)
	if err != nil {
		return nil, err
	}
	if err := c.modelFunc.ToResource(resource, data); err != nil {
		return nil, err
	}
	if err := c.modelFunc.FromResource(resource, data); err != nil {
		return nil, err
	}

	return &resource.ID, nil
}

func (c *RESTController) GetCategories() (models.Categories, error) {
	data, err := c.Get(GetCategoriesURI, "")
	if err != nil {
		return nil, err
	}

	categories := make(models.Categories, 0)
	for _, dataItem := range data.([]interface{}) {
		bytesData, err := json.Marshal(dataItem)
		if err != nil {
			return nil, err
		}

		category := models.Category{}
		if err := json.Unmarshal(bytesData, &category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *RESTController) GetResource(id string, parsedData interface{}) (*models.Resource, error) {
	params := fmt.Sprintf("?id=%s", id)
	data, err := c.Get(GetResourceURI, params)
	if err != nil {
		return nil, err
	}

	resource := &models.Resource{}
	if err := c.modelFunc.ToResource(resource, data); err != nil {
		return nil, err
	}

	if err := c.modelFunc.FromResource(resource, parsedData); err != nil {
		return nil, err
	}

	return resource, nil
}

func (c *RESTController) GetResources(ids []string) ([]models.Resource, error) {
	params := "?"
	for _, id := range ids {
		params = fmt.Sprintf("%sid=%s&", params, id)
	}
	data, err := c.Get(GetResourcesByIDsURI, params[0:len(params)-1])
	if err != nil {
		return nil, err
	}

	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resources := []models.Resource{}
	if err := json.Unmarshal(bytesData, &resources); err != nil {
		return nil, err
	}
	return resources, nil
}

func (c *RESTController) GetResourcesByCategory(id int) ([]models.Resource, error) {
	params := fmt.Sprintf("?category-id=%d", id)
	data, err := c.Get(GetResourcesByCategoryURI, params)
	if err != nil {
		return nil, err
	}

	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resources := []models.Resource{}
	if err := json.Unmarshal(bytesData, &resources); err != nil {
		return nil, err
	}
	return resources, nil
}

func (c *RESTController) UpdateResource(resource *models.Resource, data interface{}) error {
	if err := c.modelFunc.ToResourceContent(resource, data); err != nil {
		return err
	}

	_, err := c.Post(UpdateResourceURI, resource)
	if err != nil {
		return err
	}
	return nil
}
