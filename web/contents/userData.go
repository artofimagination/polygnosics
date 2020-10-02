package contents

import (
	"aiplayground/app/models"
	"aiplayground/app/utils/page"
)

// UserData stores the data that is used after the user is logged in
type UserData struct {
	CurrentProject models.Project
	User           models.User
}

// CreateUserData stores the structure needed represent logged in user data, settings.
func CreateUserData(user models.User) error {
	name := "user_data"
	var data UserData
	data.User = user

	return page.Save(name, data)
}
