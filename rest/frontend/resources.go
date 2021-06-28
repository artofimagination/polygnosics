package frontend

import (
	"fmt"
	"log"
	"net/http"

	httpModels "github.com/artofimagination/polygnosics/models/http"
	"github.com/artofimagination/polygnosics/rest"
)

func (c *RESTController) getSingleItem(w rest.ResponseWriter, r *rest.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusBadRequest)
		return
	}

	faq, err := c.BackendContext.GetItem(r.FormValue("id"))
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := httpModels.ResponseData{Data: faq}
	w.EncodeResponse(response)
}

func (c *RESTController) getFAQs(w rest.ResponseWriter, r *rest.Request) {
	faqs, err := c.BackendContext.GetAllItemsByCategory("FAQ", r)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := httpModels.ResponseData{Data: faqs}
	w.EncodeResponse(response)
}

func (c *RESTController) getFAQGroups(w rest.ResponseWriter, r *rest.Request) {
	faqGroups, err := c.BackendContext.GetFAQGroups(r)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := httpModels.ResponseData{Data: faqGroups}
	w.EncodeResponse(response)
}

func (c *RESTController) addFAQ(w rest.ResponseWriter, r *rest.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if _, err := c.BackendContext.AddHandler("FAQ", r, c.BackendContext.UpdateFAQ); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := httpModels.ResponseData{Data: "OK"}
	w.EncodeResponse(response)
}

func (c *RESTController) updateFAQ(w rest.ResponseWriter, r *rest.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := c.BackendContext.UpdateHandler("FAQ", r, c.BackendContext.UpdateFAQ); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := httpModels.ResponseData{Data: "OK"}
	w.EncodeResponse(response)
}

func (c *RESTController) getTutorials(w rest.ResponseWriter, r *rest.Request) {
	tutorials, err := c.BackendContext.GetAllItemsByCategory("Tutorial", r)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := httpModels.ResponseData{Data: tutorials}
	w.EncodeResponse(response)
}

func (c *RESTController) getNewsFeed(w rest.ResponseWriter, r *rest.Request) {
	tutorials, err := c.BackendContext.GetAllItemsByCategory("News feed", r)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := httpModels.ResponseData{Data: tutorials}
	w.EncodeResponse(response)
}

func (c *RESTController) addNewsFeedEntry(w rest.ResponseWriter, r *rest.Request) {
	if err := c.BackendContext.AddNewsFeedEntry(r); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := httpModels.ResponseData{Data: "OK"}
	w.EncodeResponse(response)
}

func (c *RESTController) updateNewsEntry(w rest.ResponseWriter, r *rest.Request) {
	if err := c.BackendContext.UpdateNewsEntry(r); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := httpModels.ResponseData{Data: "OK"}
	w.EncodeResponse(response)
}

func (c *RESTController) getCategoriesMap(w rest.ResponseWriter, r *rest.Request) {
	respondData, err := r.ForwardRequest(rest.UserDBAddress)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := httpModels.ResponseData{Data: respondData}
	w.EncodeResponse(response)
}

func (c *RESTController) updateTutorial(w rest.ResponseWriter, r *rest.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := c.BackendContext.UpdateHandler("Tutorial", r, c.BackendContext.UpdateTutorial); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := httpModels.ResponseData{Data: "OK"}
	w.EncodeResponse(response)
}

func (c *RESTController) addTutorial(w rest.ResponseWriter, r *rest.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if _, err := c.BackendContext.AddHandler("Tutorial", r, c.BackendContext.AddTutorial); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := httpModels.ResponseData{Data: "OK"}
	w.EncodeResponse(response)
}

func (c *RESTController) getFiles(w rest.ResponseWriter, r *rest.Request) {
	files, err := c.BackendContext.GetAllItemsByCategory("Files", r)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := httpModels.ResponseData{Data: files}
	w.EncodeResponse(response)
}

func (c *RESTController) addFileSection(w rest.ResponseWriter, r *rest.Request) {
	log.Println("Add files section")
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if _, err := c.BackendContext.AddHandler("FilesSection", r, c.BackendContext.AddFileSection); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := httpModels.ResponseData{Data: "OK"}
	w.EncodeResponse(response)
}

func (c *RESTController) updateFileSection(w rest.ResponseWriter, r *rest.Request) {
	log.Println("Update files section")
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := c.BackendContext.UpdateHandler("FilesSection", r, c.BackendContext.UpdateFileSection, r.FormValue("id")); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := httpModels.ResponseData{Data: "OK"}
	w.EncodeResponse(response)
}
