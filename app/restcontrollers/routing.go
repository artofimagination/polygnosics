package restcontrollers

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"polygnosics/app/restcontrollers/auth"
	"polygnosics/app/utils/page"

	"github.com/gorilla/mux"
)

// CreateRouter creates the page path structure.
func CreateRouter() *mux.Router {
	var dir string
	flag.StringVar(&dir, "dir", "./web/assets", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/auth_signup", page.MakeHandler(auth.SignupHandler, r))
	r.HandleFunc("/auth_login", page.MakeHandler(auth.LoginHandler, r))
	r.HandleFunc("/logout", page.MakeHandler(IndexHandler, r))
	r.HandleFunc("/user-main", page.MakeHandler(UserMain, r))
	r.HandleFunc("/user-settings", page.MakeHandler(UserSettings, r))
	userMain := r.PathPrefix("/user-main").Subrouter()
	userMain.HandleFunc("/new", page.MakeHandler(NewProject, r))
	userMain.HandleFunc("/{project}/run", page.MakeHandler(RunProject, r))
	userMain.HandleFunc("/{project}/webrtc", page.MakeHandler(StartWebRTC, r))
	userMain.HandleFunc("/resume", page.MakeHandler(NewProject, r))
	r.HandleFunc("/about", page.MakeHandler(AboutUsHandler, r))
	r.HandleFunc("/index", page.MakeHandler(IndexHandler, r))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(dir))))

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
