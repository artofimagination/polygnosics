package restcontrollers

import (
	"fmt"
	"net/http"
	"strings"

	"polygnosics/app/utils/page"
	"polygnosics/web/contents"
)

func StartWebRTC(w http.ResponseWriter, r *http.Request) {

}

func RunProject(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		name := "user_data"
		content := &contents.UserData{}
		if err := page.Load(name, content); err != nil {
			errorStr := fmt.Sprintf("Failed to load %s page content. %s", name, err.Error())
			page.RenderTemplate(w, "error", contents.CreateError(errorStr))
		}

		p := page.CreatePage("project")
		if err := page.Load(name, p); err != nil {
			errorStr := fmt.Sprintf("Failed to load %s page content. %s", name, err.Error())
			page.RenderTemplate(w, "error", contents.CreateError(errorStr))
		}

		p.Data["project_id"] = strings.Replace(content.CurrentProject.ID.String(), "-", "", -1)
		page.RenderTemplate(w, "run", p)
	}
}
