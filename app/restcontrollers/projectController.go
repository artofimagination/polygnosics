package restcontrollers

import (
	"fmt"
	"log"
	"net/http"
	"polygnosics/app/businesslogic/project"
	"polygnosics/app/restcontrollers/contents"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *RESTController) MyProjects(w http.ResponseWriter, r *http.Request) {
	pUser := c.ContentController.GetUserContent()
	pProduct, err := c.ContentController.GetUserProjectContent(&c.ContentController.UserData.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get project content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	for k, v := range pProduct {
		pUser[k] = v
	}
	c.RenderTemplate(w, "my-projects", pUser)
}

func (c *RESTController) ProjectDetails(w http.ResponseWriter, r *http.Request) {
	p := c.ContentController.GetUserContent()
	name := UserMain
	if err := r.ParseForm(); err != nil {
		p["message"] = contents.ErrFailedToParseForm
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
		p["message"] = "Failed to get project content"
		c.RenderTemplate(w, name, p)
		return
	}

	for k, v := range pProject {
		p[k] = v
	}
	log.Println(p)
	c.RenderTemplate(w, "project-details", p)
}

func (c *RESTController) CreateProject(w http.ResponseWriter, r *http.Request) {
	p := c.ContentController.GetUserContent()
	name := UserMain
	if err := r.ParseForm(); err != nil {
		p["message"] = contents.ErrFailedToParseForm
		c.RenderTemplate(w, name, p)
		return
	}
	productID, err := uuid.Parse(r.FormValue("product"))
	if err != nil {
		p["message"] = "Failed to parse product id"
		c.RenderTemplate(w, name, p)
		return
	}

	pProduct, err := c.ContentController.GetProductContent(&productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	for k, v := range pProduct {
		p[k] = v
	}

	if r.Method == GET {
		c.RenderTemplate(w, "new-project-wizard", p)
	} else {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			p["message"] = contents.ErrFailedToParseForm
			c.RenderTemplate(w, name, p)
			return
		}

		if err := contents.ValidateVisibility(r.FormValue("visibility")); err != nil {
			p["message"] = err
			c.RenderTemplate(w, name, p)
			return
		}

		productID, err := uuid.Parse(r.FormValue("product"))
		if err != nil {
			p["message"] = "Failed to parse project id"
			c.RenderTemplate(w, name, p)
			return
		}

		projectData, err := c.UserDBController.CreateProject(
			r.FormValue("projectName"),
			r.FormValue("visibility"),
			&c.ContentController.UserData.ID,
			&productID,
			c.ContentController.GeneratePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		err = c.ContentController.UploadFile(projectData.Assets, contents.ProjectAvatar, contents.DefaultProjectAvatarPath, "project-avatar", r)
		if err != nil {
			if errDelete := c.UserDBController.DeleteProject(&projectData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload avatar. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		projectData.Details.SetField(contents.ProjectState, project.NotRunning)
		projectData.Details.SetField(contents.ProjectVisibility, r.FormValue("visibility"))
		projectData.Details.SetField(contents.ProjectServerLogging, contents.GetBooleanString(r.FormValue("serverLogging")))
		projectData.Details.SetField(contents.ProjectClientLogging, contents.GetBooleanString(r.FormValue("clientLogging")))
		if err := c.UserDBController.UpdateProjectDetails(projectData); err != nil {
			if errDelete := c.UserDBController.DeleteProject(&projectData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to update product details. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := c.UserDBController.UpdateProjectAssets(projectData); err != nil {
			if errDelete := c.UserDBController.DeleteProject(&projectData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to update product assets. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		c.RenderTemplate(w, name, p)
	}
}

func StartWebRTC(w http.ResponseWriter, r *http.Request) {

}

func (c *RESTController) RunProject(w http.ResponseWriter, r *http.Request) {
	log.Println("Running")
}
