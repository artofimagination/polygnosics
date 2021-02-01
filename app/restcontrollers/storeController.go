package restcontrollers

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func (c *RESTController) StoreHandler(w http.ResponseWriter, r *http.Request) {
	content, err := c.ContentController.BuildStoreContent()
	if err != nil {
		errString := fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err))
		c.RenderTemplate(w, "store", c.ContentController.BuildErrorContent(errString))
		return
	}
	c.RenderTemplate(w, "store", content)
}
