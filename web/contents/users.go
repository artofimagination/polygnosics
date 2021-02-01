package contents

import (
	"polygnosics/app/businesslogic"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
)

// Details and assets field keys
const (
	UserAvatar     = "user_avatar"
	UserBackground = "user_background"
)

// GetUserContent fills a string nested map with all user details and assets info
func (c *ContentController) GetUserContent(user *models.UserData) map[string]interface{} {
	p := make(map[string]interface{})
	p["assets"] = make(map[string]interface{})
	path := c.UserDBController.ModelFunctions.GetFilePath(user.Assets, UserAvatar, businesslogic.DefaultUserAvatarPath)
	p["assets"].(map[string]interface{})[UserAvatar] = path
	p["texts"] = make(map[string]interface{})
	p["texts"].(map[string]interface{})["avatar-upload"] = "Upload your avatar"
	p["texts"].(map[string]interface{})["username"] = user.Name
	return p
}
