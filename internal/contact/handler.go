package contact

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var contactBody Contact

	if err := json.NewDecoder(r.Body).Decode(&contactBody); err != nil {
		fmt.Println(err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Validate request
	if errs := ValidateCreateContactRequest(&contactBody); len(errs) > 0 {
		respondValidationError(w, errs)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	createdId, err := h.repo.CreateContactOrUpsertTags(&contactBody)
	if err != nil {
		resp := map[string]string{
			"message": "Internal server error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := map[string]int64{
		"id": createdId,
	}
	json.NewEncoder(w).Encode(resp)
}

// respondValidationError sends a 400 Bad Request response with validation error details
func respondValidationError(w http.ResponseWriter, errors []ContactValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := ValidationErrorResponse{
		Error:   "validation failed",
		Details: errors,
	}

	json.NewEncoder(w).Encode(response)
}
