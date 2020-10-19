package restcontrollers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"aiplayground/app/models"
	"aiplayground/app/utils/jsonutils"
	"aiplayground/app/utils/page"
	"aiplayground/web/contents"
)

// NewProject is the handler for the page that is responsible for creating a new project.
func NewProject(w http.ResponseWriter, r *http.Request) {
	name := "user_data"
	content := &contents.UserData{}
	if err := page.Load(name, content); err != nil {
		errorStr := fmt.Sprintf("Failed to load %s content. %s", name, err.Error())
		page.RenderTemplate(w, "error", contents.CreateError(errorStr))
	}

	if r.Method == "GET" {
		features, err := models.GetAllFeatures()
		if err != nil {
			errorStr := "Failed to get project feature config. " + err.Error()
			page.RenderTemplate(w, "error", contents.CreateError(errorStr))
		}

		featuresMap := make(map[string]interface{})
		name := "new-project"
		p := page.CreatePage(name)
		if p.Data["features"] == nil {
			for _, feature := range features {
				featuresMap[strconv.Itoa(feature.ID)] = feature.Name
			}
			p.Data["features"] = featuresMap
			if err := page.Save(name, p); err != nil {
				errorStr := fmt.Sprintf("Failed to save %s page. %s", name, err.Error())
				page.RenderTemplate(w, "error", contents.CreateError(errorStr))
			}
		}

		page.RenderTemplate(w, name, p)
	} else {
		if err := r.ParseForm(); err != nil {
			page.HandleError("user-main", "Failed to parse form", w)
			return
		}
		if r.FormValue("submitButton") == "Select" {
			selection, err := strconv.Atoi(r.Form["select"][0])
			if err != nil {
				errorStr := fmt.Sprintf("Feature selection failed. %s", err.Error())
				page.RenderTemplate(w, "error", contents.CreateError(errorStr))
			}

			features, err := models.GetAllFeatures()
			if err != nil {
				errorStr := fmt.Sprintf("Failed to get feature config. %s", err.Error())
				page.RenderTemplate(w, "error", contents.CreateError(errorStr))
			}

			content.CurrentProject.UserID = content.User.ID
			content.CurrentProject.FeatureID = selection
			content.CurrentProject.Config = features[selection-1].Config

			name := "new-project"
			p := &page.Page{}
			if err := page.Load(name, p); err != nil {
				errorStr := fmt.Sprintf("Failed to load %s page content. %s", name, err.Error())
				page.RenderTemplate(w, "error", contents.CreateError(errorStr))
			}
			config, err := jsonutils.ProcessJSON(features[selection-1].Config)
			if err != nil {
				errorStr := fmt.Sprintf("Failed to process feature config. %s", err.Error())
				page.RenderTemplate(w, "error", contents.CreateError(errorStr))
			}
			p.Data["config"] = config["config"]

			if err := page.Save("user_data", content); err != nil {
				errorStr := fmt.Sprintf("Failed to save user data content. %s", err.Error())
				page.RenderTemplate(w, "error", contents.CreateError(errorStr))
			}

			page.RenderTemplate(w, "new-project", p)
		} else if r.FormValue("submitButton") == "Start" {
			projectName := r.FormValue("projectName")
			if len(projectName) == 0 {
				projectName = "Default_name"
			}
			content.CurrentProject.Name = projectName
			config, err := jsonutils.ProcessJSON(content.CurrentProject.Config)
			if err != nil {
				errorStr := fmt.Sprintf("Failed to process feature config. %s", err.Error())
				page.RenderTemplate(w, "error", contents.CreateError(errorStr))
			}

			for k := range config["config"].(map[string]interface{}) {
				config["config"].(map[string]interface{})[k] = r.FormValue(k)
			}
			content.CurrentProject.Config, err = json.Marshal(config)
			if err != nil {
				errorStr := fmt.Sprintf("Failed to process config data %s", err.Error())
				page.RenderTemplate(w, "error", contents.CreateError(errorStr))
			}

			if err := models.AddProject(content.CurrentProject); err != nil {
				errorStr := fmt.Sprintf("Failed to add project %s to database. %s", content.CurrentProject.Name, err.Error())
				page.RenderTemplate(w, "error", contents.CreateError(errorStr))
			}

			if content.CurrentProject, err = models.GetProjectByName(content.CurrentProject.Name); err != nil {
				errorStr := fmt.Sprintf("Failed to get project %s from database. %s", content.CurrentProject.Name, err.Error())
				page.RenderTemplate(w, "error", contents.CreateError(errorStr))
			}

			if err := page.Save("user_data", content); err != nil {
				errorStr := fmt.Sprintf("Failed to save user data content. %s", err.Error())
				page.RenderTemplate(w, "error", contents.CreateError(errorStr))
			}
			id := strings.Replace(content.CurrentProject.ID.String(), "-", "", -1)
			http.Redirect(w, r, fmt.Sprintf("/user-main/%s/run", id), http.StatusSeeOther)
		}
	}
}

func UserSettings(w http.ResponseWriter, r *http.Request) {
	name := "user-settings"
	p := page.CreatePage(name)
	page.RenderTemplate(w, name, p)
}

// UserMain renders the main page after login.
func UserMain(w http.ResponseWriter, r *http.Request) {
	name := "user-main"
	p := &page.Page{}
	if err := page.Load(name, p); err != nil {
		errorStr := fmt.Sprintf("Failed to load %s page content. %s", name, err.Error())
		page.RenderTemplate(w, "error", contents.CreateError(errorStr))
	}

	if err := contents.CreateNewProjectConfig(); err != nil {
		errorStr := fmt.Sprintf("Failed to create new project config page content. %s", err.Error())
		page.RenderTemplate(w, "error", contents.CreateError(errorStr))
	}
	page.RenderTemplate(w, name, p)
}
