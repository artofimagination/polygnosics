package restcontrollers

import (
	"fmt"
	"net/http"
	"polygnosics/app/businesslogic"
	"polygnosics/web/contents"

	"github.com/artofimagination/golang-docker/docker"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *RESTController) MyProjects(w http.ResponseWriter, r *http.Request) {
	content, err := c.ContentController.BuildMyProjectsContent()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get project content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	c.RenderTemplate(w, "my-projects", content)
}

func (c *RESTController) ProjectDetails(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseItemID(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse project id. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	content, err := c.ContentController.BuildProjectDetailsContent(projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get project content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	c.RenderTemplate(w, "project-details", content)
}

func (c *RESTController) CreateProject(w http.ResponseWriter, r *http.Request) {
	productID, err := parseItemID(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse product id. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	if r.Method == GET {
		content, err := c.ContentController.BuildProjectWizardContent(productID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get project wizard content content. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}
		c.RenderTemplate(w, ProjectWizard, content)
	} else {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse avatar. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := contents.ValidateVisibility(r.FormValue("visibility")); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse visibility. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		projectData, err := c.UserDBController.CreateProject(
			r.FormValue("projectName"),
			r.FormValue("visibility"),
			&c.ContentController.UserData.ID,
			productID,
			businesslogic.GeneratePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create project. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		err = c.BackendContext.UploadFile(projectData.Assets, businesslogic.ProjectAvatar, businesslogic.DefaultProjectAvatarPath, r)
		if err != nil {
			if errDelete := c.UserDBController.DeleteProject(&projectData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete project. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload avatar. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		containerID, err := c.BackendContext.CreateDockerContainer(&c.ContentController.UserData.ID, &projectData.ProductID)
		if err != nil {
			if errDelete := c.BackendContext.DeleteProject(projectData); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete project. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to create project container. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		c.BackendContext.UpdateProjectData(projectData, containerID, r)
		if err != nil {
			if errDelete := c.BackendContext.DeleteProject(projectData); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete project. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to update project data. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		c.MyProjects(w, r)
	}
}

func (c *RESTController) RunProject(w http.ResponseWriter, r *http.Request) {
	name := UserMain
	p := c.ContentController.GetUserContent(c.ContentController.UserData)
	if err := r.ParseForm(); err != nil {
		p["message"] = ErrFailedToParseForm
		c.RenderTemplate(w, name, p)
		return
	}

	projectID, err := uuid.Parse(r.FormValue("project"))
	if err != nil {
		p["message"] = "Failed to parse project id"
		c.RenderTemplate(w, name, p)
		return
	}

	pProject, err := c.ContentController.GetProjectContent(&projectID)
	if err != nil {
		p["message"] = "Failed to get project details"
		c.RenderTemplate(w, name, p)
		return
	}

	for k, v := range pProject {
		p[k] = v
	}

	if err := docker.StartContainer(p[businesslogic.ProjectContainerID].(string)); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start project container. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	c.RenderTemplate(w, "run", p)
}

func (c *RESTController) DeleteProject(w http.ResponseWriter, r *http.Request) {
	if r.Method == POST {
		projectID, err := parseItemID(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse project id. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		project, err := c.UserDBController.GetProject(projectID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get project. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := c.BackendContext.DeleteProject(project); err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete project. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		c.MyProjects(w, r)
	}
}
