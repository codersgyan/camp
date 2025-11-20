package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codersgyan/camp/internal/contact"
	"github.com/codersgyan/camp/internal/database"
	graceful "github.com/codersgyan/camp/internal/shutdown"
)

func main() {
	db, err := database.Connect("./camp_data/camp.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := database.RunMigration(db); err != nil {
		log.Fatal(err)
	}

	shutdown := graceful.New(10 * time.Second)

	contactRepository := contact.NewRepository(db)
	contactHandler := contact.NewHandler(contactRepository)

	srv := &http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("POST /api/contacts", contactHandler.Create)

	// Register cleanup
	shutdown.Register(srv.Shutdown)
	shutdown.Register(db.Close)

	// Start server in goroutine
	go func() {
		log.Println("Server running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for signal
	exitCode, err := shutdown.Wait(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)

}
