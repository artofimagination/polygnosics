package restcontrollers

import (
	"fmt"
	"net/http"

	"polygnosics/app/restcontrollers/contents"

	"github.com/pkg/errors"
)

func (c *RESTController) MyProductsHandler(w http.ResponseWriter, r *http.Request) {
	pUser := c.ContentController.GetUserContent()
	pProduct, err := c.ContentController.GetUserProductContent(&c.ContentController.UserData.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	for k, v := range pProduct {
		pUser[k] = v
	}
	c.RenderTemplate(w, "my-products", pUser)
}

func (c *RESTController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	p := c.ContentController.GetUserContent()
	if r.Method == GET {
		c.RenderTemplate(w, "new-product-wizard", p)
	} else {
		name := "user-main"

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			p["message"] = "Failed to parse form"
			c.RenderTemplate(w, name, p)
			return
		}

		isPublic := false
		if r.FormValue("publicProduct") == "set" {
			isPublic = true
		}

		product, err := c.UserDBController.CreateProduct(
			r.FormValue("productName"),
			isPublic,
			&c.ContentController.UserData.ID,
			c.ContentController.GeneratePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}
		c.ContentController.ProductData = product

		err = c.ContentController.UploadProductFile(contents.ProductAvatar, contents.DefaultUserAvatarPath, "product-avatar", r)
		if err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&c.ContentController.ProductData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload avatar. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		err = c.ContentController.UploadProductFile(contents.ProductMainApp, "", "main-app", r)
		if err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&c.ContentController.ProductData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload main app. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		err = c.ContentController.UploadProductFile(contents.ProductClientApp, "", "client-app", r)
		if err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&c.ContentController.ProductData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload client app. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		c.ContentController.ProductData.Details.SetURL(contents.ProductName, r.FormValue("productName"))
		c.ContentController.ProductData.Details.SetURL(contents.ProductRequires3D, r.FormValue("requires3D"))
		c.ContentController.ProductData.Details.SetURL(contents.ProductPublic, r.FormValue("publicProduct"))
		c.ContentController.ProductData.Details.SetURL(contents.ProductURL, r.FormValue("productUrl"))
		if err := c.UserDBController.UpdateProductDetails(c.ContentController.ProductData); err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&c.ContentController.ProductData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to update product details. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := c.UserDBController.UpdateProductAssets(c.ContentController.ProductData); err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&c.ContentController.ProductData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to update product assets. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		c.RenderTemplate(w, name, p)
	}
}
