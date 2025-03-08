package stone

import (
	"os"
	"testing"
)

func TestStore(t *testing.T) {
	path := "test.db"
	os.Remove(path) // Clean up before test

	// Create a new store
	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Test Set and Get
	err = store.Set([]byte("key1"), []byte("value1"))
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}
	value, err := store.Get([]byte("key1"))
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if string(value) != "value1" {
		t.Errorf("expected 'value1', got '%s'", value)
	}

	// Test update
	err = store.Set([]byte("key1"), []byte("value2"))
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	value, err = store.Get([]byte("key1"))
	if err != nil {
		t.Fatalf("get after update failed: %v", err)
	}
	if string(value) != "value2" {
		t.Errorf("expected 'value2', got '%s'", value)
	}

	// Test Delete
	err = store.Delete([]byte("key1"))
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	_, err = store.Get([]byte("key1"))
	if err == nil {
		t.Error("expected error on get after delete, got nil")
	}

	// Test non-existent key
	_, err = store.Get([]byte("key2"))
	if err == nil {
		t.Error("expected error on get for non-existent key, got nil")
	}
}

func TestPersistence(t *testing.T) {
	path := "test.db"
	os.Remove(path) // Clean up before test

	// Create store and set a value
	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	err = store.Set([]byte("key1"), []byte("value1"))
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}
	store.Close()

	// Reopen store and verify persistence
	store, err = NewStore(path)
	if err != nil {
		t.Fatalf("failed to reopen store: %v", err)
	}
	defer store.Close()

	value, err := store.Get([]byte("key1"))
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if string(value) != "value1" {
		t.Errorf("expected 'value1', got '%s'", value)
	}
}