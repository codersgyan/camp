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

// Create Contact godoc
// @Summary Create or update a contact with tags
// @Description Upsert a contact record. If the email exists, tags are updated; otherwise, a new contact is created.
// @Tags Contacts
// @Accept json
// @Produce json
// @Param contact body Contact true "Contact payload"
// @Success 201 {object} map[string]int64
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/contacts [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	// todo: request validation
	var contactBody Contact

	if err := json.NewDecoder(r.Body).Decode(&contactBody); err != nil {
		fmt.Println(err)
		http.Error(w, "invalid json", http.StatusBadRequest)
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
