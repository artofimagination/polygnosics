package restcontrollers

import (
	"fmt"
	"net/http"
	"polygnosics/app/businesslogic"
	"polygnosics/web/contents"

	"github.com/pkg/errors"
)

func (c *RESTController) MyProjects(w http.ResponseWriter, r *http.Request) {
	content, err := c.ContentController.BuildMyProjectsContent()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get project content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	c.RenderTemplate(w, MyProjects, content)
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

	c.RenderTemplate(w, UserMainProjectDetails, content)
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

		if err := contents.ValidateVisibility(r.FormValue(businesslogic.ProjectVisibilityKey)); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse visibility. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		projectData, err := c.UserDBController.CreateProject(
			r.FormValue(businesslogic.ProjectNameKey),
			r.FormValue(businesslogic.ProjectVisibilityKey),
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

		if err := c.BackendContext.UpdateProjectData(projectData, containerID, r); err != nil {
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

func (c *RESTController) HandleStatusRequest(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseItemID(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse project id. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	reachable, err := c.BackendContext.CheckProject(projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to access project. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	if !reachable {
		http.Error(w, fmt.Sprintf("Failed to access project"), http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *RESTController) RunProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseItemID(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse project id. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	if err := c.BackendContext.RunProject(&c.ContentController.UserData.ID, projectID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to run project. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	content, err := c.ContentController.BuildProjectRunContent(projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get project content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	c.RenderTemplate(w, "show", content)
}

func (c *RESTController) ShowProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseItemID(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse project id. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	content, err := c.ContentController.BuildProjectRunContent(projectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get project content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	c.RenderTemplate(w, "show", content)
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

func (c *RESTController) EditProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseItemID(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse project id. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	if r.Method == GET {
		content, err := c.ContentController.BuildProjectEditContent(projectID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to build project edit content. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		c.RenderTemplate(w, "project-edit", content)
	} else {
		project, err := c.UserDBController.GetProject(projectID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get project. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := c.BackendContext.UploadFiles(project.Assets, r); err != nil {
			http.Error(w, fmt.Sprintf("Failed to upload assets. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := c.BackendContext.EditProjectData(project, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		content, err := c.ContentController.BuildMyProjectsContent()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get my projects content. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}
		c.RenderTemplate(w, MyProjects, content)
	}
}
