package contact

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateContactOrUpsertTags(c *Contact) (int64, error) {
	existingContact, err := r.GetByEmail(c.Email)
	if err != nil {
		return 0, err
	}

	if existingContact.Email != "" && len(c.Tags) == 0 {
		return 0, fmt.Errorf("email already exists")
	}

	if existingContact.Email != "" && len(c.Tags) > 0 {
		// only tags upserts
		txn, err := r.db.Begin()
		if err != nil {
			return 0, fmt.Errorf("failed to start the transaction")
		}
		defer txn.Rollback()

		// get tags
		for _, tag := range c.Tags {
			_, err := txn.Exec(`INSERT OR IGNORE INTO tags (text) VALUES (?)`, tag.Text)
			if err != nil {
				return 0, fmt.Errorf("failed to insert tag: %w", err)
			}
		}

		placeholders := make([]string, len(c.Tags))
		args := make([]interface{}, len(c.Tags))

		for i, tag := range c.Tags {
			placeholders[i] = "?"
			args[i] = tag.Text
		}

		query := fmt.Sprintf(`
			SELECT id, text
			FROM tags
			WHERE text IN (%s)`, strings.Join(placeholders, ","))

		rows, err := txn.Query(query, args...)
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, nil
			}
			return 0, fmt.Errorf("failed to select tags: %w", err)
		}

		defer rows.Close()

		var tags []Tag
		for rows.Next() {
			var tag Tag
			if err := rows.Scan(&tag.ID, &tag.Text); err != nil {
				return 0, fmt.Errorf("failed to scan tag: %w", err)
			}
			tags = append(tags, tag)
		}

		if err := rows.Err(); err != nil {
			return 0, fmt.Errorf("failed to scan rows while selecting tags: %w", err)
		}

		// add pivot table records
		valueStrings := make([]string, 0, len(tags))
		valueArgs := make([]interface{}, 0, len(tags)*2)

		for _, tag := range tags {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, existingContact.ID, tag.ID)
		}

		pivotQuery := fmt.Sprintf(`
    	INSERT OR IGNORE INTO contact_tag (contact_id, tag_id)
    	VALUES %s`, strings.Join(valueStrings, ","))

		_, err = txn.Exec(pivotQuery, valueArgs...)
		if err != nil {
			return 0, fmt.Errorf("failed to add record in pivot table: %w", err)
		}

		err = txn.Commit()
		if err != nil {
			return 0, err
		}

		// todo: think what to return 0 is problematic, since we have created tags but zero was for created contact.
		return 0, nil
	}

	// if we reach here means we are creating contact newly.
	query := `
		INSERT INTO contacts (fname, lname, email, phone)
		VALUES(?, ?, ?, ?)
	`
	txn, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer txn.Rollback()

	createdContact, err := txn.Exec(query, c.FirstName, c.LastName, c.Email, c.Phone)
	if err != nil {
		return 0, fmt.Errorf("failed to create a contact: %w", err)
	}
	lastId, err := createdContact.LastInsertId()
	if err != nil {
		return 0, err
	}

	// if tags exists create tags
	if len(c.Tags) > 0 {
		// get tags
		for _, tag := range c.Tags {
			_, err := txn.Exec(`INSERT OR IGNORE INTO tags (text) VALUES (?)`, tag.Text)
			if err != nil {
				return 0, fmt.Errorf("failed to insert tag: %w", err)
			}
		}

		placeholders := make([]string, len(c.Tags))
		args := make([]interface{}, len(c.Tags))

		for i, tag := range c.Tags {
			placeholders[i] = "?"
			args[i] = tag.Text
		}

		query := fmt.Sprintf(`
			SELECT id, text
			FROM tags
			WHERE text IN (%s)`, strings.Join(placeholders, ","))

		rows, err := txn.Query(query, args...)
		if err != nil {
			return 0, err
		}

		defer rows.Close()

		var tags []Tag
		for rows.Next() {
			var tag Tag
			if err := rows.Scan(&tag.ID, &tag.Text); err != nil {
				return 0, err
			}
			tags = append(tags, tag)
		}

		if err := rows.Err(); err != nil {
			return 0, err
		}

		// add pivot table records
		valueStrings := make([]string, 0, len(tags))
		valueArgs := make([]interface{}, 0, len(tags)*2)

		for _, tag := range tags {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, lastId, tag.ID)
		}

		pivotQuery := fmt.Sprintf(`
    	INSERT OR IGNORE INTO contact_tag (contact_id, tag_id)
    	VALUES %s`, strings.Join(valueStrings, ","))

		_, err = txn.Exec(pivotQuery, valueArgs...)
		if err != nil {
			return 0, fmt.Errorf("failed to create record in contaCt tag pivot table: %w", err)
		}
	}

	err = txn.Commit()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (r *Repository) GetByEmail(email string) (Contact, error) {
	query := `
		SELECT
			id,
			fname,
			lname,
			email,
			phone,
			created_at,
			updated_at
		FROM contacts
		WHERE email = ?
		LIMIT 1
	`
	var result Contact

	err := r.db.QueryRow(query, email).Scan(
		&result.ID,
		&result.FirstName,
		&result.LastName,
		&result.Email,
		&result.Phone,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Contact{}, nil
		}

		return Contact{}, err
	}

	return result, nil
}
