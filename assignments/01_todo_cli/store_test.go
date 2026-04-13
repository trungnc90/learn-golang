// How to run tests:
//   go test -v -run TestFileStore ./...    Run FileStore tests only
//   go test -v -run TestMemoryStore ./...  Run MemoryStore tests only
//   go test -v -run TestStore ./...        Run all store tests

package todo

import (
	"os"
	"testing"
)

// --- Shared test logic for any Store implementation ---

func testStore_SaveAndLoad(t *testing.T, store Store) {
	t.Helper()

	tasks := []Task{
		{Id: 1, Title: "Task A", Priority: "high"},
		{Id: 2, Title: "Task B", Priority: "low"},
	}

	if err := store.Save(tasks); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(loaded))
	}
	if loaded[0].Title != "Task A" || loaded[1].Title != "Task B" {
		t.Fatal("task titles don't match")
	}
}

func testStore_LoadEmpty(t *testing.T, store Store) {
	t.Helper()

	tasks, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}
}

func testStore_Overwrite(t *testing.T, store Store) {
	t.Helper()

	store.Save([]Task{{Id: 1, Title: "Old"}})
	store.Save([]Task{{Id: 1, Title: "New"}, {Id: 2, Title: "Extra"}})

	tasks, _ := store.Load()
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Title != "New" {
		t.Fatalf("expected 'New', got '%s'", tasks[0].Title)
	}
}

func testStore_SaveEmpty(t *testing.T, store Store) {
	t.Helper()

	store.Save([]Task{{Id: 1, Title: "Task"}})
	store.Save([]Task{})

	tasks, _ := store.Load()
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks after saving empty, got %d", len(tasks))
	}
}

func testStore_LoadReturnsCopy(t *testing.T, store Store) {
	t.Helper()

	store.Save([]Task{{Id: 1, Title: "Original"}})

	// Mutate the loaded slice
	loaded, _ := store.Load()
	loaded[0].Title = "Mutated"

	// Load again — should still be "Original"
	fresh, _ := store.Load()
	if fresh[0].Title != "Original" {
		t.Fatalf("expected 'Original', got '%s' — Load did not return a copy", fresh[0].Title)
	}
}

// --- MemoryStore tests ---

func TestMemoryStore_SaveAndLoad(t *testing.T) {
	testStore_SaveAndLoad(t, NewMemoryStore())
}

func TestMemoryStore_LoadEmpty(t *testing.T) {
	testStore_LoadEmpty(t, NewMemoryStore())
}

func TestMemoryStore_Overwrite(t *testing.T) {
	testStore_Overwrite(t, NewMemoryStore())
}

func TestMemoryStore_SaveEmpty(t *testing.T) {
	testStore_SaveEmpty(t, NewMemoryStore())
}

func TestMemoryStore_LoadReturnsCopy(t *testing.T) {
	testStore_LoadReturnsCopy(t, NewMemoryStore())
}

// --- FileStore tests ---

func newTestFileStore(t *testing.T) *FileStore {
	t.Helper()
	path := "test_tasks.json"
	os.Remove(path)
	t.Cleanup(func() { os.Remove(path) })
	return NewFileStore(path)
}

func TestFileStore_SaveAndLoad(t *testing.T) {
	testStore_SaveAndLoad(t, newTestFileStore(t))
}

func TestFileStore_LoadEmpty(t *testing.T) {
	testStore_LoadEmpty(t, newTestFileStore(t))
}

func TestFileStore_Overwrite(t *testing.T) {
	testStore_Overwrite(t, newTestFileStore(t))
}

func TestFileStore_SaveEmpty(t *testing.T) {
	testStore_SaveEmpty(t, newTestFileStore(t))
}

func TestFileStore_LoadReturnsCopy(t *testing.T) {
	testStore_LoadReturnsCopy(t, newTestFileStore(t))
}

func TestFileStore_LoadCorruptedFile(t *testing.T) {
	path := "test_corrupt.json"
	os.WriteFile(path, []byte("not valid json"), 0644)
	t.Cleanup(func() { os.Remove(path) })

	store := NewFileStore(path)
	_, err := store.Load()
	if err == nil {
		t.Fatal("expected error for corrupted JSON file")
	}
}
