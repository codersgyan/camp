package main

import (
	"log"
	"net/http"

	"github.com/codersgyan/camp/internal/api/routes"
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

	apiRoutes := routes.Register(db)

	log.Fatal(http.ListenAndServe(":8080", apiRoutes))
}
