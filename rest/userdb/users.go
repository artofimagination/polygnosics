package userdb

import (
	"encoding/json"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/polygnosics/rest"
	"github.com/google/uuid"
)

const (
	UserPathLogin          = "/login"
	userPathAdd            = "/add-user"
	UserPathDetectRootUser = "/detect-root-user"
	userPathGetUserByID    = "/get-user-by-id"
	userPathGetUserByEmail = "/get-user-by-email"
	userPathUpdateSettings = "/update-user-settings"
	userPathUpdateAssets   = "/update-user-assets"
	userPathDelete         = "/delete-user"
)

func (c *RESTController) CreateUser(
	name string,
	email string,
	password []byte) (*models.UserData, error) {

	params := make(map[string]interface{})
	params["username"] = name
	params["email"] = email
	params["password"] = string(password)
	data, err := rest.Post(rest.UserDBAddress, userPathAdd, params)
	if err != nil {
		return nil, err
	}

	userData := &models.UserData{}
	if err := json.Unmarshal(data.([]byte), &userData); err != nil {
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
	data, err := rest.Get(rest.UserDBAddress, userPathGetUserByEmail, params)
	if err != nil {
		return nil, err
	}

	userData := &models.UserData{}
	if err := json.Unmarshal(data.([]byte), &userData); err != nil {
		return nil, err
	}
	return userData, nil
}

func (c *RESTController) UpdateUserSettings(userData *models.UserData) error {
	params := make(map[string]interface{})
	userDataBytes, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	params["user-data"] = userDataBytes
	_, err = rest.Post(rest.UserDBAddress, userPathUpdateSettings, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *RESTController) UpdateUserAssets(userData *models.UserData) error {
	params := make(map[string]interface{})
	userDataBytes, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	params["user-data"] = userDataBytes
	_, err = rest.Post(rest.UserDBAddress, userPathUpdateAssets, params)
	if err != nil {
		return err
	}
	return nil
}
