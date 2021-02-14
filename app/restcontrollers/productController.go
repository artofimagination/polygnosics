package restcontrollers

import (
	"fmt"
	"net/http"

	"polygnosics/app/businesslogic"

	"github.com/pkg/errors"
)

func (c *RESTController) MyProducts(w http.ResponseWriter, r *http.Request) {
	content, err := c.ContentController.BuildMyProductsContent()
	if err != nil {
		errString := fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err))
		c.RenderTemplate(w, MyProducts, c.ContentController.BuildErrorContent(errString))
		return
	}
	c.RenderTemplate(w, MyProducts, content)
}

func (c *RESTController) MyProductDetails(w http.ResponseWriter, r *http.Request) {
	productID, err := parseItemID(r)
	if err != nil {
		c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent("Failed to parse product id"))
		return
	}

	content, err := c.ContentController.BuildProductDetailsContent(productID)
	if err != nil {
		c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent("Failed to get product content"))
		return
	}

	c.RenderTemplate(w, "details", content)
}

func (c *RESTController) ProductDetails(w http.ResponseWriter, r *http.Request) {
	productID, err := parseItemID(r)
	if err != nil {
		c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent("Failed to parse product id"))
		return
	}

	content, err := c.ContentController.BuildProductDetailsContent(productID)
	if err != nil {
		c.RenderTemplate(w, UserMain, c.ContentController.BuildErrorContent("Failed to get product content"))
		return
	}

	c.RenderTemplate(w, "details", content)
}

func (c *RESTController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method == GET {
		content := c.ContentController.BuildProductWizardContent()
		c.RenderTemplate(w, "product-wizard", content)
	} else {
		content, err := c.ContentController.BuildUserMainContent()
		if err != nil {
			content["message"] = fmt.Sprintf("Failed to load user main content. %s", errors.WithStack(err))
			c.RenderTemplate(w, UserMain, content)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			content["message"] = ErrFailedToParseForm
			c.RenderTemplate(w, UserMain, content)
			return
		}

		product, err := c.BackendContext.AddProduct(&c.ContentController.UserData.ID, r.FormValue(businesslogic.ProductNameKey), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := c.BackendContext.CreateDockerImage(product, &c.ContentController.UserData.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := c.BackendContext.UpdateProductData(product, r); err != nil {
			if errDelete := c.BackendContext.DeleteProduct(product); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
				return
			}
			http.Error(w, fmt.Sprintf("Failed to update product details. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		c.RenderTemplate(w, UserMain, content)
	}
}

func (c *RESTController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method == POST {
		productID, err := parseItemID(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse product id. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		product, err := c.UserDBController.GetProduct(productID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := c.BackendContext.DeleteProduct(product); err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		c.MyProducts(w, r)
	}
}

func (c *RESTController) EditProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := parseItemID(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse product id. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	if r.Method == GET {
		content, err := c.ContentController.BuildProductEditContent(productID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get build product edit content. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		c.RenderTemplate(w, "product-edit", content)
	} else {
		product, err := c.UserDBController.GetProduct(productID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := c.BackendContext.UploadFiles(product.Assets, r); err != nil {
			http.Error(w, fmt.Sprintf("Failed to upload assets. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}
		_, _, err = r.FormFile(businesslogic.ProductMainAppKey)
		if err == nil {
			if err := c.BackendContext.CreateDockerImage(product, &c.ContentController.UserData.ID); err != nil {
				if errDelete := c.BackendContext.DeleteProduct(product); errDelete != nil {
					http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
					return
				}
				http.Error(w, fmt.Sprintf("Failed to create main app docker image. %s", errors.WithStack(err)), http.StatusInternalServerError)
				return
			}
		}

		if err := c.BackendContext.UpdateProductData(product, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		content, err := c.ContentController.BuildMyProductsContent()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get my products content. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}
		c.RenderTemplate(w, MyProducts, content)
	}
}
