package restcontrollers

import (
	"net/http"
)

func (c *RESTController) StoreHandler(w http.ResponseWriter, r *http.Request) {
	p := c.ContentController.GetUserContent()
	c.RenderTemplate(w, "store", p)
}
