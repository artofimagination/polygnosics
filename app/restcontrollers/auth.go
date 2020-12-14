package restcontrollers

import (
	"fmt"
	"net/http"

	"polygnosics/app/restcontrollers/session"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var errIncorrectEmailOrPass = errors.New("Incorrect email or password")

func authenticate(email string, password string, user *models.User) error {
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil || email != user.Email {
		return errIncorrectEmailOrPass
	}
	return nil
}

// LoginHandler checks the user email and password.
// On success generates and stores a cookie in the session sotre and adds it to the response
func (c *RESTController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == GET {
		p := make(map[string]interface{})
		c.RenderTemplate(w, "auth_login", p)
	} else {
		name := Confirm
		p := make(map[string]interface{})
		p["route"] = "/index"
		p["button_text"] = "OK"

		if err := r.ParseForm(); err != nil {
			p["message"] = "Failed to parse form"
			c.RenderTemplate(w, name, p)
			return
		}
		email := r.FormValue("email")
		pwd := r.FormValue("psw")

		// TODO Issue#45: replace this with elastic search
		searchedUser, err := c.UserDBController.GetUserByEmail(email)
		if err == dbcontrollers.ErrUserNotFound {
			p["message"] = "Incorrect email or password"
			c.RenderTemplate(w, name, p)
			return
		} else if err != nil {
			p["message"] = fmt.Sprintf("Failed to get user data.\nDetails: %s", err.Error())
			c.RenderTemplate(w, name, p)
			return
		}

		// Get and check user and password
		err = c.UserDBController.Authenticate(&searchedUser.ID, email, pwd, authenticate)
		if err == dbcontrollers.ErrUserNotFound || err == errIncorrectEmailOrPass {
			p["message"] = errIncorrectEmailOrPass
			c.RenderTemplate(w, name, p)
			return
		} else if err != nil {
			p["message"] = fmt.Sprintf("Failed to authenticate user.\nDetails: %s", err.Error())
			c.RenderTemplate(w, name, p)
			return
		}

		user, err := c.UserDBController.GetUser(&searchedUser.ID)
		if err != nil {
			p["message"] = fmt.Sprintf("Failed to get user. %s", err.Error())
		}
		c.ContentController.UserData = user

		// Create session cookie.
		sess, err := session.Store.Get(r, "cookie-name")
		if err != nil {
			p["message"] = fmt.Sprintf("Failed to create cookie. %s", errors.WithStack(err))
			c.RenderTemplate(w, name, p)
			return
		}
		sess.Options.MaxAge = 60000
		sess.Values["authenticated"] = true
		sess.Values["user"] = user.ID.String()

		cookieKey, err := session.EncryptUserAndOrigin(user.ID, r.RemoteAddr)
		if err != nil {
			p["message"] = fmt.Sprintf("Failed to generate cookie data. %s", errors.WithStack(err))
			c.RenderTemplate(w, name, p)
			return
		}
		sess.Values["cookie_key"] = cookieKey

		if err := sess.Save(r, w); err != nil {
			p["message"] = fmt.Sprintf("Failed to save cookie. %s", errors.WithStack(err))
			c.RenderTemplate(w, name, p)
			return
		}

		http.Redirect(w, r, "/user-main", http.StatusSeeOther)
	}
}

func (*RESTController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
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

func encryptPassword(password []byte) ([]byte, error) {
	var hashedPassword []byte
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 16)
	if err != nil {
		return hashedPassword, err
	}
	return hashedPassword, nil
}

func (c *RESTController) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == GET {
		p := make(map[string]interface{})
		c.RenderTemplate(w, "auth_signup", p)
	} else {
		if err := r.ParseForm(); err != nil {
			c.HandleError("index", "Failed to parse form", w)
			return
		}
		uName := r.FormValue("username")
		email := r.FormValue("email")
		pwd := []byte(r.FormValue("psw"))

		p := make(map[string]interface{})
		p["message"] = "Registration successful."
		p["route"] = "/index"
		p["button_text"] = "OK"

		_, err := c.UserDBController.CreateUser(uName, email, pwd, c.ContentController.GeneratePath, encryptPassword)
		if err != nil {
			p["message"] = fmt.Sprintf("Failed to add user. %s", err.Error())
		}
		c.RenderTemplate(w, Confirm, p)
	}
}
