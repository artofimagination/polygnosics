package models

import (
	"aiplayground/app/services/db"
	"encoding/json"

	"github.com/google/uuid"
)

type Data struct {
	Id       int     `json:"id"`
	DataType int     `json:"type"`
	Speed    float32 `json:"speed"`
}

// Project defines the project structure.
type Project struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	UserID    uuid.UUID       `json:"user_id"`
	FeatureID int             `json:"features_id"`
	Config    json.RawMessage `json:"config"`
}

// UpdateProject updates the selected project.
func UpdateProject(project Project) error {
	query := "UPDATE projects set user_id = ?, features_id = ?, name = ?, config = ? where id = ?"
	db, err := db.ConnectSystem()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Query(query, project.UserID, project.FeatureID, project.Name, project.Config, project.ID)
	if err != nil {
		return err
	}
	return nil
}

// AddProject adds a new project to the database.
func AddProject(project Project) error {
	query := "INSERT INTO projects (id, user_id, features_id, name, config) VALUES (UUID_TO_BIN(UUID()), UUID_TO_BIN(?), ?, ?, ?)"
	db, err := db.ConnectSystem()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Query(query, project.UserID, project.FeatureID, project.Name, project.Config)
	if err != nil {
		return err
	}
	return nil
}

// InsertData will insert data into timescale db.
func InsertData(projectID int, data interface{}) error {
	query := "INSERT INTO data VALUES (NOW(), ?, ?)"
	db, err := db.ConnectData()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Query(query, projectID, data)
	if err != nil {
		return err
	}
	return nil
}

// GetProjectByName returns the project by name.
func GetProjectByName(name string) (Project, error) {
	var project Project
	queryString := "select BIN_TO_UUID(id), name, config from projects where name = ?"
	db, err := db.ConnectSystem()
	if err != nil {
		return project, err
	}
	defer db.Close()

	query, err := db.Query(queryString, name)

	if err != nil {
		return project, err
	}

	defer query.Close()

	query.Next()
	if err := query.Scan(&project.ID, &project.Name, &project.Config); err != nil {
		return project, err
	}

	return project, nil
}
