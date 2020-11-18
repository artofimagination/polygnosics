package restcontrollers

import (
	"fmt"
	"net/http"

	"polygnosics/app/restcontrollers/page"
	"polygnosics/app/restcontrollers/session"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

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
