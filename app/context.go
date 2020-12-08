package app

import (
	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
)

var ContextData *Context

type Context struct {
	UserDBController    *dbcontrollers.MYSQLController
	ProjectDBController *dbcontrollers.ProjectDBDummy
}

func NewContext() (*Context, error) {
	userDBController, err := dbcontrollers.NewDBController()
	if err != nil {
		return nil, err
	}

	context := &Context{
		UserDBController:    userDBController,
		ProjectDBController: &dbcontrollers.ProjectDBDummy{},
	}

	dbcontrollers.SetProjectDB(context.ProjectDBController)
	return context, nil
}
