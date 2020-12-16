package app

import (
	"polygnosics/app/restcontrollers"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
)

type Context struct {
	UserDBController    *dbcontrollers.MYSQLController
	ProjectDBController *dbcontrollers.ProjectDBDummy
	RESTController      *restcontrollers.RESTController
}

func NewContext() (*Context, error) {
	userDBController, err := dbcontrollers.NewDBController()
	if err != nil {
		return nil, err
	}

	context := &Context{
		UserDBController:    userDBController,
		ProjectDBController: &dbcontrollers.ProjectDBDummy{},
		RESTController:      restcontrollers.NewRESTController(userDBController),
	}

	dbcontrollers.SetProjectDB(context.ProjectDBController)
	return context, nil
}
