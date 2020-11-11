package auth

import (
	"net/http"
	"polygnosics/app/utils/page"

	"polygnosics/app/services/db/mysqldb"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := page.Page{}
		page.RenderTemplate(w, "auth_login", &p)
	} else {
		name := "confirm"
		p := page.CreatePage(name)
		p.Data["route"] = "/index"
		p.Data["button_text"] = "OK"

		if err := r.ParseForm(); err != nil {
			p.Data["message"] = "Failed to parse form"

			page.RenderTemplate(w, name, p)
			return
		}
		email := r.FormValue("email")
		pwd := r.FormValue("psw")

		err := mysqldb.CheckEmailAndPassword(email, pwd)
		if err != nil {
			p.Data["message"] = err.Error()
			page.RenderTemplate(w, name, p)
			return
		}

		http.Redirect(w, r, "/user-main", http.StatusSeeOther)
	}
}
