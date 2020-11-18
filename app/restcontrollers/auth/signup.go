package auth

import (
	"fmt"
	"net/http"

	"polygnosics/app/restcontrollers/page"

	"github.com/artofimagination/mysql-user-db-go-interface/mysqldb"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := make(map[string]interface{})
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
		p := make(map[string]interface{})
		p["message"] = "Registration successful."
		p["route"] = "/index"
		p["button_text"] = "OK"

		emailExist, err := mysqldb.EmailExists(email)
		if emailExist {
			p["message"] = fmt.Sprintf("Email address %+v already in use", email)
		}

		if err != nil {
			p["message"] = fmt.Sprintf("Failed to check email. %s", err.Error())
		}

		usernameExist, err := mysqldb.UserExists(uName)
		if usernameExist {
			p["message"] = fmt.Sprintf("Username %+v already in use", uName)
		}

		if err != nil {
			p["message"] = fmt.Sprintf("Failed to check username. %s", err.Error())
		}

		if err = mysqldb.AddUser(uName, email, pwd); err != nil {
			p["message"] = fmt.Sprintf("Failed to add user. %s", err.Error())
		}
		page.RenderTemplate(w, name, &p)
	}
}
