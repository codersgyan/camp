package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/codersgyan/camp/config"
	"github.com/codersgyan/camp/internal/contact"
	"github.com/codersgyan/camp/internal/database"
)

func main() {

	cfg, err := config.SetupENV()

	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect("./camp_data/camp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// run the migration
	if err := database.RunMigration(db); err != nil {
		log.Fatal(err)
	}

	contactRepository := contact.NewRepository(db)
	contactHandler := contact.NewHandler(contactRepository)

	http.HandleFunc("POST /api/contacts", contactHandler.Create)
	http.HandleFunc("/health", HealthCheckHandler)

	// Graceful shutdown

	go func() {
		if err := http.ListenAndServe(cfg.HttpPort, nil); err != nil {
			log.Fatalf("server error: %v", err)
		}
		log.Printf("server running on this port: %v", cfg.HttpPort)
	}()

	// Wait for signal (Ctrl+C/kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down gracefully...")

}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(resp)
}
