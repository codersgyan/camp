package contact

import "time"

type Tag struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Contact struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []Tag     `json:"tags"`
}



/**

{
  "first_name": "Krishna",
  "last_name": "Bansal",
  "email": "ava@example.com",
  "phone": "+1-555-0102",
  "tags": [
    { "text": "customers" },
    { "text": "newsletter" }
  ]
}

use this data 
  */