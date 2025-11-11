package main

import (
	"fmt"
	"log"

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
		log.Fatal("failed to run migration: %w", err)
	}

	contactRepository := contact.NewRepository(db)

	newContact := contact.Contact{
		FirstName: "Rakesh",
		LastName:  "K",
		Email:     "rakesh@codersgyan.com",
	}

	createdId, err := contactRepository.Create(&newContact)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created id: %d", createdId)
}
