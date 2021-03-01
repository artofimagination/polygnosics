package businesslogic

import (
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserAvatarKey            = "avatar"
	UserNameKey              = "username"
	UserFullNameKey          = "full_name"
	UserLocationKey          = "location"
	UserCountryKey           = "country"
	UserCityKey              = "city"
	UserEmailKey             = "email"
	UserWebsiteKey           = "website"
	UserPhoneKey             = "phone"
	UserConnectionCountKey   = "connection_count"
	UserHiddenConnectionsKey = "hidden_connections"
	UserAboutKey             = "about"
	UserFacebookKey          = "facebook_link"
	UserTwitterKey           = "twitter_link"
	UserGithubKey            = "github_link"
	UserPrivilegesKey        = "privileges"
)

const (
	UserGroupDeveloper = "developer"
	UserGroupAdmin     = "admin"
	UserGroupRoot      = "root"
	UserGroupClient    = "client"
	UserGroupVisitor   = "visitor"
)

const (
	PrivilegeShowMainDashboard = "main_dashboard"
	PrivilegeShowMisuseMetrics = "misuse_metrics"
	PrivilegeShowProductStats  = "product_stats"
	PrivilegeShowProjectStats  = "project_stats"
)

const (
	PrivilegeActionDeleteUser = "delete_user"
)

func (c *Context) setGroupPrivileges(userData *models.UserData, group string) error {
	privilegeMap := make(map[string]interface{})
	switch group {
	case UserGroupRoot:
		privilegeMap[PrivilegeShowMainDashboard] = 1
		privilegeMap[PrivilegeShowMisuseMetrics] = 1
		privilegeMap[PrivilegeShowProductStats] = 1
		privilegeMap[PrivilegeShowProjectStats] = 1

		privilegeMap[PrivilegeActionDeleteUser] = 1
		c.UserDBController.ModelFunctions.SetField(userData.Settings, UserGroupRoot, 1)
	case UserGroupAdmin:
		privilegeMap[PrivilegeShowMainDashboard] = 1
		privilegeMap[PrivilegeShowMisuseMetrics] = 1
		privilegeMap[PrivilegeShowProductStats] = 1
		privilegeMap[PrivilegeShowProjectStats] = 1

		privilegeMap[PrivilegeActionDeleteUser] = 1
		c.UserDBController.ModelFunctions.SetField(userData.Settings, UserGroupAdmin, 1)
	case UserGroupDeveloper:
		privilegeMap[PrivilegeShowMainDashboard] = 0
		privilegeMap[PrivilegeShowMisuseMetrics] = 0
		privilegeMap[PrivilegeShowProductStats] = 1
		privilegeMap[PrivilegeShowProjectStats] = 1

		privilegeMap[PrivilegeActionDeleteUser] = 0
		c.UserDBController.ModelFunctions.SetField(userData.Settings, UserGroupDeveloper, 1)
	case UserGroupClient:
		privilegeMap[PrivilegeShowMainDashboard] = 0
		privilegeMap[PrivilegeShowMisuseMetrics] = 0
		privilegeMap[PrivilegeShowProductStats] = 0
		privilegeMap[PrivilegeShowProjectStats] = 1

		privilegeMap[PrivilegeActionDeleteUser] = 0
		c.UserDBController.ModelFunctions.SetField(userData.Settings, UserGroupClient, 1)
	case UserGroupVisitor:
		privilegeMap[PrivilegeShowMainDashboard] = 0
		privilegeMap[PrivilegeShowMisuseMetrics] = 0
		privilegeMap[PrivilegeShowProductStats] = 0
		privilegeMap[PrivilegeShowProjectStats] = 0

		privilegeMap[PrivilegeActionDeleteUser] = 0
		c.UserDBController.ModelFunctions.SetField(userData.Settings, UserGroupVisitor, 1)
	default:
		return fmt.Errorf("Invalid user group: %s", group)
	}
	c.UserDBController.ModelFunctions.SetField(userData.Settings, UserPrivilegesKey, privilegeMap)
	return nil
}

func encryptPassword(password []byte) ([]byte, error) {
	var hashedPassword []byte
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 16)
	if err != nil {
		return hashedPassword, err
	}
	return hashedPassword, nil
}

func (c *Context) AddUser(uName string, email string, pwd []byte, group string) error {
	userData, err := c.UserDBController.CreateUser(uName, email, pwd, GeneratePath, encryptPassword)
	if err != nil {
		return err
	}

	if err := c.setGroupPrivileges(userData, group); err != nil {
		return err
	}

	if err := c.UserDBController.UpdateUserSettings(userData); err != nil {
		return err
	}
	return nil
}

func (c *Context) GetUserByEmail(email string) (*models.UserData, error) {
	searchedUser, err := c.UserDBController.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return searchedUser, err
}
