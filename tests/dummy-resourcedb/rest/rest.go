package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
	data, err := ioutil.ReadFile("/resources/testData.json")
	if err != nil {
		return nil, err
	}
	jsonData := make(map[string]interface{})
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	requestData := make([]map[string]interface{}, 0)
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
	r.HandleFunc("/get-categories", makeHandler(c.getCategories))
	r.HandleFunc("/get-resource-by-id", makeHandler(c.getResource))
	r.HandleFunc("/add-resource", makeHandler(c.addResource))
	r.HandleFunc("/update-resource", makeHandler(c.updateResource))
	r.HandleFunc("/delete-resource", makeHandler(c.deleteResource))

	r.HandleFunc("/clear-request-data", makeHandler(c.clearRequestData))
	r.HandleFunc("/get-request-data", makeHandler(c.getRequestData))

	return r
}

func (c *Controller) getCategories(w ResponseWriter, r *Request) {
	data := make(map[string]interface{})
	data["data"] = c.TestData["categories"]
	w.encodeResponse(data, http.StatusCreated)
}

func (c *Controller) getResource(w ResponseWriter, r *Request) {
	log.Println("Get resource")
	if err := c.ParseForm(r, "/get-resource-by-id"); err != nil {
		w.writeError(fmt.Sprintf("UserDB: %s", err.Error()), http.StatusBadRequest)
		return
	}
	id := r.FormValue("id")

	resource := c.TestData["resources"].(map[string]interface{})[id]
	data := make(map[string]interface{})
	data["data"] = resource
	w.encodeResponse(data, http.StatusCreated)
}

func (c *Controller) addResource(w ResponseWriter, r *Request) {
	log.Println("Add resource")

	requestData, err := c.decodeRequest(r, "/add-resource")
	if err != nil {
		w.writeError(fmt.Sprintf("ResourceDB -> %s", err.Error()), http.StatusBadRequest)
		return
	}

	resource := make(map[string]interface{})
	resource["id"] = "c1e6122b-7986-417d-8bf6-ddf2dd9289f2"
	resource["category"] = requestData["category"]
	resource["content"] = requestData["content"]

	c.TestData["resources"].(map[string]interface{})["c1e6122b-7986-417d-8bf6-ddf2dd9289f2"] = resource
	data := make(map[string]interface{})
	data["data"] = resource
	data["error"] = ""
	w.encodeResponse(data, http.StatusCreated)
}

func (c *Controller) updateResource(w ResponseWriter, r *Request) {
	log.Println("Update resource")
	requestData, err := c.decodeRequest(r, "/update-resource")
	if err != nil {
		w.writeError(fmt.Sprintf("ResourceDB -> %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, ok := c.TestData["resources"].(map[string]interface{})[requestData["id"].(string)]; !ok {
		w.writeError("ResourceDB -> Missing resource", http.StatusNoContent)
		return
	}

	resource := c.TestData["resources"].(map[string]interface{})[requestData["id"].(string)].(map[string]interface{})
	resource["content"] = requestData["content"]
	w.writeData("OK", http.StatusOK)
}

func (c *Controller) deleteResource(w ResponseWriter, r *Request) {
	log.Println("Delete resource")
	requestData, err := c.decodeRequest(r, "/delete-resource")
	if err != nil {
		w.writeError(fmt.Sprintf("ResourceDB -> %s", err.Error()), http.StatusBadRequest)
		return
	}

	if _, ok := c.TestData["resources"].(map[string]interface{})[requestData["id"].(string)]; !ok {
		w.writeError("ResourceDB -> Missing resource", http.StatusNoContent)
		return
	}

	w.writeData("OK", http.StatusOK)
}
