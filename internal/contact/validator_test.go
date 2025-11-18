package contact

import (
	"bytes"
	"testing"
)

func TestValidateCreateContactRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     Contact
		wantErr bool
		errMsgs []string
	}{
		{
			name: "valid contact with all fields",
			req: Contact{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "+1234567890",
			},
			wantErr: false,
		},
		{
			name: "valid contact with only first_name",
			req: Contact{
				Email:     "test@example.com",
				FirstName: "John",
			},
			wantErr: false,
		},
		{
			name: "valid contact with only last_name",
			req: Contact{
				Email:    "test@example.com",
				LastName: "Doe",
			},
			wantErr: false,
		},
		{
			name: "valid contact without phone",
			req: Contact{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: false,
		},
		{
			name: "empty email",
			req: Contact{
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: true,
			errMsgs: []string{"email is required"},
		},
		{
			name: "invalid email format - no @",
			req: Contact{
				Email:     "not-an-email",
				FirstName: "John",
			},
			wantErr: true,
			errMsgs: []string{"email must be valid format"},
		},
		{
			name: "invalid email format - no domain",
			req: Contact{
				Email:     "test@",
				FirstName: "John",
			},
			wantErr: true,
			errMsgs: []string{"email must be valid format"},
		},
		{
			name: "invalid email format - no TLD",
			req: Contact{
				Email:     "test@example",
				FirstName: "John",
			},
			wantErr: true,
			errMsgs: []string{"email must be valid format"},
		},
		{
			name: "email too long",
			req: Contact{
				Email:     "a" + string(bytes.Repeat([]byte("x"), 250)) + "@example.com",
				FirstName: "John",
			},
			wantErr: true,
			errMsgs: []string{"email must not exceed 255 characters"},
		},
		{
			name: "no name fields",
			req: Contact{
				Email: "test@example.com",
			},
			wantErr: true,
			errMsgs: []string{"at least one of first_name or last_name is required"},
		},
		{
			name: "empty first_name and last_name",
			req: Contact{
				Email:     "test@example.com",
				FirstName: "",
				LastName:  "",
			},
			wantErr: true,
			errMsgs: []string{"at least one of first_name or last_name is required"},
		},
		{
			name: "whitespace only first_name and last_name",
			req: Contact{
				Email:     "test@example.com",
				FirstName: "   ",
				LastName:  "   ",
			},
			wantErr: true,
			errMsgs: []string{"at least one of first_name or last_name is required"},
		},
		{
			name: "first_name too long",
			req: Contact{
				Email:     "test@example.com",
				FirstName: string(make([]byte, 101)),
				LastName:  "Doe",
			},
			wantErr: true,
			errMsgs: []string{"first_name must not exceed 100 characters"},
		},
		{
			name: "last_name too long",
			req: Contact{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  string(make([]byte, 101)),
			},
			wantErr: true,
			errMsgs: []string{"last_name must not exceed 100 characters"},
		},
		{
			name: "invalid phone format - too short",
			req: Contact{
				Email:     "test@example.com",
				FirstName: "John",
				Phone:     "123",
			},
			wantErr: true,
			errMsgs: []string{"phone must be valid format"},
		},
		{
			name: "invalid phone format - too long",
			req: Contact{
				Email:     "test@example.com",
				FirstName: "John",
				Phone:     string(make([]byte, 20)),
			},
			wantErr: true,
			errMsgs: []string{"phone must be valid format"},
		},
		{
			name: "valid phone formats",
			req: Contact{
				Email:     "test@example.com",
				FirstName: "John",
				Phone:     "+1234567890",
			},
			wantErr: false,
		},
		{
			name: "multiple validation errors",
			req: Contact{
				Email: "",
			},
			wantErr: true,
			errMsgs: []string{"email is required", "at least one of first_name or last_name is required"},
		},
		{
			name: "valid email with subdomain",
			req: Contact{
				Email:     "test@mail.example.com",
				FirstName: "John",
			},
			wantErr: false,
		},
		{
			name: "valid email with special characters",
			req: Contact{
				Email:     "test.user+tag@example.co.uk",
				FirstName: "John",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateCreateContactRequest(&tt.req)

			if tt.wantErr {
				if len(errors) == 0 {
					t.Errorf("ValidateCreateContactRequest() expected errors but got none")
					return
				}

				// Check if expected error messages are present
				errorMessages := make(map[string]bool)
				for _, err := range errors {
					errorMessages[err.Message] = true
				}

				for _, expectedMsg := range tt.errMsgs {
					if !errorMessages[expectedMsg] {
						t.Errorf("ValidateCreateContactRequest() expected error message '%s' but didn't find it. Got errors: %v", expectedMsg, errors)
					}
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("ValidateCreateContactRequest() unexpected errors: %v", errors)
				}
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "test@mail.example.com", true},
		{"valid email with plus", "test+tag@example.com", true},
		{"valid email with dot", "test.user@example.com", true},
		{"empty email", "", false},
		{"no @ symbol", "not-an-email", false},
		{"no domain", "test@", false},
		{"no TLD", "test@example", false},
		{"invalid format", "@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidEmail(tt.email); got != tt.want {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestIsValidPhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		{"empty phone", "", true}, // Phone is optional
		{"valid phone with plus", "+1234567890", true},
		{"valid phone with dashes", "123-456-7890", true},
		{"valid phone with spaces", "123 456 7890", true},
		{"valid phone with parentheses", "(123) 456-7890", true},
		{"valid phone digits only", "1234567890", true},
		{"too short", "123", false},
		{"too long", "12345678901234567890", false},
		{"invalid characters", "abc1234567", false},
		{"valid international", "+441234567890", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidPhone(tt.phone); got != tt.want {
				t.Errorf("isValidPhone(%q) = %v, want %v", tt.phone, got, tt.want)
			}
		})
	}
}
