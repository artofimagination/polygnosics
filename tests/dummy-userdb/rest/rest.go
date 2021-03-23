package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/gorilla/mux"
)

const (
	UsersKey         = "users"
	UserProductsKey  = "user_products"
	UsersProjectKey  = "user_projects"
	UsersUsernameKey = "username"
	UsersEmailKey    = "email"
	UsersIDKey       = "id"
	UsersPasswordKey = "password"
	UserGroupKey     = "group"

	SettingsKey   = "settings"
	SettingsIDKey = "id"
	DetailsKey    = "details"
	DetailsIDKey  = "id"
	AssetsKey     = "assets"
	AssetsIDKey   = "id"
	DataMapKey    = "datamap"
)

const (
	ProductKey                 = "products"
	ProductAvatarKey           = "avatar"
	ProductMainAppKey          = "main_app"
	ProductClientAppKey        = "client_app"
	ProductDescriptionKey      = "description"
	ProductShortDescriptionKey = "short_description"
	ProductNameKey             = "name"
	ProductRequires3DKey       = "requires_3d"
	ProductURLKey              = "url"
	ProductPublicKey           = "is_public"
	ProductPricingKey          = "pricing"
	ProductPriceKey            = "amount"
	ProductTagsKey             = "tags"
	ProductCategoriesKey       = "categories"
)

const (
	UserPathAdd            = "/add-user"
	UserPathDetectRoot     = "/detect-root-user"
	UserPathUpdateSettings = "/update-user-settings"
	UserPathUpdateAssets   = "/update-user-assets"
	UserPathGetByEmail     = "/get-user-by-email"
	UserPathGetByID        = "/get-user-by-id"
	UserPathDeleteByID     = "/delete-user"
)

const (
	ProjectKey = "projects"
)

var CategoriesKey = "categories"

const (
	UserTestUUID         = "026eede8-0b9b-4355-ad48-8a4f6cf0b49e"
	UserSettingsTestUUID = "8b683a4c-198a-4cfd-abb1-7a3715a51bbb"
	UserAssetsTestUUID   = "9f02fbd5-15b7-465a-a941-f4fdc11db23e"
	RootUserTestUUID     = "f9ebc23d-81cc-4bf2-b908-7e88c58ebe91"
)

func convertCheckboxValueToText(input string) string {
	if input == "" {
		return "unchecked"
	}
	return input
}

func NewController() (*Controller, error) {
	data, err := ioutil.ReadFile("/user-assets/testData.json")
	if err != nil {
		return nil, err
	}
	jsonData := make(map[string]interface{})
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	requestData := make(map[string]interface{})
	return &Controller{
		TestData:    jsonData,
		RequestData: requestData,
	}, nil
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi! I am a dummy server!")
}

func (c *Controller) CreateRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", sayHello)
	r.HandleFunc(UserPathAdd, makeHandler(c.addUser))
	r.HandleFunc(UserPathUpdateSettings, makeHandler(c.updateUserSettings))
	r.HandleFunc(UserPathUpdateAssets, makeHandler(c.updateUserAssets))
	r.HandleFunc(UserPathGetByID, makeHandler(c.getUserByID))
	r.HandleFunc(UserPathGetByEmail, makeHandler(c.getUserByEmail))
	r.HandleFunc(UserPathDeleteByID, makeHandler(c.deleteUserByID))

	r.HandleFunc("/clear-request-data", makeHandler(c.clearRequestData))
	r.HandleFunc("/get-request-data", makeHandler(c.getRequestData))

	return r
}

func (c *Controller) addUser(w ResponseWriter, r *Request) {
	requestData, err := c.decodeRequest(r, UserPathAdd)
	if err != nil {
		w.writeError(fmt.Sprintf("UserDB: %s", err.Error()), http.StatusBadRequest)
		return
	}

	for _, v := range c.TestData[UsersKey].(map[string]interface{}) {
		if v.(map[string]interface{})[UsersUsernameKey] == requestData[UsersUsernameKey] {
			w.writeError(dbcontrollers.ErrDuplicateNameEntry.Error(), http.StatusAccepted)
			return
		}
	}
	userData := make(map[string]interface{})
	userData[AssetsKey] = make(map[string]interface{})
	userData[AssetsKey].(map[string]interface{})[AssetsIDKey] = UserAssetsTestUUID
	userData[AssetsKey].(map[string]interface{})[DataMapKey] = make(map[string]interface{})
	userData[SettingsKey] = make(map[string]interface{})
	userData[SettingsKey].(map[string]interface{})[SettingsIDKey] = UserSettingsTestUUID
	userData[SettingsKey].(map[string]interface{})[DataMapKey] = make(map[string]interface{})
	userData[UsersUsernameKey] = requestData[UsersUsernameKey]
	userData[UsersEmailKey] = requestData[UsersEmailKey]
	userData[UsersPasswordKey] = requestData[UsersPasswordKey]
	userData[UsersIDKey] = UserTestUUID
	if requestData[UsersUsernameKey] == "root" {
		userData[UsersIDKey] = RootUserTestUUID
	}

	c.TestData[UsersKey].(map[string]interface{})[UserTestUUID] = userData
	data := make(map[string]interface{})
	byteData, err := json.Marshal(c.TestData[UsersKey].(map[string]interface{})[UserTestUUID])
	if err != nil {
		w.writeError(fmt.Sprintf("UserDB: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	data["data"] = string(byteData)
	w.encodeResponse(data, http.StatusCreated)
}

func (c *Controller) updateUserSettings(w ResponseWriter, r *Request) {
	requestData, err := c.decodeRequest(r, UserPathUpdateSettings)
	if err != nil {
		w.writeError(fmt.Sprintf("UserDB: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user := c.TestData[UsersKey].(map[string]interface{})[requestData["user-id"].(string)].(map[string]interface{})
	user[SettingsKey] = requestData["user-data"]
	w.writeData("OK", http.StatusOK)
}

func (c *Controller) updateUserAssets(w ResponseWriter, r *Request) {
	requestData, err := c.decodeRequest(r, UserPathUpdateAssets)
	if err != nil {
		w.writeError(fmt.Sprintf("UserDB: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user := c.TestData[UsersKey].(map[string]interface{})[requestData["user-id"].(string)].(map[string]interface{})
	user[AssetsKey] = requestData["user-data"]
	w.writeData("OK", http.StatusOK)
}

func (c *Controller) getUserByID(w ResponseWriter, r *Request) {
	data := make(map[string]interface{})
	if err := c.ParseForm(r, UserPathGetByID); err != nil {
		w.writeError(fmt.Sprintf("UserDB: %s", err.Error()), http.StatusBadRequest)
		return
	}

	id := r.FormValue(UsersIDKey)
	for k, v := range c.TestData[UsersKey].(map[string]interface{}) {
		if k == id {
			data["data"] = v
			break
		}
	}

	w.encodeResponse(data, http.StatusOK)
}

func (c *Controller) getUserByEmail(w ResponseWriter, r *Request) {
	data := make(map[string]interface{})
	if err := c.ParseForm(r, UserPathGetByEmail); err != nil {
		w.writeError(fmt.Sprintf("UserDB: %s", err.Error()), http.StatusBadRequest)
		return
	}
	email := r.FormValue(UsersEmailKey)
	for _, v := range c.TestData[UsersKey].(map[string]interface{}) {
		if v.(map[string]interface{})[UsersEmailKey] == email {
			data["data"] = v
			break
		}
	}
	if _, ok := data["data"]; !ok {
		w.writeError(dbcontrollers.ErrUserNotFound.Error(), http.StatusAccepted)
		return
	}

	w.encodeResponse(data, http.StatusOK)
}

func (c *Controller) deleteUserByID(w ResponseWriter, r *Request) {
	requestData, err := c.decodeRequest(r, UserPathUpdateSettings)
	if err != nil {
		w.writeError(fmt.Sprintf("UserDB: %s", err.Error()), http.StatusBadRequest)
		return
	}

	for k, _ := range c.TestData[UsersKey].(map[string]interface{}) {
		if k == requestData["id"] {
			delete(c.TestData[UsersKey].(map[string]interface{}), k)
			w.writeData("OK", http.StatusOK)
			return
		}
	}

	w.writeError(dbcontrollers.ErrUserNotFound.Error(), http.StatusAccepted)
}
