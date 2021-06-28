package frontend

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/artofimagination/polygnosics/rest"
)

const (
	UsersUsernameKey = "username"
	UsersEmailKey    = "email"
	UsersPasswordKey = "password"
	UserGroupKey     = "group"
	UsersIDKey       = "id"

	SettingsKey = "settings"
	AssetsKey   = "assets"
)

func (c *RESTController) addUser(w rest.ResponseWriter, r *rest.Request) {
	requestData, err := r.DecodeRequest()
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := c.BackendContext.AddUser(
		requestData[UsersUsernameKey].(string),
		requestData[UsersEmailKey].(string),
		[]byte(requestData[UsersPasswordKey].(string)),
		requestData[UserGroupKey].(string)); err != nil {
		w.WriteError(err.Error(), http.StatusAccepted)
		return
	}

	w.WriteData("OK", http.StatusCreated)
}

func (c *RESTController) login(w rest.ResponseWriter, r *rest.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusBadRequest)
		return
	}

	email := r.FormValue(UsersEmailKey)
	password := []byte(r.FormValue(UsersPasswordKey))
	user, err := c.BackendContext.Login(email, password)
	if err != nil {
		w.WriteError(err.Error(), http.StatusAccepted)
		return
	}
	userMap := make(map[string]interface{})
	response := make(map[string]interface{})
	data, err := json.Marshal(user)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
	}

	err = json.Unmarshal(data, &userMap)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
	}
	delete(userMap, "password")
	response["data"] = userMap

	w.EncodeResponse(response)
}

func (c *RESTController) detectRootUser(w rest.ResponseWriter, r *rest.Request) {
	data := make(map[string]interface{})
	rootFound, err := c.BackendContext.DetectRootUser()
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	data["data"] = rootFound
	w.EncodeResponse(data)
}

func (c *RESTController) getUserByID(w rest.ResponseWriter, r *rest.Request) {
	data := make(map[string]interface{})
	user, err := r.ForwardRequest(rest.UserDBAddress)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	data["data"] = user

	w.EncodeResponse(data)
}
