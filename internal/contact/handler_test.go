package contact

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Create(t *testing.T) {
	db := setupDB(t)
	repo := NewRepository(db)
	handler := NewHandler(repo)

	tests := []struct {
		name             string
		body             interface{}
		expectedStatus   int
		validateResponse func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "valid contact creation",
			body: map[string]interface{}{
				"email":      "test@example.com",
				"first_name": "John",
				"last_name":  "Doe",
			},
			expectedStatus: http.StatusCreated,
			validateResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				if resp.Code != http.StatusCreated {
					t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.Code)
				}
				var result map[string]interface{}
				if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				if _, ok := result["id"]; !ok {
					t.Error("Response should contain 'id' field")
				}
			},
		},
		{
			name: "empty email - should fail",
			body: map[string]interface{}{
				"email":      "",
				"first_name": "John",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var result ValidationErrorResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				if result.Error != "validation failed" {
					t.Errorf("Expected error 'validation failed', got '%s'", result.Error)
				}
				if len(result.Details) == 0 {
					t.Error("Expected validation error details")
				}
				foundEmailError := false
				for _, detail := range result.Details {
					if detail.Field == "email" && detail.Message == "email is required" {
						foundEmailError = true
						break
					}
				}
				if !foundEmailError {
					t.Error("Expected email validation error")
				}
			},
		},
		{
			name: "invalid email format - should fail",
			body: map[string]interface{}{
				"email":      "not-an-email",
				"first_name": "John",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var result ValidationErrorResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				foundEmailError := false
				for _, detail := range result.Details {
					if detail.Field == "email" && detail.Message == "email must be valid format" {
						foundEmailError = true
						break
					}
				}
				if !foundEmailError {
					t.Error("Expected email format validation error")
				}
			},
		},
		{
			name: "no name fields - should fail",
			body: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var result ValidationErrorResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				foundNameError := false
				for _, detail := range result.Details {
					if detail.Field == "name" {
						foundNameError = true
						break
					}
				}
				if !foundNameError {
					t.Error("Expected name validation error")
				}
			},
		},
		{
			name: "empty first_name and last_name - should fail",
			body: map[string]interface{}{
				"email":      "test@example.com",
				"first_name": "",
				"last_name":  "",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var result ValidationErrorResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				foundNameError := false
				for _, detail := range result.Details {
					if detail.Field == "name" {
						foundNameError = true
						break
					}
				}
				if !foundNameError {
					t.Error("Expected name validation error")
				}
			},
		},
		{
			name: "valid contact with only first_name",
			body: map[string]interface{}{
				"email":      "test2@example.com",
				"first_name": "Jane",
			},
			expectedStatus: http.StatusCreated,
			validateResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				if resp.Code != http.StatusCreated {
					t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.Code)
				}
			},
		},
		{
			name: "valid contact with only last_name",
			body: map[string]interface{}{
				"email":     "test3@example.com",
				"last_name": "Smith",
			},
			expectedStatus: http.StatusCreated,
			validateResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				if resp.Code != http.StatusCreated {
					t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.Code)
				}
			},
		},
		{
			name: "email too long - should fail",
			body: map[string]interface{}{
				"email":      "a" + string(bytes.Repeat([]byte("x"), 250)) + "@example.com",
				"first_name": "John",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var result ValidationErrorResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				foundEmailError := false
				for _, detail := range result.Details {
					if detail.Field == "email" && detail.Message == "email must not exceed 255 characters" {
						foundEmailError = true
						break
					}
				}
				if !foundEmailError {
					t.Errorf("Expected email length validation error, got: %v", result.Details)
				}
			},
		},
		{
			name: "invalid phone format - should fail",
			body: map[string]interface{}{
				"email":      "test@example.com",
				"first_name": "John",
				"phone":      "123",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var result ValidationErrorResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				foundPhoneError := false
				for _, detail := range result.Details {
					if detail.Field == "phone" {
						foundPhoneError = true
						break
					}
				}
				if !foundPhoneError {
					t.Error("Expected phone validation error")
				}
			},
		},
		{
			name:           "invalid JSON - should fail",
			body:           "invalid json string",
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				// Should return "invalid json" error
				if resp.Body.String() != "invalid json\n" {
					t.Errorf("Expected 'invalid json' error, got: %s", resp.Body.String())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			var err error

			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/contacts", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.Create(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			if tt.validateResponse != nil {
				tt.validateResponse(t, rr)
			}
		})
	}
}
