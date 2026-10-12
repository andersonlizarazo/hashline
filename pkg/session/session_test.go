package session_test

import (
	"os"
	"path/filepath"
	"testing"

	"code.andersonlizarazo.com/hashline/pkg/database"
	"code.andersonlizarazo.com/hashline/pkg/session"
)

func TestLineTokens_smoke(t *testing.T) {
	testFile := filepath.Join(t.TempDir(), "test.txt")
	testTokens := []string{"t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8"}

	testdata := "one\ntwo"
	if err := os.WriteFile(testFile, []byte(testdata), 0o644); err != nil {
		t.Fatal(err)
	}

	db, err := database.New(t.Context(), database.InMemory, testTokens)
	if err != nil {
		t.Fatal(err)
	}

	tokens, err := session.LineTokens(t.Context(), db, testFile, 0, 2)
	if err != nil {
		t.Fatal(err)
	}

	got := len(tokens)
	want := 2
	if got != want {
		t.Errorf("session.LineTokens(%q) = %d, want %d", testdata, got, want)
	}

	testdata = "One\nTwo"
	if err := os.WriteFile(testFile, []byte(testdata), 0o644); err != nil {
		t.Fatal(err)
	}

	tokens, err = session.LineTokens(t.Context(), db, testFile, 0, 2)
	if err != nil {
		t.Fatal(err)
	}

	got = len(tokens)
	want = 2
	if got != want {
		t.Errorf("session.LineTokens(%q) = %d, want %d", testdata, got, want)
	}
}
