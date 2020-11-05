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

	"polygnosics/app/restcontrollers"
	"polygnosics/app/services/db/mysqldb"
	"polygnosics/app/services/db/timescaledb"
	"polygnosics/app/utils/configloader"

	"github.com/pkg/errors"
)

func main() {
	// Load DB configuration
	config, err := configloader.LoadDBConfigFromEnv("MYSQL")
	if err != nil {
		log.Fatalf("Failed to load MYSQl DB config. %s\n", errors.WithStack(err))
	}
	mysqldb.DBConnection = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/user_database?parseTime=true",
		config.Username,
		config.Password,
		config.Address,
		config.Port)

	config, err = configloader.LoadDBConfigFromEnv("Timescale")
	if err != nil {
		log.Fatalf("Failed to load Timescale DB config. %s\n", errors.WithStack(err))
	}

	timescaledb.DBConnection = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/data?sslmode=disable",
		config.Username,
		config.Password,
		config.Address,
		config.Port)

	// Run DB migration
	if err := mysqldb.BootstrapSystem(); err != nil {
		log.Fatalf("System bootstrap failed. %s\n", errors.WithStack(err))
	}
	if err := timescaledb.BootstrapData(); err != nil {
		log.Fatalf("Data bootstrap failed. %s\n", errors.WithStack(err))
	}

	// Create Server and Route Handlers
	srv := &http.Server{
		Handler:      restcontrollers.CreateRouter(),
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
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Shutting down")
	os.Exit(0)
}
