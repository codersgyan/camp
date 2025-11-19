package contact

import (
	"fmt"
	"testing"
)

func TestContactCreateRequest_Validate(t *testing.T) {
	t.Run("valid minimal request", func(t *testing.T) {
		req := ContactCreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
		}

		errs := req.Validate()
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %d: %+v", len(errs), errs)
		}
	})

	t.Run("valid with phone gets normalized", func(t *testing.T) {
		req := ContactCreateRequest{
			FirstName: "Alice",
			LastName:  "Smith",
			Email:     "alice@test.com",
			Phone:     "  +1 (415) 555-2671  ",
		}

		errs := req.Validate()
		if len(errs) != 0 {
			t.Fatalf("valid phone should pass, got errors: %+v", errs)
		}
		if req.Phone != "+14155552671" {
			t.Errorf("phone not normalized → got %q, want +14155552671", req.Phone)
		}
	})

	t.Run("missing first_name", func(t *testing.T) {
		req := ContactCreateRequest{LastName: "Doe", Email: "x@x.com"}
		errs := req.Validate()
		if len(errs) != 1 || errs[0].Field != "first_name" {
			t.Fatalf("expected first_name required error, got %+v", errs)
		}
	})

	t.Run("missing last_name", func(t *testing.T) {
		req := ContactCreateRequest{FirstName: "John", Email: "x@x.com"}
		errs := req.Validate()
		if len(errs) != 1 || errs[0].Field != "last_name" {
			t.Fatalf("expected first_name required error, got %+v", errs)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		req := ContactCreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "this-is-not-an-email",
		}
		errs := req.Validate()
		if len(errs) != 1 || errs[0].Field != "email" {
			t.Fatalf("expected email error, got %+v", errs)
		}
	})

	t.Run("invalid phone", func(t *testing.T) {
		req := ContactCreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "ok@x.com",
			Phone:     "12345",
		}
		errs := req.Validate()
		if len(errs) != 1 || errs[0].Field != "phone" {
			t.Fatalf("expected phone error, got %+v", errs)
		}
	})

	t.Run("first_name too long uses dynamic message", func(t *testing.T) {
		req := ContactCreateRequest{
			FirstName: stringWithLength(NameMaxLength + 1),
			LastName:  "Doe",
			Email:     "x@x.com",
		}
		errs := req.Validate()
		if len(errs) != 1 || errs[0].Message != fmt.Sprintf("max %d characters", NameMaxLength) {
			t.Fatalf("wrong dynamic message: %+v", errs[0])
		}
	})

	t.Run("last_name too long uses dynamic message", func(t *testing.T) {
		req := ContactCreateRequest{
			FirstName: "John",
			LastName:  stringWithLength(NameMaxLength + 1),
			Email:     "x@x.com",
		}
		errs := req.Validate()
		if len(errs) != 1 || errs[0].Message != fmt.Sprintf("max %d characters", NameMaxLength) {
			t.Fatalf("wrong dynamic message: %+v", errs[0])
		}
	})

	t.Run("valid request with tags", func(t *testing.T) {
		req := ContactCreateRequest{
			FirstName: "Rakesh",
			LastName:  "K",
			Email:     "rakesh@codersgyan.com",
			Tags: []struct {
				Text string `json:"text"`
			}{
				{Text: "purchased:mern"},
				{Text: "purchased:devops"},
			},
		}
		errs := req.Validate()

		if len(errs) != 0 {
			t.Fatalf("expected no errors with valid tags: %+v", errs)
		}
	})

}

func stringWithLength(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}
