package contents

import (
	"fmt"

	"polygnosics/app/businesslogic"

	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

// Details and assets field keys
const (
	ProjectMapKey               = "project"
	NewProject                  = "new_project"
	RunProject                  = "run_project"
	ShowProject                 = "show_project"
	ProjectDetailsPageKey       = "detail_path"
	ProjectStateBadge           = "state_badge"
	ProjectDeleteIDKey          = "project_to_delete"
	ProjectDeleteTextKey        = "delete_text"
	ProjectDeleteSuccessTextKey = "delete_success_text"
	ProjectDeleteURLKey         = "delete_url"
	ProjectEditPageKey          = "edit_path"
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
	case businesslogic.NotRunning:
		state.badge = "badge-warning" // orange
	case businesslogic.Paused:
		state.badge = "badge-primary" // lightblue
	case businesslogic.Running:
		state.badge = "badge-success" // green
	case businesslogic.Stopped:
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
	content[businesslogic.ProjectAvatar] = c.UserDBController.ModelFunctions.GetFilePath(projectData.Assets, businesslogic.ProjectAvatar, businesslogic.DefaultProjectAvatarPath)
	content[businesslogic.ProjectNameKey] = c.UserDBController.ModelFunctions.GetField(projectData.Details, businesslogic.ProjectNameKey, "")
	content[businesslogic.ProjectVisibilityKey] = c.UserDBController.ModelFunctions.GetField(projectData.Details, businesslogic.ProjectVisibilityKey, "")
	content[businesslogic.ProjectContainerID] = c.UserDBController.ModelFunctions.GetField(projectData.Details, businesslogic.ProjectContainerID, "")
	content[businesslogic.ProjectState] = c.UserDBController.ModelFunctions.GetField(projectData.Details, businesslogic.ProjectState, "")
	content[businesslogic.ProductCategoriesKey] = c.UserDBController.ModelFunctions.GetField(projectData.Details, businesslogic.ProductCategoriesKey, "")
	content[businesslogic.ProjectServerLogging] = convertToCheckboxValue(c.UserDBController.ModelFunctions.GetField(projectData.Details, businesslogic.ProjectServerLogging, "").(string))
	content[businesslogic.ProjectClientLogging] = convertToCheckboxValue(c.UserDBController.ModelFunctions.GetField(projectData.Details, businesslogic.ProjectClientLogging, "").(string))

	content[ProjectDetailsPageKey] = fmt.Sprintf("/user-main/my-projects/details?item-id=%s", projectData.ID.String())
	content[ProjectStateBadge] = GetProjectStateContent(c.UserDBController.ModelFunctions.GetField(projectData.Details, businesslogic.ProjectState, "").(string)).badge
	content[ProjectEditPageKey] = fmt.Sprintf("/user-main/my-projects/edit?item-id=%s", projectData.ID.String())
	content[RunProject] = fmt.Sprintf("/user-main/my-projects/run?item-id=%s", projectData.ID.String())
	content[ShowProject] = fmt.Sprintf("/user-main/my-projects/show?item-id=%s", projectData.ID.String())
	content[ProjectDeleteIDKey] = projectData.ID.String()
	content[ProjectDeleteTextKey] = "The project will not be accessible anymore"
	content[ProjectDeleteSuccessTextKey] = "Your project has been deleted"
	content[ProjectDeleteURLKey] = "/user-main/my-projects/delete"
	return content
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

// GetProjectsByCategory organizes project contents by categories. This is just a placeholder solution until proper search is introduced.
func (c *ContentController) GetProjectsByCategory(userID *uuid.UUID) (map[string]interface{}, error) {
	projects, err := c.UserDBController.GetProjectsByUserID(userID)
	if err != nil {
		return nil, err
	}

	categorizedProjects := make(map[string]interface{})
	categories := businesslogic.CreateCategoriesMap()
	for k, v := range categories {
		if categorizedProjects[k] == nil {
			categorizedProjects[k] = make(map[string]interface{})
		}
		categorizedProjects[k].(map[string]interface{})["name"] = v
		for _, project := range projects {
			for _, categoryKey := range c.UserDBController.ModelFunctions.GetField(project.ProjectData.Details, businesslogic.ProductCategoriesKey, "").([]interface{}) {
				if categoryKey.(string) == k {
					_, ok := categorizedProjects[k].(map[string]interface{})["projects"]
					if !ok {
						categorizedProjects[k].(map[string]interface{})["projects"] = make([]map[string]interface{}, 0)
					}
					categorizedProjects[k].(map[string]interface{})["projects"] = append(categorizedProjects[k].(map[string]interface{})["projects"].([]map[string]interface{}), c.generateProjectContent(project.ProjectData))
					break
				}
			}
		}
	}
	return categorizedProjects, nil
}

// GetRecentProjectsContent gathers the content of the latest 4 projects
func (c *ContentController) GetRecentProjectsContent(userID *uuid.UUID) ([]map[string]interface{}, error) {
	projects, err := c.UserDBController.GetProjectsByUserID(userID)
	if err != nil {
		return nil, err
	}

	limit := 4
	if len(projects) < limit {
		limit = len(projects)
	}
	projectContent := make([]map[string]interface{}, limit)
	for i, project := range projects {
		if limit == 0 {
			break
		}
		limit--
		projectContent[i] = c.generateProjectContent(project.ProjectData)
		projectContent[i][ProductOwnerNameKey] = c.UserData.Name
		projectContent[i][ProductOwnerPageNameKey] = fmt.Sprintf("/user-main/profile?user=%s", c.UserData.ID)
		projectContent[i][ProjectDetailsPageKey] = fmt.Sprintf("/user-main/project?item-id=%s", project.ProjectData.ID)
	}

	return projectContent, nil
}

func (c *ContentController) GetRecommendedProjectsContent(userID *uuid.UUID) ([]map[string]interface{}, error) {
	projects, err := c.UserDBController.GetProjectsByUserID(userID)
	if err != nil {
		return nil, err
	}

	limit := 4
	if len(projects) < limit {
		limit = len(projects)
	}
	projectContent := make([]map[string]interface{}, limit)
	for i, project := range projects {
		if limit == 0 {
			break
		}
		limit--
		projectContent[i] = c.generateProjectContent(project.ProjectData)
		projectContent[i][ProductOwnerNameKey] = c.UserData.Name
		projectContent[i][ProductOwnerPageNameKey] = fmt.Sprintf("/user-main/profile?user=%s", c.UserData.ID)
		projectContent[i][ProjectDetailsPageKey] = fmt.Sprintf("/user-main/project?item-id=%s", project.ProjectData.ID)
	}

	return projectContent, nil
}
