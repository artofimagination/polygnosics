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
	UserProfilePathKey         = "profile"
	UserProfileEditPathKey     = "profile_edit"
	UserProfileAvatarUploadKey = "avatar_upload"
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
	content[FutureFeature] = 1
	content[UserMapKey] = make(map[string]interface{})
	userContent := content[UserMapKey].(map[string]interface{})
	path := c.UserDBController.ModelFunctions.GetFilePath(user.Assets, businesslogic.UserAvatarKey, businesslogic.DefaultUserAvatarPath)
	userContent[businesslogic.UserAvatarKey] = path
	userContent[businesslogic.UserNameKey] = user.Name
	userContent[businesslogic.UserFullNameKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, businesslogic.UserFullNameKey, "")
	country := c.UserDBController.ModelFunctions.GetField(user.Settings, businesslogic.UserCountryKey, "").(string)
	city := c.UserDBController.ModelFunctions.GetField(user.Settings, businesslogic.UserCityKey, "").(string)
	userContent[businesslogic.UserCountryKey] = country
	userContent[businesslogic.UserCityKey] = city
	userContent[businesslogic.UserLocationKey] = setLocationString(country, city)
	userContent[businesslogic.UserEmailKey] = user.Email
	userContent[businesslogic.UserPhoneKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, businesslogic.UserPhoneKey, "")
	userContent[businesslogic.UserConnectionCountKey] = 20
	userContent[businesslogic.UserHiddenConnectionsKey] = 15
	userContent[businesslogic.UserAboutKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, businesslogic.UserAboutKey, "")
	userContent[businesslogic.UserWebsiteKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, businesslogic.UserWebsiteKey, "#")
	userContent[businesslogic.UserFacebookKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, businesslogic.UserFacebookKey, "#")
	userContent[businesslogic.UserTwitterKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, businesslogic.UserTwitterKey, "#")
	userContent[businesslogic.UserGithubKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, businesslogic.UserGithubKey, "#")
	userContent[businesslogic.UserPrivilegesKey] = c.UserDBController.ModelFunctions.GetField(user.Settings, businesslogic.UserPrivilegesKey, "#")

	userContent[UserProfileAvatarUploadKey] = "Upload your avatar"
	userContent[UserProfilePathKey] = fmt.Sprintf("/user-main/profile?user=%s", user.ID.String())
	userContent[UserProfileEditPathKey] = "/user-main/profile-edit"
	return content
}

func (c *ContentController) StoreUserInfo(r *http.Request) error {
	c.UserData.Name = r.FormValue(businesslogic.UserNameKey)
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, businesslogic.UserNameKey, r.FormValue(businesslogic.UserNameKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, businesslogic.UserFullNameKey, r.FormValue(businesslogic.UserFullNameKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, businesslogic.UserCountryKey, r.FormValue(businesslogic.UserCountryKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, businesslogic.UserCityKey, r.FormValue(businesslogic.UserCityKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, businesslogic.UserPhoneKey, r.FormValue(businesslogic.UserPhoneKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, businesslogic.UserWebsiteKey, r.FormValue(businesslogic.UserWebsiteKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, businesslogic.UserAboutKey, r.FormValue(businesslogic.UserAboutKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, businesslogic.UserFacebookKey, r.FormValue(businesslogic.UserFacebookKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, businesslogic.UserTwitterKey, r.FormValue(businesslogic.UserTwitterKey))
	c.UserDBController.ModelFunctions.SetField(c.UserData.Settings, businesslogic.UserGithubKey, r.FormValue(businesslogic.UserGithubKey))
	if err := c.UserDBController.UpdateUserSettings(c.UserData); err != nil {
		return err
	}
	return nil
}
