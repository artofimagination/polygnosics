package restcontrollers

import (
	"fmt"
	"net/http"

	"polygnosics/app"
	"polygnosics/app/restcontrollers/auth"
	"polygnosics/app/restcontrollers/page"

	"github.com/pkg/errors"
)

// NewProject is the handler for the page that is responsible for creating a new project.
func NewProject(w http.ResponseWriter, r *http.Request) {
}

func UserSettings(w http.ResponseWriter, r *http.Request) {
}

// UserMainHandler renders the main page after login.
func UserMainHandler(w http.ResponseWriter, r *http.Request) {
	p := getUserContent()
	page.RenderTemplate(w, "user-main", p)
}

// ProfileHandler renders the profile page template.
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	p := getUserContent()
	page.RenderTemplate(w, "profile", p)
}

// UploadAvatarHandler processes avatar upload request.
// Stores the image in the location defined by the asset ID and avatar ID.
// The file is named by the avatar ID and the folder is determined by the asset ID.
func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	p := getUserContent()

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form. %s", errors.WithStack(err)), http.StatusInternalServerError)
	}

	path, err := uploadUserFile(UserAvatar, DefaultUserAvatarPath, "asset", r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload asset. %s", errors.WithStack(err)), http.StatusInternalServerError)
	}

	if err := app.ContextData.UserDBController.UpdateUserAssets(auth.UserData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update asset. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	p["assets"].(map[string]interface{})[UserAvatar] = path
	http.Redirect(w, r, "/user-main/profile", http.StatusSeeOther)
}
