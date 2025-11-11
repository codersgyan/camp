package contact

import "time"

type Tag struct {
	ID        int64
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Contact struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Tags      []Tag
}
