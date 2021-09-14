package frontend

import (
	"fmt"
	"net/http"

	"github.com/artofimagination/polygnosics/businesslogic"
	httpModels "github.com/artofimagination/polygnosics/models/http"
	"github.com/artofimagination/polygnosics/rest"
)

func (c *RESTController) delete(category string, handler func(rest.RequestInterface, ...interface{}) error, w rest.ResponseWriter, r *rest.Request) {
	requestData, err := r.DecodeRequest()
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusBadRequest)
		return
	}
	if err := c.BackendContext.DeleteHandler(category, r, handler, requestData["id"].(string)); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := httpModels.ResponseData{Data: "OK"}
	w.EncodeResponse(response)
}

func (c *RESTController) add(category string, handler func(rest.RequestInterface, ...interface{}) error, w rest.ResponseWriter, r *rest.Request) {
	if _, err := c.BackendContext.AddHandler(category, r, handler); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := httpModels.ResponseData{Data: "OK"}
	w.EncodeResponse(response)
}

func (c *RESTController) addMultipart(category string, handler func(rest.RequestInterface, ...interface{}) error, w rest.ResponseWriter, r *rest.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c.add(category, handler, w, r)
}

func (c *RESTController) addForm(category string, handler func(rest.RequestInterface, ...interface{}) error, w rest.ResponseWriter, r *rest.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusBadRequest)
		return
	}

	c.add(category, handler, w, r)
}

func (c *RESTController) update(category string, handler func(rest.RequestInterface, ...interface{}) error, w rest.ResponseWriter, r *rest.Request) {
	if err := c.BackendContext.UpdateHandler(category, r, handler, r.FormValue("id")); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := httpModels.ResponseData{Data: "OK"}
	w.EncodeResponse(response)
}

func (c *RESTController) updateMultipart(category string, handler func(rest.RequestInterface, ...interface{}) error, w rest.ResponseWriter, r *rest.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c.update(category, handler, w, r)
}

func (c *RESTController) updateForm(category string, handler func(rest.RequestInterface, ...interface{}) error, w rest.ResponseWriter, r *rest.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusBadRequest)
		return
	}

	c.update(category, handler, w, r)
}

func (c *RESTController) getSingleItem(w rest.ResponseWriter, r *rest.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusBadRequest)
		return
	}

	item, err := c.BackendContext.GetItem(r.FormValue("id"))
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := httpModels.ResponseData{Data: item}
	w.EncodeResponse(response)
}

func (c *RESTController) getAll(category string, w rest.ResponseWriter, r *rest.Request) {
	items, err := c.BackendContext.GetAllItemsByCategory(category, r)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := httpModels.ResponseData{Data: items}
	w.EncodeResponse(response)
}

func (c *RESTController) getFAQs(w rest.ResponseWriter, r *rest.Request) {
	c.getAll(businesslogic.CategoryFAQ, w, r)
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
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusBadRequest)
		return
	}
	c.addForm(businesslogic.CategoryFAQ, c.BackendContext.AddFAQ, w, r)
}

func (c *RESTController) deleteFAQ(w rest.ResponseWriter, r *rest.Request) {
	c.delete(businesslogic.CategoryFAQ, c.BackendContext.DeleteFAQ, w, r)
}

func (c *RESTController) updateFAQ(w rest.ResponseWriter, r *rest.Request) {
	c.updateForm(businesslogic.CategoryFAQ, c.BackendContext.UpdateFAQ, w, r)
}

func (c *RESTController) getTutorials(w rest.ResponseWriter, r *rest.Request) {
	c.getAll(businesslogic.CategoryTutorial, w, r)
}

func (c *RESTController) getNewsFeed(w rest.ResponseWriter, r *rest.Request) {
	c.getAll(businesslogic.CategoryNews, w, r)
}

func (c *RESTController) addNewsFeedEntry(w rest.ResponseWriter, r *rest.Request) {
	c.addForm(businesslogic.CategoryNews, c.BackendContext.AddNewsFeedEntry, w, r)
}

func (c *RESTController) updateNewsEntry(w rest.ResponseWriter, r *rest.Request) {
	c.addForm(businesslogic.CategoryNews, c.BackendContext.UpdateNewsEntry, w, r)
}

func (c *RESTController) getCategoriesMap(w rest.ResponseWriter, r *rest.Request) {
	respondData, err := c.ForwardResourceDBRequest(r)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := httpModels.ResponseData{Data: respondData}
	w.EncodeResponse(response)
}

func (c *RESTController) updateTutorial(w rest.ResponseWriter, r *rest.Request) {
	c.updateMultipart(businesslogic.CategoryTutorial, c.BackendContext.UpdateTutorial, w, r)
}

func (c *RESTController) addTutorial(w rest.ResponseWriter, r *rest.Request) {
	c.addMultipart(businesslogic.CategoryTutorial, c.BackendContext.AddTutorial, w, r)
}

func (c *RESTController) deleteTutorial(w rest.ResponseWriter, r *rest.Request) {
	c.delete(businesslogic.CategoryTutorial, c.BackendContext.DeleteTutorial, w, r)
}

func (c *RESTController) getFiles(w rest.ResponseWriter, r *rest.Request) {
	c.getAll(businesslogic.CategoryFileContent, w, r)
}

func (c *RESTController) addFileSection(w rest.ResponseWriter, r *rest.Request) {
	c.addMultipart(businesslogic.CategoryFileSection, c.BackendContext.AddFileSection, w, r)
}

func (c *RESTController) updateFileSection(w rest.ResponseWriter, r *rest.Request) {
	c.updateMultipart(businesslogic.CategoryFileSection, c.BackendContext.UpdateFileSection, w, r)
}

func (c *RESTController) deleteFileSection(w rest.ResponseWriter, r *rest.Request) {
	c.delete(businesslogic.CategoryFileSection, c.BackendContext.DeleteFileSection, w, r)
}
