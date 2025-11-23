package contact

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	// allowed max lengths for fields
	emailMaxLength = 255
	// used for both first name and last name
	NameMaxLength = 100
)

func (r *ContactCreateRequest) Validate() []ValidationError {
	var errs []ValidationError

	r.trimFields()

	errs = append(errs, r.validateFirstName()...)
	errs = append(errs, r.validateLastName()...)
	errs = append(errs, r.validateEmail()...)
	errs = append(errs, r.validatePhone()...)

	return errs
}

// Private helper methods
func (r *ContactCreateRequest) trimFields() {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.TrimSpace(r.Email)
	r.Phone = strings.TrimSpace(r.Phone)
}

func (r *ContactCreateRequest) validateFirstName() []ValidationError {
	if r.FirstName == "" {
		return []ValidationError{{Field: "first_name", Message: "required"}}
	}
	if len(r.FirstName) > NameMaxLength {
		return []ValidationError{{Field: "first_name", Message: fmt.Sprintf("max %d characters", NameMaxLength)}}
	}
	return nil
}

func (r *ContactCreateRequest) validateLastName() []ValidationError {
	if r.LastName == "" {
		return []ValidationError{{Field: "last_name", Message: "required"}}
	}
	if len(r.LastName) > NameMaxLength {
		return []ValidationError{{Field: "last_name", Message: fmt.Sprintf("max %d characters", NameMaxLength)}}
	}
	return nil
}

func (r *ContactCreateRequest) validateEmail() []ValidationError {
	if r.Email == "" {
		return []ValidationError{{Field: "email", Message: "required"}}
	}
	if len(r.Email) > emailMaxLength || !emailRegex.MatchString(r.Email) {
		return []ValidationError{{Field: "email", Message: "invalid email format"}}
	}
	return nil
}

func (r *ContactCreateRequest) validatePhone() []ValidationError {
	if r.Phone == "" {
		// optional field
		return nil
	}
	r.Phone = strings.TrimSpace(r.Phone)

	if !phoneRegex.MatchString(r.Phone) {
		return []ValidationError{{Field: "phone", Message: "must be valid E.164 format (e.g. +919876543210)"}}
	}

	return nil
}
