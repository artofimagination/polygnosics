package mysqldb

import (
	"polygnosics/app/models"

	"github.com/google/uuid"
)

func AddSettings() (*uuid.UUID, error) {
	queryString := "INSERT INTO user_settings (id, settings) VALUES (UUID_TO_BIN(?), ?)"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	newID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	settings := models.UserSetting{}
	query, err := db.Query(queryString, newID, settings.Settings)
	if err != nil {
		return nil, err
	}

	defer query.Close()
	return &newID, nil
}

func GetSettings(settingsID *uuid.UUID) (*models.UserSetting, error) {
	settings := models.UserSetting{}
	queryString := "SELECT settings FROM user_settings WHERE id = UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := db.QueryRow(queryString, *settingsID)

	if err != nil {
		return nil, err
	}

	if err := query.Scan(&settings.Settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func DeleteSettings(settingsID *uuid.UUID) error {
	query := "DELETE FROM user_settings WHERE id=UUID_TO_BIN(?)"
	db, err := ConnectSystem()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query, *settingsID)
	if err != nil {
		return err
	}
	return nil
}
