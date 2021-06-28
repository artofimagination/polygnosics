package resourcesdb

import (
	"github.com/artofimagination/polygnosics/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ResourceDBInterface interface {
	GetCategories() (models.Categories, error)
	DeleteResource(id string) error
	AddResource(categoryID int, data interface{}) (*uuid.UUID, error)
	GetResource(id string, parsedData interface{}) (*models.Resource, error)
	GetResources(ids []string) ([]models.Resource, error)
	GetResourcesByCategory(id int) ([]models.Resource, error)
	UpdateResource(resource *models.Resource, data interface{}) error
}

type RESTController struct {
	modelFunc models.ResourceModels
}

func NewRESTController() *RESTController {
	controller := &RESTController{
		modelFunc: &models.ResourceModelImpl{},
	}
	return controller
}

// AddRouting adds front end endpoints.
func (c *RESTController) AddRouting(r *mux.Router) {

}
