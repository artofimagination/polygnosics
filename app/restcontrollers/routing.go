package restcontrollers

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"polygnosics/app/restcontrollers/auth"
	"polygnosics/app/restcontrollers/page"

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
func CreateRouter() *mux.Router {
	r := mux.NewRouter()
	// Publicly accessable pages
	r.HandleFunc("/auth_signup", page.MakeHandler(auth.SignupHandler, r, true))
	r.HandleFunc("/auth_login", page.MakeHandler(auth.LoginHandler, r, true))
	r.HandleFunc("/about", page.MakeHandler(AboutUsHandler, r, true))
	r.HandleFunc("/index", page.MakeHandler(IndexHandler, r, true))
	r.HandleFunc("/", page.MakeHandler(IndexHandler, r, true))

	// Authenticated pages
	r.HandleFunc("/auth_logout", page.MakeHandler(auth.LogoutHandler, r, false))
	r.HandleFunc("/user-main", page.MakeHandler(UserMainHandler, r, false))
	r.HandleFunc("/user-settings", page.MakeHandler(UserSettings, r, false))
	userMain := r.PathPrefix("/user-main").Subrouter()
	userMain.HandleFunc("/upload-avatar", page.MakeHandler(UploadAvatarHandler, r, false))
	userMain.HandleFunc("/store", page.MakeHandler(StoreHandler, r, false))
	userMain.HandleFunc("/my-products", page.MakeHandler(MyProductsHandler, r, false))
	userMain.HandleFunc("/new-product-wizard", page.MakeHandler(CreateProduct, r, false))
	userMain.HandleFunc("/profile", page.MakeHandler(ProfileHandler, r, false))
	userMain.HandleFunc("/new", page.MakeHandler(NewProject, r, false))
	userMain.HandleFunc("/{project}/run", page.MakeHandler(RunProject, r, false))
	userMain.HandleFunc("/{project}/webrtc", page.MakeHandler(StartWebRTC, r, false))
	userMain.HandleFunc("/resume", page.MakeHandler(NewProject, r, false))

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
