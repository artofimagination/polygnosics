package ipresolver

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/artofimagination/polygnosics/initialization"
	"github.com/artofimagination/polygnosics/rest"
	"github.com/artofimagination/polygnosics/rest/frontend"
	"github.com/artofimagination/polygnosics/rest/resourcesdb"
	"github.com/artofimagination/polygnosics/rest/userdb"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type IPResolver struct {
	UserDBAddress     *rest.Server
	ResourceDBAddress *rest.Server
}

func NewIPResolver(frontend *frontend.RESTController, userdb *userdb.RESTController, resourcedb *resourcesdb.RESTController, cfg *initialization.Config) *IPResolver {
	userServer := &rest.Server{
		IP:   cfg.UserDBAddress,
		Port: cfg.UserDBPort,
		Name: cfg.UserDBName,
	}
	resourceServer := &rest.Server{
		IP:   cfg.ResourceDBAddress,
		Port: cfg.ResourceDBPort,
		Name: cfg.ResourceDBName,
	}
	frontend.UserDBAddress = userServer
	frontend.ResourceDBAddress = resourceServer
	userdb.ServerAddress = userServer
	resourcedb.ServerAddress = resourceServer

	ipresolver := &IPResolver{
		UserDBAddress:     userServer,
		ResourceDBAddress: resourceServer,
	}

	return ipresolver
}

// DetectValidAddresses waits 5 seconds to allow the IP resolver server to set the user and resource db addresses
// Only waits if the address was not set through env vars before.
// Returns error if the address is not set either by env. vars or the IP resolver server.
func (c *IPResolver) DetectValidAddresses() error {
	userDBSet := false
	for retryCount := 5; retryCount > 0; retryCount-- {
		if c.UserDBAddress.Name != "Unknown" {
			userDBSet = true
			break
		}
		time.Sleep(1 * time.Second)
		log.Printf("No valid user db address")
	}

	if !userDBSet {
		return errors.New("No valid user db address detected")
	}

	resourceDBSet := false
	for retryCount := 5; retryCount > 0; retryCount-- {
		if c.ResourceDBAddress.Name != "Unknown" {
			resourceDBSet = true
			break
		}
		time.Sleep(1 * time.Second)
		log.Printf("No valid resource db address")
	}

	if !resourceDBSet {
		return errors.New("No valid resource db address detected")
	}
	return nil
}

func (c *IPResolver) setUserDBAddress(w rest.ResponseWriter, r *rest.Request) {
	requestData, err := r.DecodeRequest()
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusBadRequest)
		return
	}

	c.UserDBAddress.IP = requestData["ip"].(string)
	c.UserDBAddress.Port = requestData["port"].(int)
	c.UserDBAddress.Name = requestData["name"].(string)

	w.WriteData("OK", http.StatusCreated)
}

func (c *IPResolver) setResourceDBAddress(w rest.ResponseWriter, r *rest.Request) {
	requestData, err := r.DecodeRequest()
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusBadRequest)
		return
	}

	c.ResourceDBAddress.IP = requestData["ip"].(string)
	c.ResourceDBAddress.Port = requestData["port"].(int)
	c.ResourceDBAddress.Name = requestData["name"].(string)

	w.WriteData("OK", http.StatusCreated)
}

func (c *IPResolver) AddRouting(r *mux.Router) {
	r.HandleFunc("/set-userdb-address", rest.MakeHandler(c.setUserDBAddress))
	r.HandleFunc("/set-resourcedb-address", rest.MakeHandler(c.setResourceDBAddress))
}
