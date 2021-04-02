package frontend

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/artofimagination/polygnosics/businesslogic"
	"github.com/artofimagination/polygnosics/rest"
	"github.com/gorilla/mux"
)

var notFoundFile, notFoundErr = http.Dir("dummy").Open("does-not-exist")

type FileSystem struct {
	http.Dir
}

type RESTController struct {
	BackendContext *businesslogic.Context
}

func NewRESTController(backend *businesslogic.Context) *RESTController {
	controller := &RESTController{
		BackendContext: backend,
	}
	return controller
}

// Open is a custom implementation for the static file server
// that prevents the server from listing the static files, when accessing the path in a browser
func (m FileSystem) Open(name string) (result http.File, err error) {
	f, err := m.Dir.Open(name)
	if err != nil {
		return
	}

	fi, err := f.Stat()
	if err != nil {
		return
	}
	if fi.IsDir() {
		return notFoundFile, notFoundErr
	}
	return f, nil
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi! I am the backend server!")
}

// AddRouting adds front end endpoints.
func (c *RESTController) AddRouting(r *mux.Router) {
	r.HandleFunc("/", sayHello)
	r.HandleFunc("/detect-root-user", rest.MakeHandler(c.detectRootUser))
	r.HandleFunc("/add-user", rest.MakeHandler(c.addUser))
	r.HandleFunc("/get-user-by-id", rest.MakeHandler(c.getUserByID))
	r.HandleFunc("/auth_login", rest.MakeHandler(c.login))

	// Static file servers
	var dirUserAssets string
	flag.StringVar(&dirUserAssets, "dirUserAssets", os.Getenv("USER_STORE_DOCKER"), "the directory to serve user asset files from. Defaults to the current dir")
	flag.Parse()
	handlerUserAssets := http.FileServer(FileSystem{http.Dir(dirUserAssets)})
	r.PathPrefix("/user-assets/").Handler(http.StripPrefix("/user-assets/", handlerUserAssets))
}
