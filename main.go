package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "github.com/artofimagination/polygnosics/context"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func main() {
	appContext, err := app.NewContext()
	if err != nil {
		log.Fatalf("Failed to initiate context. %s\n", errors.WithStack(err))
	}

	r := mux.NewRouter()
	appContext.RESTFrontend.AddRouting(r)
	appContext.RESTUserDB.AddRouting(r)

	// Create Server and Route Handlers
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8184",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
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
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Shutting down")
	os.Exit(0)
}
