package businesslogic

import (
	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
)

type Context struct {
	UserDBController *dbcontrollers.MYSQLController
}
