package app

import (
	dbModels "github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/polygnosics/businesslogic"
	"github.com/artofimagination/polygnosics/initialization"
	"github.com/artofimagination/polygnosics/models"
	"github.com/artofimagination/polygnosics/rest"
	"github.com/artofimagination/polygnosics/rest/frontend"
	"github.com/artofimagination/polygnosics/rest/ipresolver"
	"github.com/artofimagination/polygnosics/rest/resourcesdb"
	"github.com/artofimagination/polygnosics/rest/userdb"
	"github.com/gorilla/mux"
)

type Context struct {
	RESTUserDB     *userdb.RESTController
	RESTFrontend   *frontend.RESTController
	RESTResouceDB  *resourcesdb.RESTController
	IPResolver     *ipresolver.IPResolver
	BackendContext *businesslogic.Context
	Config         *initialization.Config
	Router         *mux.Router
}

func NewContext() (*Context, error) {
	cfg := &initialization.Config{}
	initialization.InitConfig(cfg)
	rest.PrettyPrint(cfg)

	userdbController := userdb.NewRESTController()
	resourcedbController := resourcesdb.NewRESTController()

	uuidImpl := &dbModels.RepoUUID{}

	backend := &businesslogic.Context{
		UserDBController:      userdbController,
		ResourcesDBController: resourcedbController,
		DBModelFunctions: &dbModels.RepoFunctions{
			UUIDImpl: uuidImpl,
		},
		FileProcessor: &businesslogic.FileProcessorImpl{
			FileIO: businesslogic.IoImpl{},
			OsFunc: businesslogic.OsImpl{},
		},
		ResourceModelFunctions: &models.ResourceModelImpl{},
	}

	frontendController := frontend.NewRESTController(backend)
	ipresolver := ipresolver.NewIPResolver(frontendController, userdbController, resourcedbController, cfg)

	if err := ipresolver.DetectValidAddresses(); err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	frontendController.AddRouting(r)
	userdbController.AddRouting(r)
	resourcedbController.AddRouting(r)
	ipresolver.AddRouting(r)

	context := &Context{
		RESTFrontend:   frontendController,
		RESTUserDB:     userdbController,
		RESTResouceDB:  resourcedbController,
		IPResolver:     ipresolver,
		BackendContext: backend,
		Config:         cfg,
		Router:         r,
	}

	rest.PrettyPrint(context)
	return context, nil
}
