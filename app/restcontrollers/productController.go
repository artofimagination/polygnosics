package restcontrollers

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"

	"github.com/artofimagination/golang-docker/docker"
	"github.com/artofimagination/polygnosics/app/businesslogic"
	"github.com/artofimagination/polygnosics/app/restcontrollers/contents"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *RESTController) MyProducts(w http.ResponseWriter, r *http.Request) {
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

func (c *RESTController) ProductDetails(w http.ResponseWriter, r *http.Request) {
	pUser := c.ContentController.GetUserContent()
	name := UserMain
	if err := r.ParseForm(); err != nil {
		pUser["message"] = contents.ErrFailedToParseForm
		c.RenderTemplate(w, name, pUser)
		return
	}
	productID, err := uuid.Parse(r.FormValue("product"))
	if err != nil {
		pUser["message"] = "Failed to parse product id"
		c.RenderTemplate(w, name, pUser)
		return
	}

	pProduct, err := c.ContentController.GetProductContent(&productID)
	if err != nil {
		pUser["message"] = "Failed to get product content"
		c.RenderTemplate(w, name, pUser)
		return
	}

	for k, v := range pProduct {
		pUser[k] = v
	}
	c.RenderTemplate(w, "details", pUser)
}

func (c *RESTController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	p := c.ContentController.GetUserContent()
	if r.Method == GET {
		c.RenderTemplate(w, "new-product-wizard", p)
	} else {
		name := UserMain

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			p["message"] = contents.ErrFailedToParseForm
			c.RenderTemplate(w, name, p)
			return
		}

		isPublic := false
		if r.FormValue("publicProduct") == "Yes" {
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

		err = c.ContentController.UploadFile(product.Assets, contents.ProductAvatar, contents.DefaultProductAvatarPath, "product-avatar", r)
		if err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload avatar. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		err = c.ContentController.UploadFile(product.Assets, contents.ProductMainApp, "", "main-app", r)
		if err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload main app. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := c.ContentController.UploadFile(product.Assets, contents.ProductClientApp, "", "client-app", r); err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload client app. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		pathString := product.Assets.GetField(contents.ProductMainApp, "")
		if err := businesslogic.Untar(pathString); err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to decompress main app. %s", errors.WithStack(err)), http.StatusInternalServerError)
		}

		imageName := fmt.Sprintf("%s/%s", c.ContentController.UserData.ID.String(), product.ID.String())
		sourcePath := path.Join(filepath.Dir(pathString), "build")
		if err := docker.CreateImage(sourcePath, imageName); err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to create product image. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		product.Details.SetField(contents.ProductName, r.FormValue("productName"))
		product.Details.SetField(contents.ProductDescription, r.FormValue("productDescription"))
		product.Details.SetField(contents.ProductRequires3D, r.FormValue("requires3D"))
		product.Details.SetField(contents.ProductPublic, contents.GetBooleanString(r.FormValue("publicProduct")))
		product.Details.SetField(contents.ProductURL, r.FormValue("productUrl"))
		if err := c.UserDBController.UpdateProductDetails(product); err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to update product details. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := c.UserDBController.UpdateProductAssets(product); err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to update product assets. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		c.RenderTemplate(w, name, p)
	}

}
