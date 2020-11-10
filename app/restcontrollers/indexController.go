package restcontrollers

import (
	"net/http"

	"polygnosics/app/utils/page"
)

func AboutUsHandler(w http.ResponseWriter, r *http.Request) {
	name := "about"
	p := page.CreatePage(name)
	p.Data["title"] = "About Us"
	p.Data["body"] = "We are awesome"
	page.RenderTemplate(w, name, p)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	name := "index"
	p := page.CreatePage(name)
	p.Data["title"] = "Welcome!"
	p.Data["body"] = "Welcome to AI Playground"
	page.RenderTemplate(w, name, p)
}
