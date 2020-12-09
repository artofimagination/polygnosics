package restcontrollers

import (
	"fmt"
	"net/http"
	"polygnosics/app"
	"polygnosics/app/restcontrollers/auth"
	"polygnosics/app/restcontrollers/page"

	"github.com/pkg/errors"
)

func MyProductsHandler(w http.ResponseWriter, r *http.Request) {
	pUser := getUserContent()
	pProduct, err := getUserProductContent(&auth.UserData.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	for k, v := range pProduct {
		pUser[k] = v
	}
	page.RenderTemplate(w, "my-products", pUser)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	p := getUserContent()
	if r.Method == "GET" {
		page.RenderTemplate(w, "new-product-wizard", p)
	} else {
		name := "user-main"

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			p["message"] = "Failed to parse form"
			page.RenderTemplate(w, name, p)
			return
		}

		isPublic := false
		if r.FormValue("publicProduct") == "set" {
			isPublic = true
		}

		product, err := app.ContextData.UserDBController.CreateProduct(r.FormValue("productName"), isPublic, &auth.UserData.ID, auth.GeneratePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}
		ProductData = product

		err = uploadProductFile(ProductAvatar, DefaultUserAvatarPath, "product-avatar", r)
		if err != nil {
			if errDelete := app.ContextData.UserDBController.DeleteProduct(&ProductData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload avatar. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		err = uploadProductFile(ProductMainApp, "", "main-app", r)
		if err != nil {
			if errDelete := app.ContextData.UserDBController.DeleteProduct(&ProductData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload main app. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		err = uploadProductFile(ProductClientApp, "", "client-app", r)
		if err != nil {
			if errDelete := app.ContextData.UserDBController.DeleteProduct(&ProductData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to upload client app. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		ProductData.Details.SetURL(ProductName, r.FormValue("productName"))
		ProductData.Details.SetURL(ProductRequires3D, r.FormValue("requires3D"))
		ProductData.Details.SetURL(ProductPublic, r.FormValue("publicProduct"))
		ProductData.Details.SetURL(ProductURL, r.FormValue("productUrl"))
		if err := app.ContextData.UserDBController.UpdateProductDetails(ProductData); err != nil {
			if errDelete := app.ContextData.UserDBController.DeleteProduct(&ProductData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to update product details. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		if err := app.ContextData.UserDBController.UpdateProductAssets(ProductData); err != nil {
			if errDelete := app.ContextData.UserDBController.DeleteProduct(&ProductData.ID); errDelete != nil {
				err = errors.Wrap(errors.WithStack(err), errDelete.Error())
				http.Error(w, fmt.Sprintf("Failed to delete product. %s", errors.WithStack(err)), http.StatusInternalServerError)
			}
			http.Error(w, fmt.Sprintf("Failed to update product assets. %s", errors.WithStack(err)), http.StatusInternalServerError)
			return
		}

		page.RenderTemplate(w, name, p)
	}
}
