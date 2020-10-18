package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"polygnosics/app/restControllers"
	"polygnosics/app/services/db"

	"github.com/pkg/errors"
)

func main() {

	if err := db.BootstrapSystem(); err != nil {
		log.Fatal(fmt.Sprintf("System bootstrap failed. %s", errors.WithStack(err)))
	}
	if err := db.BootstrapData(); err != nil {
		log.Fatal(fmt.Sprintf("Data bootstrap failed. %s", errors.WithStack(err)))
	}

	// Create Server and Route Handlers
	srv := &http.Server{
		Handler:      restControllers.CreateRouter(),
		Addr:         ":8081",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
