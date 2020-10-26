package models

import (
	"github.com/google/uuid"
)

// User defines the user structures. Each user must have an associated settings entry.
type User struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Email      string    `json:"email,omitempty"`
	Password   string    `json:"password,omitempty"`
	SettingsID uuid.UUID `json:"user_settings_id,omitempty"`
}
