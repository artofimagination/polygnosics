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
	ProjectStateBadge    = "state_badge"
	ProjectContainerID   = "project_container_id"
)

// Visibility values of a project
const (
	Public    = "Public"
	Protected = "Protected"
	Private   = "Private"
)

type StateContent struct {
	text  string
	badge string
}

// GetProjectStateContent returns UI color of the project state based on the state value.
func GetProjectStateContent(stateString string) *StateContent {
	state := &StateContent{
		text: stateString,
	}
	switch stateString {
	case project.NotRunning:
		state.badge = "badge-warning" // orange
	case project.Paused:
		state.badge = "badge-primary" // lightblue
	case project.Running:
		state.badge = "badge-success" // green
	case project.Stopped:
		state.badge = "badge-danger" // red
	default:
		state.badge = "badge-secondary" // lightgray
	}
	return state
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
	content[ProjectStateBadge] = GetProjectStateContent(c.UserDBController.ModelFunctions.GetField(projectData.Details, ProjectState, "").(string)).badge
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
func (c *ContentController) GetUserProjectContent(userID *uuid.UUID, limit int) ([]map[string]interface{}, error) {
	projects, err := c.UserDBController.GetProjectsByUserID(userID)
	if err != nil {
		return nil, err
	}

	if limit > len(projects) {
		limit = len(projects)
	}

	projectContent := make([]map[string]interface{}, len(projects))
	if limit != -1 {
		projectContent = make([]map[string]interface{}, limit)
	}

	for i, project := range projects {
		if limit == 0 {
			break
		}
		limit--
		projectContent[i] = c.generateProjectContent(project.ProjectData)
	}

	return projectContent, nil
}
