package businesslogic

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
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
	UserGroupKey             = "group"
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

var errIncorrectEmailOrPass = errors.New("Incorrect email or password")

func (c *Context) setGroupPrivileges(userData *models.UserData, group string) error {
	privilegeMap := make(map[string]interface{})
	switch group {
	case UserGroupRoot:
		privilegeMap[PrivilegeShowMainDashboard] = 1
		privilegeMap[PrivilegeShowMisuseMetrics] = 1
		privilegeMap[PrivilegeShowProductStats] = 1
		privilegeMap[PrivilegeShowProjectStats] = 1

		privilegeMap[PrivilegeActionDeleteUser] = 1
		c.ModelFunctions.SetField(userData.Settings, UserGroupKey, UserGroupRoot)
	case UserGroupAdmin:
		privilegeMap[PrivilegeShowMainDashboard] = 1
		privilegeMap[PrivilegeShowMisuseMetrics] = 1
		privilegeMap[PrivilegeShowProductStats] = 1
		privilegeMap[PrivilegeShowProjectStats] = 1

		privilegeMap[PrivilegeActionDeleteUser] = 1
		c.ModelFunctions.SetField(userData.Settings, UserGroupKey, UserGroupAdmin)
	case UserGroupDeveloper:
		privilegeMap[PrivilegeShowMainDashboard] = 0
		privilegeMap[PrivilegeShowMisuseMetrics] = 0
		privilegeMap[PrivilegeShowProductStats] = 1
		privilegeMap[PrivilegeShowProjectStats] = 1

		privilegeMap[PrivilegeActionDeleteUser] = 0
		c.ModelFunctions.SetField(userData.Settings, UserGroupKey, UserGroupDeveloper)
	case UserGroupClient:
		privilegeMap[PrivilegeShowMainDashboard] = 0
		privilegeMap[PrivilegeShowMisuseMetrics] = 0
		privilegeMap[PrivilegeShowProductStats] = 0
		privilegeMap[PrivilegeShowProjectStats] = 1

		privilegeMap[PrivilegeActionDeleteUser] = 0
		c.ModelFunctions.SetField(userData.Settings, UserGroupKey, UserGroupClient)
	case UserGroupVisitor:
		privilegeMap[PrivilegeShowMainDashboard] = 0
		privilegeMap[PrivilegeShowMisuseMetrics] = 0
		privilegeMap[PrivilegeShowProductStats] = 0
		privilegeMap[PrivilegeShowProjectStats] = 0

		privilegeMap[PrivilegeActionDeleteUser] = 0
		c.ModelFunctions.SetField(userData.Settings, UserGroupKey, UserGroupVisitor)
	default:
		return fmt.Errorf("Invalid user group: %s", group)
	}
	c.ModelFunctions.SetField(userData.Settings, UserPrivilegesKey, privilegeMap)
	return nil
}

func encryptPassword(password []byte) (string, error) {
	var hashedPassword []byte
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 16)
	if err != nil {
		return "", err
	}
	hashedPasswordBase64 := base64.URLEncoding.EncodeToString(hashedPassword)
	return hashedPasswordBase64, nil
}

func authenticate(password []byte, storedPassword []byte) error {
	if err := bcrypt.CompareHashAndPassword(storedPassword, password); err != nil {
		return errIncorrectEmailOrPass
	}
	return nil
}

func (c *Context) Login(email string, password []byte) (*models.UserData, error) {
	// TODO Issue#45: replace this with elastic search
	user, err := c.UserDBController.GetUserByEmail(email)
	if err != nil && err.Error() == dbcontrollers.ErrUserNotFound.Error() {
		return nil, errIncorrectEmailOrPass
	} else if err != nil {
		return nil, err
	}

	// Get and check user and password
	err = authenticate(password, user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Context) AddUser(uName string, email string, password []byte, group string) error {
	hashedPassword, err := encryptPassword(password)
	if err != nil {
		return err
	}

	userData, err := c.UserDBController.CreateUser(uName, email, hashedPassword)
	if err != nil {
		return err
	}

	if err := GeneratePath(userData.Assets); err != nil {
		return err
	}

	if err := c.setGroupPrivileges(userData, group); err != nil {
		return err
	}

	if err := c.UserDBController.UpdateUserSettings(userData); err != nil {
		return err
	}

	if err := c.UserDBController.UpdateUserAssets(userData); err != nil {
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

func (c *Context) DetectRootUser() (bool, error) {
	// TODO Issue#107: Replace this with proper way of detecting if root has already been created.
	_, err := c.UserDBController.GetUserByEmail("root@test.com")
	if err != nil && err.Error() == dbcontrollers.ErrUserNotFound.Error() {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
