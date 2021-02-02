package restcontrollers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// ProfileHandler renders the profile page template.
func (c *RESTController) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent(ErrFailedToParseForm))
		return
	}

	userIDString := r.FormValue("user")
	userUUID, err := uuid.Parse(userIDString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user id. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	content, err := c.ContentController.BuildProfileContent(&userUUID)
	if err != nil {
		errString := fmt.Sprintf("Failed to get profile page content. %s", errors.WithStack(err))
		c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent(errString))
		return
	}
	c.RenderTemplate(w, "profile", content)
}

func (c *RESTController) ProfileEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == GET {
		content, err := c.ContentController.BuildProfileContent(&c.ContentController.UserData.ID)
		if err != nil {
			errString := fmt.Sprintf("Failed to get profile page content. %s", errors.WithStack(err))
			c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent(errString))
			return
		}
		c.RenderTemplate(w, "profile-edit", content)
	} else {
		if err := r.ParseForm(); err != nil {
			c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent(ErrFailedToParseForm))
			return
		}

		if err := c.ContentController.StoreUserInfo(r); err != nil {
			errString := fmt.Sprintf("Failed to get store user info. %s", errors.WithStack(err))
			c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent(errString))
			return
		}

		content, err := c.ContentController.BuildProfileContent(&c.ContentController.UserData.ID)
		if err != nil {
			errString := fmt.Sprintf("Failed to get profile page content. %s", errors.WithStack(err))
			c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent(errString))
			return
		}
		c.RenderTemplate(w, "profile", content)
	}

}
