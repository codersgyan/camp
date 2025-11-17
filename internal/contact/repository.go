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
	// check if email already exists
	existingContact, err := r.GetByEmail(c.Email)
	if err != nil {
		return 0, err
	}

	if existingContact != nil {

		if len(c.Tags) == 0 {
			return 0, fmt.Errorf("email already exists")
		}

		// only tags upserts
		txn, err := r.db.Begin()
		if err != nil {
			return 0, fmt.Errorf("failed to start the transaction: %w", err)
		}
		defer txn.Rollback()

		err = insertTagIfNotExist(txn, c.Tags)
		if err != nil {
			return 0, err
		}

		// get tags
		tags, err := getTagsByTexts(txn, c.Tags)
		if err != nil {
			return 0, err
		}

		// add pivot table records
		if err := linkTagsToContact(txn, existingContact.ID, tags); err != nil {
			return 0, err
		}

		if err := txn.Commit(); err != nil {
			return 0, err
		}

		return existingContact.ID, nil

	}

	// create contact
	txn, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer txn.Rollback()

	lastId, err := createContact(txn, c)
	if err != nil {
		return 0, err
	}

	// if tags exists create tags
	if len(c.Tags) > 0 {
		err = insertTagIfNotExist(txn, c.Tags)
		if err != nil {
			return 0, err
		}

		tags, err := getTagsByTexts(txn, c.Tags)
		if err != nil {
			return 0, err
		}

		if err := linkTagsToContact(txn, lastId, tags); err != nil {
			return 0, err
		}
	}

	if err := txn.Commit(); err != nil {
		return 0, err
	}

	return lastId, nil
}

func (r *Repository) GetByEmail(email string) (*Contact, error) {
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
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

func insertTagIfNotExist(txn *sql.Tx, tags []Tag) error {
	// todo: implement bulk inserts here
	for _, tag := range tags {
		_, err := txn.Exec(`INSERT OR IGNORE INTO tags (text) VALUES (?)`, tag.Text)
		if err != nil {
			return fmt.Errorf("failed to insert tag: %w", err)
		}
	}

	return nil
}

func getTagsByTexts(txn *sql.Tx, tags []Tag) ([]Tag, error) {

	placeholders := make([]string, len(tags))
	args := make([]interface{}, len(tags))

	for i, tag := range tags {
		placeholders[i] = "?"
		args[i] = tag.Text
	}

	query := fmt.Sprintf(`
			SELECT id, text
			FROM tags
			WHERE text IN (%s)`, strings.Join(placeholders, ","))

	rows, err := txn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select tags: %w", err)
	}

	defer rows.Close()

	var resultTags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Text); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		resultTags = append(resultTags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan rows while selecting tags: %w", err)
	}

	return resultTags, nil
}

func linkTagsToContact(txn *sql.Tx, contactID int64, tags []Tag) error {

	valueStrings := make([]string, 0, len(tags))
	valueArgs := make([]interface{}, 0, len(tags)*2)

	for _, tag := range tags {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, contactID, tag.ID)
	}

	query := fmt.Sprintf(`
    	INSERT OR IGNORE INTO contact_tag (contact_id, tag_id)
    	VALUES %s`, strings.Join(valueStrings, ","))

	_, err := txn.Exec(query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to add record in pivot table: %w", err)
	}

	return nil
}

func createContact(txn *sql.Tx, c *Contact) (int64, error) {
	query := `
		INSERT INTO contacts (fname, lname, email, phone)
		VALUES(?, ?, ?, ?)
	`

	result, err := txn.Exec(query, c.FirstName, c.LastName, c.Email, c.Phone)
	if err != nil {
		return 0, fmt.Errorf("failed to create a contact: %w", err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (r *Repository) GetByID(id int64) (*Contact, error) {
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
		WHERE id = ?
		LIMIT 1
`

	var contact Contact
	err := r.db.QueryRow(query, id).Scan(
		&contact.ID,
		&contact.FirstName,
		&contact.LastName,
		&contact.Email,
		&contact.Phone,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get contact: %w", err)
	}

	tagsQuery := `
		SELECT
			t.id,
			t.text,
			t.created_at,
			t.updated_at
		FROM tags t
		INNER JOIN contact_tag ct
		ON t.id = ct.tag_id
		WHERE ct.contact_id = ?
`

	rows, err := r.db.Query(tagsQuery, id)

	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Text, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tags: %w", err)
	}

	contact.Tags = tags

	return &contact, nil
}
