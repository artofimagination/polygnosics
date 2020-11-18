package restcontrollers

import (
	"fmt"
	"net/http"

	"polygnosics/app/restcontrollers/page"

	"github.com/artofimagination/mysql-user-db-go-interface/models"

	"github.com/pkg/errors"
)

// NewProject is the handler for the page that is responsible for creating a new project.
func NewProject(w http.ResponseWriter, r *http.Request) {
}

func UserSettings(w http.ResponseWriter, r *http.Request) {
}

// UserMainHandler renders the main page after login.
func UserMainHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := getContent(w, r)
	page.RenderTemplate(w, "user-main", p)
}

// ProfileHandler renders the profile page template.
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := getContent(w, r)
	page.RenderTemplate(w, "profile", p)
}

// UploadAvatarHandler processes avatar upload request.
// Stores the image in the location defined by the asset ID and avatar ID.
// The file is named by the avatar ID and the folder is determined by the asset ID.
func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	p, data := getContent(w, r)

	asset := (*data)["user-assets"].(*models.Asset)
	if err := asset.SetID(models.Avatar); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update avatar asset. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	path, err := asset.GetPath(models.Avatar)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load assets path. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	if err := page.UploadFile(path, r); err != nil {
		if err2 := asset.ClearID(models.Avatar); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		http.Error(w, fmt.Sprintf("Failed to upload file. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	if err := page.SaveAssetReferences(asset); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save asset. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	(*p)["assets"].(map[string]interface{})[models.Avatar] = path
	http.Redirect(w, r, "/user-main/profile", http.StatusSeeOther)
}
