package contact

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(c *Contact) (int64, error) {
	query := `
		INSERT INTO contacts (fname, lname, email, phone)
		VALUES(?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, c.FirstName, c.LastName, c.Email, c.Phone)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}
