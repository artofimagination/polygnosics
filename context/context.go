package app

import (
	dbModels "github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/polygnosics/businesslogic"
	"github.com/artofimagination/polygnosics/models"
	"github.com/artofimagination/polygnosics/rest/frontend"
	"github.com/artofimagination/polygnosics/rest/resourcesdb"
	"github.com/artofimagination/polygnosics/rest/userdb"
)

type Context struct {
	RESTUserDB     *userdb.RESTController
	RESTFrontend   *frontend.RESTController
	BackendContext *businesslogic.Context
}

func NewContext() (*Context, error) {
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

	context := &Context{
		RESTFrontend:   frontend.NewRESTController(backend),
		RESTUserDB:     userdb.NewRESTController(),
		BackendContext: backend,
	}

	return context, nil
}
