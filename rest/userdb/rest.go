package userdb

import (
	"github.com/artofimagination/polygnosics/rest"
	"github.com/gorilla/mux"
)

type RESTController struct {
	ServerAddress *rest.Server
}

func (c *RESTController) Post(path string, parameters interface{}) (interface{}, error) {
	return rest.Post(c.ServerAddress.GetAddress(), path, parameters)
}

func (c *RESTController) Get(path string, parameters string) (interface{}, error) {
	return rest.Get(c.ServerAddress.GetAddress(), path, parameters)
}

func NewRESTController() *RESTController {
	controller := &RESTController{}
	return controller
}

// AddRouting adds front end endpoints.
func (c *RESTController) AddRouting(r *mux.Router) {

}
