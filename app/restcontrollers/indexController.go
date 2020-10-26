package restcontrollers

import (
	"fmt"
	"net/http"

	"polygnosics/app/services/db/mysqldb"
	"polygnosics/app/utils/page"
	"polygnosics/web/contents"
)

func AboutUsHandler(w http.ResponseWriter, r *http.Request) {
	name := "about"
	p := page.CreatePage(name)
	p.Data["title"] = "About Us"
	p.Data["body"] = "We are awesome"
	page.RenderTemplate(w, name, p)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	name := "index"
	p := page.CreatePage(name)
	p.Data["title"] = "Welcome!"
	p.Data["body"] = "Weolcome to AI Playground"
	page.RenderTemplate(w, name, p)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		page.HandleError("index", "Failed to parse form", w)
		return
	}
	email := r.FormValue("email")
	pwd := r.FormValue("psw")

	err := mysqldb.CheckEmailAndPassword(email, pwd)
	if err != nil {
		page.HandleError("index", err.Error(), w)
		return
	}

	user, err := mysqldb.GetUserByEmail(email)
	if err != nil {
		page.HandleError("index", fmt.Sprintf("Failed to get user. %s", err.Error()), w)
		return
	}
	if err := contents.CreateHome(user.Name); err != nil {
		page.HandleError("index", fmt.Sprintf("Failed to create home page. %s", err.Error()), w)
		return
	}
	if err := contents.CreateUserData(user); err != nil {
		page.HandleError("index", fmt.Sprintf("Failed to create user data. %s", err.Error()), w)
		return
	}

	http.Redirect(w, r, "/user-main", http.StatusSeeOther)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		page.HandleError("index", "Failed to parse form", w)
		return
	}

	uName := r.FormValue("username")
	email := r.FormValue("email")
	pwd := r.FormValue("psw")

	emailExist, err := mysqldb.EmailExists(email)
	if emailExist {
		page.HandleError("index", fmt.Sprintf("Email address %+v already in use", email), w)
		return
	}

	if err != nil {
		page.HandleError("index", fmt.Sprintf("Failed to check email. %s", err.Error()), w)
		return
	}

	usernameExist, err := mysqldb.UserExists(uName)
	if usernameExist {
		page.HandleError("index", fmt.Sprintf("Username %+v already in use", uName), w)
		return
	}

	if err != nil {
		page.HandleError("index", fmt.Sprintf("Failed to check username. %s", err.Error()), w)
		return
	}

	if err = mysqldb.AddUser(uName, email, pwd); err != nil {
		page.HandleError("index", fmt.Sprintf("Failed to add user. %s", err.Error()), w)
		return
	}
	name := "confirm"
	p := page.CreatePage(name)
	p.Data["message"] = "Registration successful."
	p.Data["route"] = "/index"
	page.RenderTemplate(w, name, p)
}
