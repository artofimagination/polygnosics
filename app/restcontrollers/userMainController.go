package restcontrollers

import (
	"fmt"
	"log"
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
	p := getContent()
	log.Println(p)
	page.RenderTemplate(w, "user-main", p)
}

// ProfileHandler renders the profile page template.
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	p := getContent()
	page.RenderTemplate(w, "profile", p)
}

// UploadAvatarHandler processes avatar upload request.
// Stores the image in the location defined by the asset ID and avatar ID.
// The file is named by the avatar ID and the folder is determined by the asset ID.
func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	p := getContent()

	if err := auth.UserData.Assets.SetImagePath(UserAvatar); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update avatar asset. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	path := auth.UserData.Assets.GetImagePath(UserAvatar, DefaultAvatarPath)

	if err := uploadFile(path, r); err != nil {
		if err2 := auth.UserData.Assets.ClearAsset(UserAvatar); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		http.Error(w, fmt.Sprintf("Failed to upload file. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	if err := app.ContextData.UserDBController.UpdateUserAssets(auth.UserData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update asset. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	p["assets"].(map[string]interface{})[UserAvatar] = path
	http.Redirect(w, r, "/user-main/profile", http.StatusSeeOther)
}
