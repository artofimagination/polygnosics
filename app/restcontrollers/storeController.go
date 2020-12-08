package restcontrollers

import (
	"net/http"
	"polygnosics/app/restcontrollers/page"
)

func StoreHandler(w http.ResponseWriter, r *http.Request) {
	p := getContent()
	page.RenderTemplate(w, "store", p)
}
