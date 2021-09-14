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
	BackendContext    *businesslogic.Context
	UserDBAddress     *rest.Server
	ResourceDBAddress *rest.Server
}

func (c *RESTController) ForwardResourceDBRequest(r *rest.Request) (interface{}, error) {
	return r.ForwardRequest(c.ResourceDBAddress.GetAddress())
}

func (c *RESTController) ForwardUserDBRequest(r *rest.Request) (interface{}, error) {
	return r.ForwardRequest(c.UserDBAddress.GetAddress())
}

func NewRESTController(backend *businesslogic.Context) *RESTController {
	controller := &RESTController{
		BackendContext:    backend,
		UserDBAddress:     &rest.Server{},
		ResourceDBAddress: &rest.Server{},
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
	// User endpoints
	r.HandleFunc("/detect-root-user", rest.MakeHandler(c.detectRootUser))
	r.HandleFunc("/add-user", rest.MakeHandler(c.addUser))
	r.HandleFunc("/get-user-by-id", rest.MakeHandler(c.getUserByID))
	r.HandleFunc("/auth_login", rest.MakeHandler(c.login))
	r.HandleFunc("/get-categories", rest.MakeHandler(c.getCategoriesMap))

	// Resource endpoints
	r.HandleFunc("/get-tutorials", rest.MakeHandler(c.getTutorials))
	r.HandleFunc("/get-tutorial", rest.MakeHandler(c.getSingleItem))
	r.HandleFunc("/get-files", rest.MakeHandler(c.getFiles))
	r.HandleFunc("/get-files-section", rest.MakeHandler(c.getSingleItem))
	r.HandleFunc("/get-news-feed", rest.MakeHandler(c.getNewsFeed))
	r.HandleFunc("/get-news-item", rest.MakeHandler(c.getSingleItem))
	r.HandleFunc("/get-faqs", rest.MakeHandler(c.getFAQs))
	r.HandleFunc("/get-faq", rest.MakeHandler(c.getSingleItem))
	r.HandleFunc("/get-faq-groups", rest.MakeHandler(c.getFAQGroups))
	resources := r.PathPrefix("/resources").Subrouter()
	resources.HandleFunc("/create-news-item", rest.MakeHandler(c.addNewsFeedEntry))
	resources.HandleFunc("/edit-news-item", rest.MakeHandler(c.updateNewsEntry))
	resources.HandleFunc("/create-files-item", rest.MakeHandler(c.addFileSection))
	resources.HandleFunc("/edit-files-item", rest.MakeHandler(c.updateFileSection))
	resources.HandleFunc("/delete-files-item", rest.MakeHandler(c.deleteFileSection))
	resources.HandleFunc("/create-tutorial-item", rest.MakeHandler(c.addTutorial))
	resources.HandleFunc("/edit-tutorial-item", rest.MakeHandler(c.updateTutorial))
	resources.HandleFunc("/delete-tutorial-item", rest.MakeHandler(c.deleteTutorial))
	resources.HandleFunc("/create-faq-item", rest.MakeHandler(c.addFAQ))
	resources.HandleFunc("/edit-faq-item", rest.MakeHandler(c.updateFAQ))
	resources.HandleFunc("/delete-faq-item", rest.MakeHandler(c.deleteFAQ))
	resources.HandleFunc("/get-article", rest.MakeHandler(c.getSingleItem))

	// Static file servers
	var dirUserAssets string
	flag.StringVar(&dirUserAssets, "dirUserAssets", os.Getenv("USER_STORE_DOCKER"), "the directory to serve user asset files from. Defaults to the current dir")
	flag.Parse()
	handlerUserAssets := http.FileServer(FileSystem{http.Dir(dirUserAssets)})
	r.PathPrefix("/user-assets/").Handler(http.StripPrefix("/user-assets/", handlerUserAssets))

	var dirResources string
	flag.StringVar(&dirResources, "dirResources", os.Getenv("RESOURCES_DOCKER"), "the directory to serve public resources files from. Defaults to the current dir")
	flag.Parse()
	handlerResources := http.FileServer(FileSystem{http.Dir(dirResources)})
	r.PathPrefix(businesslogic.ResourcesPath).Handler(http.StripPrefix(businesslogic.ResourcesPath, handlerResources))
}
