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

// UserMain renders the main page after login.
func UserMain(w http.ResponseWriter, r *http.Request) {
	session, err := session.Store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get cookie. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	user, err := mysqldb.GetUserByID(uuid.MustParse(session.Values["user"].(string)))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load user. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	// Assets are introduced later, some users will have no assets associated yet.
	// Second solution could be to do this during migration, but with a large database it would take long time
	if user.AssetsID == models.NullUUID {
		if err := mysqldb.UpdateAssetID(user); err != nil {
			http.Error(w, fmt.Sprintf("Failed to update user assets. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}
	}

	asset, err := page.LoadAssetReferences(&user.AssetsID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load user assets. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	p := make(map[string]interface{})
	p["assets"] = make(map[string]interface{})
	p["assets"].(map[string]interface{})["avatar"] = asset.GetAvatarPath()
	p["username"] = user.Name
	page.RenderTemplate(w, "user-main", &p)
}
