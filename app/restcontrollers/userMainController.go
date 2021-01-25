package restcontrollers

import (
	"fmt"
	"net/http"

	"polygnosics/app/restcontrollers/contents"

	"github.com/pkg/errors"
)

// NewProject is the handler for the page that is responsible for creating a new project.
func NewProject(w http.ResponseWriter, r *http.Request) {
}

func UserSettings(w http.ResponseWriter, r *http.Request) {
}

// UserMainHandler renders the main page after login.
func (c *RESTController) UserMainHandler(w http.ResponseWriter, r *http.Request) {
	p := c.ContentController.GetUserContent()
	c.RenderTemplate(w, UserMain, p)
}

// ProfileHandler renders the profile page template.
func (c *RESTController) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	p := c.ContentController.GetUserContent()
	c.RenderTemplate(w, "profile", p)
}

// UploadAvatarHandler processes avatar upload request.
// Stores the image in the location defined by the asset ID and avatar ID.
// The file is named by the avatar ID and the folder is determined by the asset ID.
func (c *RESTController) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	p := c.ContentController.GetUserContent()

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form. %s", errors.WithStack(err)), http.StatusInternalServerError)
	}

	if c.ContentController.UserData == nil {
		http.Error(w, "User is not configured", http.StatusInternalServerError)
	}

	err := c.ContentController.UploadFile(c.ContentController.UserData.Assets, contents.UserAvatar, contents.DefaultUserAvatarPath, "asset", r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload asset. %s", errors.WithStack(err)), http.StatusInternalServerError)
	}

	if err := c.UserDBController.UpdateUserAssets(c.ContentController.UserData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update asset. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	p["assets"].(map[string]interface{})[contents.UserAvatar] =
		c.UserDBController.ModelFunctions.GetFilePath(
			c.ContentController.UserData.Assets,
			contents.UserAvatar,
			contents.DefaultUserAvatarPath)
	http.Redirect(w, r, "/user-main/profile", http.StatusSeeOther)
}
