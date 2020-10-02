package contents

import (
	"aiplayground/app/utils/page"
)

// CreateHome creates the content structure for the user home page.
func CreateHome(user string) error {
	name := "user-main"
	p1 := page.CreatePage(name)
	p1.Data["title"] = "Main"
	p1.Data["body"] = "Welcome " + user
	return page.Save(name, p1)
}

// CreateNewProjectConfig creates the content structure for new projects config page.
func CreateNewProjectConfig() error {
	name := "new-project"
	p1 := page.CreatePage(name)
	p1.Data["title"] = "New project Config"
	p1.Data["body"] = "Create new"
	p1.Data["features"] = nil
	p1.Data["config"] = nil
	return page.Save(name, p1)
}
