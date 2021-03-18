package businesslogic

import (
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/polygnosics/rest/userdb"
)

type Context struct {
	UserDBController *userdb.RESTController
	ModelFunctions   *models.RepoFunctions
}
