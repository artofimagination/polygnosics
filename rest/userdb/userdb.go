package userdb

import (
	"github.com/gorilla/mux"
)

type RESTController struct {
}

func NewRESTController() *RESTController {
	controller := &RESTController{}
	return controller
}

// AddRouting adds front end endpoints.
func (c *RESTController) AddRouting(r *mux.Router) {

}
