package auth

import (
	"fmt"
	"net/http"

	"polygnosics/app/services/db/mysqldb"
	"polygnosics/app/utils/page"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := page.Page{}
		page.RenderTemplate(w, "auth_signup", &p)
	} else {
		if err := r.ParseForm(); err != nil {
			page.HandleError("index", "Failed to parse form", w)
			return
		}
		uName := r.FormValue("username")
		email := r.FormValue("email")
		pwd := r.FormValue("psw")

		name := "confirm"
		p := page.CreatePage(name)
		p.Data["message"] = "Registration successful."
		p.Data["route"] = "/index"
		p.Data["button_text"] = "OK"

		emailExist, err := mysqldb.EmailExists(email)
		if emailExist {
			p.Data["message"] = fmt.Sprintf("Email address %+v already in use", email)
			page.RenderTemplate(w, name, p)
			return
		}

		if err != nil {
			p.Data["message"] = fmt.Sprintf("Failed to check email. %s", err.Error())
			page.RenderTemplate(w, name, p)
			return
		}

		usernameExist, err := mysqldb.UserExists(uName)
		if usernameExist {
			p.Data["message"] = fmt.Sprintf("Username %+v already in use", uName)
			page.RenderTemplate(w, name, p)
			return
		}

		if err != nil {
			p.Data["message"] = fmt.Sprintf("Failed to check username. %s", err.Error())
			page.RenderTemplate(w, name, p)
			return
		}

		if err = mysqldb.AddUser(uName, email, pwd); err != nil {
			p.Data["message"] = fmt.Sprintf("Failed to add user. %s", err.Error())
			page.RenderTemplate(w, name, p)
			return
		}
		page.RenderTemplate(w, name, p)
	}
}
