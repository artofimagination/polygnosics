package businesslogic

import (
	"fmt"
	"net/http"

	"github.com/artofimagination/golang-docker/docker"
	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
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
	ProjectNameKey       = "name"
	ProjectVisibilityKey = "visibility"
	ProjectServerLogging = "server_logging"
	ProjectClientLogging = "client_logging"
	ProjectState         = "state"
	ProjectContainerID   = "project_container_id"
)

func (c *Context) DeleteProject(project *models.ProjectData) error {
	if err := c.UserDBController.DeleteProject(&project.ID); err != nil {
		return err
	}

	folder := c.ModelFunctions.GetFilePath(project.Assets, models.BaseAssetPath, "")
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

func (c *Context) getProjectState(details *models.Asset) string {
	state := NotRunning
	containerID := c.ModelFunctions.GetField(details, ProjectContainerID, "").(string)
	if err := docker.ContainerExists(containerID); err != nil {
		state = NotRunning
	}
	return state
}

func (c *Context) CheckProject(id *uuid.UUID) (bool, error) {
	// TODO Issue#109: Add status end point to the project servers. Also need to find a way to
	// Have predefined IP addresses for the project docker.
	// (To handle: G107: Potential HTTP request made with variable url)
	// project, err := c.UserDBController.GetProject(id)
	// if err != nil {
	// 	return false, err
	// }
	// containerID := c.UserDBController.ModelFunctions.GetField(project.Details, ProjectContainerID, "")

	// ip, err := docker.GetIPAddress(containerID.(string), "polygnosics_poly_backend")
	// if err != nil {
	// 	return false, err
	// }
	// ip = strings.Split(ip, "/")[0]
	// ipString := fmt.Sprintf("http://%s:10000/", ip)
	// response, err := http.Get(ip)
	// if err != nil {
	// 	return false, err
	// }

	// if response.StatusCode == 200 {
	// 	return true, nil
	// }
	return true, nil
}

func (c *Context) SetProjectDetails(details *models.Asset, productDetails *models.Asset, r *http.Request, containerID string) {
	c.ModelFunctions.SetField(details, ProjectContainerID, containerID)
	c.ModelFunctions.SetField(details, ProjectState, c.getProjectState(details))
	c.ModelFunctions.SetField(details, ProjectVisibilityKey, r.FormValue(ProjectVisibilityKey))
	c.ModelFunctions.SetField(details, ProjectServerLogging, getBooleanString(r.FormValue(ProjectServerLogging)))
	c.ModelFunctions.SetField(details, ProjectClientLogging, getBooleanString(r.FormValue(ProjectClientLogging)))
	categories := c.ModelFunctions.GetField(productDetails, ProductCategoriesKey, "")
	c.ModelFunctions.SetField(details, ProductCategoriesKey, categories.([]interface{}))
}

func (c *Context) EditProjectData(project *models.ProjectData, r *http.Request) error {
	c.ModelFunctions.SetField(project.Details, ProjectNameKey, r.FormValue(ProjectNameKey))
	c.ModelFunctions.SetField(project.Details, ProjectVisibilityKey, r.FormValue(ProjectVisibilityKey))
	c.ModelFunctions.SetField(project.Details, ProjectServerLogging, getBooleanString(r.FormValue(ProjectServerLogging)))
	c.ModelFunctions.SetField(project.Details, ProjectClientLogging, getBooleanString(r.FormValue(ProjectClientLogging)))
	if err := c.UserDBController.UpdateProjectDetails(project); err != nil && err != dbcontrollers.ErrNoProjectDetailsUpdate {
		return err
	}

	if err := c.UserDBController.UpdateProjectAssets(project); err != nil && err != dbcontrollers.ErrNoProjectAssetsUpdate {
		return err
	}
	return nil
}

func (c *Context) UpdateProjectData(project *models.ProjectData, containerID string, r *http.Request) error {
	product, err := c.UserDBController.GetProduct(&project.ProductID)
	if err != nil {
		return err
	}

	c.SetProjectDetails(project.Details, product.Details, r, containerID)

	if err := c.UserDBController.UpdateProjectDetails(project); err != nil && err != dbcontrollers.ErrNoProjectDetailsUpdate {
		return err
	}

	if err := c.UserDBController.UpdateProjectAssets(project); err != nil && err != dbcontrollers.ErrNoProjectAssetsUpdate {
		return err
	}
	return nil
}

func (c *Context) RunProject(userID *uuid.UUID, projectID *uuid.UUID) error {
	project, err := c.UserDBController.GetProject(projectID)
	if err != nil {
		return err
	}

	containerID := c.ModelFunctions.GetField(project.Details, ProjectContainerID, "").(string)
	if err := docker.ContainerExists(containerID); err != nil {
		containerID, err = c.CreateDockerContainer(userID, &project.ProductID)
		if err != nil {
			c.ModelFunctions.SetField(project.Details, ProjectState, NotRunning)
			return err
		}
	}

	if err := docker.StartContainer(containerID, "polygnosics_poly_backend"); err != nil {
		c.ModelFunctions.SetField(project.Details, ProjectState, NotRunning)
		return err
	}
	c.ModelFunctions.SetField(project.Details, ProjectContainerID, containerID)
	c.ModelFunctions.SetField(project.Details, ProjectState, Running)

	if err := c.UserDBController.UpdateProjectDetails(project); err != nil && err != dbcontrollers.ErrNoProjectDetailsUpdate {
		return err
	}
	return nil
}
