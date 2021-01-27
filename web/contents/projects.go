package contents

import (
	"fmt"
	"net/http"

	"polygnosics/app/businesslogic"
	"polygnosics/app/businesslogic/project"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

// Details and assets field keys
const (
	ProjectAvatar        = "project_avatar"
	ProjectPath          = "project_path"
	ProjectName          = "name"
	ProjectVisibility    = "visibility"
	ProjectServerLogging = "server_logging"
	ProjectClientLogging = "client_logging"
	NewProject           = "new_project"
	RunProject           = "run_project"
	ProjectState         = "project_state"
	ProjectStateColor    = "project_state_color"
	ProjectContainerID   = "project_container_id"
)

// Visibility values of a project
const (
	Public    = "Public"
	Protected = "Protected"
	Private   = "Private"
)

// GetProjectStateColorString returns UI color of the project state based on the state value.
func GetProjectStateColorString(state string) string {
	switch state {
	case project.NotRunning:
		return "#f5cf0a" // orange
	case project.Running:
		return "#00ff00" // green
	case project.Stopped:
		return "#ff0000" // red
	default:
		return "#e0dfd6" // lightgray
	}
}

// ValidateVisibility validates the visibility string
func ValidateVisibility(value string) error {
	if value != Public && value != Protected && value != Private {
		return fmt.Errorf("Invalid visibility: %s", value)
	}
	return nil
}

// generateProjectContent fills a string nested map with all project details and assets info
func (c *ContentController) generateProjectContent(projectData *models.ProjectData) map[string]interface{} {
	content := make(map[string]interface{})
	content[ProjectAvatar] = c.UserDBController.ModelFunctions.GetFilePath(projectData.Assets, ProjectAvatar, businesslogic.DefaultProjectAvatarPath)
	content[ProjectName] = c.UserDBController.ModelFunctions.GetField(projectData.Details, ProjectName, "")
	content[ProjectVisibility] = c.UserDBController.ModelFunctions.GetField(projectData.Details, ProjectVisibility, "")
	content[ProjectContainerID] = c.UserDBController.ModelFunctions.GetField(projectData.Details, ProjectContainerID, "")
	content[ProjectState] = c.UserDBController.ModelFunctions.GetField(projectData.Details, ProjectState, "")

	content[ProjectPath] = fmt.Sprintf("/user-main/my-projects/details?project=%s", projectData.ID.String())
	content[ProjectStateColor] = GetProjectStateColorString(c.UserDBController.ModelFunctions.GetField(projectData.Details, ProjectState, ""))
	content[RunProject] = fmt.Sprintf("/user-main/my-projects/run?project=%s", projectData.ID.String())
	return content
}

func (c *ContentController) SetProjectDetails(details *models.Asset, r *http.Request, containerID string) {
	c.UserDBController.ModelFunctions.SetField(details, ProjectContainerID, containerID)
	c.UserDBController.ModelFunctions.SetField(details, ProjectState, project.NotRunning)
	c.UserDBController.ModelFunctions.SetField(details, ProjectVisibility, r.FormValue("visibility"))
	c.UserDBController.ModelFunctions.SetField(details, ProjectServerLogging, getBooleanString(r.FormValue("serverLogging")))
	c.UserDBController.ModelFunctions.SetField(details, ProjectClientLogging, getBooleanString(r.FormValue("clientLogging")))
}

// GetProjectContent returns the selected project details and assets info.
func (c *ContentController) GetProjectContent(projectID *uuid.UUID) (map[string]interface{}, error) {
	project, err := c.UserDBController.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	return c.generateProjectContent(project), nil
}

// GetUserProjectContent gathers the contents of all projects belonging to the specified user.
func (c *ContentController) GetUserProjectContent(userID *uuid.UUID) (map[string]interface{}, error) {
	projects, err := c.UserDBController.GetProjectsByUserID(userID)
	if err != nil {
		return nil, err
	}

	p := make(map[string]interface{})

	projectContent := make([]map[string]interface{}, len(projects))
	for i, project := range projects {
		projectContent[i] = c.generateProjectContent(project.ProjectData)
	}
	p["project"] = projectContent

	return p, nil
}
