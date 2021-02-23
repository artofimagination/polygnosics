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
	r.HandleFunc("/index", c.MakeHandler(c.IndexHandler, r, true))
	r.HandleFunc("/", c.MakeHandler(c.IndexHandler, r, true))
	r.HandleFunc("/news", c.MakeHandler(c.News, r, true))
	resources := r.PathPrefix("/resources").Subrouter()
	resources.HandleFunc("/news", c.MakeHandler(c.News, r, true))
	resources.HandleFunc("/docs", c.MakeHandler(c.Docs, r, true))
	resources.HandleFunc("/tutorials", c.MakeHandler(c.Tutorials, r, true))
	resources.HandleFunc("/faq", c.MakeHandler(c.FAQ, r, true))
	resources.HandleFunc("/instructions", c.MakeHandler(c.Instructions, r, true))
	resources.HandleFunc("/examples", c.MakeHandler(c.Examples, r, true))
	resources.HandleFunc("/files", c.MakeHandler(c.Files, r, true))
	about := r.PathPrefix("/about").Subrouter()
	about.HandleFunc("/who-we-are", c.MakeHandler(c.About, r, true))
	about.HandleFunc("/contact", c.MakeHandler(c.Contact, r, true))
	about.HandleFunc("/career", c.MakeHandler(c.Career, r, true))

	resources.HandleFunc("/news", c.MakeHandler(c.News, r, true))
	resources.HandleFunc("/docs", c.MakeHandler(c.Docs, r, true))
	resources.HandleFunc("/tutorials", c.MakeHandler(c.Tutorials, r, true))
	resources.HandleFunc("/faq", c.MakeHandler(c.FAQ, r, true))
	resources.HandleFunc("/instructions", c.MakeHandler(c.Instructions, r, true))
	resources.HandleFunc("/examples", c.MakeHandler(c.Examples, r, true))

	// Authenticated pages
	r.HandleFunc("/auth_logout", c.MakeHandler(c.LogoutHandler, r, false))
	r.HandleFunc("/user-main", c.MakeHandler(c.UserMainHandler, r, false))
	r.HandleFunc("/user-settings", c.MakeHandler(UserSettings, r, false))
	userMain := r.PathPrefix("/user-main").Subrouter()
	userMain.HandleFunc("/upload-avatar", c.MakeHandler(c.UploadAvatarHandler, r, false))
	userMain.HandleFunc("/store", c.MakeHandler(c.StoreHandler, r, false))
	userMain.HandleFunc("/my-products", c.MakeHandler(c.MyProducts, r, false))
	userMain.HandleFunc("/my-projects", c.MakeHandler(c.MyProjects, r, false))
	userMain.HandleFunc("/project-browser", c.MakeHandler(c.BrowseProjects, r, false))
	userMain.HandleFunc("/product", c.MakeHandler(c.ProductDetails, r, false))
	userMain.HandleFunc("/product-wizard", c.MakeHandler(c.CreateProduct, r, false))
	userMain.HandleFunc("/profile", c.MakeHandler(c.ProfileHandler, r, false))
	userMain.HandleFunc("/profile-edit", c.MakeHandler(c.ProfileEdit, r, false))
	userMain.HandleFunc("/project-stats", c.MakeHandler(c.ProjectStats, r, false))
	userMain.HandleFunc("/product-stats", c.MakeHandler(c.ProductStats, r, false))
	userMain.HandleFunc("/user-stats", c.MakeHandler(c.UserStats, r, false))
	userMain.HandleFunc("/accounting", c.MakeHandler(c.AccountingStats, r, false))
	userMain.HandleFunc("/system-health", c.MakeHandler(c.SystemHealthStats, r, false))
	userMain.HandleFunc("/products-projects", c.MakeHandler(c.ProductsProjectsStats, r, false))
	userMain.HandleFunc("/ui-stats", c.MakeHandler(c.UIStats, r, false))
	userMain.HandleFunc("/misuse-metrics", c.MakeHandler(c.MisuseMetrics, r, false))
	userMain.HandleFunc("/mail-inbox", c.MakeHandler(c.MailInbox, r, false))
	userMain.HandleFunc("/mail-compose", c.MakeHandler(c.MailCompose, r, false))
	userMain.HandleFunc("/mail-read", c.MakeHandler(c.MailRead, r, false))
	stats := userMain.PathPrefix("/stats").Subrouter()
	stats.HandleFunc("/webrtc", c.MakeHandler(c.StatsWebRTC, r, false))
	myProducts := userMain.PathPrefix("/my-products").Subrouter()
	myProducts.HandleFunc("/details", c.MakeHandler(c.MyProductDetails, r, false))
	myProducts.HandleFunc("/delete", c.MakeHandler(c.DeleteProduct, r, false))
	myProducts.HandleFunc("/edit", c.MakeHandler(c.EditProduct, r, false))
	myProducts.HandleFunc("/project-wizard", c.MakeHandler(c.CreateProject, r, false))
	myProjects := userMain.PathPrefix("/my-projects").Subrouter()
	myProjects.HandleFunc("/details", c.MakeHandler(c.ProjectDetails, r, false))
	myProjects.HandleFunc("/delete", c.MakeHandler(c.DeleteProject, r, false))
	myProjects.HandleFunc("/edit", c.MakeHandler(c.EditProject, r, false))
	myProjects.HandleFunc("/run", c.MakeHandler(c.RunProject, r, false))
	myProjects.HandleFunc("/show", c.MakeHandler(c.ShowProject, r, false))

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
