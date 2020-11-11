package models

import (
	"github.com/google/uuid"
)

type Asset struct {
	ID         uuid.UUID `json:"id" validation:"required"`
	References uuid.UUID `json:"refs" validation:"required"`
}

// Assets structure contains the identification of all user related documents images.
type References struct {
	AvatarID uuid.UUID `json:"avatar_id,omitempty"`
}
