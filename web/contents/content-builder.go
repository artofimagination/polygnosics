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
	ProjectsPageBrowserName    = "Browser"
	ProjectsPageEditName       = "Edit"
	ProjectsPageRunName        = "Run"
)

const (
	UserPageName         = "User"
	UserPageProfileName  = "Profile"
	UserPageMainPageName = "Info board"
)

const (
	ResourcesPageName              = "Resources"
	ResourcesPageNewsName          = "News"
	ResourcesPageDocumentationName = "Documentation"
	ResourcesPageTutorialsName     = "Tutorials"
	ResourcesPageFAQsName          = "FAQs"
	ResourcesPageFilesName         = "Files"
)

const (
	StatsPageName     = "Statistics"
	StatsUsers        = "User Statistics"
	StatsItems        = "Overall Product and Project Stats"
	StatsProducts     = "Product Statistics"
	StatsProjects     = "Project Statistics"
	StatsAccounting   = "Accounting Statistics"
	StatsUI           = "UI Statistics"
	StatsSystemHealth = "System Health"
)

var FutureFeature = "future_feature"

// TODO Issue#40: Replace  user/product/project data with redis storage.
type ContentController struct {
	UserData         *models.UserData
	ProductData      *models.ProductData
	ProjectData      *models.ProjectData
	UserDBController *dbcontrollers.MYSQLController
}

func convertToCheckboxValue(input string) string {
	if input == businesslogic.CheckBoxUnChecked {
		return ""
	}
	return input
}

func convertCheckedToYesNo(input string) string {
	if input == businesslogic.CheckBoxUnChecked {
		return "No"
	}
	return "Yes"
}

func (c *ContentController) BuildProductWizardContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProductsPageName, ProductsPageCreateName)
	content["categories"] = businesslogic.CreateCategoriesMap()
	return content
}

func (c *ContentController) BuildProjectWizardContent(productID *uuid.UUID) (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProjectsPageName, ProjectsPageCreateName)
	product, err := c.GetProductContent(productID)
	if err != nil {
		return nil, err
	}
	content[ProductMapKey] = product
	return content, nil
}

func (c *ContentController) BuildProductEditContent(productID *uuid.UUID) (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProductsPageName, ProductsPageEditName)
	productContent, err := c.GetProductContent(productID)
	if err != nil {
		return nil, err
	}
	content[ProductMapKey] = productContent
	content["categories"] = businesslogic.CreateCategoriesMap()
	return content, err
}

func (c *ContentController) BuildProjectEditContent(projectID *uuid.UUID) (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProjectsPageName, ProjectsPageEditName)
	projectContent, err := c.GetProjectContent(projectID)
	if err != nil {
		return nil, err
	}
	content[ProjectMapKey] = projectContent
	content["categories"] = businesslogic.CreateCategoriesMap()

	return content, err
}

func (c *ContentController) BuildProjectRunContent(projectID *uuid.UUID) (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProjectsPageName, ProjectsPageRunName)
	projectContent, err := c.GetProjectContent(projectID)
	if err != nil {
		return nil, err
	}
	content[ProjectMapKey] = projectContent
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
	content[ProductMapKey] = productContent
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
	projectContent, err := c.GetProjectContent(projectID)
	if err != nil {
		return nil, err
	}
	content[ProjectMapKey] = projectContent
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

func (c *ContentController) BuildUserStatsContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, StatsPageName, StatsUsers)
	return content
}

func (c *ContentController) BuildProductStatsContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, StatsPageName, StatsProducts)
	return content
}

func (c *ContentController) BuildProjectStatsContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, StatsPageName, StatsProjects)
	return content
}

func (c *ContentController) BuildAccountingStatsContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, StatsPageName, StatsAccounting)
	return content
}

func (c *ContentController) BuildSystemHealthContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, StatsPageName, StatsSystemHealth)
	return content
}

func (c *ContentController) BuildUIStatsContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, StatsPageName, StatsUI)
	return content
}

func (c *ContentController) BuildItemStatsContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, StatsPageName, StatsItems)
	return content
}

func (c *ContentController) BuildMailInboxContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	return content
}

func (c *ContentController) BuildMailComposeContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	return content
}

func (c *ContentController) BuildDocsContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ResourcesPageName, ResourcesPageDocumentationName)
	return content
}

func (c *ContentController) BuildFilesContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ResourcesPageName, ResourcesPageFilesName)
	return content
}

func (c *ContentController) BuildTutorialsContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ResourcesPageName, ResourcesPageDocumentationName)
	return content
}

func (c *ContentController) BuildFAQContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ResourcesPageName, ResourcesPageFAQsName)
	return content
}

func (c *ContentController) BuildMailReadContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	return content
}

func (c *ContentController) BuildContactContent() map[string]interface{} {
	content := c.GetUserContent(c.UserData)
	return content
}

func (c *ContentController) BuildStoreContent() (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProductsPageName, ProductsPageStoreName)
	categorizedProducts, err := c.GetProductsByCategory(&c.UserData.ID)
	if err != nil {
		return nil, err
	}
	recent, err := c.GetRecentProductsContent(&c.UserData.ID)
	if err != nil {
		return nil, err
	}
	recommended, err := c.GetRecommendedProductsContent(&c.UserData.ID)
	if err != nil {
		return nil, err
	}
	content[businesslogic.ProductCategoriesKey] = categorizedProducts
	content["recent"] = recent
	content["recommended"] = recommended
	return content, nil
}

func (c *ContentController) BuildProjectBrowserContent() (map[string]interface{}, error) {
	content := c.GetUserContent(c.UserData)
	content = c.prepareContentHeader(content, ProductsPageName, ProjectsPageBrowserName)
	categorizedProducts, err := c.GetProjectsByCategory(&c.UserData.ID)
	if err != nil {
		return nil, err
	}
	recent, err := c.GetRecentProjectsContent(&c.UserData.ID)
	if err != nil {
		return nil, err
	}
	recommended, err := c.GetRecommendedProjectsContent(&c.UserData.ID)
	if err != nil {
		return nil, err
	}

	content[businesslogic.ProductCategoriesKey] = categorizedProducts
	content["recent"] = recent
	content["recommended"] = recommended
	return content, nil
}
