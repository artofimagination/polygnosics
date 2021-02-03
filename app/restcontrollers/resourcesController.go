package restcontrollers

import (
	"net/http"
)

func (c *RESTController) News(w http.ResponseWriter, r *http.Request) {
	content := c.ContentController.BuildNewsContent()
	c.RenderTemplate(w, "news", content)
}
