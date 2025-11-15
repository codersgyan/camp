package contact

import "time"

type Tag struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text" validate:"required,min=1,max=255"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Contact struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string    `json:"last_name" validate:"required,min=1,max=100"`
	Email     string    `json:"email" validate:"required,email"`
	Phone     string    `json:"phone" validate:"omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []Tag     `json:"tags" validate:"omitempty,dive"`
}
