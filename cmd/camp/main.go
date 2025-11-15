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

	contactRepository := contact.NewRepository(db)
	contactHandler := contact.NewHandler(contactRepository)

<<<<<<< HEAD
	http.HandleFunc("POST /api/contacts", contactHandler.Create)
=======
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/contacts", contactHandler.Create)
	mux.HandleFunc("GET /api/contacts", contactHandler.List)
>>>>>>> 67b9709 (Add GET /api/contacts with pagination support)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
