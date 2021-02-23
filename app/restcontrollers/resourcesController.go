package restcontrollers

import (
	"net/http"
)

func (c *RESTController) News(w http.ResponseWriter, r *http.Request) {
	content := c.ContentController.BuildNewsContent()
	c.RenderTemplate(w, ResourcesNews, content)
}

func (c *RESTController) Docs(w http.ResponseWriter, r *http.Request) {
	content := c.ContentController.BuildDocsContent()
	c.RenderTemplate(w, ResourcesDocs, content)
}

func (c *RESTController) Tutorials(w http.ResponseWriter, r *http.Request) {
	content := c.ContentController.BuildTutorialsContent()
	c.RenderTemplate(w, ResourcesTutorials, content)
}

func (c *RESTController) FAQ(w http.ResponseWriter, r *http.Request) {
	content := c.ContentController.BuildFAQContent()
	c.RenderTemplate(w, ResourcesFAQ, content)
}

func (c *RESTController) Instructions(w http.ResponseWriter, r *http.Request) {
	content := make(map[string]interface{})
	c.RenderTemplate(w, ResourcesInstructions, content)
}

func (c *RESTController) Examples(w http.ResponseWriter, r *http.Request) {
	content := make(map[string]interface{})
	c.RenderTemplate(w, ResourcesExamples, content)
}
