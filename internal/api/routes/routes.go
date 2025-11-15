package routes

import (
	"database/sql"
	"net/http"

	"github.com/codersgyan/camp/internal/contact"
)

func route(method string, prefix, path string) string {
	return method + " " + prefix + path
}

func ContactRoutes(prefix string, mux *http.ServeMux, db *sql.DB) {
	contactRepository := contact.NewRepository(db)
	contactHandler := contact.NewHandler(contactRepository)
	mux.HandleFunc(route("POST", prefix, ""), contactHandler.Create)
}

func Register(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	ContactRoutes("/api/v1/contacts", mux, db)
	return mux
}
