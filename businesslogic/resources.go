package businesslogic

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/artofimagination/polygnosics/models"
	"github.com/artofimagination/polygnosics/rest"
	"github.com/google/uuid"
)

const ResourcesPath = "/resources/"
const TutorialPath = "/tutorials/"
const FilesPath = "/files/"
const FAQPath = "/faq/"

const (
	FileSectionFileType = "file"
	FileSectionLinkType = "link"
)

const (
	CategoryFAQ         = "FAQ"
	CategoryFAQGroups   = "FAQGroups"
	CategoryTutorial    = "Tutorial"
	CategoryNews        = "NewsFeed"
	CategoryFileContent = "FileContent"
	CategoryFileSection = "FilesSection"
)

var MaxFileAttachements = 50

func getDataModel(category string) (interface{}, error) {
	switch category {
	case CategoryFAQ:
		return &models.FAQ{}, nil
	case CategoryTutorial:
		return &models.Tutorial{}, nil
	case CategoryFileSection:
		return &models.FilesSection{}, nil
	case CategoryFileContent:
		return &models.FileContent{}, nil
	case CategoryNews:
		return &models.NewsEntry{}, nil
	default:
		return nil, fmt.Errorf("Invalid resource category: %s", category)
	}
}

func (c *Context) getCategoryID(categoryName string) (int, error) {
	categoryID := -1
	categories, err := c.ResourcesDBController.GetCategories()
	if err != nil {
		return categoryID, err
	}

	for _, category := range categories {
		if category.Name == categoryName {
			categoryID = category.ID
			return categoryID, err
		}
	}

	return categoryID, fmt.Errorf("%s category not found", categoryName)
}

func (c *Context) UpdateHandler(category string, r rest.RequestInterface, handler func(rest.RequestInterface, ...interface{}) error, parameters ...interface{}) error {
	dataModel, err := getDataModel(category)
	if err != nil {
		return err
	}
	resource, err := c.ResourcesDBController.GetResource(parameters[0].(string), dataModel)
	if err != nil {
		return err
	}

	if err := handler(r, dataModel, parameters); err != nil {
		return err
	}

	if err := c.ResourcesDBController.UpdateResource(resource, dataModel); err != nil {
		return err
	}

	return nil
}

func (c *Context) AddHandler(category string, r rest.RequestInterface, handler func(rest.RequestInterface, ...interface{}) error, parameters ...interface{}) (*uuid.UUID, error) {
	dataModel, err := getDataModel(category)
	if err != nil {
		return nil, err
	}

	if err := handler(r, dataModel, parameters); err != nil {
		return nil, err
	}

	categoryID, err := c.getCategoryID(category)
	if err != nil {
		return nil, err
	}

	id, err := c.ResourcesDBController.AddResource(categoryID, dataModel)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (c *Context) DeleteHandler(category string, r rest.RequestInterface, handler func(rest.RequestInterface, ...interface{}) error, parameters ...interface{}) error {
	dataModel, err := getDataModel(category)
	if err != nil {
		return err
	}

	if _, err := c.ResourcesDBController.GetResource(parameters[0].(string), dataModel); err != nil {
		return err
	}

	if err := handler(r, dataModel); err != nil {
		return err
	}

	if err := c.ResourcesDBController.DeleteResource(parameters[0].(string)); err != nil {
		return err
	}

	return nil
}

func (c *Context) setTutorialArticle(tutorial *models.Tutorial, r rest.RequestInterface) error {
	if r.FormValue("article") == "" {
		if tutorial.Content != "" {
			if err := c.FileProcessor.RemoveFile(tutorial.Content); err != nil {
				return err
			}
		}
		tutorial.Content = ""
		return nil
	}

	if tutorial.Content == "" && r.FormValue("article") != "" {
		articleFileName := fmt.Sprintf("%s.%s", c.FileProcessor.GenerateID(), "txt")
		tutorial.Content = filepath.Join(ResourcesPath, TutorialPath, articleFileName)
	}
	if err := c.FileProcessor.WriteToFile(tutorial.Content, r.FormValue("article")); err != nil {
		return err
	}

	return nil
}

func (c *Context) setTutorialAvatar(tutorial *models.Tutorial, r rest.RequestInterface) error {
	newType := r.FormValue("avatar_type")
	if newType == "image" {
		if tutorial.AvatarType != newType {
			avatarFileName := fmt.Sprintf("%s%s", c.FileProcessor.GenerateID(), filepath.Ext(r.FormValue("avatar")))
			tutorial.AvatarSource = filepath.Join(ResourcesPath, TutorialPath, avatarFileName)
		}

		if err := c.FileProcessor.UploadFile("avatar_image", tutorial.AvatarSource, r); err != nil {
			return err
		}
	} else {
		if tutorial.AvatarType != newType {
			if err := c.FileProcessor.RemoveFile(tutorial.AvatarSource); err != nil {
				return err
			}
		}
		tutorial.AvatarSource = r.FormValue("avatar_video")
	}
	tutorial.AvatarType = newType
	return nil
}

func (c *Context) UpdateTutorial(r rest.RequestInterface, input ...interface{}) error {
	tutorial := input[0].(*models.Tutorial)

	if err := c.setTutorialArticle(tutorial, r); err != nil {
		return err
	}

	if err := c.setTutorialAvatar(tutorial, r); err != nil {
		return err
	}

	tutorial.Title = r.FormValue("title")
	tutorial.ShortDesc = r.FormValue("short")
	tutorial.LastUpdated = time.Now()

	return nil
}

func (c *Context) AddTutorial(r rest.RequestInterface, input ...interface{}) error {
	tutorial := input[0].(*models.Tutorial)

	if err := c.setTutorialArticle(tutorial, r); err != nil {
		return err
	}

	if err := c.setTutorialAvatar(tutorial, r); err != nil {
		return err
	}
	tutorial.Title = r.FormValue("title")
	tutorial.ShortDesc = r.FormValue("short")
	tutorial.LastUpdated = time.Now()

	return nil
}

func (c *Context) DeleteTutorial(r rest.RequestInterface, input ...interface{}) error {
	id := input[0].(string)
	if err := c.DeleteHandler(CategoryTutorial, r, func(r rest.RequestInterface, parameters ...interface{}) error {
		tutorial := parameters[0].(*models.Tutorial)
		if tutorial.AvatarType == "image" {
			if err := c.FileProcessor.RemoveFile(tutorial.AvatarSource); err != nil {
				return err
			}
		}
		if tutorial.Content != "" {
			if err := c.FileProcessor.RemoveFile(tutorial.Content); err != nil {
				return err
			}
		}
		return nil
	}, id); err != nil {
		return err
	}

	return nil
}

func (c *Context) setFileItem(r rest.RequestInterface, input ...interface{}) error {
	fileItem := input[0].(*models.FileContent)
	index := input[1].([]interface{})[0].(string)
	if index == "" {
		index = input[1].([]interface{})[1].(string)
	}

	fileItem.Type = r.FormValue(fmt.Sprintf("type_%s", index))
	fileItem.RefName = r.FormValue(fmt.Sprintf("ref_name_%s", index))

	if fileItem.Type == FileSectionFileType {
		formName := fmt.Sprintf("upload_file_%s", index)
		_, handler, err := r.FormFile(formName)
		if err != nil && err == http.ErrMissingFile {
			return nil
		}
		fileItem.OriginalFileName = handler.Filename

		fileName := fmt.Sprintf("%s%s", c.FileProcessor.GenerateID(), filepath.Ext(fileItem.OriginalFileName))
		fileItem.Ref = filepath.Join(ResourcesPath, FilesPath, fileName)
		if err := c.FileProcessor.UploadFile(formName, fileItem.Ref, r); err != nil {
			return err
		}
	} else {
		fileItem.Ref = r.FormValue(fmt.Sprintf("repo_link_%s", index))
	}
	return nil
}

func (c *Context) AddFileSection(r rest.RequestInterface, input ...interface{}) error {
	fileSection := input[0].(*models.FilesSection)
	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		return err
	}

	if count > MaxFileAttachements {
		count = MaxFileAttachements
	}

	fileSection.Title = r.FormValue("title")
	fileSection.ShortDescription = r.FormValue("short")
	fileSection.ContentIDList = make([]string, count)
	for i := 0; i < count; i++ {
		resourceID, err := c.AddHandler(CategoryFileContent, r, c.setFileItem, "", fmt.Sprintf("%d", i))
		if err != nil {
			return err
		}
		fileSection.ContentIDList[i] = resourceID.String()
	}

	return nil
}

func (c *Context) removeFileItem(id string, r rest.RequestInterface, fileSectionContentIDs []string) ([]string, error) {
	formName := fmt.Sprintf("upload_file_%s", id)
	fileName := fmt.Sprintf("%s%s", id, filepath.Ext(r.FormValue(formName)))

	if err := c.FileProcessor.RemoveFile(fileName); err != nil {
		return nil, err
	}

	if err := c.ResourcesDBController.DeleteResource(id); err != nil {
		return nil, err
	}

	newContentList := make([]string, len(fileSectionContentIDs)-1)
	newIndex := 0
	for i, contentID := range fileSectionContentIDs {
		if id != contentID {
			newContentList[newIndex] = fileSectionContentIDs[i]
			newIndex++
		}
	}

	return newContentList, nil
}

func (c *Context) UpdateFileSection(r rest.RequestInterface, input ...interface{}) error {
	fileSection := input[0].(*models.FilesSection)
	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		return err
	}

	if count > MaxFileAttachements {
		count = MaxFileAttachements
	}

	fileSection.Title = r.FormValue("title")
	fileSection.ShortDescription = r.FormValue("short")

	for _, resourceIDString := range fileSection.ContentIDList {
		resourceIDString := resourceIDString
		if r.FormValue(fmt.Sprintf("remove_%s", resourceIDString)) == "checked" {
			newContentList, err := c.removeFileItem(resourceIDString, r, fileSection.ContentIDList)
			if err != nil {
				return err
			}
			fileSection.ContentIDList = newContentList
			count--
			continue
		}
		if err := c.UpdateHandler(CategoryFileContent, r, c.setFileItem, resourceIDString, ""); err != nil {
			return err
		}
	}

	for i := len(fileSection.ContentIDList); i < count; i++ {
		resourceID, err := c.AddHandler(CategoryFileContent, r, c.setFileItem, "", fmt.Sprintf("%d", i))
		if err != nil {
			return err
		}
		fileSection.ContentIDList = append(fileSection.ContentIDList, resourceID.String())
	}

	return nil
}

func (c *Context) DeleteFileSection(r rest.RequestInterface, input ...interface{}) error {
	fileSection := input[0].(*models.FilesSection)

	for _, resourceIDString := range fileSection.ContentIDList {
		resourceIDString := resourceIDString
		if err := c.DeleteHandler(CategoryFileContent, r, func(r rest.RequestInterface, parameters ...interface{}) error {
			fileItem := parameters[0].(*models.FileContent)
			if fileItem.Type != FileSectionFileType {
				return nil
			}
			if err := c.FileProcessor.RemoveFile(fileItem.RefName); err != nil {
				return err
			}
			return nil
		}, resourceIDString); err != nil {
			return err
		}
	}

	return nil
}

func (c *Context) AddFAQ(r rest.RequestInterface, input ...interface{}) error {
	faq := input[0].(*models.FAQ)
	log.Println("GGGG", r.FormValue("group"))
	faq.Group = r.FormValue("group")
	pairPath := c.FileProcessor.GenerateID()
	path := filepath.Join(ResourcesPath, FAQPath, pairPath)
	faq.Answer = filepath.Join(path, "answer.txt")
	faq.Question = filepath.Join(path, "question.txt")

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	if err := c.FileProcessor.WriteToFile(faq.Answer, r.FormValue("answer")); err != nil {
		return err
	}

	if err := c.FileProcessor.WriteToFile(faq.Question, r.FormValue("question")); err != nil {
		return err
	}

	return nil
}

func (c *Context) UpdateFAQ(r rest.RequestInterface, input ...interface{}) error {
	faq := input[0].(*models.FAQ)

	if err := c.FileProcessor.WriteToFile(faq.Answer, r.FormValue("answer")); err != nil {
		return err
	}

	if err := c.FileProcessor.WriteToFile(faq.Question, r.FormValue("question")); err != nil {
		return err
	}

	return nil
}

func (c *Context) DeleteFAQ(r rest.RequestInterface, input ...interface{}) error {
	faq := input[0].(*models.FAQ)

	if err := c.FileProcessor.RemoveFile(faq.Answer); err != nil {
		return err
	}

	if err := c.FileProcessor.RemoveFile(faq.Question); err != nil {
		return err
	}

	return nil
}

func (c *Context) GetItem(id string) (map[string]string, error) {
	itemMap := make(map[string]string)
	if _, err := c.ResourcesDBController.GetResource(id, itemMap); err != nil {
		return nil, err
	}

	itemMap["id"] = id
	return itemMap, nil
}

func (c *Context) GetAllItemsByCategory(categoryName string, r *rest.Request) ([]map[string]string, error) {
	categoryID, err := c.getCategoryID(categoryName)
	if err != nil {
		return nil, err
	}

	resources, err := c.ResourcesDBController.GetResourcesByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	itemMapList := make([]map[string]string, len(resources))
	for i, resource := range resources {
		resource := resource
		itemMap := make(map[string]string)
		if err := c.ResourceModelFunctions.FromResource(&resource, itemMap); err != nil {
			return nil, err
		}
		itemMap["id"] = resource.ID.String()
		itemMapList[i] = itemMap
	}

	return itemMapList, nil
}

func (c *Context) GetFAQGroups(r *rest.Request) ([]string, error) {
	categoryID, err := c.getCategoryID(CategoryFAQGroups)
	if err != nil {
		return nil, err
	}

	resources, err := c.ResourcesDBController.GetResourcesByCategory(categoryID)
	if err != nil {
		return nil, err
	}

	if len(resources) > 0 {
		return nil, errors.New("To many FAQ group entries")
	}

	faqGroups := make([]string, 0)
	if err := c.ResourceModelFunctions.FromResource(&resources[0], faqGroups); err != nil {
		return nil, err
	}

	return faqGroups, nil
}

func (c *Context) AddNewsFeedEntry(r rest.RequestInterface, input ...interface{}) error {
	newsEntry := input[0].(*models.NewsEntry)

	dt, err := time.Parse("Mon Jan 02 2006 15:04:05.0000 GMT-0700", time.Now().String())
	if err != nil {
		return err
	}
	year := strconv.Itoa(dt.Year())
	day := strconv.Itoa(dt.Day())

	entryTextFile := fmt.Sprintf("%s.%s", uuid.New().String(), "txt")
	newsEntry.Text = filepath.Join(ResourcesPath, FAQPath, uuid.New().String(), entryTextFile)
	newsEntry.Day = day
	newsEntry.Month = dt.Month().String()[0:3]
	newsEntry.Year = year

	if err := c.FileProcessor.WriteToFile(newsEntry.Text, r.FormValue("news_text")); err != nil {
		return err
	}

	return nil
}

func (c *Context) UpdateNewsEntry(r rest.RequestInterface, input ...interface{}) error {
	newsEntry := input[0].(*models.NewsEntry)

	if err := c.FileProcessor.WriteToFile(newsEntry.Text, r.FormValue("news_text")); err != nil {
		return err
	}

	return nil
}
