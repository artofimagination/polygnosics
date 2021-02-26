package restcontrollers

import (
	"net/http"
)

func (c *RESTController) Contact(w http.ResponseWriter, r *http.Request) {
	content := c.ContentController.BuildContactContent()
	c.RenderTemplate(w, AboutContact, content)
}

func (c *RESTController) Career(w http.ResponseWriter, r *http.Request) {
	content := make(map[string]interface{})
	c.RenderTemplate(w, AboutCareer, content)
}

func (c *RESTController) About(w http.ResponseWriter, r *http.Request) {
	content := make(map[string]interface{})
	c.RenderTemplate(w, AboutWhoWeAre, content)
}

func (c *RESTController) GeneralContact(w http.ResponseWriter, r *http.Request) {
	content := make(map[string]interface{})
	c.RenderTemplate(w, IndexContact, content)
}

func (c *RESTController) GeneralNews(w http.ResponseWriter, r *http.Request) {
	content := make(map[string]interface{})
	c.RenderTemplate(w, IndexNews, content)
}

func (c *RESTController) IndexHandler(w http.ResponseWriter, r *http.Request) {
	content := make(map[string]interface{})
	c.RenderTemplate(w, IndexPage, content)
}
