package auth

import (
	"fmt"
	"net/http"

	"polygnosics/app"
	"polygnosics/app/restcontrollers/page"
	"polygnosics/app/restcontrollers/session"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// TODO Issue#40: Replace with redis storage.

var UserData *models.UserData

var errIncorrectEmailOrPass = errors.New("Incorrect email or password")

func authenticate(email string, password string, user *models.User) error {
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil || email != user.Email {
		return errIncorrectEmailOrPass
	}
	return nil
}

// LoginHandler checks the user email and password.
// On success generates and stores a cookie in the session sotre and adds it to the response
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := make(map[string]interface{})
		page.RenderTemplate(w, "auth_login", p)
	} else {
		name := "confirm"
		p := make(map[string]interface{})
		p["route"] = "/index"
		p["button_text"] = "OK"

		if err := r.ParseForm(); err != nil {
			p["message"] = "Failed to parse form"
			page.RenderTemplate(w, name, p)
			return
		}
		email := r.FormValue("email")
		pwd := r.FormValue("psw")

		// TODO Issue#45: replace this with elastic search
		searchedUser, err := app.ContextData.UserDBController.GetUserByEmail(email)
		if err == dbcontrollers.ErrUserNotFound {
			p["message"] = "Incorrect email or password"
			page.RenderTemplate(w, name, p)
			return
		} else if err != nil {
			p["message"] = fmt.Sprintf("Failed to get user data.\nDetails: %s", err.Error())
			page.RenderTemplate(w, name, p)
			return
		}

		// Get and check user and password
		err = app.ContextData.UserDBController.Authenticate(&searchedUser.ID, email, pwd, authenticate)
		if err == dbcontrollers.ErrUserNotFound || err == errIncorrectEmailOrPass {
			p["message"] = errIncorrectEmailOrPass
			page.RenderTemplate(w, name, p)
			return
		} else if err != nil {
			p["message"] = fmt.Sprintf("Failed to authenticate user.\nDetails: %s", err.Error())
			page.RenderTemplate(w, name, p)
			return
		}

		user, err := app.ContextData.UserDBController.GetUser(&searchedUser.ID)
		if err != nil {
			p["message"] = fmt.Sprintf("Failed to get user. %s", err.Error())
		}
		UserData = user

		// Create session cookie.
		sess, err := session.Store.Get(r, "cookie-name")
		if err != nil {
			p["message"] = fmt.Sprintf("Failed to create cookie. %s", errors.WithStack(err))
			page.RenderTemplate(w, name, p)
			return
		}
		sess.Options.MaxAge = 60000
		sess.Values["authenticated"] = true
		sess.Values["user"] = user.ID.String()

		cookieKey, err := session.EncryptUserAndOrigin(user.ID, r.RemoteAddr)
		if err != nil {
			p["message"] = fmt.Sprintf("Failed to generate cookie data. %s", errors.WithStack(err))
			page.RenderTemplate(w, name, p)
			return
		}
		sess.Values["cookie_key"] = cookieKey

		if err := sess.Save(r, w); err != nil {
			p["message"] = fmt.Sprintf("Failed to save cookie. %s", errors.WithStack(err))
			page.RenderTemplate(w, name, p)
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
