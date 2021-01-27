package contents

import (
	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
)

// TODO Issue#40: Replace  user/product/project data with redis storage.
type ContentController struct {
	UserData         *models.UserData
	ProductData      *models.ProductData
	ProjectData      *models.ProjectData
	UserDBController *dbcontrollers.MYSQLController
}

// getBooleanString returns a check box stat Yes/No string
func getBooleanString(input string) string {
	if input == "" {
		return "No"
	}
	return input
}
