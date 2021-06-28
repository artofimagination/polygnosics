package resourcesdb

import (
	"github.com/artofimagination/polygnosics/models"
	"github.com/google/uuid"
)

type Mock struct {
	Categories      models.Categories
	CategoriesError error
	AddedResource   interface{}
	UpdatedResource interface{}
}

func (r *Mock) GetCategories() (models.Categories, error) {
	return r.Categories, r.CategoriesError
}

func (r *Mock) DeleteResource(id string) error {
	return nil
}

func (r *Mock) AddResource(categoryID int, data interface{}) (*uuid.UUID, error) {
	r.AddedResource = data
	uuid := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	return &uuid, nil
}

func (r *Mock) GetResource(id string, data interface{}) (*models.Resource, error) {
	return nil, nil
}

func (r *Mock) GetResources(ids []string) ([]models.Resource, error) {
	return nil, nil
}

func (r *Mock) GetResourcesByCategory(id int) ([]models.Resource, error) {
	return nil, nil
}

func (r *Mock) UpdateResource(resource *models.Resource, data interface{}) error {
	r.UpdatedResource = data
	return nil
}
