package userdb

import (
	"encoding/json"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/polygnosics/rest"
	"github.com/google/uuid"
)

const (
	userPathAdd            = "/add-user"
	UserPathGetUserByEmail = "/get-user-by-email"
	userPathUpdateSettings = "/update-user-settings"
	userPathUpdateAssets   = "/update-user-assets"
	userPathDelete         = "/delete-user"
)

func (c *RESTController) CreateUser(
	name string,
	email string,
	password string) (*models.UserData, error) {

	params := make(map[string]interface{})
	params["username"] = name
	params["email"] = email
	params["password"] = password
	data, err := rest.Post(rest.UserDBAddress, userPathAdd, params)
	if err != nil {
		return nil, err
	}

	userData := &models.UserData{}
	if err := json.Unmarshal([]byte(data.(string)), &userData); err != nil {
		return nil, err
	}

	return userData, nil
}

func (c *RESTController) DeleteUser(ID *uuid.UUID, nominatedOwners map[uuid.UUID]uuid.UUID) error {
	params := make(map[string]interface{})
	params["id"] = ID.String()
	_, err := rest.Post(rest.UserDBAddress, userPathDelete, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *RESTController) GetUserByEmail(email string) (*models.UserData, error) {
	params := fmt.Sprintf("?email=%s", email)
	data, err := rest.Get(rest.UserDBAddress, UserPathGetUserByEmail, params)
	if err != nil {
		return nil, err
	}

	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	userData := &models.UserData{}
	if err := json.Unmarshal(bytesData, &userData); err != nil {
		return nil, err
	}
	return userData, nil
}

func (c *RESTController) UpdateUserSettings(userData *models.UserData) error {
	params := make(map[string]interface{})
	params["user-id"] = userData.ID
	params["user-data"] = userData.Settings
	_, err := rest.Post(rest.UserDBAddress, userPathUpdateSettings, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *RESTController) UpdateUserAssets(userData *models.UserData) error {
	params := make(map[string]interface{})
	params["user-id"] = userData.ID
	params["user-data"] = userData.Assets
	_, err := rest.Post(rest.UserDBAddress, userPathUpdateAssets, params)
	if err != nil {
		return err
	}
	return nil
}
