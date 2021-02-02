package contents

import (
	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

const (
	ProductsPageName           = "Products"
	ProductsPageCreateName     = "Product Wizard"
	ProductsPageMyProductsName = "My Products"
	ProductsPageDetailsName    = "Details"
)

const (
	ProjectsPageName           = "Projects"
	ProjectsPageCreateName     = "Project Wizard"
	ProjectsPageMyProjectsName = "My Projects"
	ProjectsPageDetailsName    = "Details"
)

const (
	UserPageName         = "User"
	UserPageProfileName  = "Profile"
	UserPageMainPageName = "Info board"
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

func (c *ContentController) BuildProductWizardContent() map[string]interface{} {
	content := c.GetUserContent()
	content = c.prepareContentHeader(content, ProductsPageName, ProductsPageCreateName)
	return content
}

func (c *ContentController) BuildMyProductsContent() (map[string]interface{}, error) {
	content := c.GetUserContent()
	content = c.prepareContentHeader(content, ProductsPageName, ProductsPageMyProductsName)
	productsContent, err := c.GetUserProductContent(&c.UserData.ID)
	if err != nil {
		return nil, err
	}
	for k, v := range productsContent {
		content[k] = v
	}
	return content, nil
}

func (c *ContentController) BuildProductDetailsContent(productID *uuid.UUID) (map[string]interface{}, error) {
	content := c.GetUserContent()
	content = c.prepareContentHeader(content, ProductsPageName, ProductsPageDetailsName)
	productContent, err := c.GetProductContent(productID)
	if err != nil {
		return nil, err
	}
	for k, v := range productContent {
		content[k] = v
	}
	return content, nil
}

func (c *ContentController) BuildMyProjectsContent() (map[string]interface{}, error) {
	content := c.GetUserContent()
	content = c.prepareContentHeader(content, ProjectsPageName, ProjectsPageMyProjectsName)
	productsContent, err := c.GetUserProjectContent(&c.UserData.ID)
	if err != nil {
		return nil, err
	}
	for k, v := range productsContent {
		content[k] = v
	}
	return content, nil
}

func (c *ContentController) BuildProjectDetailsContent(projectID *uuid.UUID) (map[string]interface{}, error) {
	content := c.GetUserContent()
	content = c.prepareContentHeader(content, ProjectsPageName, ProjectsPageDetailsName)
	productContent, err := c.GetProjectContent(projectID)
	if err != nil {
		return nil, err
	}
	for k, v := range productContent {
		content[k] = v
	}
	return content, nil
}

func (c *ContentController) BuildProfileContent() map[string]interface{} {
	content := c.GetUserContent()
	content = c.prepareContentHeader(content, UserPageName, UserPageProfileName)
	return content
}

func (c *ContentController) BuildUserMainContent() map[string]interface{} {
	content := c.GetUserContent()
	content = c.prepareContentHeader(content, UserPageName, UserPageMainPageName)
	return content
}

func (c *ContentController) BuildErrorContent(errString string) map[string]interface{} {
	content := c.GetUserContent()
	content["message"] = errString
	return content
}
