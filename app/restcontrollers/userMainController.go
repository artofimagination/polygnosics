package restcontrollers

import (
	"net/http"

	"polygnosics/app/restcontrollers/page"
)

// NewProject is the handler for the page that is responsible for creating a new project.
func NewProject(w http.ResponseWriter, r *http.Request) {
}

func UserSettings(w http.ResponseWriter, r *http.Request) {
}

// UserMain renders the main page after login.
func UserMain(w http.ResponseWriter, r *http.Request) {
	p := make(map[string]interface{})
	page.RenderTemplate(w, "user-main", &p)
}
