package page

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

type Page struct {
	Name string
	Data map[string]interface{}
}

var htmls = []string{
	"/web/templates/about.html",
	"/web/templates/error.html",
	"/web/templates/index.html",
	"/web/templates/confirm.html",
	"/web/templates/user/user-main.html",
	"/web/templates/user/user-settings.html",
	"/web/templates/user/new-project.html",
	"/web/templates/project/run.html"}
var paths = []string{}
var lock sync.Mutex

// MakeHandler creates the page handler and check the route validity.
func MakeHandler(fn func(http.ResponseWriter, *http.Request), router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Security-Policy", "default-src 'self' data: 'unsafe-inline' 'unsafe-eval'")
		routeMatch := mux.RouteMatch{}
		if matched := router.Match(r, &routeMatch); !matched {
			http.Error(w, "Url does not exist", http.StatusInternalServerError)
		}
		fn(w, r)
	}
}

// RenderTemplate renders html.
func RenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(paths) == 0 {
		for i := 0; i < len(htmls); i++ {
			paths = append(paths, wd+htmls[i])
		}
	}

	t := template.Must(template.ParseFiles(paths...))

	err = t.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Save stores the data in binary file after json marshalling.
func Save(name string, data interface{}) error {
	filename := fmt.Sprintf("%s.bin", name)

	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = binary.Write(f, binary.LittleEndian, r)
	if err != nil {
		return err
	}
	return nil
}

// CreatePage cretaes a page instance.
func CreatePage(name string) *Page {
	p := &Page{Name: name}
	p.Data = make(map[string]interface{})
	return p
}

// Load loads data from binary file.
func Load(name string, data interface{}) error {
	filename := fmt.Sprintf("%s.bin", name)
	lock.Lock()
	defer lock.Unlock()

	dataBinary, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dataBinary, data)
	if err != nil {
		return err
	}
	return nil
}

// HandleError creates page details and renders html template for an error modal.
func HandleError(route string, errorStr string, w http.ResponseWriter) {
	name := "confirm"
	p := CreatePage(name)
	p.Data["message"] = errorStr
	p.Data["route"] = fmt.Sprintf("/%s", route)
	RenderTemplate(w, name, p)
}
