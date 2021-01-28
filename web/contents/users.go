package contents

import (
	"polygnosics/app/businesslogic"
)

// Details and assets field keys
const (
	UserAvatar     = "user_avatar"
	UserBackground = "user_background"
)

// GetUserContent fills a string nested map with all user details and assets info
func (c *ContentController) GetUserContent() map[string]interface{} {
	p := make(map[string]interface{})
	p["assets"] = make(map[string]interface{})
	path := c.UserDBController.ModelFunctions.GetFilePath(c.UserData.Assets, UserAvatar, businesslogic.DefaultUserAvatarPath)
	p["assets"].(map[string]interface{})[UserAvatar] = path
	p["texts"] = make(map[string]interface{})
	p["texts"].(map[string]interface{})["avatar-upload"] = "Upload your avatar"
	p["texts"].(map[string]interface{})["username"] = c.UserData.Name
	return p
}
