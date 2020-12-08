package restcontrollers

import (
	"net/http"
	"polygnosics/app/restcontrollers/page"
)

func MyProductsHandler(w http.ResponseWriter, r *http.Request) {
	p := getContent()
	page.RenderTemplate(w, "my-products", p)
}
