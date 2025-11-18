package contact

import (
	"database/sql"
	"os"
	"testing"

	"github.com/codersgyan/camp/internal/database"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

func TestContactRepositoryCreate(t *testing.T) {
	testDB = setupDB(t)
	repo := NewRepository(testDB)

	contact := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
	}

	_, err := repo.CreateContactOrUpsertTags(contact)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	var count int
	testDB.QueryRow("SELECT COUNT(*) FROM contacts WHERE email = $1",
		contact.Email).Scan(&count)

	if count != 1 {
		t.Errorf("Expected 1 contact, got %d", count)
	}
}

func TestContactRepositoryCreateWithTags(t *testing.T) {
	testDB = setupDB(t)
	repo := NewRepository(testDB)

	contact := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
		Tags:      []Tag{{Text: "purchase:golang"}, {Text: "subscribed:platform"}},
	}

	_, err := repo.CreateContactOrUpsertTags(contact)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	var count int
	testDB.QueryRow("SELECT COUNT(*) FROM contacts WHERE email = $1",
		contact.Email).Scan(&count)

	if count != 1 {
		t.Errorf("Expected 1 contact, got %d", count)
	}

	var tagsCount int
	testDB.QueryRow("SELECT COUNT(*) FROM contact_tag").Scan(&tagsCount)

	if tagsCount != 2 {
		t.Errorf("Expected 2 tags, got %d", tagsCount)
	}
}

func TestContactRepositoryUpsertWithTagsIfEmailExists(t *testing.T) {
	testDB = setupDB(t)
	repo := NewRepository(testDB)

	contact1 := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
		Tags:      []Tag{{Text: "purchase:golang"}, {Text: "subscribed:platform"}},
	}

	_, err := repo.CreateContactOrUpsertTags(contact1)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	contact2 := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
		Tags:      []Tag{{Text: "joined:annual"}},
	}

	_, err = repo.CreateContactOrUpsertTags(contact2)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	var count int
	testDB.QueryRow("SELECT COUNT(*) FROM contacts WHERE email = $1",
		contact1.Email).Scan(&count)

	if count != 1 {
		t.Errorf("Expected 1 contact, got %d", count)
	}

	var tagsCount int
	testDB.QueryRow("SELECT COUNT(*) FROM contact_tag").Scan(&tagsCount)

	if tagsCount != 3 {
		t.Errorf("Expected 3 tags, got %d", tagsCount)
	}
}

func TestContactRepositoryThrowErrorIfTagsNotSent(t *testing.T) {
	testDB = setupDB(t)
	repo := NewRepository(testDB)

	contact1 := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
		Tags:      []Tag{{Text: "purchase:golang"}, {Text: "subscribed:platform"}},
	}

	_, err := repo.CreateContactOrUpsertTags(contact1)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	// here we are not sending the tags
	contact2 := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
	}

	_, err = repo.CreateContactOrUpsertTags(contact2)
	if err == nil {
		t.Fatalf("Expected error when tags are not sent, but got nil")
	}

	var count int
	testDB.QueryRow("SELECT COUNT(*) FROM contacts WHERE email = $1",
		contact1.Email).Scan(&count)

	if count != 1 {
		t.Errorf("Expected 1 contact, got %d", count)
	}

	var tagsCount int

	testDB.QueryRow("SELECT COUNT(*) FROM contact_tag").Scan(&tagsCount)

	if tagsCount != 2 {
		t.Errorf("Expected 2 tags, got %d", tagsCount)
	}
}

func TestContactRepositoryGetByID(t *testing.T) {
	testDB = setupDB(t)
	repo := NewRepository(testDB)

	t.Run("get existing contact with all fields", func(t *testing.T) {
		contact := &Contact{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@codersgyan.com",
			Tags:      []Tag{{Text: "Tag1"}, {Text: "Tag2"}},
			Phone:     "000-1234",
		}

		id, err := repo.CreateContactOrUpsertTags(contact)
		if err != nil {
			t.Fatalf("Failed to create contact: %v", err)
		}

		retrievedContact, err := repo.GetContactByID(id)
		if err != nil {
			t.Fatalf("Failed to get contact by ID: %v", err)
		}

		if retrievedContact == nil {
			t.Fatal("Expected contact, got nil")
		}

		if retrievedContact.ID != id {
			t.Errorf("Expected ID %d, got %d", id, retrievedContact.ID)
		}
		if retrievedContact.FirstName != contact.FirstName {
			t.Errorf("Expected FirstName '%s', got '%s'", contact.FirstName, retrievedContact.FirstName)
		}
		if retrievedContact.LastName != contact.LastName {
			t.Errorf("Expected LastName '%s', got '%s'", contact.LastName, retrievedContact.LastName)
		}
		if retrievedContact.Email != contact.Email {
			t.Errorf("Expected Email '%s', got '%s'", contact.Email, retrievedContact.Email)
		}
		if retrievedContact.Phone != contact.Phone {
			t.Errorf("Expected Phone '%s', got '%s'", contact.Phone, retrievedContact.Phone)
		}
		if len(retrievedContact.Tags) != len(contact.Tags) {
			t.Fatalf("Expected %d tags, got %d", len(contact.Tags), len(retrievedContact.Tags))
		}

		tagTexts := make(map[string]bool)
		for _, tag := range retrievedContact.Tags {
			tagTexts[tag.Text] = true
		}
		for _, expectedTag := range contact.Tags {
			if !tagTexts[expectedTag.Text] {
				t.Errorf("Expected tag '%s' not found in retrieved tags", expectedTag.Text)
			}
		}
	})

	t.Run("non-existing contact with negative ID", func(t *testing.T) {
		notFound, err := repo.GetContactByID(-1)
		if err != nil {
			t.Fatalf("Expected no error for non-existing contact, got: %v", err)
		}
		if notFound != nil {
			t.Error("Expected nil for non-existing contact")
		}
	})

	t.Run("non-existing contact with zero ID", func(t *testing.T) {
		notFound, err := repo.GetContactByID(0)
		if err != nil {
			t.Fatalf("Expected no error for non-existing contact, got: %v", err)
		}
		if notFound != nil {
			t.Error("Expected nil for non-existing contact")
		}
	})

	t.Run("non-existing contact with large ID", func(t *testing.T) {
		notFound, err := repo.GetContactByID(999999)
		if err != nil {
			t.Fatalf("Expected no error for non-existing contact, got: %v", err)
		}
		if notFound != nil {
			t.Error("Expected nil for non-existing contact")
		}
	})
}

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err := database.RunMigration(db); err != nil {
		t.Fatal(err)
	}

	return db
}
