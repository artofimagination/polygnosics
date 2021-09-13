package resourcesdb

import (
	"github.com/artofimagination/polygnosics/models"
	"github.com/artofimagination/polygnosics/rest"
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
	modelFunc     models.ResourceModels
	ServerAddress *rest.Server
}

func (c *RESTController) Post(path string, parameters interface{}) (interface{}, error) {
	return rest.Post(c.ServerAddress.GetAddress(), path, parameters)
}

func (c *RESTController) Get(path string, parameters string) (interface{}, error) {
	return rest.Get(c.ServerAddress.GetAddress(), path, parameters)
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
