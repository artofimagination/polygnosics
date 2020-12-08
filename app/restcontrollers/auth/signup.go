package auth

import (
	"fmt"
	"net/http"
	"strings"

	"polygnosics/app"
	"polygnosics/app/restcontrollers/page"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func encryptPassword(password []byte) ([]byte, error) {
	var hashedPassword []byte
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 16)
	if err != nil {
		return hashedPassword, err
	}
	return hashedPassword, nil
}

func generatePath(assetID *uuid.UUID) string {
	path := ""
	increment := 4
	assetIDString := assetID.String()
	assetIDString = strings.Replace(assetIDString, "-", "", -1)
	for i := 0; i < 31; i = i + increment {
		path = fmt.Sprintf("%s/%s", path, assetIDString[i:i+increment])
	}
	return path
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		p := make(map[string]interface{})
		page.RenderTemplate(w, "auth_signup", p)
	} else {
		if err := r.ParseForm(); err != nil {
			page.HandleError("index", "Failed to parse form", w)
			return
		}
		uName := r.FormValue("username")
		email := r.FormValue("email")
		pwd := []byte(r.FormValue("psw"))

		name := "confirm"
		p := make(map[string]interface{})
		p["message"] = "Registration successful."
		p["route"] = "/index"
		p["button_text"] = "OK"

		_, err := app.ContextData.UserDBController.CreateUser(uName, email, pwd, generatePath, encryptPassword)
		if err != nil {
			p["message"] = fmt.Sprintf("Failed to add user. %s", err.Error())
		}
		page.RenderTemplate(w, name, p)
	}
}
