package contents

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"polygnosics/app/businesslogic/project"
	"regexp"
	"strings"

	"github.com/artofimagination/mysql-user-db-go-interface/dbcontrollers"
	"github.com/artofimagination/mysql-user-db-go-interface/models"
	"github.com/pkg/errors"

	"github.com/google/uuid"
)

// TODO Issue#40: Replace with redis storage.

const (
	DefaultUserAvatarPath    = "/assets/images/avatar.jpg"
	DefaultProductAvatarPath = "/assets/images/avatar.jpg"
	DefaultProjectAvatarPath = "/assets/images/avatar.jpg"
)

const (
	UserAvatar     = "user_avatar"
	UserBackground = "user_background"

	ProductAvatar      = "product_avatar"
	ProductMainApp     = "main_app"
	ProductClientApp   = "client-app"
	ProductDescription = "product_description"
	ProductName        = "product_name"
	ProductVisibility  = "product_visibility"
	ProductPath        = "product_path"
	ProductFolder      = "product_folder"
	ProductRequires3D  = "requires_3d"
	ProductURL         = "product_url"
	ProductPublic      = "is_public"

	ProjectAvatar        = "project_avatar"
	ProjectPath          = "project_path"
	ProjectName          = "name"
	ProjectVisibility    = "visibility"
	ProjectServerLogging = "server_logging"
	ProjectClientLogging = "client_logging"
	NewProject           = "new_project"
	RunProject           = "run_project"
	ProjectState         = "project_state"
	ProjectStateColor    = "project_state_color"
	ProjectContainerID   = "project_container_id"
)

const (
	Product = iota
	User
)

const (
	Public    = "Public"
	Protected = "Protected"
	Private   = "Private"
)

var ErrFailedToParseForm = "Failed to parse form"

type ContentController struct {
	UserData         *models.UserData
	ProductData      *models.ProductData
	ProjectData      *models.ProjectData
	UserDBController *dbcontrollers.MYSQLController
}

func GetBooleanString(input string) string {
	if input == "" {
		return "No"
	}
	return input
}

func GetProjectStateColorString(state string) string {
	switch state {
	case project.NotRunning:
		return "#f5cf0a" // orange
	case project.Running:
		return "#00ff00" // green
	case project.Stopped:
		return "#ff0000" // red
	default:
		return "#e0dfd6" // lightgray
	}
}

func ValidateVisibility(value string) error {
	if value != Public && value != Protected && value != Private {
		return fmt.Errorf("Invalid visibility: %s", value)
	}
	return nil
}

func (c *ContentController) GetUserContent() map[string]interface{} {
	p := make(map[string]interface{})
	p["assets"] = make(map[string]interface{})
	path := c.UserDBController.ModelFunctions.GetFilePath(c.UserData.Assets, UserAvatar, DefaultUserAvatarPath)
	p["assets"].(map[string]interface{})[UserAvatar] = path
	p["texts"] = make(map[string]interface{})
	p["texts"].(map[string]interface{})["avatar-upload"] = "Upload your avatar"
	p["texts"].(map[string]interface{})["username"] = c.UserData.Name

	return p
}

func (c *ContentController) generateProductContent(productData *models.ProductData) map[string]interface{} {
	content := make(map[string]interface{})
	content[ProductAvatar] = c.UserDBController.ModelFunctions.GetFilePath(c.UserData.Assets, ProductAvatar, DefaultProductAvatarPath)
	content[ProductVisibility] = c.UserDBController.ModelFunctions.GetField(productData.Details, ProductVisibility, "")
	content[ProductName] = c.UserDBController.ModelFunctions.GetField(productData.Details, ProductName, "")
	content[ProductPublic] = c.UserDBController.ModelFunctions.GetField(productData.Details, ProductPublic, "")
	content[ProductDescription] = c.UserDBController.ModelFunctions.GetField(productData.Details, ProductDescription, "")
	content[ProductPath] = fmt.Sprintf("/user-main/my-products/details?product=%s", productData.ID.String())
	content[NewProject] = fmt.Sprintf("/user-main/my-products/new-project-wizard?product=%s", productData.ID.String())
	return content
}

func (c *ContentController) generateProjectContent(projectData *models.ProjectData) map[string]interface{} {
	content := make(map[string]interface{})
	content[ProjectAvatar] = c.UserDBController.ModelFunctions.GetFilePath(projectData.Assets, ProjectAvatar, DefaultProjectAvatarPath)
	content[ProjectName] = c.UserDBController.ModelFunctions.GetField(projectData.Details, ProjectName, "")
	content[ProjectVisibility] = c.UserDBController.ModelFunctions.GetField(projectData.Details, ProjectVisibility, "")
	content[ProjectContainerID] = c.UserDBController.ModelFunctions.GetField(projectData.Details, ProjectContainerID, "")
	content[ProjectPath] = fmt.Sprintf("/user-main/my-projects/details?project=%s", projectData.ID.String())
	content[ProjectState] = c.UserDBController.ModelFunctions.GetField(projectData.Details, ProjectState, "")
	content[ProjectStateColor] = GetProjectStateColorString(c.UserDBController.ModelFunctions.GetField(projectData.Details, ProjectState, ""))
	content[RunProject] = fmt.Sprintf("/user-main/my-projects/run?project=%s", projectData.ID.String())
	return content
}

func (c *ContentController) GetProductContent(productID *uuid.UUID) (map[string]interface{}, error) {
	product, err := c.UserDBController.GetProduct(productID)
	if err != nil {
		return nil, err
	}
	return c.generateProductContent(product), nil
}

func (c *ContentController) GetProjectContent(projectID *uuid.UUID) (map[string]interface{}, error) {
	project, err := c.UserDBController.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	return c.generateProjectContent(project), nil
}

func (c *ContentController) GetUserProductContent(userID *uuid.UUID) (map[string]interface{}, error) {
	products, err := c.UserDBController.GetProductsByUserID(userID)
	if err != nil {
		return nil, err
	}

	p := make(map[string]interface{})

	productContent := make([]map[string]interface{}, len(products))
	for i, product := range products {
		productContent[i] = c.generateProductContent(&product.ProductData)
	}
	p["product"] = productContent

	return p, nil
}

func (c *ContentController) GetUserProjectContent(userID *uuid.UUID) (map[string]interface{}, error) {
	projects, err := c.UserDBController.GetProjectsByUserID(userID)
	if err != nil {
		return nil, err
	}

	p := make(map[string]interface{})

	projectContent := make([]map[string]interface{}, len(projects))
	for i, project := range projects {
		projectContent[i] = c.generateProjectContent(project.ProjectData)
	}
	p["project"] = projectContent

	return p, nil
}

func (c *ContentController) UploadFile(asset *models.Asset, fileType string, defaultPath string, formName string, r *http.Request) error {
	file, handler, err := r.FormFile(formName)
	if err == http.ErrMissingFile {
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	defer file.Close()

	if err := c.UserDBController.ModelFunctions.SetFilePath(asset, fileType, filepath.Ext(handler.Filename)); err != nil {
		return err
	}
	path := c.UserDBController.ModelFunctions.GetFilePath(asset, fileType, defaultPath)

	// Create file
	dst, err := os.Create(path)
	if err != nil {
		if err2 := dst.Close(); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		if err2 := c.UserDBController.ModelFunctions.ClearAsset(asset, fileType); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return err
	}

	// Copy the uploaded file to the created file on the file system.
	if _, err := io.Copy(dst, file); err != nil {
		if err2 := dst.Close(); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		if err2 := c.UserDBController.ModelFunctions.ClearAsset(asset, fileType); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return err
	}

	return nil
}

var splitRegexp = regexp.MustCompile(`(\S{4})`)

func (*ContentController) GeneratePath(assetID *uuid.UUID) (string, error) {
	assetIDString := strings.Replace(assetID.String(), "-", "", -1)
	assetStringSplit := splitRegexp.FindAllString(assetIDString, -1)
	assetPath := path.Join(assetStringSplit...)
	rootPath := os.Getenv("USER_STORE_DOCKER")
	assetPath = path.Join(rootPath, assetPath)
	if err := os.MkdirAll(assetPath, os.ModePerm); err != nil {
		return "", err
	}
	return assetPath, nil
}
