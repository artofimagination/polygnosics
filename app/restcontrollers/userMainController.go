package restcontrollers

import (
	"fmt"
	"net/http"

	"polygnosics/app/models"
	"polygnosics/app/restcontrollers/page"
	"polygnosics/app/restcontrollers/session"
	"polygnosics/app/services/db/mysqldb"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// NewProject is the handler for the page that is responsible for creating a new project.
func NewProject(w http.ResponseWriter, r *http.Request) {
}

func UserSettings(w http.ResponseWriter, r *http.Request) {
}

func getContent(w http.ResponseWriter, r *http.Request) (*map[string]interface{}, *map[string]interface{}) {
	session, err := session.Store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get cookie. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return nil, nil
	}

	user, err := mysqldb.GetUserByID(uuid.MustParse(session.Values["user"].(string)))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load user. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return nil, nil
	}

	// Assets are introduced later, some users will have no assets associated yet.
	// Second solution could be to do this during migration, but with a large database it would take long time
	if user.AssetsID == models.NullUUID {
		if err := mysqldb.UpdateAssetID(user); err != nil {
			http.Error(w, fmt.Sprintf("Failed to update user assets. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return nil, nil
		}
	}

	asset, err := page.LoadAssetReferences(&user.AssetsID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load user assets. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return nil, nil
	}

	p := make(map[string]interface{})
	p["assets"] = make(map[string]interface{})
	path, err := asset.GetPath(models.Avatar)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load assets path. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return nil, nil
	}
	p["assets"].(map[string]interface{})[models.Avatar] = path
	p["texts"] = make(map[string]interface{})
	p["texts"].(map[string]interface{})["avatar-upload"] = "Upload your avatar"
	p["texts"].(map[string]interface{})["username"] = user.Name

	data := make(map[string]interface{})
	data["user-data"] = user
	data["user-assets"] = asset
	return &p, &data
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
