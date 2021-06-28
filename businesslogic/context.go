package businesslogic

import (
	dbModels "github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/artofimagination/polygnosics/models"
	"github.com/artofimagination/polygnosics/rest/resourcesdb"
	"github.com/artofimagination/polygnosics/rest/userdb"
)

type Context struct {
	UserDBController       *userdb.RESTController
	ResourcesDBController  resourcesdb.ResourceDBInterface
	DBModelFunctions       *dbModels.RepoFunctions
	FileProcessor          FileProcessor
	ResourceModelFunctions models.ResourceModels
}
