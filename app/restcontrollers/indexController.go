package restcontrollers

import (
	"net/http"

	"polygnosics/app/restcontrollers/page"
)

func AboutUsHandler(w http.ResponseWriter, r *http.Request) {
	name := "about"
	p := make(map[string]interface{})
	p["title"] = "About Us"
	p["body"] = "We are awesome"
	page.RenderTemplate(w, name, p)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	name := "index"
	p := make(map[string]interface{})
	p["title"] = "Welcome!"
	p["body"] = "Welcome to AI Playground"
	page.RenderTemplate(w, name, p)
}
