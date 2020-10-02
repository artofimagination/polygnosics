package contents

import (
	"aiplayground/app/utils/page"
)

// CreateError returns the error page content, determined by the errorStr input parameter.
func CreateError(errorStr string) *page.Page {
	p := page.CreatePage("error")
	p.Data["title"] = "Something went wrong"
	p.Data["body"] = errorStr
	return p
}
