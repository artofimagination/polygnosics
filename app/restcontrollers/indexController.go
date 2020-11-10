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
