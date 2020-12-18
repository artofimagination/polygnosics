package restcontrollers

import (
	"fmt"
	"net/http"
	"os"
	"polygnosics/app/restcontrollers/contents"
	"polygnosics/app/restcontrollers/session"
	"text/template"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type RESTController struct {
	UserDBController  *dbcontrollers.MYSQLController
	ContentController *contents.ContentController
}

var htmls = []string{
	"/web/templates/about.html",
	"/web/templates/error.html",
	"/web/templates/index.html",
	"/web/templates/confirm.html",
	"/web/templates/user/user-main.html",
	"/web/templates/user/profile.html",
	"/web/templates/user/user-settings.html",
	"/web/templates/user/new-project.html",
	"/web/templates/project/run.html",
	"/web/templates/project/project-details.html",
	"/web/templates/project/my-projects.html",
	"/web/templates/project/new-project-wizard.html",
	"/web/templates/auth_signup.html",
	"/web/templates/auth_login.html",
	"/web/templates/products/store.html",
	"/web/templates/products/new-product-wizard.html",
	"/web/templates/products/my-products.html",
	"/web/templates/products/details.html",
	"/web/templates/components/side-bar.html",
	"/web/templates/components/content-header.html",
}
var paths = []string{}

const (
	GET      = "GET"
	Confirm  = "confirm"
	UserMain = "user-main"
)

func NewRESTController(userDB *dbcontrollers.MYSQLController) *RESTController {
	controller := &RESTController{
		UserDBController: userDB,
		ContentController: &contents.ContentController{
			UserDBController: userDB,
		},
	}
	return controller
}

// MakeHandler creates the page handler and check the route validity.
func (c *RESTController) MakeHandler(fn func(http.ResponseWriter, *http.Request), router *mux.Router, isPublicPage bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Security-Policy", "default-src 'self' data: 'unsafe-inline' 'unsafe-eval'")

		routeMatch := mux.RouteMatch{}
		if matched := router.Match(r, &routeMatch); !matched {
			http.Error(w, "Url does not exist", http.StatusInternalServerError)
			return
		}

		if !isPublicPage {
			sess, err := session.Store.Get(r, "cookie-name")
			if err != nil {
				http.Error(w, "Unable to retrieve session cookie.", http.StatusForbidden)
				return
			}

			uuidString, ok := sess.Values["user"].(string)
			if !ok {
				http.Error(w, "Unable to decode session cookie.", http.StatusInternalServerError)
				return
			}

			userUUID, err := uuid.Parse(uuidString)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to get user id. %s", errors.WithStack(err)), http.StatusInternalServerError)
				return
			}

			user, err := c.UserDBController.GetUser(&userUUID)
			if err != nil {
				http.Error(w, "Unable to retrieve user info", http.StatusInternalServerError)
				return
			}

			match, err := session.IsAuthenticated(user.ID, sess, r)
			if err != nil {
				errorString := fmt.Sprintf("Unable to check session cookie:\n%s\n", errors.WithStack(err))
				http.Error(w, errorString, http.StatusInternalServerError)
				return
			}

			if !match {
				http.Error(w, "Forbidden access", http.StatusForbidden)
				return
			}
		}
		fn(w, r)
	}
}

// RenderTemplate renders html.
func (c *RESTController) RenderTemplate(w http.ResponseWriter, tmpl string, p map[string]interface{}) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(paths) == 0 {
		for i := 0; i < len(htmls); i++ {
			paths = append(paths, wd+htmls[i])
		}
	}

	t := template.Must(template.ParseFiles(paths...))

	err = t.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleError creates page details and renders html template for an error modal.
func (c *RESTController) HandleError(route string, errorStr string, w http.ResponseWriter) {
	p := make(map[string]interface{})
	p["message"] = errorStr
	p["route"] = fmt.Sprintf("/%s", route)
	c.RenderTemplate(w, Confirm, p)
}
