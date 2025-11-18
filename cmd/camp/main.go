package main

import (
	"log"
	"net/http"

	"github.com/codersgyan/camp/internal/contact"
	"github.com/codersgyan/camp/internal/database"
)

func main() {
	db, err := database.Connect("./camp_data/camp.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// run the migration
	if err := database.RunMigration(db); err != nil {
		log.Fatal(err)
	}

	// this is the implementation of graceful shutdown

	contactRepository := contact.NewRepository(db)
	contactHandler := contact.NewHandler(contactRepository)

	http.HandleFunc("POST /api/contacts", contactHandler.Create)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
