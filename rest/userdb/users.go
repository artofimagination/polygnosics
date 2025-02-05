package userdb

import (
	"encoding/json"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	dbrest "github.com/artofimagination/mysql-user-db-go-interface/restcontrollers"
	"github.com/google/uuid"
)

func (c *RESTController) CreateUser(
	name string,
	email string,
	password string) (*models.UserData, error) {

	params := make(map[string]interface{})
	params["username"] = name
	params["email"] = email
	params["password"] = password
	data, err := c.Post(dbrest.UserPathAdd, params)
	if err != nil {
		return nil, err
	}

	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	userData := &models.UserData{}
	if err := json.Unmarshal(bytesData, userData); err != nil {
		return nil, err
	}
	return userData, nil
}

func (c *RESTController) DeleteUser(ID *uuid.UUID, nominatedOwners map[uuid.UUID]uuid.UUID) error {
	params := make(map[string]interface{})
	params["id"] = ID.String()
	_, err := c.Post(dbrest.UserPathDeleteByID, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *RESTController) GetUserByEmail(email string) (*models.UserData, error) {
	params := fmt.Sprintf("?email=%s", email)
	data, err := c.Get(dbrest.UserPathGetByEmail, params)
	if err != nil {
		return nil, err
	}
	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	userData := &models.UserData{}
	if err := json.Unmarshal(bytesData, userData); err != nil {
		return nil, err
	}

	return userData, nil
}

func (c *RESTController) UpdateUserSettings(userData *models.UserData) error {
	params := make(map[string]interface{})
	params["user-id"] = userData.ID
	params["user-data"] = userData.Settings
	_, err := c.Post(dbrest.UserPathUpdateSettings, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *RESTController) UpdateUserAssets(userData *models.UserData) error {
	params := make(map[string]interface{})
	params["user-id"] = userData.ID
	params["user-data"] = userData.Assets
	_, err := c.Post(dbrest.UserPathUpdateAssets, params)
	if err != nil {
		return err
	}
	return nil
}
