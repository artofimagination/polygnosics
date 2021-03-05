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
	content := c.ContentController.BuildLoginContent()
	if r.Method == GET {
		c.RenderTemplate(w, "auth_login", content)
	} else {
		name := "confirm"
		p := make(map[string]interface{})

		if err := r.ParseForm(); err != nil {
			p["message"] = ErrFailedToParseForm
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

		http.Redirect(w, r, UserMainPath, http.StatusSeeOther)
	}
}

func (c *RESTController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := session.Store.Get(r, "cookie-name")
	if err != nil {
		c.HandleError(w, fmt.Sprintf("Failed to get cookie. %s", errors.WithStack(err)), http.StatusInternalServerError, IndexPath)
		return
	}

	// Revoke users authentication
	session.Values["authenticated"] = false
	if err := session.Save(r, w); err != nil {
		c.HandleError(w, fmt.Sprintf("Failed to save cookie. %s", errors.WithStack(err)), http.StatusInternalServerError, IndexPath)
		return
	}

	http.Redirect(w, r, IndexPath, http.StatusSeeOther)
}

func (c *RESTController) SignupHandler(w http.ResponseWriter, r *http.Request) {
	content := c.ContentController.BuildSignupContent()
	if r.Method == GET {
		c.RenderTemplate(w, "auth_signup", content)
	} else {
		if err := r.ParseForm(); err != nil {
			c.HandleError(w, fmt.Sprintf("Failed to parse form. %s", errors.WithStack(err)), http.StatusInternalServerError, IndexPath)
			return
		}
		uName := r.FormValue("username")
		email := r.FormValue("email")
		pwd := []byte(r.FormValue("psw"))
		group := r.FormValue("developer")
		if group == "" {
			group = "client"
		}

		if err := c.BackendContext.AddUser(uName, email, pwd, group); err != nil {
			c.HandleError(w, fmt.Sprintf("Failed to add user. %s", errors.WithStack(err)), http.StatusInternalServerError, IndexPath)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Registration successful")
	}
}
