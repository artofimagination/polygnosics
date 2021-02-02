package contents

import (
	"fmt"
	"net/http"
	"polygnosics/app/businesslogic"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
)

// Details and assets field keys
const (
	UserMapKey                 = "user"
	UserAvatarKey              = "avatar"
	UserProfilePathKey         = "profile"
	UserProfileEditPathKey     = "profile_edit"
	UserProfileAvatarUploadKey = "avatar_upload"
	UserNameKey                = "username"
	UserFullNameKey            = "full_name"
	UserLocationKey            = "location"
	UserCountryKey             = "country"
	UserCityKey                = "city"
	UserEmailKey               = "email"
	UserWebsiteKey             = "website"
	UserPhoneKey               = "phone"
	UserConnectionCountKey     = "connection_count"
	UserHiddenConnectionsKey   = "hidden_connections"
	UserAboutKey               = "about"
	UserFacebookKey            = "facebook_link"
	UserTwitterKey             = "twitter_link"
	UserGithubKey              = "github_link"
)

func setLocationString(country string, city string) string {
	if country == "" && city == "" {
		return "Not specified"
	} else if country != "" && city == "" {
		return country
	} else if country == "" && city != "" {
		return city
	}
	return fmt.Sprintf("%s, %s", city, country)
}

// GetUserContent fills a string nested map with all user details and assets info
func (c *ContentController) GetUserContent(user *models.UserData) map[string]interface{} {
	content := make(map[string]interface{})
	content[UserMapKey] = make(map[string]interface{})
	userContent := content[UserMapKey].(map[string]interface{})
	path := c.UserDBController.ModelFunctions.GetFilePath(user.Assets, UserAvatarKey, businesslogic.DefaultUserAvatarPath)
	userContent[UserAvatarKey] = path
	userContent[UserProfileAvatarUploadKey] = "Upload your avatar"
	userContent[UserNameKey] = user.Name
	userContent[UserFullNameKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, UserFullNameKey, "")

	country := c.UserDBController.ModelFunctions.GetField(user.Settings, UserCountryKey, "")
	city := c.UserDBController.ModelFunctions.GetField(user.Settings, UserCityKey, "")
	userContent[UserCountryKey] = country
	userContent[UserCityKey] = city
	userContent[UserLocationKey] = setLocationString(country, city)
	userContent[UserEmailKey] = user.Email
	userContent[UserPhoneKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, UserPhoneKey, "")
	userContent[UserConnectionCountKey] = 20
	userContent[UserHiddenConnectionsKey] = 15
	userContent[UserAboutKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, UserAboutKey, "")
	userContent[UserProfilePathKey] = fmt.Sprintf("/user-main/profile?user=%s", user.ID.String())
	userContent[UserProfileEditPathKey] = "/user-main/profile-edit"
	userContent[UserWebsiteKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, UserWebsiteKey, "#")
	userContent[UserFacebookKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, UserFacebookKey, "#")
	userContent[UserTwitterKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, UserTwitterKey, "#")
	userContent[UserGithubKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, UserGithubKey, "#")
	return content
}

func (c *ContentController) StoreUserInfo(r *http.Request) error {
	c.UserData.Name = r.FormValue(UserNameKey)
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, UserNameKey, r.FormValue(UserNameKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, UserFullNameKey, r.FormValue(UserFullNameKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, UserCountryKey, r.FormValue(UserCountryKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, UserCityKey, r.FormValue(UserCityKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, UserPhoneKey, r.FormValue(UserPhoneKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, UserWebsiteKey, r.FormValue(UserWebsiteKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, UserAboutKey, r.FormValue(UserAboutKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, UserFacebookKey, r.FormValue(UserFacebookKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, UserTwitterKey, r.FormValue(UserTwitterKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, UserGithubKey, r.FormValue(UserGithubKey))
	if err := c.UserDBController.UpdateUserSettings(c.UserData); err != nil {
		return err
	}
	return nil
}
