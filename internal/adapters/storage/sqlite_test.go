package storage

import (
	"os"
	"testing"
	"time"

	"bountyos-v8/internal/core"
)

func TestSQLiteStorage(t *testing.T) {
	// Create a temp file for the database
	tmpfile, err := os.CreateTemp("", "testdb-*.sqlite")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up
	dbPath := tmpfile.Name()

	// Initialize storage
	store, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer store.Close()

	// Test Data
	bounty := core.Bounty{
		URL:         "https://example.com/bounty/1",
		Title:       "Test Bounty",
		Platform:    "TEST",
		Reward:      "100 USDC",
		Currency:    "USDC",
		CreatedAt:   time.Now(),
		Score:       85,
		Description: "This is a test bounty",
		Tags:        []string{"test", "urgent"},
		PaymentType: "crypto",
	}

	// 1. Test IsNew (should be true initially)
	isNew, err := store.IsNew(bounty.URL)
	if err != nil {
		t.Errorf("IsNew() error = %v", err)
	}
	if !isNew {
		t.Errorf("IsNew() = false, want true for new bounty")
	}

	// 2. Test Save
	if err := store.Save(bounty); err != nil {
		t.Errorf("Save() error = %v", err)
	}

	// 3. Test IsNew (should be false now)
	isNew, err = store.IsNew(bounty.URL)
	if err != nil {
		t.Errorf("IsNew() error = %v", err)
	}
	if isNew {
		t.Errorf("IsNew() = true, want false for existing bounty")
	}

	// 4. Test GetRecent
	recent, err := store.GetRecent(10)
	if err != nil {
		t.Errorf("GetRecent() error = %v", err)
	}
	if len(recent) != 1 {
		t.Errorf("GetRecent() count = %d, want 1", len(recent))
	}
	if len(recent) > 0 {
		got := recent[0]
		if got.URL != bounty.URL {
			t.Errorf("GetRecent() URL = %s, want %s", got.URL, bounty.URL)
		}
		if got.Title != bounty.Title {
			t.Errorf("GetRecent() Title = %s, want %s", got.Title, bounty.Title)
		}
		// Check tags
		if len(got.Tags) != 2 {
			t.Errorf("GetRecent() Tags count = %d, want 2", len(got.Tags))
		}
	}
}
