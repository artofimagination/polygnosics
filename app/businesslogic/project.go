package businesslogic

import (
	"fmt"
	"net/http"

	"github.com/artofimagination/golang-docker/docker"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	NotRunning = "Not running"
	Running    = "Running"
	Stopped    = "Stopped"
	Paused     = "Paused"
)

// Details and assets field keys
const (
	ProjectAvatar        = "avatar"
	ProjectName          = "name"
	ProjectVisibility    = "visibility"
	ProjectServerLogging = "server_logging"
	ProjectClientLogging = "client_logging"
	ProjectState         = "state"
	ProjectContainerID   = "project_container_id"
)

func (c *Context) DeleteProject(project *models.ProjectData) error {
	if err := c.UserDBController.DeleteProject(&project.ID); err != nil {
		return err
	}

	folder := c.UserDBController.ModelFunctions.GetFilePath(project.Assets, models.BaseAssetPath, "")
	if err := removeContents(folder); err != nil {
		return fmt.Errorf("Failed to delete project. %s", errors.WithStack(err))
	}
	return nil
}

func (c *Context) CreateDockerContainer(userID *uuid.UUID, produtID *uuid.UUID) (string, error) {
	containerName := fmt.Sprintf("%s/%s", userID.String(), produtID.String())
	containerID, err := docker.CreateNewContainer(containerName, "0.0.0.0", "10000")
	if err != nil {
		return "", err
	}
	return containerID, nil
}

func (c *Context) SetProjectDetails(details *models.Asset, productDetails *models.Asset, r *http.Request, containerID string) {
	c.UserDBController.ModelFunctions.SetField(details, ProjectContainerID, containerID)
	c.UserDBController.ModelFunctions.SetField(details, ProjectState, NotRunning)
	c.UserDBController.ModelFunctions.SetField(details, ProjectVisibility, r.FormValue(ProjectVisibility))
	c.UserDBController.ModelFunctions.SetField(details, ProjectServerLogging, getBooleanString(r.FormValue(ProjectServerLogging)))
	c.UserDBController.ModelFunctions.SetField(details, ProjectClientLogging, getBooleanString(r.FormValue(ProjectClientLogging)))
	categories := c.UserDBController.ModelFunctions.GetField(productDetails, ProductCategoriesKey, "")
	c.UserDBController.ModelFunctions.SetField(details, ProductCategoriesKey, categories.([]interface{}))
}

func (c *Context) UpdateProjectData(project *models.ProjectData, containerID string, r *http.Request) error {
	product, err := c.UserDBController.GetProduct(&project.ProductID)
	if err != nil {
		return err
	}

	c.SetProjectDetails(project.Details, product.Details, r, containerID)

	if err := c.UserDBController.UpdateProjectDetails(project); err != nil {
		return err
	}

	if err := c.UserDBController.UpdateProjectAssets(project); err != nil {
		return err
	}
	return nil
}
