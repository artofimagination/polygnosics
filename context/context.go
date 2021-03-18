package app

import (
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/polygnosics/businesslogic"
	"github.com/artofimagination/polygnosics/rest/frontend"
	"github.com/artofimagination/polygnosics/rest/userdb"
)

type Context struct {
	RESTUserDB     *userdb.RESTController
	RESTFrontend   *frontend.RESTController
	BackendContext *businesslogic.Context
}

func NewContext() (*Context, error) {
	userdbController := userdb.NewRESTController()

	uuidImpl := &models.RepoUUID{}

	backend := &businesslogic.Context{
		UserDBController: userdbController,
		ModelFunctions: &models.RepoFunctions{
			UUIDImpl: uuidImpl,
		},
	}

	context := &Context{
		RESTFrontend:   frontend.NewRESTController(backend),
		RESTUserDB:     userdb.NewRESTController(),
		BackendContext: backend,
	}

	return context, nil
}
