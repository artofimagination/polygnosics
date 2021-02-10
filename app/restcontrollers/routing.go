package restcontrollers

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var notFoundFile, notFoundErr = http.Dir("dummy").Open("does-not-exist")

type FileSystem struct {
	http.Dir
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

// CreateRouter creates the page path structure.
func CreateRouter(c *RESTController) *mux.Router {
	r := mux.NewRouter()
	// Publicly accessable pages
	r.HandleFunc("/auth_signup", c.MakeHandler(c.SignupHandler, r, true))
	r.HandleFunc("/auth_login", c.MakeHandler(c.LoginHandler, r, true))
	r.HandleFunc("/about", c.MakeHandler(c.AboutUsHandler, r, true))
	r.HandleFunc("/index", c.MakeHandler(c.IndexHandler, r, true))
	r.HandleFunc("/", c.MakeHandler(c.IndexHandler, r, true))
	r.HandleFunc("/news", c.MakeHandler(c.News, r, true))

	// Authenticated pages
	r.HandleFunc("/auth_logout", c.MakeHandler(c.LogoutHandler, r, false))
	r.HandleFunc("/user-main", c.MakeHandler(c.UserMainHandler, r, false))
	r.HandleFunc("/user-settings", c.MakeHandler(UserSettings, r, false))
	resources := r.PathPrefix("/resources").Subrouter()
	resources.HandleFunc("/news", c.MakeHandler(c.News, r, false))
	userMain := r.PathPrefix("/user-main").Subrouter()
	userMain.HandleFunc("/upload-avatar", c.MakeHandler(c.UploadAvatarHandler, r, false))
	userMain.HandleFunc("/store", c.MakeHandler(c.StoreHandler, r, false))
	userMain.HandleFunc("/my-products", c.MakeHandler(c.MyProducts, r, false))
	userMain.HandleFunc("/my-projects", c.MakeHandler(c.MyProjects, r, false))
	userMain.HandleFunc("/product", c.MakeHandler(c.ProductDetails, r, false))
	userMain.HandleFunc("/product-wizard", c.MakeHandler(c.CreateProduct, r, false))
	userMain.HandleFunc("/profile", c.MakeHandler(c.ProfileHandler, r, false))
	userMain.HandleFunc("/profile-edit", c.MakeHandler(c.ProfileEdit, r, false))
	myProducts := userMain.PathPrefix("/my-products").Subrouter()
	myProducts.HandleFunc("/details", c.MakeHandler(c.MyProductDetails, r, false))
	myProducts.HandleFunc("/delete", c.MakeHandler(c.DeleteProduct, r, false))
	myProducts.HandleFunc("/edit", c.MakeHandler(c.EditProduct, r, false))
	myProducts.HandleFunc("/new-project-wizard", c.MakeHandler(c.CreateProject, r, false))
	myProjects := userMain.PathPrefix("/my-projects").Subrouter()
	myProjects.HandleFunc("/details", c.MakeHandler(c.ProjectDetails, r, false))
	myProjects.HandleFunc("/run", c.MakeHandler(c.RunProject, r, false))

	// Static file servers
	// Default web assets
	var dirDefaultAssets string
	var dirUserAssets string
	var dirTemplates string
	flag.StringVar(&dirDefaultAssets, "dirDefaultAssets", "./web/assets", "the directory to serve default web assets from. Defaults to the current dir")
	handlerDefaultAssets := http.FileServer(FileSystem{http.Dir(dirDefaultAssets)})
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", handlerDefaultAssets))
	flag.StringVar(&dirTemplates, "dirTemplates", "./web/templates", "the directory to serve default web assets from. Defaults to the current dir")
	handlerTemplates := http.FileServer(FileSystem{http.Dir(dirTemplates)})
	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", handlerTemplates))
	flag.StringVar(&dirUserAssets, "dirUserAssets", os.Getenv("USER_STORE_DOCKER"), "the directory to serve user asset files from. Defaults to the current dir")
	flag.Parse()
	handlerUserAssets := http.FileServer(FileSystem{http.Dir(dirUserAssets)})
	r.PathPrefix("/user-assets/").Handler(http.StripPrefix("/user-assets/", handlerUserAssets))

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	http.Handle("/", r)

	return r
}
