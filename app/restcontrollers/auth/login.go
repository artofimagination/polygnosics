package auth

import (
	"database/sql"
	"fmt"
	"net/http"

	"polygnosics/app/restcontrollers/page"
	"polygnosics/app/restcontrollers/session"
	"polygnosics/app/services/db/mysqldb"

	"github.com/pkg/errors"
)

// LoginHandler checks the user email and password.
// On success generates and stores a cookie in the session sotre and adds it to the response
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := make(map[string]interface{})
		page.RenderTemplate(w, "auth_login", &p)
	} else {
		name := "confirm"
		p := make(map[string]interface{})
		p["route"] = "/index"
		p["button_text"] = "OK"

		if err := r.ParseForm(); err != nil {
			p["message"] = "Failed to parse form"
			page.RenderTemplate(w, name, &p)
			return
		}
		email := r.FormValue("email")
		pwd := r.FormValue("psw")

		// Get and check user and password
		user, err := mysqldb.GetUserByEmail(email)
		if err == sql.ErrNoRows {
			p["message"] = "Incorrect email or password"
			page.RenderTemplate(w, name, &p)
			return
		}

		if !mysqldb.IsPasswordCorrect(pwd, user) {
			p["message"] = "Incorrect email or password"
			page.RenderTemplate(w, name, &p)
			return
		}

		// Create session cookie.
		sess, err := session.Store.Get(r, "cookie-name")
		if err != nil {
			p["message"] = fmt.Sprintf("Failed to create cookie. %s", errors.WithStack(err))
			page.RenderTemplate(w, name, &p)
			return
		}
		sess.Options.MaxAge = 600
		sess.Values["authenticated"] = true
		sess.Values["user"] = user.ID.String()

		cookieKey, err := session.EncryptUserAndOrigin(user.ID, r.RemoteAddr)
		if err != nil {
			p["message"] = fmt.Sprintf("Failed to generate cookie data. %s", errors.WithStack(err))
			page.RenderTemplate(w, name, &p)
			return
		}
		sess.Values["cookie_key"] = cookieKey

		if err := sess.Save(r, w); err != nil {
			p["message"] = fmt.Sprintf("Failed to save cookie. %s", errors.WithStack(err))
			page.RenderTemplate(w, name, &p)
			return
		}

		http.Redirect(w, r, "/user-main", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := session.Store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get cookie. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	// Revoke users authentication
	session.Values["authenticated"] = false
	if err := session.Save(r, w); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save cookie. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/index", http.StatusSeeOther)
}
