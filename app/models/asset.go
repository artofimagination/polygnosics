package models

import (
	"fmt"

	"github.com/google/uuid"
)

var NullUUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

type Asset struct {
	ID         uuid.UUID  `json:"id" validation:"required"`
	References References `json:"refs" validation:"required"`
	Path       string
}

// Assets structure contains the identification of all user related documents images.
type References struct {
	AvatarID uuid.UUID `json:"avatar_id,omitempty"`
}

func (r Asset) GetAvatarPath() string {
	if r.References.AvatarID == NullUUID {
		return "/assets/images/avatar/avatar.jpg"
	}
	return fmt.Sprintf("%s/%s.jpg", r.Path, r.References.AvatarID)
}
