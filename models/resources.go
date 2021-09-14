package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ResourceModels interface {
	ToResource(resource *Resource, data interface{}) error
	ToResourceContent(resource *Resource, data interface{}) error
	FromResource(resource *Resource, data interface{}) error
}

type Resource struct {
	ID       uuid.UUID              `json:"id" validate:"required"`
	Category int                    `json:"category" validate:"required"`
	Content  map[string]interface{} `json:"content" validate:"required"`
}

type ResourceModelImpl struct {
}

func (r *ResourceModelImpl) ToResource(resource *Resource, data interface{}) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &resource)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResourceModelImpl) ToResourceContent(resource *Resource, data interface{}) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &resource.Content)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResourceModelImpl) FromResource(resource *Resource, data interface{}) error {
	resourceBytes, err := json.Marshal(resource.Content)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resourceBytes, &data); err != nil {
		return err
	}
	return nil
}

type Tutorial struct {
	Title        string    `json:"title" validate:"required"`
	ShortDesc    string    `json:"short" validate:"required"`
	AvatarSource string    `json:"avatar" validate:"required"`
	AvatarType   string    `json:"avatar_type" validate:"required"`
	Content      string    `json:"content" validate:"required"`
	LastUpdated  time.Time `json:"last_updated" validate:"required"`
}

type FAQ struct {
	Group    string `json:"group" validate:"required"`
	Question string `json:"question" validate:"required"`
	Answer   string `json:"answer" validate:"required"`
}

type FilesSection struct {
	Title            string   `json:"title" validate:"required"`
	ShortDescription string   `json:"short" validate:"required"`
	ContentIDList    []string `json:"files" validate:"required"`
}

type FileContent struct {
	Type             string `json:"type" validate:"required"`
	Ref              string `json:"ref" validate:"required"`
	OriginalFileName string `json:"orig_file_name" validate:"required"`
	RefName          string `json:"ref_name" validate:"required"`
}

type Category struct {
	ID          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type Categories []Category

type NewsEntry struct {
	Text  string `json:"news_text" validate:"required"`
	Year  string `json:"news_year" validate:"required"`
	Month string `json:"news_month" validate:"required"`
	Day   string `json:"news_day" validate:"required"`
}
