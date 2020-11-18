package restcontrollers

import (
	"net/http"
	"polygnosics/app/restcontrollers/page"
)

func MyProductsHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := getContent(w, r)
	page.RenderTemplate(w, "my-products", p)
}
