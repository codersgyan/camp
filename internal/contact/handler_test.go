package contact

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codersgyan/camp/internal/database"
)

func TestCreateContactHandler_Validation(t *testing.T) {
	// Real in-memory DB + migrations — exactly like production
	db, _ := database.Connect(":memory:")
	if err := database.RunMigration(db); err != nil {
		t.Fatalf("migration failed: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	handler := NewHandler(repo)

	tests := []struct {
		name       string
		payload    map[string]any
		wantStatus int
		wantBody   string // substring check
	}{
		// Success case — we only check that we get 201 and a valid JSON with "id"
		{
			name: "valid request → 201 with id",
			payload: map[string]any{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "john@example.com",
			},
			wantStatus: http.StatusCreated,
			wantBody:   `"id":`,
		},
		{
			name: "valid with phone → 201",
			payload: map[string]any{
				"first_name": "Alice",
				"last_name":  "Smith",
				"email":      "alice@test.com",
				"phone":      "+911472583695",
			},
			wantStatus: http.StatusCreated,
			wantBody:   `"id":`,
		},
		// Validation failures
		{
			name: "missing email → 400",
			payload: map[string]any{
				"first_name": "John",
				"last_name":  "Doe",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `"field":"email","message":"required"`,
		},
		{
			name: "invalid email → 400",
			payload: map[string]any{
				"first_name": "John",
				"last_name":  "Doe",
				"email":      "bademail",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `"invalid email format"`,
		},
		{
			name: "invalid phone → 400",
			payload: map[string]any{
				"first_name": "A",
				"last_name":  "B",
				"email":      "a@b.com",
				"phone":      "12345",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `"field":"phone"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/contacts", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Create(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf("status: got %d, want %d\nbody: %s", w.Code, tt.wantStatus, w.Body.String())
			}

			if tt.wantBody != "" && !bytes.Contains(w.Body.Bytes(), []byte(tt.wantBody)) {
				t.Errorf("body missing expected substring %q\ngot: %s", tt.wantBody, w.Body.String())
			}
		})
	}
}
