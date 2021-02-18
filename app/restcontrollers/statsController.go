package restcontrollers

import (
	"fmt"
	"net/http"

	"polygnosics/app/utils/webrtc"

	"github.com/pkg/errors"
)

func (c *RESTController) StatsWebRTC(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse frontend webrtc offer. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}

	offer := r.FormValue("offer")
	if err := webrtc.SetupFrontend(w, r, offer, c.BackendContext.ProvideUserStats); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start frontend webrtc. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
}

func (c *RESTController) ProductStats(w http.ResponseWriter, r *http.Request) {
	content, err := c.ContentController.BuildStoreContent()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	c.RenderTemplate(w, ProductStats, content)
}

func (c *RESTController) ProjectStats(w http.ResponseWriter, r *http.Request) {
	content, err := c.ContentController.BuildStoreContent()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	c.RenderTemplate(w, ProjectStats, content)
}

func (c *RESTController) UserStats(w http.ResponseWriter, r *http.Request) {
	content, err := c.ContentController.BuildStoreContent()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	c.RenderTemplate(w, UserStats, content)
}

func (c *RESTController) ProductsProjectsStats(w http.ResponseWriter, r *http.Request) {
	content, err := c.ContentController.BuildStoreContent()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	c.RenderTemplate(w, ProductProject, content)
}

func (c *RESTController) AccountingStats(w http.ResponseWriter, r *http.Request) {
	content, err := c.ContentController.BuildStoreContent()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	c.RenderTemplate(w, Accounting, content)
}

func (c *RESTController) SystemHealthStats(w http.ResponseWriter, r *http.Request) {
	content, err := c.ContentController.BuildStoreContent()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get product content. %s", errors.WithStack(err)), http.StatusInternalServerError)
		return
	}
	c.RenderTemplate(w, SystemHealth, content)
}
