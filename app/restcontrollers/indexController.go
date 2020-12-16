package restcontrollers

import (
	"net/http"
)

func (c *RESTController) AboutUsHandler(w http.ResponseWriter, r *http.Request) {
	name := "about"
	p := make(map[string]interface{})
	p["title"] = "About Us"
	p["body"] = "We are awesome"
	c.RenderTemplate(w, name, p)
}

func (c *RESTController) IndexHandler(w http.ResponseWriter, r *http.Request) {
	name := "index"
	p := make(map[string]interface{})
	p["title"] = "Welcome!"
	p["body"] = "Welcome to AI Playground"
	c.RenderTemplate(w, name, p)
}
