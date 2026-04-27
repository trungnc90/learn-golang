// How to run tests:
//   go test -v ./...                    Run all tests
//   go test -run TestTokenize           Run Tokenize tests only
//   go test -run TestParseCommand       Run ParseCommand tests only

package todo

import (
	"testing"
)

// --- Tokenize ---

func TestTokenize_Simple(t *testing.T) {
	tokens := Tokenize("add task1 --priority high")
	expected := []string{"add", "task1", "--priority", "high"}
	assertTokens(t, tokens, expected)
}

func TestTokenize_QuotedStrings(t *testing.T) {
	tokens := Tokenize(`add "Buy groceries" --desc "Milk and eggs"`)
	expected := []string{"add", "Buy groceries", "--desc", "Milk and eggs"}
	assertTokens(t, tokens, expected)
}

func TestTokenize_Empty(t *testing.T) {
	tokens := Tokenize("")
	if len(tokens) != 0 {
		t.Fatalf("expected 0 tokens, got %d", len(tokens))
	}
}

func TestTokenize_ExtraSpaces(t *testing.T) {
	tokens := Tokenize("  add   task1  ")
	expected := []string{"add", "task1"}
	assertTokens(t, tokens, expected)
}

func assertTokens(t *testing.T, got, expected []string) {
	t.Helper()
	if len(got) != len(expected) {
		t.Fatalf("expected %d tokens, got %d: %v", len(expected), len(got), got)
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("token[%d]: expected %q, got %q", i, expected[i], got[i])
		}
	}
}

// --- ParseCommand: add ---

func TestParseCommand_Add(t *testing.T) {
	cmd, err := ParseCommand(`add "Buy groceries" --desc "Milk, eggs" --priority high`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Add == nil {
		t.Fatal("expected Add command")
	}
	if cmd.Add.Title != "Buy groceries" {
		t.Fatalf("expected title 'Buy groceries', got %q", cmd.Add.Title)
	}
	if cmd.Add.Description != "Milk, eggs" {
		t.Fatalf("expected desc 'Milk, eggs', got %q", cmd.Add.Description)
	}
	if cmd.Add.Priority != "high" {
		t.Fatalf("expected priority 'high', got %q", cmd.Add.Priority)
	}
}

func TestParseCommand_Add_MissingTitle(t *testing.T) {
	_, err := ParseCommand("add")
	if err == nil {
		t.Fatal("expected error for missing title")
	}
}

func TestParseCommand_Add_DefaultFlags(t *testing.T) {
	cmd, err := ParseCommand("add task1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Add.Description != "" {
		t.Fatalf("expected empty description, got %q", cmd.Add.Description)
	}
	if cmd.Add.Priority != "" {
		t.Fatalf("expected empty priority, got %q", cmd.Add.Priority)
	}
}

// --- ParseCommand: list ---

func TestParseCommand_List_NoFilter(t *testing.T) {
	cmd, err := ParseCommand("list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.List == nil {
		t.Fatal("expected List command")
	}
	if cmd.List.Filter != "" {
		t.Fatalf("expected empty filter, got %q", cmd.List.Filter)
	}
}

func TestParseCommand_List_WithFilter(t *testing.T) {
	cmd, err := ParseCommand("list done")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.List.Filter != "done" {
		t.Fatalf("expected filter 'done', got %q", cmd.List.Filter)
	}
}

// --- ParseCommand: delete ---

func TestParseCommand_Delete(t *testing.T) {
	cmd, err := ParseCommand("delete 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Delete == nil {
		t.Fatal("expected Delete command")
	}
	if cmd.Delete.Id != 3 {
		t.Fatalf("expected id 3, got %d", cmd.Delete.Id)
	}
}

func TestParseCommand_Delete_InvalidId(t *testing.T) {
	_, err := ParseCommand("delete abc")
	if err == nil {
		t.Fatal("expected error for invalid id")
	}
}

func TestParseCommand_Delete_MissingId(t *testing.T) {
	_, err := ParseCommand("delete")
	if err == nil {
		t.Fatal("expected error for missing id")
	}
}

// --- ParseCommand: update ---

func TestParseCommand_Update(t *testing.T) {
	cmd, err := ParseCommand(`update 1 --title "New title" --desc "New desc" --priority low`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Update == nil {
		t.Fatal("expected Update command")
	}
	if cmd.Update.Id != 1 {
		t.Fatalf("expected id 1, got %d", cmd.Update.Id)
	}
	if cmd.Update.Title != "New title" {
		t.Fatalf("expected title 'New title', got %q", cmd.Update.Title)
	}
	if cmd.Update.Description != "New desc" {
		t.Fatalf("expected desc 'New desc', got %q", cmd.Update.Description)
	}
	if cmd.Update.Priority != "low" {
		t.Fatalf("expected priority 'low', got %q", cmd.Update.Priority)
	}
}

func TestParseCommand_Update_PartialFlags(t *testing.T) {
	cmd, err := ParseCommand("update 2 --title Updated")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Update.Title != "Updated" {
		t.Fatalf("expected title 'Updated', got %q", cmd.Update.Title)
	}
	if cmd.Update.Description != "" {
		t.Fatalf("expected empty description, got %q", cmd.Update.Description)
	}
}

// --- ParseCommand: done ---

func TestParseCommand_Done(t *testing.T) {
	cmd, err := ParseCommand("done 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Done == nil {
		t.Fatal("expected Done command")
	}
	if cmd.Done.Id != 5 {
		t.Fatalf("expected id 5, got %d", cmd.Done.Id)
	}
}

func TestParseCommand_Done_InvalidId(t *testing.T) {
	_, err := ParseCommand("done xyz")
	if err == nil {
		t.Fatal("expected error for invalid id")
	}
}

// --- ParseCommand: help, exit, unknown ---

func TestParseCommand_Help(t *testing.T) {
	cmd, err := ParseCommand("help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cmd.Help {
		t.Fatal("expected Help to be true")
	}
}

func TestParseCommand_Exit(t *testing.T) {
	cmd, err := ParseCommand("exit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cmd.Exit {
		t.Fatal("expected Exit to be true")
	}
}

func TestParseCommand_Quit(t *testing.T) {
	cmd, err := ParseCommand("quit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cmd.Exit {
		t.Fatal("expected Exit to be true")
	}
}

func TestParseCommand_Unknown(t *testing.T) {
	_, err := ParseCommand("foobar")
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestParseCommand_Empty(t *testing.T) {
	_, err := ParseCommand("")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}
