package database_test

import (
	"reflect"
	"testing"

	"code.andersonlizarazo.com/hashline/pkg/database"
)

func TestMintTokens(t *testing.T) {
	testTokens := []string{"t1", "t2", "t3", "t4"}
	db, err := database.New(t.Context(), database.InMemory, testTokens)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		inAmount  int64
		wantToken []string
	}{
		{
			inAmount:  1,
			wantToken: []string{"t1"},
		},
		{
			inAmount:  2,
			wantToken: []string{"t2", "t3"},
		},
		{
			inAmount:  3,
			wantToken: []string{"t4", "t1", "t2"},
		},
		{
			inAmount:  4,
			wantToken: []string{"t3", "t4", "t1", "t2"},
		},
	}

	for _, tt := range tests {
		// Isn't wrapped in [t.Run] because it depends on the state of the previous run.
		gotToken, err := db.MintTokens(t.Context(), tt.inAmount)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(gotToken, tt.wantToken) {
			t.Errorf("MintTokens(%d) = %v, want %v", tt.inAmount, gotToken, tt.wantToken)
		}
	}
}

func TestWriteReadTokens(t *testing.T) {
	testTokens := []string{"t1", "t2", "t3", "t4"}
	db, err := database.New(t.Context(), database.InMemory, testTokens)
	if err != nil {
		t.Fatal(err)
	}

	testFilename := "file.txt"
	testFileHash := uint64(1)
	testFileTokens := []string{"t1", "t2"}
	testFilePos := int64(0)
	testFileSpan := int64(2)

	if err := db.WriteTokens(t.Context(), testFilename, testFileHash, testFileTokens, testFilePos); err != nil {
		t.Fatal(err)
	}

	gotFileHash, gotFileTokens, err := db.ReadTokens(t.Context(), testFilename, testFilePos, testFileSpan)
	if err != nil {
		t.Fatal(err)
	}

	if gotFileHash != testFileHash {
		t.Errorf("ReadTokens(%q, %d, %d) filehash = %v, want %v", testFilename, testFilePos, testFileSpan, gotFileHash, testFileHash)
	}

	if !reflect.DeepEqual(gotFileTokens, testFileTokens) {
		t.Errorf("ReadTokens(%q, %d, %d) tokens = %v, want %v", testFilename, testFilePos, testFileSpan, gotFileTokens, testFileTokens)
	}
}
