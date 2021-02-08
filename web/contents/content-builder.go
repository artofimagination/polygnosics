package contents

import (
	"polygnosics/app/businesslogic"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/google/uuid"
)

const (
	ProductsPageName           = "Products"
	ProductsPageCreateName     = "Product Wizard"
	ProductsPageMyProductsName = "My Products"
	ProductsPageDetailsName    = "Details"
	ProductsPageStoreName      = "Marketplace"
	ProductsPageEditName       = "Edit"
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

const (
	ResourcesPageName     = "Resources"
	ResourcesPageNewsName = "News"
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
	if input == "" || input == businesslogic.CheckBoxUnChecked {
		return businesslogic.CheckBoxUnChecked
	}
	return businesslogic.CheckBoxChecked
}

func convertToCheckboxValue(input string) string {
	if input == businesslogic.CheckBoxUnChecked {
		return ""
	}
	return input
}

func (c *ContentController) BuildProductWizardContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProductsPageName, ProductsPageCreateName)
	return content
}

func (c *ContentController) BuildProductEditContent(productID *uuid.UUID) (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProductsPageName, ProductsPageEditName)
	productContent, err := c.GetProductContent(productID)
	if err != nil {
		return nil, err
	}
	content[ProductMapKey] = productContent[ProductMapKey]
	return content, err
}

func (c *ContentController) BuildMyProductsContent() (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProductsPageName, ProductsPageMyProductsName)
	productsContent, err := c.GetUserProductContent(&c.UserData.ID)
	if err != nil {
		return nil, err
	}
	content[ProductMapKey] = productsContent
	return content, nil
}

func (c *ContentController) BuildProductDetailsContent(productID *uuid.UUID) (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProductsPageName, ProductsPageDetailsName)
	productContent, err := c.GetProductContent(productID)
	if err != nil {
		return nil, err
	}
	content[ProductMapKey] = productContent[ProductMapKey]
	return content, nil
}

func (c *ContentController) BuildMyProjectsContent() (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProjectsPageName, ProjectsPageMyProjectsName)
	projectsContent, err := c.GetUserProjectContent(&c.UserData.ID, -1)
	if err != nil {
		return nil, err
	}
	content["project"] = projectsContent
	return content, nil
}

func (c *ContentController) BuildProjectDetailsContent(projectID *uuid.UUID) (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
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

func (c *ContentController) BuildProfileContent(id *uuid.UUID) (map[string]interface{}, error) {
	user, err := c.UserDBController.GetUser(id)
	if err != nil {
		return nil, err
	}
	content := c.GetUserContent(user)
	content = c.prepareContentHeader(content, UserPageName, UserPageProfileName)
	return content, err
}

func (c *ContentController) BuildUserMainContent() (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, UserPageName, UserPageMainPageName)
	content = c.prepareNewsFeed(content)
	productsContent, err := c.GetRecentProductsContent(&c.UserData.ID)
	if err != nil {
		return nil, err
	}
	content["product"] = productsContent
	projectsContent, err := c.GetUserProjectContent(&c.UserData.ID, 4)
	if err != nil {
		return nil, err
	}
	content["project"] = projectsContent
	return content, nil
}

func (c *ContentController) BuildErrorContent(errString string) map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content["message"] = errString
	return content
}

func (c *ContentController) BuildNewsContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ResourcesPageName, ResourcesPageNewsName)
	return content
}

func (c *ContentController) BuildStoreContent() (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProductsPageName, ProductsPageStoreName)
	productsContent, err := c.GetUserProductContent(&c.UserData.ID)
	if err != nil {
		return nil, err
	}
	content["product"] = productsContent
	return content, nil
}
