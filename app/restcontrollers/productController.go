package restcontrollers

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"polygnosics/app/businesslogic"
	"polygnosics/web/contents"

	"github.com/artofimagination/golang-docker/docker"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (c *RESTController) MyProducts(w http.ResponseWriter, r *http.Request) {
	content, err := c.ContentController.BuildMyProductsContent()
	if err != nil {
		errString := fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err))
		c.RenderTemplate(w, "my-products", c.ContentController.BuildErrorContent(errString))
		return
	}
	c.RenderTemplate(w, "my-products", content)
}

func (c *RESTController) ProductDetails(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent(ErrFailedToParseForm))
		return
	}
	productID, err := uuid.Parse(r.FormValue("product"))
	if err != nil {
		c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent("Failed to parse product id"))
		return
	}

	content, err := c.ContentController.BuildProductDetailsContent(&productID)
	if err != nil {
		c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent("Failed to get product content"))
		return
	}

	c.RenderTemplate(w, "details", content)
}

func (c *RESTController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	content := c.ContentController.BuildProductWizardContent()
	if r.Method == GET {
		c.RenderTemplate(w, "new-product-wizard", content)
	} else {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			content["message"] = ErrFailedToParseForm
			c.RenderTemplate(w, UserMain, content)
			return
		}

		product, err := c.UserDBController.CreateProduct(
			r.FormValue("productName"),
			&c.ContentController.UserData.ID,
			businesslogic.GeneratePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		err = c.BackendContext.UploadFile(product.Assets, contents.ProductAvatar, businesslogic.DefaultProductAvatarPath, "product-avatar", r)
		if err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload avatar. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		err = c.BackendContext.UploadFile(product.Assets, contents.ProductMainApp, "", "main-app", r)
		if err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload main app. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := c.BackendContext.UploadFile(product.Assets, contents.ProductClientApp, "", "client-app", r); err != nil {
			if errDelete := c.UserDBController.DeleteProduct(&product.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload client app. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		pathString := c.UserDBController.ModelFunctions.GetFilePath(product.Assets, contents.ProductMainApp, "")
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

		c.ContentController.SetProductDetails(product.Details, r)

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

		c.RenderTemplate(w, UserMain, content)
	}
}
