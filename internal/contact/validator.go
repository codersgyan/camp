package contact

import (
	"regexp"
	"strings"
)

const (
	maxEmailLength = 255
	maxNameLength  = 100
)

// ContactValidationError represents a validation error for a specific field
type ContactValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse represents the error response structure
type ValidationErrorResponse struct {
	Error   string                   `json:"error"`
	Details []ContactValidationError `json:"details"`
}

var (
	// emailRegex is a simple regex for email validation
	// This matches most common email formats but is not fully RFC 5322 compliant
	// For production, consider using a more robust library
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// phoneRegex matches common phone number formats
	// Supports: +1234567890, (123) 456-7890, 123-456-7890, 1234567890
	phoneRegex = regexp.MustCompile(`^[\+]?[(]?[0-9]{1,4}[)]?[-\s\.]?[(]?[0-9]{1,4}[)]?[-\s\.]?[0-9]{1,9}$`)
)

// isValidEmail checks if the email format is valid
func isValidEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	return emailRegex.MatchString(email)
}

// isValidPhone checks if the phone number format is valid (optional field)
func isValidPhone(phone string) bool {
	if len(phone) == 0 {
		return true // Phone is optional
	}
	// Remove common separators for validation
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	// Check if it's a reasonable length (at least 7 digits, max 15)
	if len(cleaned) < 7 || len(cleaned) > 15 {
		return false
	}

	// Check if remaining characters are digits or +
	hasDigit := false
	for _, char := range cleaned {
		if char >= '0' && char <= '9' {
			hasDigit = true
		} else if char != '+' {
			return false
		}
	}
	return hasDigit
}

// ValidateCreateContactRequest validates a contact creation request
func ValidateCreateContactRequest(req *Contact) []ContactValidationError {
	var errors []ContactValidationError

	// Email validation - required
	if req.Email == "" {
		errors = append(errors, ContactValidationError{
			Field:   "email",
			Message: "email is required",
		})
	} else if !isValidEmail(req.Email) {
		errors = append(errors, ContactValidationError{
			Field:   "email",
			Message: "email must be valid format",
		})
	} else if len(req.Email) > maxEmailLength {
		errors = append(errors, ContactValidationError{
			Field:   "email",
			Message: "email must not exceed 255 characters",
		})
	}

	// Name validation - at least one of first_name or last_name is required
	firstNameTrimmed := strings.TrimSpace(req.FirstName)
	lastNameTrimmed := strings.TrimSpace(req.LastName)

	if firstNameTrimmed == "" && lastNameTrimmed == "" {
		errors = append(errors, ContactValidationError{
			Field:   "name",
			Message: "at least one of first_name or last_name is required",
		})
	} else {
		// Validate first_name length if provided
		if firstNameTrimmed != "" && len(req.FirstName) > maxNameLength {
			errors = append(errors, ContactValidationError{
				Field:   "first_name",
				Message: "first_name must not exceed 100 characters",
			})
		}

		// Validate last_name length if provided
		if lastNameTrimmed != "" && len(req.LastName) > maxNameLength {
			errors = append(errors, ContactValidationError{
				Field:   "last_name",
				Message: "last_name must not exceed 100 characters",
			})
		}
	}

	// Phone validation - optional but must be valid if provided
	if req.Phone != "" && !isValidPhone(req.Phone) {
		errors = append(errors, ContactValidationError{
			Field:   "phone",
			Message: "phone must be valid format",
		})
	}

	return errors
}
