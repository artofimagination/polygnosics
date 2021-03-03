package restcontrollers

import (
	"fmt"
	"net/http"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/pkg/errors"
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
	// TODO Issue#107: Replace this with proper way of detecting if root has already been created.
	_, err := c.BackendContext.GetUserByEmail("root@test.com")
	name := IndexPage
	if err == dbcontrollers.ErrUserNotFound {
		if r.Host == "polygnosics.localhost" {
			name = "auth_signup"
		} else {
			c.HandleError(w, "Server is not configured yet", http.StatusInternalServerError, IndexPath)
			return
		}

	} else if err != nil {
		c.HandleError(w, fmt.Sprintf("Failed to get root user. %s", errors.WithStack(err)), http.StatusInternalServerError, IndexPath)
		return
	}

	c.RenderTemplate(w, name, content)
}
