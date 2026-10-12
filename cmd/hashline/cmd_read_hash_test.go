package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"code.andersonlizarazo.com/hashline/pkg/database"
)

func Test_cmdReadHash(t *testing.T) {
	testTokens := []string{"t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8", "t9", "t10", "t11", "t12", "t13", "t14", "t15", "t16"}
	db, err := database.New(t.Context(), database.InMemory, testTokens)
	if err != nil {
		t.Fatal(err)
	}

	testFilePath := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(testFilePath, []byte("A\nB\nC\nD\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var buffer bytes.Buffer
	if err := cmdReadHash(t.Context(), db, &buffer, testFilePath); err != nil {
		t.Fatal(err)
	}

	gotOutput := buffer.String()
	wantOutput := "t1\tA\nt2\tB\nt3\tC\nt4\tD\n"

	if gotOutput != wantOutput {
		t.Errorf("cmdReadHash(%q) output = %q, want %q", testFilePath, gotOutput, wantOutput)
	}
}
