package userdb

import (
	"encoding/json"
	"fmt"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

const (
	projectPathAdd           = "/add-project"
	projectPathDelete        = "/delete-project"
	projectPathGet           = "/get-project"
	projectPathUpdateDetails = "/update-project-details"
	projectPathUpdateAssets  = "/update-project-assets"
)

func (c *RESTController) CreateProject(owner *uuid.UUID, productID *uuid.UUID) (*models.ProjectData, error) {
	params := make(map[string]interface{})
	params["owner_id"] = owner.String()
	params["product_id"] = productID.String()
	data, err := c.Post(projectPathAdd, params)
	if err != nil {
		return nil, err
	}

	project := &models.ProjectData{}
	if err := json.Unmarshal(data.([]byte), &project); err != nil {
		return nil, err
	}

	return project, nil
}

func (c *RESTController) DeleteProject(projectID *uuid.UUID) error {
	params := make(map[string]interface{})
	params["id"] = projectID.String()
	_, err := c.Post(projectPathDelete, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *RESTController) GetProject(projectID *uuid.UUID) (*models.ProjectData, error) {
	params := fmt.Sprintf("?id=%s", projectID.String())
	data, err := c.Get(projectPathGet, params)
	if err != nil {
		return nil, err
	}

	projectData := &models.ProjectData{}
	if err := json.Unmarshal(data.([]byte), &projectData); err != nil {
		return nil, err
	}
	return projectData, nil
}

func (c *RESTController) UpdateProjectDetails(projectData *models.ProjectData) error {
	params := make(map[string]interface{})
	dataBytes, err := json.Marshal(projectData)
	if err != nil {
		return err
	}
	params["project-data"] = dataBytes
	_, err = c.Post(projectPathUpdateDetails, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *RESTController) UpdateProjectAssets(projectData *models.ProjectData) error {
	params := make(map[string]interface{})
	dataBytes, err := json.Marshal(projectData)
	if err != nil {
		return err
	}
	params["project-data"] = dataBytes
	_, err = c.Post(projectPathUpdateAssets, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *RESTController) GetProjectsByProductID(productID *uuid.UUID) ([]models.ProjectData, error) {
	params := fmt.Sprintf("?id=%s", productID.String())
	data, err := c.Get(projectPathGet, params)
	if err != nil {
		return nil, err
	}

	projectData := make([]models.ProjectData, 0)
	if err := json.Unmarshal(data.([]byte), &projectData); err != nil {
		return nil, err
	}
	return projectData, nil
}
